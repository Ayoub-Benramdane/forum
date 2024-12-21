package server

import (
	database "forum/Database"
	structs "forum/Data"
	"html/template"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func LogIn(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
		Errors(w, structs.Error{Code: http.StatusNotFound, Message: "Page not found"})
		return
	}
	switch r.Method {
	case http.MethodGet:
		LogInGet(w, r)
	case http.MethodPost:
		LogInPost(w, r)
	default:
		Errors(w, structs.Error{Code: http.StatusMethodNotAllowed, Message: "Method not allowed"})
		return
	}
}

func LogInGet(w http.ResponseWriter, r *http.Request) {
	if user := database.GetUserConnected(); user != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	tmpl, tmplErr := template.ParseFiles("./Template/html/login.html")
	if tmplErr != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Failed to load login page template"})
		return
	}
	tmpl.Execute(w, nil)
}

func LogInPost(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	user, errData := database.GetUserByUsername(username)
	if errData != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		Errors(w, structs.Error{Code: http.StatusUnauthorized, Message: "Check Username Or Password"})
		return
	}
	if errData := database.CreateSession(user.Username, user.ID); errData != nil {
		Errors(w, structs.Error{Code: http.StatusUnauthorized, Message: "Error Connection"})
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
