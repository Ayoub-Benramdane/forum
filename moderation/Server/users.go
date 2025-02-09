package server

import (
	structs "forum/Data"
	database "forum/Database"
	"net/http"
	"strconv"
)

func User(w http.ResponseWriter, r *http.Request) {
	id_user, err := strconv.ParseInt(r.URL.Path[len("/user/"):], 10, 64)
	if err != nil {
		Errors(w, structs.Error{Code: http.StatusBadRequest, Message: "Invalid user ID", Page: "Admin", Path: "/"})
		return
	} else if r.Method != http.MethodPost {
		Errors(w, structs.Error{Code: http.StatusMethodNotAllowed, Message: "Method not allowed", Page: "Admin", Path: "/"})
		return
	}
	cookie, err := r.Cookie("session")
	if err != nil {
		Errors(w, structs.Error{Code: http.StatusNotFound, Message: "you cant change role of this user", Page: "Home", Path: "/"})
		return
	}
	user, err := database.GetUserConnected(cookie.Value)
	if err != nil || user.Role != "admin" {
		if err != nil {
			http.SetCookie(w, &http.Cookie{Name: "session", Value: "", MaxAge: -1})
		}
		Errors(w, structs.Error{Code: http.StatusNotFound, Message: "you cant change role of this user", Page: "Home", Path: "/"})
		return
	}
	role := r.FormValue("role")
	if database.UpdateInfo(id_user, "", "", role) != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error Updating role user", Page: "Profile edit", Path: "/profile-edit"})
		return
	}
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
