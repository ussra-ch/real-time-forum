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


	DB.Exec(`INSERT INTO categories (name, icon) VALUES('Sport', '<i class="fa-solid fa-medal"></i>'),
	('Music', '<i class="fa-solid fa-music"></i>'),
	('Science', '<i class="fa-solid fa-flask"></i>'),
	('Tecknology', '<i class="fa-solid fa-microchip"></i>'),('Culture', '<i class="fa-solid fa-person-walking"></i>');`)
	
	fmt.Println("Queries executed successfully!")
}
