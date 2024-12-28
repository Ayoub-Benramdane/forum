package server

import (
	structs "forum/Data"
	database "forum/Database"
	"html/template"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func Profile(w http.ResponseWriter, r *http.Request) {
	user := database.GetUserConnected()
	if r.URL.Path != "/profile" || user == nil {
		Errors(w, structs.Error{Code: http.StatusNotFound, Message: "Page not found", Page: "Home", Path: "/"})
		return
	} else if r.Method != http.MethodPost && r.Method != http.MethodGet {
		Errors(w, structs.Error{Code: http.StatusMethodNotAllowed, Message: "Method not allowed", Page: "Home", Path: "/"})
		return
	}
	tmpl, tmplErr := template.ParseFiles("Template/html/profile.html")
	if tmplErr != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Failed to load profil page template", Page: "Home", Path: "/"})
		return
	}
	info, errLoadInfo := database.GetInfoUser(user.UserID)
	if errLoadInfo != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error loading Info for user", Page: "Home", Path: "/"})
		return
	}
	tmpl.Execute(w, info)
}

func EditProfil(w http.ResponseWriter, r *http.Request) {
	user := database.GetUserConnected()
	if r.URL.Path != "/profile-edit" || user == nil {
		Errors(w, structs.Error{Code: http.StatusNotFound, Message: "Page not found", Page: "Home", Path: "/"})
		return
	}
	switch r.Method {
	case http.MethodGet:
		EditProfilGet(w, r)
	case http.MethodPost:
		EditProfilPost(w, r)
	default:
		Errors(w, structs.Error{Code: http.StatusMethodNotAllowed, Message: "Method not allowed", Page: "Profile", Path: "/profile"})
		return
	}
}

func EditProfilGet(w http.ResponseWriter, r *http.Request) {
	tmpl, tmplErr := template.ParseFiles("Template/html/profile-edit.html")
	if tmplErr != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error loading profil edit page", Page: "Profile", Path: "/profile"})
		return
	}
	user := database.GetUserConnected()
	info, errLoadInfo := database.GetInfoUser(user.UserID)
	if errLoadInfo != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error loading Info for user", Page: "Profile", Path: "/profile"})
		return
	}
	tmpl.Execute(w, info)
}

func EditProfilPost(w http.ResponseWriter, r *http.Request) {
	username := strings.TrimSpace(r.FormValue("username"))
	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")
	password1 := r.FormValue("new-password")
	password2 := r.FormValue("confirm-password")
	user := database.GetUserConnected()
	if password != "" {
		user, errData := database.GetUserByUsername(username)
		if errData != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
			Errors(w, structs.Error{Code: http.StatusUnauthorized, Message: "Check your Password", Page: "Profile edit", Path: "/profile-edit"})
			return
		}
	}
	if password1 == "" && password2 == "" {
		password1 = "Aa@11111"
	} else if password != password2 {
		Errors(w, structs.Error{Code: http.StatusConflict, Message: "Password not matched", Page: "Profile edit", Path: "/profile-edit"})
		return
	}
	if errSigne := validateSignupInput(username, email, password1); errSigne != nil {
		Errors(w, structs.Error{Code: http.StatusBadRequest, Message: errSigne.Error(), Page: "Profile edit", Path: "/profile-edit"})
		return
	}
	hashedPassword, errCrepting := bcrypt.GenerateFromPassword([]byte(password1), bcrypt.DefaultCost)
	if errCrepting != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error processing registration", Page: "Profile edit", Path: "/profile-edit"})
		return
	}

	if errUpdate := database.UpdateInfo(user.UserID, username, email); errUpdate != nil {
		if strings.Contains(errUpdate.Error(), "UNIQUE constraint failed") {
			Errors(w, structs.Error{Code: http.StatusConflict, Message: "Username already taken", Page: "Profile edit", Path: "/profile-edit"})
			return
		}
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error Updating user", Page: "Profile edit", Path: "/profile-edit"})
		return
	}
	if password1 != "" {
		if errUpdate := database.UpdatePass(user.UserID, string(hashedPassword)); errUpdate != nil {
			Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error Updating Password", Page: "Profile edit", Path: "/profile-edit"})
			return
		}
	}
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}
