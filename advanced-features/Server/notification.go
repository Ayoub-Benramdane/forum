package server

import (
	"html/template"
	"net/http"

	structs "forum/Data"
	database "forum/Database"

	"golang.org/x/crypto/bcrypt"
)

func Notification(w http.ResponseWriter, r *http.Request) {
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
	switch r.Method {
	case http.MethodGet:
		NotificationGet(w, r, user)
	case http.MethodPost:
		NotificationPost(w, r)
	default:
		Errors(w, structs.Error{Code: http.StatusMethodNotAllowed, Message: "Method not allowed", Page: "Home", Path: "/"})
		return
	}
}

func NotificationGet(w http.ResponseWriter, r *http.Request, user *structs.User) {
	tmpl, tmplErr := template.ParseFiles("Template/html/notification.html")
	if tmplErr != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Failed to load login page template", Page: "Home", Path: "/"})
		return
	}
	data := struct {
		User          *structs.User
		Notifications []structs.Post
		Categories    []structs.Category
		Pagination    []int64
	}{
		User:          user,
		Notifications: nil,
		Categories:    nil,
		Pagination:    nil,
	}
	tmpl.Execute(w, data)
}

func NotificationPost(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	user, errData := database.GetUserByUsername(username)
	if errData != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		Errors(w, structs.Error{Code: http.StatusUnauthorized, Message: "Check Username Or Password", Page: "Login", Path: "/login"})
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
