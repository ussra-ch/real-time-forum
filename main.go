package main

import (
	"fmt"
	"log"
	"net/http"

	"handlers/databases"
	"handlers/handlers"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func Test(w http.ResponseWriter, r *http.Request) {
	fs := http.FileServer(http.Dir("static"))
	if r.URL.Path == "/static/" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	http.StripPrefix("/static/", fs).ServeHTTP(w, r)
}

func main() {
	databases.InitDB("forum.db")
	defer databases.DB.Close()

	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/static/", Test)
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/login", handlers.RateLimitLoginMiddleware(handlers.LoginHandler))
	http.HandleFunc("/logout", handlers.LogoutHandler)
	http.HandleFunc("/api/logout", handlers.LogoutHandler)
	http.HandleFunc("/api/authenticated", handlers.IsAuthenticated)
	http.HandleFunc("/api/post", handlers.RateLimitPostsMiddleware(handlers.PostHandler))
	http.HandleFunc("/api/fetch_posts", handlers.FetchPostsHandler)
	http.HandleFunc("/comment", handlers.CommentsRatelimitMiddleware(handlers.CommentHandler))
	http.HandleFunc("/api/fetch_comments", handlers.FetchCommentsHandler)
	http.HandleFunc("/user", handlers.FetchUsers)
	http.HandleFunc("/chat", handlers.WebSocketHandler)
	http.HandleFunc("/api/fetchMessages", handlers.FetchMessages)
	http.HandleFunc("/delete", handlers.DeletePost)
	http.HandleFunc("/edit", handlers.EditPost)
	http.HandleFunc("/editProfile", handlers.EditProfile)
	fmt.Println("Server started at http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
