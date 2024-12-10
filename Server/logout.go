package server

import (
	database "forum/Database"
	structs "forum/Structs"
	"net/http"
)

func LogOut(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		Errors(w, structs.Error{Code: http.StatusMethodNotAllowed, Message: "Method not allowed"})
		return
	}
	database.DeleteSession()
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
