package server

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"

	structs "forum/Data"
	database "forum/Database"
)

func Post(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		Errors(w, structs.Error{Code: http.StatusMethodNotAllowed, Message: "Method not allowed"})
		return
	}
	id_post, err := strconv.ParseInt(r.URL.Path[len("/post/"):], 10, 64)
	if err != nil {
		Errors(w, structs.Error{Code: http.StatusBadRequest, Message: "Invalid post ID"})
		return
	}
	user := database.GetUserConnected()
	if user == nil {
		user = &structs.Session{Status: "Disconnected"}
	}
	post, errLoadPost := database.GetPostByID(id_post)
	if errLoadPost != nil {
		Errors(w, structs.Error{Code: http.StatusNotFound, Message: "Post not found"})
		return
	}
	tmpl, err := template.ParseFiles("Template/html/post.html")
	if err != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Failed to load post page template"})
		return
	}
	if r.Method == http.MethodPost {
		content := strings.TrimSpace(r.FormValue("content"))
		if content == "" {
			Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Check your input"})
			return
		}
		if errCrePost := database.CreateComment(content, user.UserID, id_post); errCrePost != nil {
			Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error Creating Comment"})
			return
		}
	}
	comments, errLoadComment := database.GetAllComments(id_post, user.Status)
	if errLoadComment != nil {
		Errors(w, structs.Error{Code: http.StatusNotFound, Message: "Comments not found"})
		return
	}
	data := struct {
		User     *structs.Session
		Post     *structs.Post
		Comments []structs.Comment
	}{
		User:     user,
		Post:     post,
		Comments: comments,
	}
	tmpl.Execute(w, data)
}
