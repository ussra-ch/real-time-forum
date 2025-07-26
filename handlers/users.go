package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"handlers/databases"
)

func IsLoggedIn(r *http.Request) (bool, int) {
	cookie, err := r.Cookie("session")
	if err != nil {
		return false, 0
	}

	var userID int
	err = databases.DB.QueryRow(`
		SELECT user_id FROM sessions 
		WHERE id = ? AND expires_at > DATETIME('now')
	`, cookie.Value).Scan(&userID)

	if err != nil {
		return false, 0
	}

	return true, userID
}
func FetchUsers(w http.ResponseWriter, r *http.Request) {
	loggedIn, userID := IsLoggedIn(r)
	if !loggedIn {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"loggedIn":    false,
			"nickname":    nil,
			"onlineUsers": []string{},
			"offlineUsers": []string{},
		})
		return
	}

	var myNickname string
	err := databases.DB.QueryRow("SELECT nickname FROM users WHERE id = ?", userID).Scan(&myNickname)
	if err != nil {
		log.Println("Error getting nickname:", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	rows, err := databases.DB.Query(`
		SELECT u.nickname 
		FROM users u
		JOIN sessions s ON u.id = s.user_id
		WHERE s.expires_at > DATETIME('now') AND u.id != ?
	`, userID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var onlineUsers []string
	for rows.Next() {
		var nickname string
		if err := rows.Scan(&nickname); err != nil {
			log.Fatal(err)
		}
		onlineUsers = append(onlineUsers, nickname)
	}

	row, err := databases.DB.Query(`
		SELECT nickname FROM users
		WHERE id != ? AND id NOT IN (
			SELECT user_id FROM sessions WHERE expires_at > DATETIME('now')
		)
	`, userID)
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()

	var offlineUsers []string
	for row.Next() {
		var nickname string
		if err := row.Scan(&nickname); err != nil {
			log.Fatal(err)
		}
		offlineUsers = append(offlineUsers, nickname)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"loggedIn":     true,
		"nickname":     myNickname,
		"onlineUsers":  onlineUsers,
		"offlineUsers": offlineUsers,
	})
}
