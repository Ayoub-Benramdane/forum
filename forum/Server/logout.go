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
	if err != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "No cookies", Page: "Home", Path: "/"})
		return
	}
	user_id, err := database.GetUserFromToken(cookie.Value)
	if err != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "User not connected", Page: "Home", Path: "/"})
		return
	}
	if database.DeleteSession(user_id) != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error Ending Session", Page: "Home", Path: "/"})
		return
	}
	http.SetCookie(w, &http.Cookie{Name: "session", Value: "", MaxAge: -1})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
