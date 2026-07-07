package errors

import "fmt"

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("ошибка валидации: %s", e.Message)
}

func NewValidationError(message string) *ValidationError {
	return &ValidationError{
		Message: message,
	}
}
