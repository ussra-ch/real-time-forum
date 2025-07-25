package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"handlers/databases"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	Nickname := r.FormValue("Nickname")
	password := r.FormValue("password")
	var dbPassword string
	err := databases.DB.QueryRow("SELECT password FROM users WHERE nickname = ?", Nickname).Scan(&dbPassword)
	if err == sql.ErrNoRows {
		http.Error(w, "Invalid Nickname or password", http.StatusUnauthorized)
		return
	} else if err != nil {
		log.Println(err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(password)) != nil {
		http.Error(w, "Invalid Nickname or password", http.StatusUnauthorized)
		return
	}
	// Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    Nickname,
		Path:     "/",
		HttpOnly: true,
		MaxAge: 3600,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
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

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.FormValue("Nickname")
	email := r.FormValue("email")
	gender := r.FormValue("gender")
	age := r.FormValue("age")
	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	password := r.FormValue("password")

	var exists int
	err := databases.DB.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", email).Scan(&exists)
	if err != nil {
		log.Println("Error checking email:", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if exists > 0 {
		http.Error(w, "Email already in use", http.StatusConflict)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error encrypting password", http.StatusInternalServerError)
		return
	}

	res, err := databases.DB.Exec(`
		INSERT INTO users (nickname, age, gender, first_name, last_name, email, password)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		username, age, gender, firstName, lastName, email, hashedPassword)
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
		Name:     "sessionId",
		Value:    sessionID,
		Path:     "/",
		Expires:  expiresAt,
		HttpOnly: true,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func GetUserIDFromSession(r *http.Request) (int64, error) {
	cookie, err := r.Cookie("sessionId")
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
