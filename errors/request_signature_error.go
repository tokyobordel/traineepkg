package errors

import "errors"

type RequestSignatureErrorValues string

const (
	RequestSignatureErrorNotFound    = "Client not authorized"
	RequestSignatureErrorInvalid     = "Invalid signature"
	RequestSignatureErrorDataInvalid = "Invalid signature data"
	RequestSignatureErrorApiKey      = "Api key required"
	RequestSignatureErrorSignature   = "Request signature required"
)

type RequestSignatureError struct {
	RequestSignatureErrorValue RequestSignatureErrorValues
}

func (e *RequestSignatureError) Error() string {
	return string(e.RequestSignatureErrorValue)
}

func (e *RequestSignatureError) Type() string {
	return ErrorTypeSignature
}

func NewRequestSignatureError(requestSignatureErrorValue RequestSignatureErrorValues) *RequestSignatureError {
	return &RequestSignatureError{
		RequestSignatureErrorValue: requestSignatureErrorValue,
	}
}

func IsRequestSignatureError(err error) bool {
	var domainErr DomainError
	return errors.As(err, &domainErr) && domainErr.Type() == ErrorTypeSignature
}
