package models

type User struct {
	ID    int               `json:"id" example:"1"`
	Login string            `json:"login" example:"john_doe"`
	Data  map[string]string `json:"data"`
}
