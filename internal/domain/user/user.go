package user

import "time"

type User struct {
	ID            int       `json:"id"`
	First_name    string    `json:"first_name"`
	Last_name     string    `json:"last_name"`
	Email         string    `json:"email"`
	Password      string    `json:"-"`
	Avatar        string    `json:"avatar"`
	Phone         string    `json:"phone"`
	Token         string    `json:"token"`
	Refresh_token string    `json:"refresh_token"`
	Role          string    `json:"role"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
