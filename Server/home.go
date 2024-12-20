package server

import (
	"html/template"
	"net/http"
	"strings"

	structs "forum/Data"
	database "forum/Database"
)

func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		Errors(w, structs.Error{Code: http.StatusNotFound, Message: "Page not found"})
		return
	} else if r.Method != http.MethodPost && r.Method != http.MethodGet {
		Errors(w, structs.Error{Code: http.StatusMethodNotAllowed, Message: "Method not allowed"})
		return
	}
	tmpl, tmplErr := template.ParseFiles("Template/html/home.html")
	if tmplErr != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Failed to load home page template"})
		return
	}
	user := database.GetUserConnected()
	if user == nil {
		user = &structs.Session{Status: "Disconnected"}
	}
	if r.Method == http.MethodPost {
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
	totalPosts := len(posts)
	for i := range posts {
		posts[i].TotalPosts = int64(totalPosts)
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
