package service

import "traineesheep/imageservice/pkg/models"

type IAuthService interface {
	Login(pass string, login string) (models.User, error)
	Register(pass string, login string) (models.User, error)
	GetMe(id int) (models.User, error)
}
