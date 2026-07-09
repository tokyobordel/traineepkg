package errors

import "fmt"

// UniqueConstraintError описывает нарушение уникальности данных.
type UniqueConstraintError struct {
	// Resource — имя сущности с конфликтом уникальности.
	Resource string
	// Field — поле, значение которого уже существует.
	Field string
	// Value — конфликтующее значение поля.
	Value interface{}
}

// Error возвращает текстовое представление ошибки.
func (e *UniqueConstraintError) Error() string {
	return fmt.Sprintf("нарушение уникальности для %s: поле %s со значением %v уже существует", e.Resource, e.Field, e.Value)
}

// Type возвращает тип доменной ошибки.
func (e *UniqueConstraintError) Type() string {
	return ErrorTypeUniqueConstraint
}

// NewUniqueConstraintError создаёт ошибку UniqueConstraintError.
func NewUniqueConstraintError(resource, field string, value interface{}) *UniqueConstraintError {
	return &UniqueConstraintError{
		Resource: resource,
		Field:    field,
		Value:    value,
	}
}
