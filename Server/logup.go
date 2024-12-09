package server

import (
	"fmt"
	database "forum/Database"
	structs "forum/Structs"
	"html/template"
	"net/http"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func LogUp(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/register" {
		Errors(w, structs.Error{Code: http.StatusNotFound, Message: "Page not found"})
		return
	}
	switch r.Method {
	case http.MethodGet:
		LogUpGet(w, r)
	case http.MethodPost:
		LogUpPost(w, r)
	default:
		Errors(w, structs.Error{Code: http.StatusMethodNotAllowed, Message: "Method not allowed"})
		return
	}
}

func LogUpGet(w http.ResponseWriter, r *http.Request) {
	if user := database.GetUserConnected(); user != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	tmpl, tmplErr := template.ParseFiles("Template/logup.html")
	if tmplErr != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error loading signup page"})
		return
	}
	tmpl.Execute(w, nil)
}

func LogUpPost(w http.ResponseWriter, r *http.Request) {
	username := strings.TrimSpace(r.FormValue("username"))
	password := r.FormValue("password")
	password2 := r.FormValue("confirm-password")
	if password != password2 {
		Errors(w, structs.Error{Code: http.StatusConflict, Message: "Password not matched"})
		return
	}
	if errSigne := validateSignupInput(username, password); errSigne != nil {
		Errors(w, structs.Error{Code: http.StatusBadRequest, Message: errSigne.Error()})
		return
	}
	hashedPassword, errCrepting := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if errCrepting != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error processing registration"})
		return
	}
	if errCreate := database.CreateNewUser(username, string(hashedPassword)); errCreate != nil {
		if strings.Contains(errCreate.Error(), "UNIQUE constraint failed") {
			Errors(w, structs.Error{Code: http.StatusConflict, Message: "Username already taken"})
			return
		}
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error creating user"})
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func validateSignupInput(username, password string) error {
	if len(username) < 3 || len(username) > 20 {
		return fmt.Errorf("username must be between 3 and 20 characters")
	} else if !regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(username) {
		return fmt.Errorf("username can only contain letters, numbers, and underscores")
	} else if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	} else if !regexp.MustCompile(`[A-Z]`).MatchString(password) || !regexp.MustCompile(`[a-z]`).MatchString(password) || !regexp.MustCompile(`[0-9]`).MatchString(password) || !regexp.MustCompile(`[^a-zA-Z0-9]`).MatchString(password) {
		return fmt.Errorf("password must contain at least one uppercase letter, lowercase letter, number, and special character")
	}
	return nil
}
