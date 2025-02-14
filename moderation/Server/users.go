package server

import (
	structs "forum/Data"
	database "forum/Database"
	"net/http"
	"strconv"
	"strings"
)

func Users(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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
	slc := strings.Split(r.URL.Path[len("/users/"):], "/")
	if len(slc) == 2 {
		user_id, err := strconv.ParseInt(slc[0], 10, 64)
		if err != nil {
			Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error Updating role user", Page: "Profile edit", Path: "/profile-edit"})
			return
		} else if slc[1] != "user" && slc[1] != "guest" && slc[1] != "moderateur" {
			Errors(w, structs.Error{Code: http.StatusNotFound, Message: "role not found", Page: "Profile edit", Path: "/profile-edit"})
			return
		} else if database.UpdateInfo(user_id, "", "", slc[1]) != nil {
			Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error Updating role user", Page: "Profile edit", Path: "/profile-edit"})
			return
		}
	} else {
		Errors(w, structs.Error{Code: http.StatusNotFound, Message: "you cant change role of this user", Page: "Home", Path: "/"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
