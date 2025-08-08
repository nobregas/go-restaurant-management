package main

import (
	"database/sql"
	"fmt"
	"go-restaurant-management/config"
	"go-restaurant-management/internal/app"
	"log"
	"os"
	"strconv"

	mysqlCfg "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	// Connect to MySQL server without specifying a database
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/", config.Envs.DB_USER, config.Envs.DB_PASSWORD, config.Envs.DB_ADDRESS)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Migration: " + err.Error())
	}
	defer db.Close()

	// Create the database if it doesn't exist
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", config.Envs.DB_NAME))
	if err != nil {
		log.Fatal("Migration: " + err.Error())
	}

	// Now, connect to the database
	dbWithDb, err := app.NewMySqlStorage(mysqlCfg.Config{
		User:                 config.Envs.DB_USER,
		Passwd:               config.Envs.DB_PASSWORD,
		Addr:                 config.Envs.DB_ADDRESS,
		DBName:               config.Envs.DB_NAME,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	})

	if err != nil {
		log.Fatal("Migration: " + err.Error())
	}

	driver, err := mysql.WithInstance(dbWithDb, &mysql.Config{})
	if err != nil {
		log.Fatal("Migration: " + err.Error())
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://cmd/migrate/migrations",
		"mysql",
		driver,
	)
	if err != nil {
		log.Fatal("Migration: " + err.Error())
	}

	cmd := os.Args[(len(os.Args) - 1)]
	if cmd == "up" {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal("Migration: " + err.Error())
		}
	}
	if cmd == "down" {
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal("Migration: " + err.Error())
		}
	}
	if cmd == "force" {
		if len(os.Args) < 3 {
			log.Fatal("Migration: Missing version argument for force command")
		}
		version := os.Args[2]
		v, err := strconv.Atoi(version)
		if err != nil {
			log.Fatal("Migration: Invalid version number")
		}
		if err := m.Force(v); err != nil {
			log.Fatal("Migration: " + err.Error())
		}
	}
	log.Println("Migration: Successfully executed")
}
