package internal

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var Db *sql.DB

func InitDB() {
	var err error
	Db, err = sql.Open("sqlite3", "./go_motion.db")
	if err != nil {
		log.Fatal(err)
	}

	query := `CREATE TABLE IF NOT EXISTS processmodels (
		id UUID PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT,
		version INTEGER NOT NULL,
		created_at TIMESTAMP,
		updated_at TIMESTAMP,
		is_executable BOOLEAN,
		process_data JSONB
	);`
	_, err = Db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database initialized and connected.")
}
