package user

import (
	"database/sql"
	"go-restaurant-management/internal/shared/utils"
)

type UserRepository interface {
	Save(user User) (User, error)
}

type userRepository struct {
	*sql.DB
}

func (u *userRepository) Save(user User) (User, error) {
	query := "INSERT INTO users (email, password, avatar, role) VALUES (?, ?, ?, ?)"

	result, err := u.DB.Exec(query, user.Email, user.Password, user.Avatar, user.Role)
	if err != nil {
		return User{}, err
	}

	composer := utils.NewComposerDB(u.DB, result, "users")

	return composer.Compose(user).(User), nil
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db}
}
