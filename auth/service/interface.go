// Package service определяет контракт сервиса аутентификации для HTTP-адаптеров.
package service

import "github.com/tokyobordel/traineepkg/models"

// IAuthService описывает операции аутентификации, реализуемые прикладным сервисом.
type IAuthService interface {
	// Login аутентифицирует пользователя по паролю и логину.
	Login(pass string, login string) (models.User, error)
	// Register создаёт пользователя с дополнительными данными data.
	Register(pass string, login string, data map[string]string) (models.User, error)
	// GetMe возвращает пользователя по его идентификатору.
	GetMe(id int) (models.User, error)
}
