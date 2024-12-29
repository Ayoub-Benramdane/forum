package server

import (
	"html/template"
	"net/http"
	"strconv"

	structs "forum/Data"
	database "forum/Database"
)

var Posts = &structs.PostsShowing

func Page(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.ParseInt(r.URL.Path[len("/page/"):], 10, 64)
	if err != nil {
		Errors(w, structs.Error{Code: http.StatusBadRequest, Message: "Invalid page ID", Page: "Home", Path: "/"})
		return
	}
	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		Errors(w, structs.Error{Code: http.StatusMethodNotAllowed, Message: "Method not allowed", Page: "Home", Path: "/"})
		return
	}
	tmpl, tmplErr := template.ParseFiles("Template/html/home.html")
	if tmplErr != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Failed to load home page template", Page: "Home", Path: "/"})
		return
	}
	user := database.GetUserConnected()
	if user == nil {
		user = &structs.Session{Status: "Disconnected"}
	}
	x := (page - 1) * 10
	y := x + 10
	var posts []structs.Post
	if int64(len(*Posts)) > y {
		posts = (*Posts)[x:y]
	} else if int64(len(*Posts)) > x {
		posts = (*Posts)[x:]
	} else {
		posts = *Posts
	}
	categories, errLoadPost := database.GetAllCategorys()
	if errLoadPost != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error loading categories", Page: "Home", Path: "/"})
		return
	}
	pagination, errPage := Pagination([]string{"All"}, 0)
	if errPage != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error loading pagination", Page: "Home", Path: "/"})
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
