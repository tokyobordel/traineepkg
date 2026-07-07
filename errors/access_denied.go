package errors

import "fmt"

// AuthTokenError - ошибка связанная с обработкой токена авторизации
type AccessDeniedError struct {
	Resurse interface{}
}

func (e *AccessDeniedError) Error() string {

	return fmt.Sprintf("Отказано в доступе к ресурсу %v", e.Resurse)
}

func (e *AccessDeniedError) Type() string {
	return ErrorTypeAccessDenied
}

func NewAccessDeniedError(resourse interface{}) *AccessDeniedError {
	return &AccessDeniedError{
		Resurse: resourse,
	}
}
