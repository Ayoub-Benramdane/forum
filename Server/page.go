package server

import (
	"html/template"
	"math"
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
	} else if page > int64(len(*Posts)+9)/10 {
		Errors(w, structs.Error{Code: http.StatusNotFound, Message: "Page Not Found", Page: "Home", Path: "/"})
		return
	} else if r.Method != http.MethodGet {
		Errors(w, structs.Error{Code: http.StatusMethodNotAllowed, Message: "Method not allowed", Page: "Home", Path: "/"})
		return
	}
	tmpl, tmplErr := template.ParseFiles("Template/html/home.html")
	if tmplErr != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Failed to load home page template", Page: "Home", Path: "/"})
		return
	}
	cookie, err := r.Cookie("session")
	var user *structs.Session
	if err == nil {
		user = database.GetUserConnected(cookie.Value)
		if user == nil {
			http.SetCookie(w, &http.Cookie{Name: "session", Value: "", MaxAge: -1})
			user = &structs.Session{Status: "Disconnected"}
		}
	} else {
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

func Pagination(categories []string, posts int) ([]int64, error) {
	var totalPosts float64
	var err error
	var pagination []int64
	for _, cat := range categories {
		if cat == "All" {
			totalPosts, err = database.CountPosts()
			if err != nil {
				return nil, err
			}
			break
		}
		totalPosts = float64(posts)
	}
	for i := int64(1); i <= int64(math.Ceil(totalPosts/10)); i++ {
		pagination = append(pagination, i)
	}
	return pagination, nil
}
