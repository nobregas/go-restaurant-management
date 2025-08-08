package user

import (
	"database/sql"
	"log"
)

type UserRepository interface {
	Save(user User) (User, error)
	FindByEmail(email string) (User, error)
}

type userRepository struct {
	*sql.DB
}

func (u *userRepository) Save(user User) (User, error) {
	log.Printf("saving user %s to database", user.Email)
	query := "INSERT INTO users (first_name, last_name, email, password, phone, avatar, role) VALUES (?, ?, ?, ?, ?, ?, ?)"

	result, err := u.DB.Exec(query, user.First_name, user.Last_name, user.Email, user.Password, user.Phone, user.Avatar, user.Role)
	if err != nil {
		log.Printf("error executing insert for user %s: %v", user.Email, err)
		return User{}, err
	}

	// Obter o ID do usuário inserido
	userID, err := result.LastInsertId()
	if err != nil {
		log.Printf("error getting last insert ID for user %s: %v", user.Email, err)
		return User{}, err
	}

	// Definir o ID no usuário
	user.ID = int(userID)

	log.Printf("user %s saved successfully with ID %d", user.Email, user.ID)

	return user, nil
}

func (u *userRepository) FindByEmail(email string) (User, error) {
	log.Printf("finding user %s in database", email)
	query := "SELECT id, first_name, last_name, email, phone, avatar, role FROM users WHERE email = ?"

	row := u.DB.QueryRow(query, email)
	var user User
	err := row.Scan(&user.ID, &user.First_name, &user.Last_name, &user.Email, &user.Phone, &user.Avatar, &user.Role)
	if err != nil {
		log.Printf("error finding user %s: %v", email, err)
		return User{}, err
	}

	log.Printf("user %s found successfully", email)
	return user, nil
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db}
}
