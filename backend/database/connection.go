package database

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func Connect(dbPath string) error {
	var err error
	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}

	if err = DB.Ping(); err != nil {
		return err
	}

	log.Println("Database connected successfully")
	return nil
}

func Close() {
	if DB != nil {
		DB.Close()
		log.Println("Database connection closed")
	}
}
