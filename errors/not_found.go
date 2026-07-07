package errors

import "fmt"

// NotFoundError - ошибка когда ресурс не найден
type NotFoundError struct {
	Resource string
	ID       interface{}
}

func (e *NotFoundError) Error() string {
	if e.ID != nil {
		return fmt.Sprintf("%s с ID %v не найден", e.Resource, e.ID)
	}
	return fmt.Sprintf("%s не найден", e.Resource)
}

func (e *NotFoundError) Type() string {
	return ErrorTypeNotFound
}

func NewNotFoundError(resource string, id interface{}) *NotFoundError {
	return &NotFoundError{
		Resource: resource,
		ID:       id,
	}
}
