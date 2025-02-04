package server

import (
	"net/http"
	"strings"
	"time"

	structs "forum/Data"
	database "forum/Database"

	"golang.org/x/crypto/bcrypt"
)

func RegisterPostAuth(username, email, password string) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	database.CreateNewUser(username, email, string(hashedPassword))
}

func LoginPostAuth(w http.ResponseWriter, r *http.Request, username, password string) {
	user, _ := database.GetUserByUsername(username)

	hashedUser, _ := bcrypt.GenerateFromPassword([]byte(username), bcrypt.DefaultCost)

	tkn, _ := database.GenerateToken(string(hashedUser))

	token := string(tkn)
	if err := database.CreateSession(user.Username, user.ID, token); err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			if database.DeleteSession(username) != nil {
				Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error Ending Session", Page: "Home", Path: "/"})
				return
			}
			http.SetCookie(w, &http.Cookie{Name: "session", Value: "", MaxAge: -1})
			if database.CreateSession(user.Username, user.ID, token) != nil {
				Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error Connection", Page: "Home", Path: "/"})
				return
			}
		} else {
			Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error Connection", Page: "Login", Path: "/login"})
			return
		}
	}

	cookie := &http.Cookie{
		Name:     "session",
		Value:    token,
		Expires:  time.Now().Add(5 * time.Minute),
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
