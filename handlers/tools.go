package handlers

import (
	"net/http"

	"handlers/databases"
)

func GetUserIDFromSession(r *http.Request) (int64, error) {
	cookie, err := r.Cookie("session")
	if err != nil {
		return 0, err
	}

	var userID int64
	err = databases.DB.QueryRow(`
		SELECT user_id FROM sessions
		WHERE id = ? AND expires_at > datetime('now')`, cookie.Value).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func IsAuthenticated(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil || cookie.Value == "" {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}
	w.Write([]byte(cookie.Value))
}

