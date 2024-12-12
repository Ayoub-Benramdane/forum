package server

import (
	"html/template"
	"net/http"

	database "forum/Database"
	structs "forum/Structs"
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
		category := r.FormValue("category")
		if errCrePost := database.CreatePost(title, content, category, user.UserID); errCrePost != nil {
			Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error Creating post"})
			return
		}
	}
	if user == nil {
		user = &structs.Session{ID: 1, Username: "", UserID: 1, Status: "Login"}
	}
	posts, errLoadPost := database.GetAllPosts(user.Status)
	if errLoadPost != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error loading posts"})
		return
	}
	categorys, errLoadPost := database.GetAllCategorys()
	if errLoadPost != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error loading posts"})
		return
	}
	data := struct {
		User      *structs.Session
		Posts     []structs.Post
		Categorys []structs.Category
	}{
		User:      user,
		Posts:     posts,
		Categorys: categorys,
	}
	tmpl.Execute(w, data)
}
