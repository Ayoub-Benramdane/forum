package server

import (
	"html/template"
	"net/http"

	database "forum/Database"
	structs "forum/Data"
)

func Filter(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/filter" {
		Errors(w, structs.Error{Code: http.StatusNotFound, Message: "Page not found"})
		return
	} else if r.Method != http.MethodPost {
		Errors(w, structs.Error{Code: http.StatusMethodNotAllowed, Message: "Method not allowed"})
		return
	}
	tmpl, tmplErr := template.ParseFiles("./Template/html/home.html")
	if tmplErr != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Failed to load home page template"})
		return
	}
	user := database.GetUserConnected()
	if user == nil {
		user = &structs.Session{Status: "Disconnected"}
	}
	if err := r.ParseForm(); err != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error parsing form"})
		return
	}
	categorie := r.Form["category"]
	posts, errLoadPost := database.GetFilterPosts(user, categorie)
	if errLoadPost != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error loading posts"})
		return
	}
	categories, errLoadCategory := database.GetAllCategorys()
	if errLoadCategory != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error loading posts"})
		return
	}
	data := struct {
		User       *structs.Session
		Posts      []structs.Post
		Categories []structs.Category
	}{
		User:       user,
		Posts:      posts,
		Categories: categories,
	}
	tmpl.Execute(w, data)
}
