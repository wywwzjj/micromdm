package log

import (
	"context"

	"github.com/go-kit/kit/log"
)

type key int

const loggerKey key = 0

func NewContext(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func FromContext(ctx context.Context) Logger {
	v, ok := ctx.Value(loggerKey).(Logger)
	if !ok {
		return log.NewNopLogger()
	}

	span := traceFromContext(ctx)

	return log.With(v, "trace_id", span.TraceID)
}
