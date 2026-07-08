package middleware

import (
	contextTrace "github.com/tokyobordel/traineepkg/context/trace"
	"github.com/tokyobordel/traineepkg/logger"

	"github.com/gofiber/fiber/v3"
)

type SpreadMiddleware struct {
	logger *logger.ContextLogger
}

func NewSpreadMiddleware(logger *logger.ContextLogger) *SpreadMiddleware {
	return &SpreadMiddleware{logger: logger}
}

// AddSpreadInContext автоматически добавляет spread ID в контекст каждого запроса
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
