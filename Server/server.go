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

var err structs.Error

func LogIn(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		err = structs.Error{Code: http.StatusNotFound, Message: "Page not found"}
		Errors(w, err)
		return
	}
	switch r.Method {
	case http.MethodGet:
		LogInGet(w, r)
	case http.MethodPost:
		LogInPost(w, r)
	default:
		err = structs.Error{Code: http.StatusMethodNotAllowed, Message: "Method not allowed"}
		Errors(w, err)
		return
	}
}

func LogInGet(w http.ResponseWriter, r *http.Request) {
	tmpl, tmplErr := template.ParseFiles("Template/sign in.html")
	if tmplErr != nil {
		err = structs.Error{Code: http.StatusInternalServerError, Message: "Error loading userin page"}
		Errors(w, err)
		return
	}
	tmpl.Execute(w, nil)
}

func LogInPost(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	user, errData := database.GetUserByUsername(username)
	if errData != nil || user.Password != password {
		err = structs.Error{Code: http.StatusUnauthorized, Message: "Check Username Or Password"}
		Errors(w, err)
		return
	}
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func LogUp(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/register" {
		err = structs.Error{Code: http.StatusNotFound, Message: "Page not found"}
		Errors(w, err)
		return
	}
	switch r.Method {
	case http.MethodGet:
		LogUpGet(w, r)
	case http.MethodPost:
		LogUpPost(w, r)
	default:
		err = structs.Error{Code: http.StatusMethodNotAllowed, Message: "Method not allowed"}
		Errors(w, err)
		return
	}
}

func LogUpGet(w http.ResponseWriter, r *http.Request) {
	tmpl, tmplErr := template.ParseFiles("Template/sign up.html")
	if tmplErr != nil {
		err = structs.Error{Code: http.StatusInternalServerError, Message: "Error loading signup page"}
		Errors(w, err)
		return
	}
	tmpl.Execute(w, nil)
}

func LogUpPost(w http.ResponseWriter, r *http.Request) {
	username := strings.TrimSpace(r.FormValue("username"))
	password := r.FormValue("password")
	if errSigne := validateSignupInput(username, password); errSigne != nil {
		err = structs.Error{Code: http.StatusBadRequest, Message: errSigne.Error()}
		Errors(w, err)
		return
	}
	hashedPassword, errCrepting := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if errCrepting != nil {
		err = structs.Error{Code: http.StatusInternalServerError, Message: "Error processing registration"}
		Errors(w, err)
		return
	}
	_, errCreate := database.CreateNewUser(username, string(hashedPassword))
	if errCreate != nil {
		if strings.Contains(errCreate.Error(), "UNIQUE constraint failed") {
			err = structs.Error{Code: http.StatusConflict, Message: "Username already taken"}
			Errors(w, err)
			return
		}
		err = structs.Error{Code: http.StatusInternalServerError, Message: "Error creating user"}
		Errors(w, err)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func validateSignupInput(username, password string) error {
	if len(username) < 3 || len(username) > 20 {
		return fmt.Errorf("username must be between 3 and 20 characters")
	}
	if !regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(username) {
		return fmt.Errorf("username can only contain letters, numbers, and underscores")
	}
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}
	if !regexp.MustCompile(`[A-Z]`).MatchString(password) ||
		!regexp.MustCompile(`[a-z]`).MatchString(password) ||
		!regexp.MustCompile(`[0-9]`).MatchString(password) ||
		!regexp.MustCompile(`[^a-zA-Z0-9]`).MatchString(password) {
		return fmt.Errorf("password must contain at least one uppercase letter, lowercase letter, number, and special character")
	}
	return nil
}

func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/home" {
		err = structs.Error{Code: http.StatusNotFound, Message: "Page not found"}
		Errors(w, err)
		return
	} else if r.Method != http.MethodPost {
		err = structs.Error{Code: http.StatusMethodNotAllowed, Message: "Method not allowed"}
		Errors(w, err)
		return
	}
	tmpl, tmplErr := template.ParseFiles("Template/home page.html")
	if tmplErr != nil {
		err = structs.Error{Code: http.StatusInternalServerError, Message: "Error loading home page"}
		Errors(w, err)
		return
	}
	tmpl.Execute(w, nil)
}

func Errors(w http.ResponseWriter, err structs.Error) {
	w.WriteHeader(err.Code)
	tmpl, tmplErr := template.ParseFiles("Template/errors.html")
	if tmplErr != nil {
		http.Error(w, "Status Internal Server Error 500", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, err)
}
