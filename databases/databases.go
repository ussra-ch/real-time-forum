package databases

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB(filepath string) {
	var err error
	DB, err = sql.Open("sqlite3", filepath)
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	sqlfile, err := os.ReadFile("./databases/my.sql")
	if err != nil {
		log.Fatal("read error:", err)
	}

	_, err = DB.Exec(string(sqlfile))
	if err != nil {
		log.Fatal("exec error: ", err)
	}

	fmt.Println("Queries executed successfully!")
}
