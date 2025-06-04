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
