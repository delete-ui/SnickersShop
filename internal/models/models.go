package models

import "github.com/google/uuid"

type Snickers struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Cost        float64   `json:"cost"`
}

type SnickerInput struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Cost        float64 `json:"cost"`
}

type Pagination struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type User struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"` //В реальном проекте хранить кэш
}

type UserInput struct {
	Username string `json:"username"`
	Password string `json:"password"` //В реальном проекте хранить кэш
}

type UserLogInResponse struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	Status   string `json:"status"`
}

type CostRange struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}
