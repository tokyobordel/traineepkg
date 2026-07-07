package errors

// DomainError - общий интерфейс для всех доменных ошибок
type DomainError interface {
	error
	Type() string
}

// ErrorType - типы ошибок
const (
	ErrorTypeNotFound          = "NOT_FOUND"
	ErrorTypeUniqueConstraint  = "UNIQUE_CONSTRAINT"
	ErrorTypeInternalService   = "INTERNAL_SERVICE"
	ErrorTypeInvalidParameters = "INVALID_PARAMETERS"
	ErrorTypeValidation        = "VALIDATION"
	ErrorTypeUnauthorized      = "UNAUTHORIZED"
	ErrorTypeToken             = "TOKEN_ERROR"
	ErrorTypeAccessDenied      = "ACCESS_DENIED"
	ErrorTypeSignature         = "SIGNATURE_ERROR"
	ErrprTypeIntegration       = "INTEGRATION_ERROR"
)
