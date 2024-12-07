package main

import (
	"fmt"
	fs "forum/Function"
	"net/http"
)

func main() {
	http.HandleFunc("/", fs.Login)
	http.HandleFunc("/home", fs.Home)
	fmt.Println("Server is Runing...")
	fmt.Println("Link : http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
}
