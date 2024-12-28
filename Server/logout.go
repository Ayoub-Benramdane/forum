package server

import (
	database "forum/Database"
	structs "forum/Data"
	"net/http"
)

func LogOut(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		Errors(w, structs.Error{Code: http.StatusMethodNotAllowed, Message: "Method not allowed", Page: "Home", Path: "/"})
		return
	}
	if database.DeleteSession() != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error Ending Session", Page: "Home", Path: "/"})
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
