package auth

import "github.com/tokyobordel/traineepkg/models"

// CredentialsRequest учётные данные пользователя.
type CredentialsRequest struct {
	Login string `json:"login" example:"john_doe"`
	Pass  string `json:"pass" example:"secret123"`
}

// AuthResponse данные аутентифицированного пользователя.
type AuthResponse struct {
	User models.User `json:"user"`
}
