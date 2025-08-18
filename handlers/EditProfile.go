package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"handlers/databases"
)

func EditProfile(w http.ResponseWriter, r *http.Request) {
	// PUT PATCH
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

	if len(strings.TrimSpace(nickname)) == 0 {
		_ = databases.DB.QueryRow("SELECT nickname FROM users WHERE id = ?", userID).Scan(&nickname)
	}

	if len(strings.TrimSpace(email)) == 0 {
		_ = databases.DB.QueryRow("SELECT email FROM users WHERE id = ?", userID).Scan(&email)
	} else {
		emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

		re := regexp.MustCompile(emailRegex)
		if !re.MatchString(email) {
			errorr := ErrorStruct{
				Type: "error",
				Text: "Invalid email",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errorr)
			return
		}
	}
	tmpAge, _ := strconv.Atoi(age)
	if tmpAge == 0 {
		_ = databases.DB.QueryRow("SELECT age FROM users WHERE id = ?", userID).Scan(&age)
	} else if tmpAge < 13 || tmpAge > 120 {
		errorr := ErrorStruct{
				Type: "error",
				Text: "Please provide a valid age",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errorr)
			return
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
		errorr := ErrorStruct{
			Type: "error",
			Text: "The nickname or email is not available. Please try another.",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorr)
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
