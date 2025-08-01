package handlers

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"handlers/databases"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func generateSessionID() string {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}

type a struct {
	Nickname string `json:Nickname`
	Password string `json:password`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var aa a
	err := json.NewDecoder(r.Body).Decode(&aa)
	// fmt.Println(aa.Nickname)
	var dbPassword string
	var userID int
	err = databases.DB.QueryRow("SELECT id, password FROM users WHERE nickname = ?", aa.Nickname).Scan(&userID, &dbPassword)
	if err == sql.ErrNoRows || bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(aa.Password)) != nil {
		http.Error(w, "Invalid Nickname or password", http.StatusUnauthorized)
		return
	} else if err != nil {
		log.Println(err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	sessionID := generateSessionID()
	_, err = databases.DB.Exec(`
		INSERT INTO sessions (id, user_id, expires_at)
		VALUES (?, ?, DATETIME('now', '+1 hour'))
	`, sessionID, userID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   3600,
	})
}

func IsAuthenticated(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil || cookie.Value == "" {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}
	w.Write([]byte(cookie.Value))
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	_, err = databases.DB.Exec("DELETE FROM sessions WHERE id = ?", cookie.Value)
	if err != nil {
		log.Println("Failed to delete session:", err)
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})
	w.WriteHeader(http.StatusOK)
}

type data struct {
	Nickname  string `json:Nickname`
	Email     string `json:email`
	Gender    string `json:gender`
	Age       string `json:age`
	Firstname string `json:first_name`
	Lastname  string `json:last_name`
	Password  string `json:password`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var aa data
	err := json.NewDecoder(r.Body).Decode(&aa)
	var exists int
	err = databases.DB.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", aa.Email).Scan(&exists)
	if err != nil {
		log.Println("Error checking email:", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if exists > 0 {
		http.Error(w, "Email already in use", http.StatusConflict)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(aa.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error encrypting password", http.StatusInternalServerError)
		return
	}

	res, err := databases.DB.Exec(`
		INSERT INTO users (nickname, age, gender, first_name, last_name, email, password)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		aa.Nickname, aa.Age, aa.Gender, aa.Firstname, aa.Lastname, aa.Email, hashedPassword)
	if err != nil {
		log.Println("Error inserting user:", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	userID, err := res.LastInsertId()
	if err != nil {
		log.Println("Error getting inserted user ID:", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	sessionID := uuid.New().String()
	expiresAt := time.Now().Add(24 * time.Hour)

	_, err = databases.DB.Exec(`
		INSERT INTO sessions (id, user_id, expires_at)
		VALUES (?, ?, ?)`,
		sessionID, userID, expiresAt)
	if err != nil {
		log.Println("Error inserting session:", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    sessionID,
		Path:     "/",
		Expires:  expiresAt,
		HttpOnly: true,
	})
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
