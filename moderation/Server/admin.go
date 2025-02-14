package server

import (
	"fmt"
	structs "forum/Data"
	database "forum/Database"
	"html/template"
	"net/http"
)

func Admin(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil {
		Errors(w, structs.Error{Code: http.StatusNotFound, Message: "Page not found", Page: "Home", Path: "/"})
		return
	}
	user, err := database.GetUserConnected(cookie.Value)
	if err != nil || user.Role != "admin" {
		if err != nil {
			http.SetCookie(w, &http.Cookie{Name: "session", Value: "", MaxAge: -1})
		}
		Errors(w, structs.Error{Code: http.StatusNotFound, Message: "Page not found", Page: "Home", Path: "/"})
		return
	}
	if r.Method != http.MethodGet {
		Errors(w, structs.Error{Code: http.StatusMethodNotAllowed, Message: "Method not allowed", Page: "Home", Path: "/"})
		return
	}
	tmpl, tmplErr := template.ParseFiles("Template/html/admin.html")
	if tmplErr != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Failed to load activity page template", Page: "Home", Path: "/"})
		return
	}
	users, err := database.GetAllUsers()
	if err != nil {
		fmt.Println(err)
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error loading users", Page: "Home", Path: "/"})
		return
	}
	posts, errLoadPost := database.GetAllPosts()
	if errLoadPost != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error loading posts", Page: "Home", Path: "/"})
		return
	}
	categories, errLoadCategories := database.GetAllCategories()
	if errLoadCategories != nil || database.PostsCategory(&categories) != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error loading categories", Page: "Home", Path: "/"})
		return
	}
	comments, errLoadComment := database.GetAllComments(0)
	if errLoadComment != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error loading comments", Page: "Home", Path: "/"})
		return
	}
	reports, errLoadReports := database.GetAllReports()
	if errLoadReports != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error loading reports", Page: "Home", Path: "/"})
		return
	}
	activities, err := database.GetActivities()
	if err != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error loading activities", Page: "Home", Path: "/"})
		return
	}
	stats := structs.Activity{}
	stats.TotalPosts = int64(len(posts))
	stats.TotalUsers = int64(len(users))
	stats.TotalComments = int64(len(comments))
	stats.TotalReports = int64(len(reports))
	data := struct {
		User             *structs.User
		Users            []structs.User
		Posts            []structs.Post
		Categories       []structs.Category
		Stats            structs.Activity
		Reports          []structs.Reports
		RecentActivities []structs.RecentActivities
	}{
		User:             user,
		Users:            users,
		Posts:            posts,
		Categories:       categories,
		Stats:            stats,
		Reports:          reports,
		RecentActivities: activities,
	}
	tmpl.Execute(w, data)
}
