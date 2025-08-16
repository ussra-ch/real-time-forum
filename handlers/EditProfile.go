package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"handlers/databases"
)

func EditProfile(w http.ResponseWriter, r *http.Request) {
	//PUT PATCH
	if r.Method != http.MethodPost {
		errorHandler(http.StatusMethodNotAllowed, w)
		return
	}

	mu.Lock()
	_, userID := IsLoggedIn(r)
	mu.Unlock()
	
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		errorHandler(http.StatusBadRequest, w)
		return
	}

	nickname := r.FormValue("nickname")
	email := r.FormValue("email")
	age := r.FormValue("age")

	var photoPath string
	file, handler, err := r.FormFile("photo")
	if err == nil {
		defer file.Close()
		photoPath = fmt.Sprintf("static/uploads/%d_%s", time.Now().UnixNano(), handler.Filename) /// TRAVERSAL ATTACK
		dst, err := os.Create(photoPath)
		if err != nil {
			errorHandler(http.StatusInternalServerError, w)
			return
		}
		defer dst.Close()
		_, err = io.Copy(dst, file)
		if err != nil {
			errorHandler(http.StatusInternalServerError, w)
			return
		}
	}
	state := false
	state1 := false

	if nickname == "" {
		_ = databases.DB.QueryRow("SELECT nickname FROM users WHERE id = ?", userID).Scan(&nickname)
		state = true
	}
	if email == "" {
		_ = databases.DB.QueryRow("SELECT email FROM users WHERE id = ?", userID).Scan(&email)
		state1 = true
	}
	if age == "" {
		_ = databases.DB.QueryRow("SELECT age FROM users WHERE id = ?", userID).Scan(&age)
	}
	if !state {
		var exists bool
		err = databases.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE nickname = ?)", nickname).Scan(&exists)
		if err != nil {
			log.Fatal(err)
		}
		if exists {
			errorr := ErrorStruct{
				Type: "error",
				Text: "Nickname is already in use",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errorr)
			return
		}
	}

	if !state1 {
		var exists bool
		err = databases.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", email).Scan(&exists)
		if err != nil {
			log.Fatal(err)
		}
		if exists {
			errorr := ErrorStruct{
				Type: "error",
				Text: "Email is already in use",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errorr)
			return
		}
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

	success := ErrorStruct{
		Type: "success",
		Text: "Your information has been updated",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(success)

}
