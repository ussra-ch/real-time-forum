package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"handlers/databases"
)

type Post struct {
	ID int `json:"id"`
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	var post Post

	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	query := "DELETE FROM posts WHERE id = ?"
	_, err = databases.DB.Exec(query, post.ID)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Post deleted successfully"))
}

func EditPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var data struct {
		ID      int    `json:"id"`
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	_, err := databases.DB.Exec("UPDATE posts SET title = ?, content = ? WHERE id = ?", data.Title, data.Content, data.ID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Post deleted successfully"))
}

func EditProfil(w http.ResponseWriter, r *http.Request) {
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
