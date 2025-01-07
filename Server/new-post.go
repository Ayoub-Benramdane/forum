package server

import (
	"html/template"
	"net/http"
	"strings"
	"time"

	structs "forum/Data"
	database "forum/Database"
)

func NewPost(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if r.URL.Path != "/new-post" || err != nil {
		Errors(w, structs.Error{Code: http.StatusNotFound, Message: "Page not found", Page: "Home", Path: "/"})
		return
	}
	switch r.Method {
	case http.MethodGet:
		NewPostGet(w, r)
	case http.MethodPost:
		NewPostPost(w, r, cookie)
	default:
		Errors(w, structs.Error{Code: http.StatusMethodNotAllowed, Message: "Method not allowed", Page: "Home", Path: "/"})
		return
	}
}

func NewPostGet(w http.ResponseWriter, r *http.Request) {
	tmpl, tmplErr := template.ParseFiles("Template/html/new-post.html")
	if tmplErr != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Failed to load new post page template", Page: "Home", Path: "/"})
		return
	}
	categories, errLoadPost := database.GetAllCategorys()
	if errLoadPost != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error loading categories", Page: "New-Post", Path: "/new-post"})
		return
	}
	data := struct {
		Categories []structs.Category
	}{
		Categories: categories,
	}
	tmpl.Execute(w, data)
}

func NewPostPost(w http.ResponseWriter, r *http.Request, cookie *http.Cookie) {
	user := database.GetUserConnected(cookie.Value)
	title := strings.TrimSpace(r.FormValue("title"))
	content := strings.TrimSpace(r.FormValue("content"))
	if title == "" || content == "" {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Check your input", Page: "New-Post", Path: "/new-post"})
		return
	}
	if err := r.ParseForm(); err != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error parsing form", Page: "New-Post", Path: "/new-post"})
		return
	}
	categories := r.Form["category"]
	if errCrePost := database.CreatePost(title, content, categories, user.UserID); errCrePost != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error Creating post", Page: "New-Post", Path: "/new-post"})
		return
	}
	token := cookie.Value
	cookie = &http.Cookie{
		Name:     "session",
		Value:    token,
		Expires:  time.Now().Add(5 * time.Minute),
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
