package databases

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB(filepath string) {
	var err error
	DB, err = sql.Open("sqlite3", filepath)
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	createTable := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nickname TEXT NOT NULL UNIQUE,
		age INTEGER NOT NULL,
		gender TEXT NOT NULL,
		first_name TEXT NOT NULL,
		last_name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	);`
	_, err = DB.Exec(createTable)
	if err != nil {
		log.Fatal("Failed to create users table:", err)
	}
}
