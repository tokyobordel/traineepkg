package errors

import "fmt"

type UnauthorizedError struct {
	Message string
}

func (e *UnauthorizedError) Error() string {
	return fmt.Sprintf("неавторизованный доступ: %s", e.Message)
}

func (e *UnauthorizedError) Type() string {
	return ErrorTypeUnauthorized
}

func NewUnauthorizedError(message string) *UnauthorizedError {
	return &UnauthorizedError{
		Message: message,
	}
}
