package auth

import "github.com/tokyobordel/traineepkg/models"

// CredentialsRequest содержит учётные данные пользователя для входа или регистрации.
type CredentialsRequest struct {
	Login string            `json:"login" example:"john_doe"`
	Pass  string            `json:"pass" example:"secret123"`
	Data  map[string]string `json:"data,omitempty"`
}

// AuthResponse содержит данные аутентифицированного пользователя.
type AuthResponse struct {
	User models.User `json:"user"`
}
