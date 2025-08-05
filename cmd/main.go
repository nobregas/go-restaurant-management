package main

import (
	"database/sql"
	"go-restaurant-management/config"
	"go-restaurant-management/internal/app"
	"log"

	"github.com/go-sql-driver/mysql"
)

func main() {
	cfg := mysql.Config{
		User:                 config.Envs.DB_USER,
		Passwd:               config.Envs.DB_PASSWORD,
		Addr:                 config.Envs.DB_ADDRESS,
		DBName:               config.Envs.DB_NAME,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	db, err := app.NewMySqlStorage(cfg)
	if err != nil {
		log.Fatal("DB: " + err.Error())
	}

	initStorage(db)

	server := app.NewApiServer(":8080", db)
	if err := server.Run(); err != nil {
		log.Fatal("Server: " + err.Error())
	}
}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal("DB: " + err.Error())
	}

	log.Println("DB: Successfully connected")
}
