package server

import (
	"html/template"
	"math"
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
	posts, errLoadPost := database.GetAllPosts(user.Status, 20, 0)
	if errLoadPost != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error loading posts"})
		return
	}
	categories, errLoadPost := database.GetAllCategorys()
	if errLoadPost != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error loading categories"})
		return
	}
	pagination, errPage := Pagination()
	if errPage != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error loading pagination"})
		return
	}
	data := struct {
		User       *structs.Session
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

func Pagination() ([]int64, error) {
	totalPosts, err := database.CountPosts()
	if err != nil {
		return nil, err
	}
	var pagination []int64
	for i := int64(1); i <= int64(math.Ceil(totalPosts/20)); i++ {
		pagination = append(pagination, i)
	}
	return pagination, nil
}
