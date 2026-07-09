// Package models содержит общие доменные структуры данных.
package models

// User описывает пользователя системы.
type User struct {
	// ID — уникальный идентификатор пользователя.
	ID int `json:"id" example:"1"`
	// Login — логин пользователя.
	Login string `json:"login" example:"john_doe"`
	// Data — дополнительные строковые атрибуты пользователя.
	Data map[string]string `json:"data"`
}
