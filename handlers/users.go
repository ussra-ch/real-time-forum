package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"handlers/databases"
)

func FetchUsers(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil {
		fmt.Println("No sessionId cookie:", err)
		w.Write([]byte(`{"loggedIn": false}`))
		return
	}

	query := `SELECT user_id FROM sessions WHERE id = ? AND expires_at > DATETIME('now')`
	var userID int
	err = databases.DB.QueryRow(query, cookie.Value).Scan(&userID)
	if err == sql.ErrNoRows {
		fmt.Println("No session found in DB for:", cookie.Value)
		w.Write([]byte(`{"loggedIn": false}`))
		return
	} else if err != nil {
		log.Println("DB error:", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
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
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nicknames)
}
