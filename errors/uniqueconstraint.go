package errors

import "fmt"

// UniqueConstraintError - ошибка нарушения уникальности
type UniqueConstraintError struct {
	Resource string
	Field    string
	Value    interface{}
}

func (e *UniqueConstraintError) Error() string {
	return fmt.Sprintf("нарушение уникальности для %s: поле %s со значением %v уже существует", e.Resource, e.Field, e.Value)
}

func (e *UniqueConstraintError) Type() string {
	return ErrorTypeUniqueConstraint
}

func NewUniqueConstraintError(resource, field string, value interface{}) *UniqueConstraintError {
	return &UniqueConstraintError{
		Resource: resource,
		Field:    field,
		Value:    value,
	}
}
