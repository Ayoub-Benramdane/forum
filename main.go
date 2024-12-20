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


	http.HandleFunc("/register", server.LogUp)
	http.HandleFunc("/login", server.LogIn)
	http.HandleFunc("/logout", server.LogOut)
	http.HandleFunc("/profil", server.Profil)
	http.HandleFunc("/editProfil", server.EditProfil)
	http.HandleFunc("/", server.Home)
	http.HandleFunc("/filter", server.Filter)
	http.HandleFunc("/post/", server.Post)
	http.HandleFunc("/like/", server.LikePost)
	http.HandleFunc("/dislike/", server.DislikePost)
	http.HandleFunc("/post/like_comment/", server.LikeComment)
	http.HandleFunc("/post/dislike_comment/", server.DislikeComment)
	http.HandleFunc("/post/delete/", server.DeletePost)
	http.HandleFunc("/post/edit/", server.EditPost)
	http.HandleFunc("/post/delete_comment/", server.DeleteComment)
	http.HandleFunc("/post/edit_comment/", server.EditComment)
	log.Println("Server is running...")
	log.Println("Link: http://localhost:8090")
	log.Fatal(http.ListenAndServe(":8090", nil))
}
