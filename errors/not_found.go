package errors

import "fmt"

// NotFoundError описывает ситуацию, когда запрошенный ресурс не найден.
type NotFoundError struct {
	// Resource — имя сущности, которая не была найдена.
	Resource string
	// ID — идентификатор отсутствующего ресурса.
	ID interface{}
}

// Error возвращает текстовое представление ошибки.
func (e *NotFoundError) Error() string {
	if e.ID != nil {
		return fmt.Sprintf("%s с ID %v не найден", e.Resource, e.ID)
	}
	return fmt.Sprintf("%s не найден", e.Resource)
}

// Type возвращает тип доменной ошибки.
func (e *NotFoundError) Type() string {
	return ErrorTypeNotFound
}

// NewNotFoundError создаёт ошибку NotFoundError.
func NewNotFoundError(resource string, id interface{}) *NotFoundError {
	return &NotFoundError{
		Resource: resource,
		ID:       id,
	}
}
