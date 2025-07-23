package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"handlers/databases"

	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	email := r.FormValue("email")
	password := r.FormValue("password")
	var dbPassword string
	err := databases.DB.QueryRow("SELECT password FROM users WHERE email = ?", email).Scan(&dbPassword)
	if err == sql.ErrNoRows {

		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	} else if err != nil {
		log.Println(err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if password != dbPassword {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}
	fmt.Println("password:", password, "dbPassword:", dbPassword)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// if r.Method == http.MethodPost {
	// 	// Handle registration logic here
	// 	http.Redirect(w, r, "/", http.StatusSeeOther)
	// 	return
	// }
	r.ParseForm()
	username := r.FormValue("Nickname")
	email := r.FormValue("email")
	gender := r.FormValue("gender")
	age := r.FormValue("age")
	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	password := r.FormValue("password")
	fmt.Println(username)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to encrypt password", http.StatusInternalServerError)
		return
	}
	_, err = databases.DB.Exec(`
		INSERT INTO users (nickname, age, gender, first_name, last_name, email, password)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		username, age, gender, firstName, lastName, email, hashedPassword)
	if err != nil {
		log.Println(err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
