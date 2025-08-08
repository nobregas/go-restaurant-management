package user

import (
	"go-restaurant-management/internal/shared/errors/exceptions"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(user User) (User, error)
	FindByEmail(email string) (User, error)
}

type userService struct {
	UserRepository
}

func (u *userService) Register(user User) (User, error) {
	log.Printf("starting to register user %s", user.Email)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		log.Printf("error generating password hash for user %s: %v", user.Email, err)
		return User{}, exceptions.NewInternalServerError(err.Error())
	}

	user.Password = string(hashedPassword)

	user, err = u.UserRepository.Save(user)
	if err != nil {
		log.Printf("error saving user %s: %v", user.Email, err)
		return User{}, exceptions.NewInternalServerError(err.Error())
	}

	log.Printf("user %s registered successfully in service", user.Email)
	return user, nil
}

func (u *userService) FindByEmail(email string) (User, error) {
	log.Printf("finding user %s in service", email)
	user, err := u.UserRepository.FindByEmail(email)
	if err != nil {
		log.Printf("error finding user %s: %v", email, err)
		return User{}, exceptions.NewEntityNotFound("user", email)
	}

	log.Printf("user %s found successfully in service", email)
	return user, nil
}

func NewUserService(userRepository UserRepository) UserService {
	return &userService{userRepository}
}
