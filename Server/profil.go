package server

import (
	database "forum/Database"
	structs "forum/Structs"
	"html/template"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func Profil(w http.ResponseWriter, r *http.Request) {
	user := database.GetUserConnected()
	if r.URL.Path != "/profil" || user == nil {
		Errors(w, structs.Error{Code: http.StatusNotFound, Message: "Page not found"})
		return
	} else if r.Method != http.MethodPost && r.Method != http.MethodGet {
		Errors(w, structs.Error{Code: http.StatusMethodNotAllowed, Message: "Method not allowed"})
		return
	}
	tmpl, tmplErr := template.ParseFiles("Template/profil.html")
	if tmplErr != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Failed to load profil page template"})
		return
	}
	info, errLoadInfo := database.GetInfoUser(user.UserID)
	if errLoadInfo != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error loading Info for user"})
		return
	}
	tmpl.Execute(w, info)
}

func EditProfil(w http.ResponseWriter, r *http.Request) {
	user := database.GetUserConnected()
	if r.URL.Path != "/edit" || user == nil {
		Errors(w, structs.Error{Code: http.StatusNotFound, Message: "Page not found"})
		return
	}
	switch r.Method {
	case http.MethodGet:
		EditProfilGet(w, r)
	case http.MethodPost:
		EditProfilPost(w, r)
	default:
		Errors(w, structs.Error{Code: http.StatusMethodNotAllowed, Message: "Method not allowed"})
		return
	}
}

func EditProfilGet(w http.ResponseWriter, r *http.Request) {
	tmpl, tmplErr := template.ParseFiles("Template/edit.html")
	if tmplErr != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error loading profil edit page"})
		return
	}
	user := database.GetUserConnected()
	info, errLoadInfo := database.GetInfoUser(user.UserID)
	if errLoadInfo != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error loading Info for user"})
		return
	}
	tmpl.Execute(w, info)
}

func EditProfilPost(w http.ResponseWriter, r *http.Request) {
	username := strings.TrimSpace(r.FormValue("username"))
	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")
	password2 := r.FormValue("confirm-password")
	if password == "" && password2 == "" {
		password = "Aa@11111"
	} else if password != password2 {
		Errors(w, structs.Error{Code: http.StatusConflict, Message: "Password not matched"})
		return
	}
	if errSigne := validateSignupInput(username, email, password); errSigne != nil {
		Errors(w, structs.Error{Code: http.StatusBadRequest, Message: errSigne.Error()})
		return
	}
	hashedPassword, errCrepting := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if errCrepting != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error processing registration"})
		return
	}
	user := database.GetUserConnected()
	if errUpdate := database.UpdateInfo(user.UserID, username, email); errUpdate != nil {
		if strings.Contains(errUpdate.Error(), "UNIQUE constraint failed") {
			Errors(w, structs.Error{Code: http.StatusConflict, Message: "Username already taken"})
			return
		}
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error Updating user"})
		return
	}
	if password2 != "" {
		if errUpdate := database.UpdatePass(user.UserID, string(hashedPassword)); errUpdate != nil {
			Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error Updating Password"})
			return
		}
	}
	http.Redirect(w, r, "/profil", http.StatusSeeOther)
}
