package server

import (
	"net/http"

	structs "forum/Data"
	database "forum/Database"
)

func Moderation(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil {
		Errors(w, structs.Error{Code: http.StatusNotFound, Message: "Page not found", Page: "Home", Path: "/profile"})
		return
	}
	user, err := database.GetUserConnected(cookie.Value)
	if err != nil || user.Role != "user" {
		if err != nil {
			http.SetCookie(w, &http.Cookie{Name: "session", Value: "", MaxAge: -1})
		}
		Errors(w, structs.Error{Code: http.StatusNotFound, Message: "Page not found", Page: "Home", Path: "/profile"})
		return
	}
	if r.Method != http.MethodGet {
		Errors(w, structs.Error{Code: http.StatusMethodNotAllowed, Message: "Method not allowed", Page: "Home", Path: "/profile"})
		return
	}
	if database.CreateRequest(user.ID) != nil {
		Errors(w, structs.Error{Code: http.StatusNotFound, Message: "Cannot send request", Path: "/profile"})
		return
	}
	user.Request = true
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}
