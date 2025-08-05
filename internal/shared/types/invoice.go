package types

import "time"

type Invoice struct {
	ID               uint64    `json:"id"`
	Order_id         uint64    `json:"order_id"`
	Payment_method   string    `json:"payment_method" validate:"eq=CARD|eq=CASH|eq=PIX"`
	Payment_status   string    `json:"payment_status" validate:"required,eq=PENDING|eq=PAID|eq=REFUNDED"`
	Payment_due_date time.Time `json:"payment_due_date"`
	Created_at       time.Time `json:"created_at"`
	Updated_at       time.Time `json:"updated_at"`
}
