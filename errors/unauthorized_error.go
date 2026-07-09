package errors

import "fmt"

// UnauthorizedError описывает ошибку неавторизованного доступа.
type UnauthorizedError struct {
	// Message — пояснение причины отказа в авторизации.
	Message string
}

// Error возвращает текстовое представление ошибки.
func (e *UnauthorizedError) Error() string {
	return fmt.Sprintf("неавторизованный доступ: %s", e.Message)
}

// Type возвращает тип доменной ошибки.
func (e *UnauthorizedError) Type() string {
	return ErrorTypeUnauthorized
}

// NewUnauthorizedError создаёт ошибку UnauthorizedError.
func NewUnauthorizedError(message string) *UnauthorizedError {
	return &UnauthorizedError{
		Message: message,
	}
}
