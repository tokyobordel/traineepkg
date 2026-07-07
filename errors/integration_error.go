package errors

type IntegrationError struct {
	Message string
	Cause   error
}

func (e *IntegrationError) Error() string {
	return "Ошибка во внешнем сервисе"
}

func (e *IntegrationError) Type() string {
	return ErrprTypeIntegration
}

func NewIntegrationError(message string, cause error) *IntegrationError {
	return &IntegrationError{
		Message: message,
		Cause:   cause,
	}
}
