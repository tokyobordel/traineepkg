package trace

import (
	"context"

	pkgContext "traineesheep/imageservice/pkg/context"
	"traineesheep/imageservice/pkg/trace"
)

const SpreadKey pkgContext.CtxKey = "traceSpread"

// Добавить spread в контекст
func WithSpread(ctx context.Context, spread string) context.Context {
	/*
		Добавляет spread в контекст
		- spread - не обязательный параметр, если не передан, то генерируется новый spread
	*/
	if spread == "" {
		spread = trace.GenerateSpreadId()
	}
	return context.WithValue(ctx, SpreadKey, spread)
}

func SpreadFromContext(ctx context.Context) (string, bool) {
	spread := ctx.Value(SpreadKey)

	if spread == nil {
		return "", false
	}

	// Безопасный type assertion
	if str, ok := spread.(string); ok {
		return str, true
	} else {
		return "", false
	}
}
