package authjwt

import (
	"context"

	"github.com/gofiber/fiber/v2"

	"traineesheep/imageservice/pkg/adapters/api/v1/response"
	jwtAuth "traineesheep/imageservice/pkg/authorization/jwt"
	"traineesheep/imageservice/pkg/errors"
	"traineesheep/imageservice/pkg/logger"
)

type ctxKey string

const UserIDContextKey ctxKey = "authUserID"

type Middleware struct {
	jwtService *jwtAuth.Service
	logger     *logger.ContextLogger
}

func NewMiddleware(jwtService *jwtAuth.Service, logger *logger.ContextLogger) *Middleware {
	return &Middleware{
		jwtService: jwtService,
		logger:     logger,
	}
}

func (m *Middleware) RequireAccessToken() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Cookies(jwtAuth.AccessTokenCookieName)
		if token == "" {
			response.MakeErrorResponse(c, m.logger, errors.NewAuthTokenError(errors.TokenNotFound))
			return nil
		}

		userID, err := m.jwtService.ValidateAccessToken(token)
		if err != nil {
			response.MakeErrorResponse(c, m.logger, errors.NewAuthTokenError(errors.InvalidToken))
			return nil
		}

		ctx := context.WithValue(c.UserContext(), UserIDContextKey, userID)
		c.SetUserContext(ctx)
		return c.Next()
	}
}

func UserIDFromContext(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value(UserIDContextKey).(int)
	if !ok {
		return 0, false
	}
	return userID, true
}
