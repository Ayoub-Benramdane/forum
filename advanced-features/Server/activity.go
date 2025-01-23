package server

import (
	"html/template"
	"net/http"

	structs "forum/Data"
	database "forum/Database"
)

func Activity(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil {
		Errors(w, structs.Error{Code: http.StatusNotFound, Message: "Page not found", Page: "Home", Path: "/"})
		return
	}
	user, err := database.GetUserConnected(cookie.Value)
	if err != nil {
		http.SetCookie(w, &http.Cookie{Name: "session", Value: "", MaxAge: -1})
		Errors(w, structs.Error{Code: http.StatusNotFound, Message: "Page not found", Page: "Home", Path: "/"})
		return
	}
	if r.Method != http.MethodGet {
		Errors(w, structs.Error{Code: http.StatusMethodNotAllowed, Message: "Method not allowed", Page: "Home", Path: "/"})
		return
	}
	tmpl, tmplErr := template.ParseFiles("Template/html/activity.html")
	if tmplErr != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Failed to load activity page template", Page: "Home", Path: "/"})
		return
	}
	all_data, err := database.GetData(user.ID)
	if err != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Failed to load activity data", Page: "Home", Path: "/"})
		return
	}
	data := struct {
		User            *structs.User
		All             *structs.All
		Posts           []structs.Post
		Comments        []structs.Comment
		ReactedPosts    []structs.Post
		ReactedComments []structs.Comment
	}{
		User:            user,
		All:             all_data,
		Posts:           all_data.Posts,
		Comments:        all_data.Comments,
		ReactedPosts:    all_data.ReactedPosts,
		ReactedComments: all_data.ReactedComments,
	}
	tmpl.Execute(w, data)
}