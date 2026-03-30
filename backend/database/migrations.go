package database

import (
	"log"
)

func Migrate() error {
	query := `
	CREATE TABLE IF NOT EXISTS recipes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		url TEXT,
		cuisine TEXT,
		ingredients TEXT,
		instructions TEXT,
		cooking_time INTEGER DEFAULT 0,
		servings INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := DB.Exec(query)
	if err != nil {
		log.Printf("Migration failed: %v", err)
		return err
	}

	log.Println("Migration completed successfully")
	return nil
}
