package server

import (
	"fmt"
	"html/template"
	"net/http"

	structs "forum/Data"
	database "forum/Database"
)

func Notification(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil {
		Errors(w, structs.Error{Code: http.StatusNotFound, Message: "Page not found", Page: "Home", Path: "/"})
		return
	}
	user, err := database.GetUserConnected(cookie.Value)
	if err != nil {
		http.SetCookie(w, &http.Cookie{Name: "session", Value: "", MaxAge: -1})
		Errors(w, structs.Error{Code: http.StatusNotFound, Message: "Page not found", Page: "Home", Path: "/"})
		return
	}
	switch r.Method {
	case http.MethodGet:
		NotificationGet(w, r, user)
	case http.MethodPost:
		NotificationPost(w, r)
	default:
		Errors(w, structs.Error{Code: http.StatusMethodNotAllowed, Message: "Method not allowed", Page: "Home", Path: "/"})
		return
	}
}

func NotificationGet(w http.ResponseWriter, r *http.Request, user *structs.User) {
	tmpl, tmplErr := template.ParseFiles("Template/html/notification.html")
	if tmplErr != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Failed to load login page template", Page: "Home", Path: "/"})
		return
	}
	notifications, err := database.GetNotification(user.ID)
	if err != nil {
		fmt.Println(err)
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error geting notification", Page: "Home", Path: "/"})
		return
	}
	data := struct {
		User          *structs.User
		Notifications []structs.Notification
	}{
		User:          user,
		Notifications: notifications,
	}
	tmpl.Execute(w, data)
}

func NotificationPost(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
