package main

import (
	"fmt"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"net/http"
)

func main() {

	tracer.Start(tracer.WithService("IDK-2"))
	defer tracer.Stop()

	fmt.Println("hi")
	wrapped := http.HandlerFunc(handleIt)
	http.Handle("/hit", wrapped)

	http.ListenAndServe(":9002", nil)
}

func handleIt(w http.ResponseWriter, r *http.Request) {
	sctx, err := tracer.Extract(tracer.HTTPHeadersCarrier(r.Header))
	if err != nil {
		if err == tracer.ErrSpanContextNotFound {
			fmt.Println("header not found.")
		} else {
			fmt.Println(err)
			return
		}

	}

	//_, span := trace.SpanFromContext(r.Context()).TracerProvider().Tracer("traceMis").Start(r.Context(), "dummy-name-2")
	//defer span.End()
	//fmt.Println(span.SpanContext().SpanID(), span.SpanContext().TraceID())

	fmt.Println(r.Header)

	span := tracer.StartSpan("web.request", tracer.ChildOf(sctx))
	defer span.Finish()

	fmt.Println(span.Context().SpanID(), span.Context().TraceID())

	w.Write([]byte("hit inside"))
}
