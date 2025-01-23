package main

import (
	"forum/Database"
	"forum/Server"
	"log"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) != 1 {
		return
	} else if err := database.ConnectDatabase(); err != nil {
        log.Fatalf("Failed to initialize database: %v", err)
    }

	fs := http.FileServer(http.Dir("./Template"))
	http.Handle("/Template/", http.StripPrefix("/Template/", fs))

	http.HandleFunc("/login", server.Login)
	http.HandleFunc("/register", server.Register)
	http.HandleFunc("/logout", server.Logout)
	http.HandleFunc("/profile", server.Profile)
	http.HandleFunc("/profile-edit", server.EditProfile)
	http.HandleFunc("/", server.Home)
	http.HandleFunc("/page/", server.Page)
	http.HandleFunc("/filter", server.Filter)
	http.HandleFunc("/post/", server.Post)
	http.HandleFunc("/new-post", server.NewPost)
	http.HandleFunc("/like/", server.LikePost)
	http.HandleFunc("/dislike/", server.DislikePost)
	http.HandleFunc("/like_comment/", server.LikeComment)
	http.HandleFunc("/dislike_comment/", server.DislikeComment)
	http.HandleFunc("/delete/", server.DeletePost)
	http.HandleFunc("/edit/", server.EditPost)
	http.HandleFunc("/delete_comment/", server.DeleteComment)
	http.HandleFunc("/edit_comment/", server.EditComment)
	log.Println("Server is running...")
	log.Println("Link: http://localhost:4444")
	log.Fatal(http.ListenAndServe(":4444", nil))
}
