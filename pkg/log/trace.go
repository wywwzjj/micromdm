package log

import (
	"context"

	"micromdm.io/v2/pkg/id"
)

// TraceID return a unique ID associated with a request.
// The trace ID can be used to identify a particular HTTP response and logged error.
func TraceID(ctx context.Context) string {
	return traceFromContext(ctx).TraceID
}

const traceKey key = 1

// span will eventually get replaced by http://opentelemetry.io
// or something similar.
type span struct {
	TraceID string
}

func newTraceContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, traceKey, span{TraceID: id.New()})
}

func traceFromContext(ctx context.Context) span {
	v, ok := ctx.Value(traceKey).(span)
	if !ok {
		return span{TraceID: id.New()}
	}

	return v
}
