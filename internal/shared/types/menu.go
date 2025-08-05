package types

import "time"

type Menu struct {
	ID         uint64    `json:"id"`
	Name       string    `json:"name" validate:"required,min=3,max=100"`
	Category   string    `json:"category" validate:"required,min=3,max=100"`
	Start_Date time.Time `json:"start_date"`
	End_Date   time.Time `json:"end_date"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
