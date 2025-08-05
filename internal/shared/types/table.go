package types

import "time"

type Table struct {
	ID               uint64    `json:"id"`
	Number_of_guests uint64    `json:"number_of_guests" validate:"required,min=1,max=100"`
	Table_number     uint64    `json:"table_number" validate:"required,min=1,max=100"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
