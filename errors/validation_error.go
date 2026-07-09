package errors

import "fmt"

// ValidationError описывает ошибку бизнес-валидации данных.
type ValidationError struct {
	// Message — описание нарушения правил валидации.
	Message string
}

// Error возвращает текстовое представление ошибки.
func (e *ValidationError) Error() string {
	return fmt.Sprintf("ошибка валидации: %s", e.Message)
}

// NewValidationError создаёт ошибку ValidationError.
func NewValidationError(message string) *ValidationError {
	return &ValidationError{
		Message: message,
	}
}
