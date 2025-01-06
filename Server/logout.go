package server

import (
	"net/http"

	structs "forum/Data"
	database "forum/Database"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		Errors(w, structs.Error{Code: http.StatusMethodNotAllowed, Message: "Method not allowed", Page: "Home", Path: "/"})
		return
	}
	cookie, err := r.Cookie("session")
	if err != nil || cookie.Value == "" {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "No cookies", Page: "Home", Path: "/"})
		return
	}
	if database.DeleteSession(cookie.Value) != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error Ending Session", Page: "Home", Path: "/"})
		return
	}
	http.SetCookie(w, &http.Cookie{Name: "session", Value: "", MaxAge: -1})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
