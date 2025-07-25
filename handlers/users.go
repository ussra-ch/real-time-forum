package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"handlers/databases"
)

func FetchUsers(w http.ResponseWriter, r *http.Request) {
	
	cookie, err := r.Cookie("sessionId")
	if err != nil {
		w.Write([]byte(`{"loggedIn": false}`))
		return
	}

	query1 := `SELECT user_id FROM sessions WHERE id = ? AND expires_at > DATETIME('now')`
	var userID int
	err = databases.DB.QueryRow(query1, cookie.Value).Scan(&userID)
	if err != nil {
		w.Write([]byte(`{"loggedIn": false}`))
		return
	}

	rows, err := databases.DB.Query("SELECT nickname FROM users WHERE id != ?", userID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var nicknames []string
	var nickname string
	for rows.Next() {
		err = rows.Scan(&nickname)
		if err != nil {
			log.Fatal(err)
		}
		nicknames = append(nicknames, nickname)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Println(nicknames)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nicknames)
}
