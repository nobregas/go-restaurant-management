package types

import "time"

type OrderItem struct {
	ID         uint64    `json:"id"`
	Quantity   uint64    `json:"quantity" validate:"required"`
	Unit_price float64   `json:"unit_price" validate:"required"`
	Food_id    uint64    `json:"food_id" validate:"required"`
	Order_id   uint64    `json:"order_id" validate:"required"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
}
