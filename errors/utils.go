package errors

import "errors"

// IsNotFoundError сообщает, является ли err ошибкой типа NOT_FOUND.
func IsNotFoundError(err error) bool {
	var domainErr DomainError
	return errors.As(err, &domainErr) && domainErr.Type() == ErrorTypeNotFound
}

// IsUniqueConstraintError сообщает, является ли err ошибкой типа UNIQUE_CONSTRAINT.
func IsUniqueConstraintError(err error) bool {
	var domainErr DomainError
	return errors.As(err, &domainErr) && domainErr.Type() == ErrorTypeUniqueConstraint
}

// IsInternalServiceError сообщает, является ли err ошибкой типа INTERNAL_SERVICE.
func IsInternalServiceError(err error) bool {
	var domainErr DomainError
	return errors.As(err, &domainErr) && domainErr.Type() == ErrorTypeInternalService
}

// IsInvalidParametersError сообщает, является ли err ошибкой типа INVALID_PARAMETERS.
func IsInvalidParametersError(err error) bool {
	var domainErr DomainError
	return errors.As(err, &domainErr) && domainErr.Type() == ErrorTypeInvalidParameters
}

// ConvertToServiceError преобразует произвольную ошибку в доменную InternalServiceError.
func ConvertToServiceError(err error, context string) DomainError {
	var domainErr DomainError
	if errors.As(err, &domainErr) {
		return domainErr
	}
	return NewInternalServiceError(context, err)
}

// IsDomainError сообщает, реализует ли err интерфейс DomainError.
func IsDomainError(err error) bool {
	var domainErr DomainError
	return errors.As(err, &domainErr)
}
