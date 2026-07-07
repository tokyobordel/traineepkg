package errors

// InternalServiceError - внутренняя ошибка сервиса
type InternalServiceError struct {
	Message string
	Cause   error
}

func (e *InternalServiceError) Error() string {
	return "внутренняя ошибка сервиса"
}

func (e *InternalServiceError) Type() string {
	return ErrorTypeInternalService
}

func NewInternalServiceError(message string, cause error) *InternalServiceError {
	return &InternalServiceError{
		Message: message,
		Cause:   cause,
	}
}
