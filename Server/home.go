package server

import (
	database "forum/Database"
	structs "forum/Structs"
	"html/template"
	"net/http"
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
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error loading home page"})
		return
	}
	user := database.GetUserConnected()
	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		content := r.FormValue("content")
		if errCrePost := database.CreatePost(title, content, user.UserID); errCrePost != nil {
			Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error Creating post"})
			return
		}
	}
	posts, errLoadPost := database.GetAllPosts()
	if errLoadPost != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error loading posts"})
		return
	}
	if user == nil {
		user = &structs.Session{ID: 1, Username: "", UserID: 1, Statut: "Login"}
	}
	data := struct {
		User  *structs.Session
		Posts []structs.Post
	}{
		User:  user,
		Posts: posts,
	}
	tmpl.Execute(w, data)
}
