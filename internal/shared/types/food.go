package types

import "time"

type Food struct {
	ID          uint64    `json:"id"`
	Name        string    `json:"name" validate:"required,min=3,max=100"`
	Description string    `json:"description" validate:"min=3,max=255"`
	Price       float64   `json:"price" validate:"required,min=0"`
	Image       string    `json:"image" validate:"required"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Menu_id     uint64    `json:"menu_id" validate:"required"`
}
