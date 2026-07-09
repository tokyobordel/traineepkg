// Package trace предоставляет функции работы с идентификатором spread в context.Context.
package trace

import (
	"context"

	pkgContext "github.com/tokyobordel/traineepkg/context"
	pkgTrace "github.com/tokyobordel/traineepkg/trace"
)

// SpreadKey — ключ context.Context для хранения spread ID.
const SpreadKey pkgContext.CtxKey = "traceSpread"

// WithSpread добавляет spread в контекст. Если spread пустой, генерируется новый идентификатор.
func WithSpread(ctx context.Context, spread string) context.Context {
	if spread == "" {
		spread = pkgTrace.GenerateSpreadId()
	}
	return context.WithValue(ctx, SpreadKey, spread)
}

// SpreadFromContext извлекает spread ID из контекста.
func SpreadFromContext(ctx context.Context) (string, bool) {
	spread := ctx.Value(SpreadKey)

	if spread == nil {
		return "", false
	}

	if str, ok := spread.(string); ok {
		return str, true
	}

	return "", false
}
