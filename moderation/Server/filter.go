package server

import (
	"html/template"
	"net/http"

	structs "forum/Data"
	database "forum/Database"
)

func Filter(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		Errors(w, structs.Error{Code: http.StatusMethodNotAllowed, Message: "Method not allowed", Page: "Home", Path: "/"})
		return
	}
	tmpl, tmplErr := template.ParseFiles("Template/html/home.html")
	if tmplErr != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Failed to load home page template", Page: "Home", Path: "/"})
		return
	}
	cookie, err := r.Cookie("session")
	var user *structs.User
	if err == nil {
		user, err = database.GetUserConnected(cookie.Value)
		if err != nil {
			http.SetCookie(w, &http.Cookie{Name: "session", Value: "", MaxAge: -1})
			user = &structs.User{Status: "Disconnected"}
		}
	} else {
		user = &structs.User{Status: "Disconnected"}
	}
	if err := r.ParseForm(); err != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error parsing form", Page: "Home", Path: "/"})
		return
	}
	categorie := r.Form["category"]
	posts, errLoadPost := database.GetFilterPosts(user.ID, categorie)
	if errLoadPost != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error loading posts", Page: "Home", Path: "/"})
		return
	}
	categories, errLoadCategory := database.GetAllCategories()
	if errLoadCategory != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error loading posts", Page: "Home", Path: "/"})
		return
	}
	pagination, errPage := Pagination(categorie, len(posts))
	if errPage != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error loading pagination", Page: "Home", Path: "/"})
		return
	}
	data := struct {
		User       *structs.User
		Posts      []structs.Post
		Categories []structs.Category
		Pagination []int64
	}{
		User:       user,
		Posts:      posts,
		Categories: categories,
		Pagination: pagination,
	}
	tmpl.Execute(w, data)
}
