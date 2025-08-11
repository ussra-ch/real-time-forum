package handlers

import (
	"fmt"
	"handlers/databases"
	"io"
	"net/http"
	"os"
	"time"
)

func EditProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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
	// fmt.Println(age)
	var photoPath string
	file, handler, err := r.FormFile("photo")
	if err == nil {
		defer file.Close()
		photoPath = fmt.Sprintf("static/uploads/%d_%s", time.Now().UnixNano(), handler.Filename)
		// fmt.Println(photoPath)
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

	query := "UPDATE users SET nickname = ?, email = ?, age = ?"
	args := []interface{}{nickname, email, age}

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
