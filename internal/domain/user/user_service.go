package user

import (
	"go-restaurant-management/internal/shared/errors/exceptions"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(user User) (User, error)
}

type userService struct {
	UserRepository
}

func (u *userService) Register(user User) (User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return User{}, exceptions.NewInternalServerError(err.Error())
	}

	user.Password = string(hashedPassword)

	user, err = u.UserRepository.Save(user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func NewUserService(userRepository UserRepository) UserService {
	return &userService{userRepository}
}
