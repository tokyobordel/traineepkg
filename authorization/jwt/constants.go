// Package jwt предоставляет генерацию и валидацию JWT-токенов доступа и обновления.
package jwt

// AccessTokenCookieNameDefault — имя HTTP-cookie для access-токена.
const (
	AccessTokenCookieNameDefault  = "access_token"
	RefreshTokenCookieNameDefault = "refresh_token"
)

// AccessTokenType — тип JWT access-токена.
const (
	AccessTokenType  = "access"
	RefreshTokenType = "refresh"
)
