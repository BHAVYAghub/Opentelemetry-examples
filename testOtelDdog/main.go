package main

import (
	"fmt"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	tracer.Start(tracer.WithService("IDK"))
	defer tracer.Stop()

	wrappedHandler := http.HandlerFunc(helloHandler)
	http.Handle("/hello", wrappedHandler)

	http.ListenAndServe(":8000", nil)

}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	span1, b := tracer.SpanFromContext(r.Context())
	fmt.Println("@@1", span1.Context().SpanID(), span1.Context().TraceID(), b)

	span, ctx := tracer.StartSpanFromContext(r.Context(), "web.request", tracer.ResourceName("/posts"))
	defer span.Finish()

	fmt.Println("@@2", span.Context().SpanID(), span.Context().TraceID())

	//span2, b := tracer.SpanFromContext(ctx)
	//fmt.Println("@@3", span2.Context().SpanID(), span2.Context().TraceID(), b)

	span4, ctx := tracer.StartSpanFromContext(ctx, "web.request", tracer.ResourceName("/posts"))
	defer span.Finish()

	fmt.Println("@@4", span4.Context().SpanID(), span4.Context().TraceID())

	client := http.Client{}

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:9002/hit", nil)
	if err != nil {
		fmt.Println("error", err)
		return
	}

	//err = tracer.Inject(span.Context(), tracer.HTTPHeadersCarrier(req.Header))
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}

	fmt.Println(req.Header)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error:  ", err)
		return
	}
	//
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error from ioutil: ", err)
		return
	}

	// Convert response body to string
	bodyString := string(bodyBytes)
	fmt.Println("API Response as String:\n" + bodyString)

	<-time.After(time.Millisecond * 1) // Waiting for 1ms to simulate workload
	w.Write([]byte("Hy"))

}
