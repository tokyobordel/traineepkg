// Package errors предоставляет унифицированные доменные ошибки и вспомогательные функции для их проверки.
package errors

// DomainError — общий интерфейс для всех доменных ошибок пакета.
type DomainError interface {
	error
	Type() string
}

// ErrorTypeNotFound обозначает отсутствие запрошенного ресурса.
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
