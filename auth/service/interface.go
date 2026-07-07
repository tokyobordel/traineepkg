package service

import "github.com/tokyobordel/traineepkg/models"

type IAuthService interface {
	Login(pass string, login string) (models.User, error)
	Register(pass string, login string) (models.User, error)
	GetMe(id int) (models.User, error)
}
