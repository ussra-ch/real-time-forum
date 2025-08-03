package handlers

import (
	"database/sql"
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
			"loggedIn":     false,
			"nickname":     nil,
			"onlineUsers":  []string{},
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
		SELECT u.nickname, u.id, u.photo
		FROM users u
		JOIN sessions s ON u.id = s.user_id
		WHERE s.expires_at > DATETIME('now') AND u.id != ?
	`, userID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	type User struct {
		Nickname string         `json:"nickname"`
		UserId   int            `json:"userId"`
		Photo    sql.NullString `json:"photo"`
	}
	var onlineUsers []User
	for rows.Next() {
		var photo sql.NullString
		var nickname string
		var userId int
		if err := rows.Scan(&nickname, &userId, &photo); err != nil {
			log.Fatal("error", err)
		}

		onlineUsers = append(onlineUsers, User{Nickname: nickname, UserId: userId, Photo: photo})
	}

	row, err := databases.DB.Query(`
		SELECT u.nickname, u.id
		FROM users u
		WHERE id != ? AND id NOT IN (
			SELECT user_id FROM sessions WHERE expires_at > DATETIME('now')
		)
	`, userID)
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()

	var offlineUsers []User
	for row.Next() {
		// fmt.Println("row is :", rows)
		var nickname string
		var userId int
		if err := row.Scan(&nickname, &userId); err != nil {
			// fmt.Println("121212")
			log.Fatal(err)
		}
		// fmt.Println(User{nickname: nickname, userId: userId})
		offlineUsers = append(offlineUsers, User{Nickname: nickname, UserId: userId})
	}

	w.Header().Set("Content-Type", "application/json")
	// fmt.Println(offlineUsers)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"loggedIn":     true,
		"nickname":     myNickname,
		"onlineUsers":  onlineUsers,
		"offlineUsers": offlineUsers,
		"UserId":       userID,
	})
}
