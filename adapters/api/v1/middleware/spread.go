package middleware

import (
	contextTrace "traineepkg/context/trace"
	"traineepkg/logger"

	"github.com/gofiber/fiber/v2"
)

type SpreadMiddleware struct {
	logger *logger.ContextLogger
}

func NewSpreadMiddleware(logger *logger.ContextLogger) *SpreadMiddleware {
	return &SpreadMiddleware{logger: logger}
}

// AddSpreadInContext автоматически добавляет spread ID в контекст каждого запроса
func (m *SpreadMiddleware) AddSpreadInContext() fiber.Handler {

	return func(c *fiber.Ctx) error {
		ctx := contextTrace.WithSpread(c.UserContext(), "")
		c.SetUserContext(ctx)

		spreadID, ok := contextTrace.SpreadFromContext(c.UserContext())
		if !ok {
			m.logger.Errorf(c.UserContext(), "Error adding spread identifier to context")
		}
		m.logger.Infof(c.UserContext(), "Spread identifier added %v", spreadID)
		return c.Next()
	}
}
