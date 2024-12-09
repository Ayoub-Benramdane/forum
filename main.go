package main

import (
	"forum/Database"
	"forum/Server"
	"log"
	"net/http"
)

func main() {
    if err := database.ConnectDatabase(); err != nil {
        log.Fatalf("Failed to initialize database: %v", err)
    }
	http.HandleFunc("/register", server.LogUp)
	http.HandleFunc("/login", server.LogIn)
	http.HandleFunc("/logout", server.LogOut)
	http.HandleFunc("/", server.Home)
	http.HandleFunc("/post/", server.Post)
	log.Println("Server is running...")
	log.Println("Link: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
