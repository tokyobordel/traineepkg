package errors

// IntegrationError описывает ошибку при обращении к внешнему сервису.
type IntegrationError struct {
	// Message — описание сбоя интеграции.
	Message string
	// Cause — исходная ошибка внешнего вызова.
	Cause error
}

// Error возвращает текстовое представление ошибки.
func (e *IntegrationError) Error() string {
	return "Ошибка во внешнем сервисе"
}

// Type возвращает тип доменной ошибки.
func (e *IntegrationError) Type() string {
	return ErrprTypeIntegration
}

// NewIntegrationError создаёт ошибку IntegrationError.
func NewIntegrationError(message string, cause error) *IntegrationError {
	return &IntegrationError{
		Message: message,
		Cause:   cause,
	}
}
