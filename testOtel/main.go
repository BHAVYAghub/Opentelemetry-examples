package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testOtel/apm"
	"time"

	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"

	//"github.com/prometheus/client_golang/prometheus"
	//"github.com/prometheus/client_golang/prometheus/promauto"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
)

//var (
//	cusCounter = promauto.NewCounter(prometheus.CounterOpts{
//		Name: "cus_counter_total",
//		Help: "The total number of processed custom counter.",
//	})
//)

func main() {

	// increment custom counter.
	//go func() {
	//	for {
	//		cusCounter.Inc()
	//		time.Sleep(2 * time.Second)
	//	}
	//}()

	var tp *trace.TracerProvider
	tp, err := getZipkinExporter()
	if err != nil {
		return
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()
	wrappedHandler := otelhttp.NewHandler(http.HandlerFunc(helloHandler), "/hello-operation")

	//http.Handle("/metrics", promhttp.Handler())
	http.Handle("/hello", wrappedHandler)

	otel.SetTracerProvider(tp)
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

	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:9001/hit", nil)
	if err != nil {
		fmt.Println("error", err)
		return
	}
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

func fun(ctx context.Context) {
	apm := apm.Otel{}
	traceId := apm.GetTraceId(ctx)
	spanId := apm.GetSpanId(ctx)
	fmt.Println(traceId, spanId)
}

func getZipkinExporter() (*trace.TracerProvider, error) {
	url := "http://localhost:2005/api/v2/spans"

	exporter, err := zipkin.New(url)
	if err != nil {
		return nil, err
	}

	svcName := "otel_test_service_1"
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

	if err != nil {
		panic(err)
	}

	batcher := trace.NewBatchSpanProcessor(exporter, trace.BatchSpanProcessorOption())

	tp := trace.NewTracerProvider(trace.WithSampler(), trace.WithSpanProcessor(batcher), trace.WithResource(r))

	return tp, nil
}

func AddTraceId(next http.HandlerFunc) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {

		// clone current context and append request-id

		//parentSpan := trc.SpanFromContext(r.Context())
		//fmt.Println("trc", parentSpan.SpanContext().TraceID())
		//fmt.Println("span", parentSpan.SpanContext().SpanID())

		ctx := r.Context()
		//ctx = context.WithValue(ctx, "trcID", parentSpan.SpanContext().TraceID())

		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return fn
}
