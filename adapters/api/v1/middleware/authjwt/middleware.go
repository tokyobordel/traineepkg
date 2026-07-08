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

const UserIDContextKey ctxKey = "authUserID"

type Middleware struct {
	jwtService *jwtAuth.Service
}

func NewMiddleware(jwtService *jwtAuth.Service) *Middleware {
	return &Middleware{
		jwtService: jwtService,
	}
}

func (m *Middleware) RequireAccessToken() fiber.Handler {
	return func(c fiber.Ctx) error {
		accessToken := c.Cookies(jwtAuth.AccessTokenCookieName)
		if accessToken == "" {
			response.MakeErrorResponse(c, nil, errors.NewAuthTokenError(errors.TokenNotFound))
			return nil
		}

		refreshToken := c.Cookies(jwtAuth.RefreshTokenCookieName)
		if accessToken == "" {
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

		var tokenPair jwtAuth.TokenPair
		tokenPair, err = m.jwtService.GenerateTokenPair(userID)
		if err != nil {

			response.MakeErrorResponse(c, nil, errors.NewInternalServiceError("Token generate error", err))
			return nil
		}

		expires := time.Now().Add(m.jwtService.GetAccessTTL())
		c.Cookie(&fiber.Cookie{
			Name:     jwtAuth.AccessTokenCookieName,
			Value:    tokenPair.AccessToken,
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

func UserIDFromContext(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value(UserIDContextKey).(int)
	if !ok {
		return 0, false
	}
	return userID, true
}
