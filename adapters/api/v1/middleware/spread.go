// Package middleware предоставляет общие HTTP-middleware для Fiber API.
package middleware

import (
	contextTrace "github.com/tokyobordel/traineepkg/context/trace"
	"github.com/tokyobordel/traineepkg/logger"

	"github.com/gofiber/fiber/v3"
)

// SpreadMiddleware добавляет spread ID в контекст входящих HTTP-запросов.
type SpreadMiddleware struct {
	logger *logger.ContextLogger
}

// NewSpreadMiddleware создаёт middleware для работы с spread ID.
func NewSpreadMiddleware(logger *logger.ContextLogger) *SpreadMiddleware {
	return &SpreadMiddleware{logger: logger}
}

// AddSpreadInContext возвращает Fiber-handler, который записывает spread ID в контекст запроса.
func (m *SpreadMiddleware) AddSpreadInContext() fiber.Handler {
	return func(c fiber.Ctx) error {
		ctx := contextTrace.WithSpread(c.Context(), "")
		c.SetContext(ctx)

		spreadID, ok := contextTrace.SpreadFromContext(c.Context())
		if !ok {
			m.logger.Errorf(c.Context(), "Error adding spread identifier to context")
		}
		m.logger.Infof(c.Context(), "Spread identifier added %v", spreadID)
		return c.Next()
	}
}
