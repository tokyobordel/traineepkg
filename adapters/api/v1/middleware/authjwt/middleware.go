// Package authjwt предоставляет JWT-middleware для защиты HTTP-маршрутов.
package authjwt

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v3"

	"github.com/tokyobordel/traineepkg/adapters/api/v1/response"
	jwtAuth "github.com/tokyobordel/traineepkg/authorization/jwt"
	"github.com/tokyobordel/traineepkg/errors"
)

type ctxKey string

// UserIDContextKey — ключ context.Context для идентификатора аутентифицированного пользователя.
const UserIDContextKey ctxKey = "authUserID"

// Middleware проверяет JWT-токены и обновляет access cookie при необходимости.
type Middleware struct {
	jwtService *jwtAuth.Service
}

// NewMiddleware создаёт JWT-middleware с указанным сервисом токенов.
func NewMiddleware(jwtService *jwtAuth.Service) *Middleware {
	return &Middleware{
		jwtService: jwtService,
	}
}

// RequireAccessToken возвращает Fiber-handler, требующий валидный access- или refresh-токен.
func (m *Middleware) RequireAccessToken() fiber.Handler {
	return func(c fiber.Ctx) error {
		accessToken := c.Cookies(m.jwtService.GetAccessTokenCookieName())
		refreshToken := c.Cookies(m.jwtService.GetRefreshTokenCookieName())

		if refreshToken == "" {
			response.MakeErrorResponse(c, nil, errors.NewAuthTokenError(errors.TokenNotFound))
			return nil
		}

		userID, err := m.jwtService.ValidateAccessToken(accessToken)
		if err != nil {
			userID, err = m.jwtService.ValidateRefreshToken(refreshToken)
			if err != nil {
				response.MakeErrorResponse(c, nil, errors.NewAuthTokenError(errors.InvalidToken))
				return nil
			}
		}

		newAccessToken, err := m.jwtService.GenerateAccess(userID)
		if err != nil {
			response.MakeErrorResponse(c, nil, errors.NewInternalServiceError("Token generate error", err))
			return nil
		}

		expires := time.Now().Add(m.jwtService.GetAccessTTL())
		c.Cookie(&fiber.Cookie{
			Name:     m.jwtService.GetAccessTokenCookieName(),
			Value:    newAccessToken,
			Expires:  expires,
			HTTPOnly: true,
			Secure:   true,
			SameSite: "Lax",
		})

		ctx := context.WithValue(c.Context(), UserIDContextKey, userID)
		c.SetContext(ctx)
		return c.Next()
	}
}

// UserIDFromContext извлекает идентификатор пользователя из контекста запроса.
func UserIDFromContext(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value(UserIDContextKey).(int)
	if !ok {
		return 0, false
	}
	return userID, true
}
