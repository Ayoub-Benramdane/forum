package server

import (
	"html/template"
	"net/http"
	"time"

	structs "forum/Data"
	database "forum/Database"

	"golang.org/x/crypto/bcrypt"
)

func LogIn(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
		Errors(w, structs.Error{Code: http.StatusNotFound, Message: "Page not found", Page: "Home", Path: "/"})
		return
	}
	switch r.Method {
	case http.MethodGet:
		LogInGet(w, r)
	case http.MethodPost:
		LogInPost(w, r)
	default:
		Errors(w, structs.Error{Code: http.StatusMethodNotAllowed, Message: "Method not allowed", Page: "Home", Path: "/"})
		return
	}
}

func LogInGet(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err == nil && database.GetUserConnected(cookie.Value) != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	tmpl, tmplErr := template.ParseFiles("Template/html/login.html")
	if tmplErr != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Failed to load login page template", Page: "Home", Path: "/"})
		return
	}
	tmpl.Execute(w, nil)
}

func LogInPost(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	user, errData := database.GetUserByUsername(username)
	if errData != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		Errors(w, structs.Error{Code: http.StatusUnauthorized, Message: "Check Username Or Password", Page: "Login", Path: "/login"})
		return
	}
	token := database.GenerateToken(username)
	if database.CreateSession(user.Username, user.ID, token) != nil {
		Errors(w, structs.Error{Code: http.StatusUnauthorized, Message: "Error Connection", Page: "Login", Path: "/login"})
		return
	}
	cookie := &http.Cookie{
		Name:     "session",
		Value:    token,
		Expires:  time.Now().Add(2 * time.Minute),
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
