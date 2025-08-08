package user

import "go-restaurant-management/internal/shared/types"

func ToUser(req types.RegisterUserRequest) User {
	return User{
		First_name: req.First_name,
		Last_name:  req.Last_name,
		Email:      req.Email,
		Password:   req.Password,
		Phone:      req.Phone,
	}
}
