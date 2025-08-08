package types

type RegisterUserRequest struct {
	First_name string `json:"first_name" validate:"required,min=2,max=100"`
	Last_name  string `json:"last_name" validate:"required,min=2,max=100"`
	Password   string `json:"password" validate:"required,min=6,max=100"`
	Email      string `json:"email" validate:"required,email|unique"`
	Phone      string `json:"phone" validate:"required|unique|numeric|len=11"`
}
