// Package jwt предоставляет генерацию и валидацию JWT-токенов доступа и обновления.
package jwt

// AccessTokenCookieName — имя HTTP-cookie для access-токена.
const (
	AccessTokenCookieName  = "access_token"
	RefreshTokenCookieName = "refresh_token"
)

// AccessTokenType — тип JWT access-токена.
const (
	AccessTokenType  = "access"
	RefreshTokenType = "refresh"
)
