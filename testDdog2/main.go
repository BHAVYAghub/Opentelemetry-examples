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
		//if err == tracer.ErrSpanContextNotFound {
		//	fmt.Println("header not found.")
		//	span = tracer.StartSpan("post.filter", tracer.ChildOf(sctx))
		//	defer span.Finish()
		//	fmt.Println("span created inside: ", span.Context().TraceID(), span.Context().TraceID())
		//
		//} else {
		//	fmt.Println(err)
		//	return
		//}

		// TODO: log the error
	}

	if sctx != nil {
		fmt.Println("in header: ", sctx.SpanID(), sctx.TraceID())
	}
	span := tracer.StartSpan("post.filter", tracer.ChildOf(sctx))
	defer span.Finish()
	fmt.Println("after span start: ", span.Context().SpanID(), span.Context().TraceID())

	ctx := tracer.ContextWithSpan(r.Context(), span)

	fmt.Println("after context start: ", span.Context().SpanID(), span.Context().TraceID())
	//span11, b := tracer.SpanFromContext(r.Context())
	//fmt.Println("r.context ", span11.Context().TraceID(), span11.Context().TraceID(), b)

	//span12, b := tracer.SpanFromContext(ctx)
	//fmt.Println("r.context ctx", span12.Context().TraceID(), span12.Context().TraceID(), b)

	*r = *r.WithContext(ctx)

	//_, span := trace.SpanFromContext(r.Context()).TracerProvider().Tracer("traceMis").Start(r.Context(), "dummy-name-2")
	//defer span.End()
	//fmt.Println(span.SpanContext().SpanID(), span.SpanContext().TraceID())

	span1, b := tracer.SpanFromContext(r.Context())
	fmt.Println(span1.Context().SpanID(), span1.Context().TraceID(), b)

	fmt.Println(r.Header)

	//span := tracer.StartSpan("web.request", tracer.ChildOf(sctx))
	//defer span.Finish()

	//fmt.Println(span.Context().SpanID(), span.Context().TraceID())

	w.Write([]byte("hit inside"))
}
