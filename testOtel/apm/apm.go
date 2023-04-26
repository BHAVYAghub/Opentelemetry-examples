package apm

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

type APM interface {
	GetSpanFromContext(ctx context.Context) trace.Span
	CreateSpanFromContext(ctx context.Context, spanName string) (context.Context, interface{})
	GetTraceId(ctx context.Context) string
	GetSpanId(ctx context.Context) string
	CloseSpan(span interface{})
}
