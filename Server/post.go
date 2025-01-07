package server

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	structs "forum/Data"
	database "forum/Database"
)

func Post(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		Errors(w, structs.Error{Code: http.StatusMethodNotAllowed, Message: "Method not allowed", Page: "Home", Path: "/"})
		return
	}
	id_post, err := strconv.ParseInt(r.URL.Path[len("/post/"):], 10, 64)
	if err != nil {
		Errors(w, structs.Error{Code: http.StatusBadRequest, Message: "Invalid post ID", Page: "Home", Path: "/"})
		return
	}
	cookie, err := r.Cookie("session")
	var user *structs.Session
	if err == nil {
		user = database.GetUserConnected(cookie.Value)
	} else {
		if database.DeleteSession(cookie.Value) != nil {
			Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error Ending Session", Page: "Home", Path: "/"})
			return
		}
		user = &structs.Session{Status: "Disconnected"}
	}
	post, errLoadPost := database.GetPostByID(id_post)
	if errLoadPost != nil {
		Errors(w, structs.Error{Code: http.StatusNotFound, Message: "Post not found", Page: "Home", Path: "/"})
		return
	}
	tmpl, err := template.ParseFiles("Template/html/post.html")
	if err != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Failed to load post page template", Page: "Home", Path: "/"})
		return
	}
	if r.Method == http.MethodPost {
		content := strings.TrimSpace(r.FormValue("content"))
		if content == "" {
			Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Check your input", Page: "Post", Path: fmt.Sprintf("/post/%d", id_post)})
			return
		}
		if errCrePost := database.CreateComment(content, user.UserID, id_post); errCrePost != nil {
			Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error Creating Comment", Path: fmt.Sprintf("/post/%d", id_post)})
			return
		}
		token := cookie.Value
		cookie = &http.Cookie{
			Name:     "session",
			Value:    token,
			Expires:  time.Now().Add(5 * time.Minute),
			HttpOnly: true,
			Path:     "/",
		}
		http.SetCookie(w, cookie)
	}
	comments, errLoadComment := database.GetAllComments(id_post)
	if errLoadComment != nil {
		Errors(w, structs.Error{Code: http.StatusNotFound, Message: "Comments not found", Path: fmt.Sprintf("/post/%d", id_post)})
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
