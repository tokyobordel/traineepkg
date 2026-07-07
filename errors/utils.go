package errors

import "errors"

func IsNotFoundError(err error) bool {
	var domainErr DomainError
	return errors.As(err, &domainErr) && domainErr.Type() == ErrorTypeNotFound
}

func IsUniqueConstraintError(err error) bool {
	var domainErr DomainError
	return errors.As(err, &domainErr) && domainErr.Type() == ErrorTypeUniqueConstraint
}

func IsInternalServiceError(err error) bool {
	var domainErr DomainError
	return errors.As(err, &domainErr) && domainErr.Type() == ErrorTypeInternalService
}

func IsInvalidParametersError(err error) bool {
	var domainErr DomainError
	return errors.As(err, &domainErr) && domainErr.Type() == ErrorTypeInvalidParameters
}

// ConvertToServiceError - конвертирует обычную ошибку в доменную
func ConvertToServiceError(err error, context string) DomainError {
	var domainErr DomainError
	if errors.As(err, &domainErr) {
		return domainErr
	}
	return NewInternalServiceError(context, err)
}

func IsDomainError(err error) bool {
	var domainErr DomainError
	return errors.As(err, &domainErr)
}
