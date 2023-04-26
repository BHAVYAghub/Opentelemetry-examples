package apm

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

type Otel struct {
}

// getSpanFromContext returns span set within the context.
func (o Otel) GetSpanFromContext(ctx context.Context) trace.Span {
	span := trace.SpanFromContext(ctx)

	return span
}

// createSpanFromContext creates a new span which should be closed.
func (o Otel) CreateSpanFromContext(ctx context.Context, spanName string) (context.Context, interface{}) {
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("").Start(ctx, spanName)

	return ctx, span
}

// getTraceId returns trace Id within the current span
func (o Otel) GetTraceId(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)

	return span.SpanContext().TraceID().String()
}

// getSpanId returns span Id within the current span
func (o Otel) GetSpanId(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)

	return span.SpanContext().SpanID().String()
}

func (o Otel) CloseSpan(span interface{}) {
	span2, ok := span.(trace.Span)
	if !ok {
		// TODO: unable to close span.
	}
	span2.End()
}
