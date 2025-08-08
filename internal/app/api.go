package app

import (
	"database/sql"
	"go-restaurant-management/internal/app/handler"
	"go-restaurant-management/internal/domain/user"
	"log"
	"net/http"
)

type ApiServer struct {
	addr string
	db   *sql.DB
}

func NewApiServer(addr string, db *sql.DB) *ApiServer {
	return &ApiServer{
		addr: addr,
		db:   db,
	}

}

func (s *ApiServer) Run() error {
	// User
	userRepository := user.NewUserRepository(s.db)
	userService := user.NewUserService(userRepository)

	http.HandleFunc("/users", handler.UserHandler(userService))

	log.Printf("Server has started, listening on %s", s.addr)
	return http.ListenAndServe(s.addr, nil)
}
