package types

import "time"

type Order struct {
	ID         uint64    `json:"id"`
	Order_Date time.Time `json:"order_date" validate:"required"`
	Table_id   uint64    `json:"table_id" validate:"required"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
}
