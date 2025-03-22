package server

import (
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"strings"

	structs "forum/Data"
	database "forum/Database"

	"golang.org/x/crypto/bcrypt"
)

func Register(w http.ResponseWriter, r *http.Request) {
	if _, err := r.Cookie("session"); err == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	switch r.Method {
	case http.MethodGet:
		RegisterGet(w, r)
	case http.MethodPost:
		RegisterPost(w, r)
	default:
		Errors(w, structs.Error{Code: http.StatusMethodNotAllowed, Message: "Method not allowed", Page: "Home", Path: "/"})
		return
	}
}

func RegisterGet(w http.ResponseWriter, r *http.Request) {
	tmpl, tmplErr := template.ParseFiles("Template/html/register.html")
	if tmplErr != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Failed to load Register page template", Page: "Home", Path: "/"})
		return
	}
	tmpl.Execute(w, nil)
}

func RegisterPost(w http.ResponseWriter, r *http.Request) {
	username := strings.TrimSpace(r.FormValue("username"))
	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")
	password2 := r.FormValue("confirm-password")
	if password != password2 {
		Errors(w, structs.Error{Code: http.StatusConflict, Message: "Password not matched", Page: "Register", Path: "/register"})
		return
	}
	if errSigne := validateSignupInput(username, email, password); errSigne != nil {
		Errors(w, structs.Error{Code: http.StatusBadRequest, Message: errSigne.Error(), Page: "Register", Path: "/register"})
		return
	}
	hashedPassword, errCrepting := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if errCrepting != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error processing registration", Page: "Register", Path: "/register"})
		return
	}
	if errCreate := database.CreateNewUser(username, email, string(hashedPassword)); errCreate != nil {
		fmt.Println(errCreate)
		if strings.Contains(errCreate.Error(), "UNIQUE constraint failed") {
			Errors(w, structs.Error{Code: http.StatusConflict, Message: "Username or email already taken", Page: "Register", Path: "/register"})
			return
		}
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error creating user", Page: "Register", Path: "/register"})
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func validateSignupInput(username, email, password string) error {
	if len(username) < 3 || len(username) > 20 {
		return fmt.Errorf("username must be between 3 and 20 characters")
	} else if !regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(username) {
		return fmt.Errorf("username can only contain letters, numbers, and underscores")
	} else if !regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`).MatchString(email) {
		return fmt.Errorf("please enter a valid email address")
	} else if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	} else if !regexp.MustCompile(`[A-Z]`).MatchString(password) || !regexp.MustCompile(`[a-z]`).MatchString(password) || !regexp.MustCompile(`[0-9]`).MatchString(password) || !regexp.MustCompile(`[^a-zA-Z0-9]`).MatchString(password) {
		return fmt.Errorf("password must contain at least one uppercase letter, lowercase letter, number, and special character")
	}
	return nil
}
