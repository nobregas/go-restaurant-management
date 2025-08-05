package app

import (
	"database/sql"
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
	return http.ListenAndServe(s.addr, nil)
}
