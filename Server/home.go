package server

import (
	"html/template"
	"net/http"

	database "forum/Database"
	structs "forum/Data"
)

func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		Errors(w, structs.Error{Code: http.StatusNotFound, Message: "Page not found"})
		return
	} else if r.Method != http.MethodPost && r.Method != http.MethodGet {
		Errors(w, structs.Error{Code: http.StatusMethodNotAllowed, Message: "Method not allowed"})
		return
	}
	tmpl, tmplErr := template.ParseFiles("Template/home.html")
	if tmplErr != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Failed to load home page template"})
		return
	}
	user := database.GetUserConnected()
	if user == nil {
		user = &structs.Session{Status: "Disconnected"}
	}
	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		content := r.FormValue("content")
		if err := r.ParseForm(); err != nil {
			Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error parsing form"})
			return
		}
		categories := r.Form["category"]
		if errCrePost := database.CreatePost(title, content, categories, user.UserID); errCrePost != nil {
			Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error Creating post"})
			return
		}
	}
	posts, errLoadPost := database.GetAllPosts(user.Status)
	if errLoadPost != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error loading posts"})
		return
	}
	categories, errLoadPost := database.GetAllCategorys()
	if errLoadPost != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error loading categories"})
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
