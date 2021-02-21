package frontend

import "context"

type key int

const dataKey key = 0

func AddFormData(ctx context.Context, data Data) context.Context {
	return context.WithValue(ctx, dataKey, data)
}

func dataFromContext(ctx context.Context) Data {
	v, ok := ctx.Value(dataKey).(Data)
	if !ok {
		return Data{}
	}
	return v
}
