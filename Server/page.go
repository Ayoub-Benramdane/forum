package server

import (
	"html/template"
	"net/http"
	"strconv"

	structs "forum/Data"
	database "forum/Database"
)

func Page(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.ParseInt(r.URL.Path[len("/page/"):], 10, 64)
	if err != nil {
		Errors(w, structs.Error{Code: http.StatusBadRequest, Message: "Invalid page ID"})
		return
	}
	if r.Method != http.MethodPost && r.Method != http.MethodGet {
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
	posts, errLoadPost := database.GetAllPosts(user.Status, 20, page-1)
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
