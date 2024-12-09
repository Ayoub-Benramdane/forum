package server

import (
	database "forum/Database"
	"net/http"
)

func LogOut(w http.ResponseWriter, r *http.Request) {
	database.DeleteSession()
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
