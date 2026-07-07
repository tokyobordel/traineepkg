package errors

import (
	"errors"
	"fmt"
)

type TokenError string

const (
	TokenNotFound TokenError = "Токен не найден"
	TokenExpires  TokenError = "Токен истек"
	InvalidToken  TokenError = "Невалидный токен"
)

// AuthTokenError - ошибка связанная с обработкой токена авторизации
type AuthTokenError struct {
	TokenError TokenError
}

func (e *AuthTokenError) Error() string {
	return fmt.Sprintf("Ошибка обработки токена %s", e.TokenError)
}

func (e *AuthTokenError) Type() string {
	return ErrorTypeToken
}

func NewAuthTokenError(tokenError TokenError) *AuthTokenError {
	return &AuthTokenError{
		TokenError: tokenError,
	}
}

func IsAutnTokenError(err error) bool {
	var domainErr DomainError
	return errors.As(err, &domainErr) && domainErr.Type() == ErrorTypeToken
}
