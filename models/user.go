package models

type User struct {
	ID    int         `json:"id"`
	Login string      `json:"login"`
	Data  interface{} `json:"data"`
}
