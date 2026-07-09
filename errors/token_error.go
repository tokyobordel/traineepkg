package errors

import (
	"errors"
	"fmt"
)

// TokenError описывает конкретную причину ошибки JWT-токена.
type TokenError string

// TokenNotFound означает отсутствие токена в запросе.
const (
	TokenNotFound TokenError = "Токен не найден"
	TokenExpires  TokenError = "Токен истек"
	InvalidToken  TokenError = "Невалидный токен"
)

// AuthTokenError описывает ошибку обработки токена авторизации.
type AuthTokenError struct {
	// TokenError — причина ошибки токена.
	TokenError TokenError
}

// Error возвращает текстовое представление ошибки.
func (e *AuthTokenError) Error() string {
	return fmt.Sprintf("Ошибка обработки токена %s", e.TokenError)
}

// Type возвращает тип доменной ошибки.
func (e *AuthTokenError) Type() string {
	return ErrorTypeToken
}

// NewAuthTokenError создаёт ошибку AuthTokenError.
func NewAuthTokenError(tokenError TokenError) *AuthTokenError {
	return &AuthTokenError{
		TokenError: tokenError,
	}
}

// IsAutnTokenError сообщает, является ли err ошибкой типа TOKEN_ERROR.
func IsAutnTokenError(err error) bool {
	var domainErr DomainError
	return errors.As(err, &domainErr) && domainErr.Type() == ErrorTypeToken
}
