package errors

import "fmt"

// AccessDeniedError описывает отказ в доступе к ресурсу.
type AccessDeniedError struct {
	// Resurse — ресурс, к которому доступ запрещён.
	Resurse interface{}
}

// Error возвращает текстовое представление ошибки.
func (e *AccessDeniedError) Error() string {
	return fmt.Sprintf("Отказано в доступе к ресурсу %v", e.Resurse)
}

// Type возвращает тип доменной ошибки.
func (e *AccessDeniedError) Type() string {
	return ErrorTypeAccessDenied
}

// NewAccessDeniedError создаёт ошибку AccessDeniedError.
func NewAccessDeniedError(resourse interface{}) *AccessDeniedError {
	return &AccessDeniedError{
		Resurse: resourse,
	}
}
