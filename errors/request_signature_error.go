package errors

import "errors"

// RequestSignatureErrorValues содержит предопределённые сообщения ошибок проверки подписи запроса.
type RequestSignatureErrorValues string

// RequestSignatureErrorNotFound означает, что клиент не авторизован.
const (
	RequestSignatureErrorNotFound    = "Client not authorized"
	RequestSignatureErrorInvalid     = "Invalid signature"
	RequestSignatureErrorDataInvalid = "Invalid signature data"
	RequestSignatureErrorApiKey      = "Api key required"
	RequestSignatureErrorSignature   = "Request signature required"
)

// RequestSignatureError описывает ошибку проверки подписи HTTP-запроса.
type RequestSignatureError struct {
	// RequestSignatureErrorValue — конкретное значение ошибки подписи.
	RequestSignatureErrorValue RequestSignatureErrorValues
}

// Error возвращает текстовое представление ошибки.
func (e *RequestSignatureError) Error() string {
	return string(e.RequestSignatureErrorValue)
}

// Type возвращает тип доменной ошибки.
func (e *RequestSignatureError) Type() string {
	return ErrorTypeSignature
}

// NewRequestSignatureError создаёт ошибку RequestSignatureError.
func NewRequestSignatureError(requestSignatureErrorValue RequestSignatureErrorValues) *RequestSignatureError {
	return &RequestSignatureError{
		RequestSignatureErrorValue: requestSignatureErrorValue,
	}
}

// IsRequestSignatureError сообщает, является ли err ошибкой типа SIGNATURE_ERROR.
func IsRequestSignatureError(err error) bool {
	var domainErr DomainError
	return errors.As(err, &domainErr) && domainErr.Type() == ErrorTypeSignature
}
