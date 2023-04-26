package main

import (
	"context"
	"fmt"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"otelCol1/apm"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"log"
	"net/http"
	"time"
)

func main() {

	// increment custom counter.
	//go func() {
	//	for {
	//		cusCounter.Inc()
	//		time.Sleep(2 * time.Second)
	//	}
	//}()

	ctx := context.Background()
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String("test-service"),
		),
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	traceExporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithEndpoint("localhost:4318"),
	)

	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)

	defer func() {
		if err := tracerProvider.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()
	wrappedHandler := otelhttp.NewHandler(http.HandlerFunc(helloHandler), "/hello-operation")

	//http.Handle("/metrics", promhttp.Handler())
	http.Handle("/hello", wrappedHandler)

	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	http.ListenAndServe(":8000", nil)

}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	//span1, b := tracer.SpanFromContext(r.Context())
	//fmt.Println(span1.Context().SpanID(), span1.Context().TraceID(), b)
	//
	//span, ctx := tracer.StartSpanFromContext(r.Context(), "web.request", tracer.ResourceName("/posts"))
	//defer span.Finish()
	//
	//fmt.Println(span.Context().SpanID(), span.Context().TraceID())
	//
	//span, ctx2 := tracer.StartSpanFromContext(ctx, "web.request", tracer.ResourceName("/posts"))
	//defer span.Finish()

	//fmt.Println(span.Context().SpanID(), span.Context().TraceID())

	//var apm apm.APM = apm.Otel{}
	//traceId := apm.GetTraceId(r.Context())
	//spanId := apm.GetSpanId(r.Context())
	//
	//fmt.Println(traceId, spanId)

	//ctx, span := apm.CreateSpanFromContext(r.Context(), "trName", "spName")
	//defer span.End()
	//fun(ctx)
	//fun(r.Context())

	fun(r.Context())

	//ctx, span := trc.SpanFromContext(r.Context()).TracerProvider().Tracer("gofr-context").Start(r.Context(), "dummy-name")
	//defer span.End()

	apm := apm.Otel{}
	ctx, span := apm.CreateSpanFromContext(r.Context(), "er")
	defer apm.CloseSpan(span)
	fun(ctx)

	//log.Println("trace id in logs", r.Context().Value("trcID"))

	//client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	//
	//req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:9001/hit", nil)
	//if err != nil {
	//	fmt.Println("error", err)
	//	return
	//}
	//resp, err := client.Do(req)
	//if err != nil {
	//	fmt.Println("error:  ", err)
	//	return
	//}
	////
	//defer resp.Body.Close()
	//bodyBytes, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	fmt.Println("error from ioutil: ", err)
	//	return
	//}
	//
	//// Convert response body to string
	//bodyString := string(bodyBytes)
	//fmt.Println("API Response as String:\n" + bodyString)

	<-time.After(time.Millisecond * 1)
	//w.WriteHeader(http.StatusInternalServerError) // Waiting for 1ms to simulate workload
	w.Write([]byte("Hy"))

}

func fun(ctx context.Context) {
	apm := apm.Otel{}
	traceId := apm.GetTraceId(ctx)
	spanId := apm.GetSpanId(ctx)
	fmt.Println(traceId, spanId)
}
