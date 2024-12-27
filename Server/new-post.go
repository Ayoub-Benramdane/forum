package server

import (
	structs "forum/Data"
	database "forum/Database"
	"html/template"
	"net/http"
	"strings"
)

func NewPost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/new-post" {
		Errors(w, structs.Error{Code: http.StatusNotFound, Message: "Page not found"})
		return
	}
	switch r.Method {
	case http.MethodGet:
		NewPostGet(w, r)
	case http.MethodPost:
		NewPostPost(w, r)
	default:
		Errors(w, structs.Error{Code: http.StatusMethodNotAllowed, Message: "Method not allowed"})
		return
	}
}

func NewPostGet(w http.ResponseWriter, r *http.Request) {
	if user := database.GetUserConnected(); user == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	tmpl, tmplErr := template.ParseFiles("Template/html/new-post.html")
	if tmplErr != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Failed to load new post page template"})
		return
	}
	categories, errLoadPost := database.GetAllCategorys()
	if errLoadPost != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error loading categories"})
		return
	}
	data := struct {
		Categories []structs.Category
	}{
		Categories: categories,
	}
	tmpl.Execute(w, data)
}

func NewPostPost(w http.ResponseWriter, r *http.Request) {
	user := database.GetUserConnected()
	if user == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	title := strings.TrimSpace(r.FormValue("title"))
	content := strings.TrimSpace(r.FormValue("content"))
	if title == "" || content == "" {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Check your input"})
		return
	}
	if err := r.ParseForm(); err != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error parsing form"})
		return
	}
	categories := r.Form["category"]
	if errCrePost := database.CreatePost(title, content, categories, user.UserID); errCrePost != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error Creating post"})
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
