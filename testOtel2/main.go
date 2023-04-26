package main

import (
	"context"
	"fmt"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	trc "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	fmt.Println("Hi")
	var tp *trc.TracerProvider
	tp, err := getZipkinExporter()
	if err != nil {
		return
	}

	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	wrapped := otelhttp.NewHandler(http.HandlerFunc(handleIt), "/hit")
	http.Handle("/hit", wrapped)

	http.ListenAndServe(":9001", nil)
}

func handleIt(w http.ResponseWriter, r *http.Request) {
	fun(r.Context())
	ctx, span := trace.SpanFromContext(r.Context()).TracerProvider().Tracer("traceMis").Start(r.Context(), "dummy-name-2")
	defer span.End()
	fun(ctx)

	//span, _ := tracer.StartSpanFromContext(r.Context(), "web.request", tracer.ResourceName("/posts"))
	//defer span.Finish()
	//
	//fmt.Println(span.Context().SpanID(), span.Context().TraceID())

	fmt.Println(r.Header)
	w.Write([]byte("hit inside"))
}

func getZipkinExporter() (*trc.TracerProvider, error) {
	url := "http://localhost:2005/api/v2/spans"

	exporter, err := zipkin.New(url)
	if err != nil {
		return nil, err
	}

	svcName := "otel_test_service_2"
	executable, err := os.Executable()
	if err != nil {
		svcName += ":go"
	} else {
		svcName = fmt.Sprintf("%s:%s", svcName, filepath.Base(executable))
	}

	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(svcName),
		),
	)

	batcher := trc.NewBatchSpanProcessor(exporter)

	tp := trc.NewTracerProvider(trc.WithSampler(trc.AlwaysSample()), trc.WithSpanProcessor(batcher), trc.WithResource(r))

	return tp, nil
}

func fun(ctx context.Context) {
	span := trace.SpanFromContext(ctx)

	if !span.SpanContext().TraceID().IsValid() {
		fmt.Println("not valid")
	}

	fmt.Println(span.SpanContext().TraceID().String(), span.SpanContext().SpanID())
}
