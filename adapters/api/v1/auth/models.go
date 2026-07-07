package auth

import "traineepkg/models"

type credentialsRequest struct {
	Login string `json:"login"`
	Pass  string `json:"pass"`
}

type authResponse struct {
	User models.User `json:"user"`
}
