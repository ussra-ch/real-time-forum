package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"handlers/databases"
)

func EditProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", 301)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	_, userID := IsLoggedIn(r)

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		fmt.Println(err)
		http.Error(w, "Cannot parse form", http.StatusBadRequest)
		return
	}

	nickname := r.FormValue("nickname")
	email := r.FormValue("email")
	age := r.FormValue("age")

	var photoPath string
	file, handler, err := r.FormFile("photo")
	if err == nil {
		defer file.Close()
		photoPath = fmt.Sprintf("static/uploads/%d_%s", time.Now().UnixNano(), handler.Filename)
		dst, err := os.Create(photoPath)
		if err != nil {
			http.Error(w, "Failed to save photo", http.StatusInternalServerError)
			return
		}
		defer dst.Close()
		_, err = io.Copy(dst, file)
		if err != nil {
			http.Error(w, "Failed to save photo", http.StatusInternalServerError)
			return
		}
	}

	if nickname == "" {
		_ = databases.DB.QueryRow("SELECT nickname FROM users WHERE id = ?", userID).Scan(&nickname)
	}
	if email == "" {
		_ = databases.DB.QueryRow("SELECT email FROM users WHERE id = ?", userID).Scan(&email)
	}
	if age == "" {
		_ = databases.DB.QueryRow("SELECT age FROM users WHERE id = ?", userID).Scan(&age)
	}
	args := []interface{}{nickname, email, age}
	query := "UPDATE users SET nickname = ?, email = ?, age = ?"

	if photoPath != "" {
		query += ", photo = ?"
		args = append(args, photoPath)
	}
	query += " WHERE id = ?"
	args = append(args, userID)

	_, err = databases.DB.Exec(query, args...)
	if err != nil {
		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
