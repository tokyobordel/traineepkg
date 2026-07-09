package errors

import "fmt"

// InvalidParametersError описывает ошибку неверных входных параметров.
type InvalidParametersError struct {
	// Parameter — имя параметра с некорректным значением.
	Parameter string
	// Value — переданное значение параметра.
	Value interface{}
	// Reason — пояснение причины ошибки.
	Reason string
}

// Error возвращает текстовое представление ошибки.
func (e *InvalidParametersError) Error() string {
	if e.Reason != "" {
		return fmt.Sprintf("неверный параметр %s со значением %v: %s", e.Parameter, e.Value, e.Reason)
	}
	return fmt.Sprintf("неверный параметр %s со значением %v", e.Parameter, e.Value)
}

// Type возвращает тип доменной ошибки.
func (e *InvalidParametersError) Type() string {
	return ErrorTypeInvalidParameters
}

// NewInvalidParametersError создаёт ошибку InvalidParametersError.
func NewInvalidParametersError(parameter string, value interface{}, reason string) *InvalidParametersError {
	return &InvalidParametersError{
		Parameter: parameter,
		Value:     value,
		Reason:    reason,
	}
}
