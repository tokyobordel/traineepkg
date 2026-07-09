package errors

// InternalServiceError описывает внутреннюю ошибку сервиса, скрытую от клиента.
type InternalServiceError struct {
	// Message — внутреннее описание ошибки для логирования.
	Message string
	// Cause — исходная ошибка, вызвавшая сбой.
	Cause error
}

// Error возвращает безопасное для клиента сообщение об ошибке.
func (e *InternalServiceError) Error() string {
	return "внутренняя ошибка сервиса"
}

// Type возвращает тип доменной ошибки.
func (e *InternalServiceError) Type() string {
	return ErrorTypeInternalService
}

// NewInternalServiceError создаёт ошибку InternalServiceError.
func NewInternalServiceError(message string, cause error) *InternalServiceError {
	return &InternalServiceError{
		Message: message,
		Cause:   cause,
	}
}
