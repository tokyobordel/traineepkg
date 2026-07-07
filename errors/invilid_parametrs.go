package errors

import "fmt"

// InvalidParametersError - ошибка неверных параметров
type InvalidParametersError struct {
	Parameter string
	Value     interface{}
	Reason    string
}

func (e *InvalidParametersError) Error() string {
	if e.Reason != "" {
		return fmt.Sprintf("неверный параметр %s со значением %v: %s", e.Parameter, e.Value, e.Reason)
	}
	return fmt.Sprintf("неверный параметр %s со значением %v", e.Parameter, e.Value)
}

func (e *InvalidParametersError) Type() string {
	return ErrorTypeInvalidParameters
}

func NewInvalidParametersError(parameter string, value interface{}, reason string) *InvalidParametersError {
	return &InvalidParametersError{
		Parameter: parameter,
		Value:     value,
		Reason:    reason,
	}
}
