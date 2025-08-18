package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"handlers/databases"
)

func errorHandler(errorType int, w http.ResponseWriter) {
	errorr := ErrorStruct{
		Type: "error",
		Text: http.StatusText(errorType),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errorType)
	json.NewEncoder(w).Encode(errorr)
}

func generateSessionID() string {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}

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

func ProtectStaticDir(w http.ResponseWriter, r *http.Request) {
	fs := http.FileServer(http.Dir("static"))
	path := r.URL.Path
	if path == "/static/" || path == "/static/uploads/" {
		errorHandler(http.StatusForbidden, w)
		return
	}

	http.StripPrefix("/static/", fs).ServeHTTP(w, r)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}
