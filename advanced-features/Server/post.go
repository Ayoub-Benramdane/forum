package server

import (
	"encoding/json"
	"fmt"
	"html"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	structs "forum/Data"
	database "forum/Database"
)

func Post(w http.ResponseWriter, r *http.Request) {
	id_post, err := strconv.ParseInt(strings.TrimPrefix(r.URL.Path, "/post/"), 10, 64)
	if err != nil {
		Errors(w, structs.Error{Code: http.StatusBadRequest, Message: "Invalid post ID", Page: "Home", Path: "/"})
		return
	}
	cookie, err := r.Cookie("session")
	var user *structs.User
	if err == nil {
		user, err = database.GetUserConnected(cookie.Value)
		if err != nil {
			http.SetCookie(w, &http.Cookie{Name: "session", Value: "", MaxAge: -1})
			user = &structs.User{Status: "Disconnected"}
		}
	} else {
		user = &structs.User{Status: "Disconnected"}
	}
	post, errLoadPost := database.GetPostByID(id_post)
	if errLoadPost != nil {
		Errors(w, structs.Error{Code: http.StatusNotFound, Message: "Post not found", Page: "Home", Path: "/"})
		return
	}
	switch r.Method {
	case http.MethodGet:
		PostGet(w, post, user)
	case http.MethodPost:
		PostComment(w, r, id_post, user, cookie)
	default:
		Errors(w, structs.Error{Code: http.StatusMethodNotAllowed, Message: "Method not allowed", Page: "Home", Path: "/"})
		return
	}
}

func PostGet(w http.ResponseWriter, post *structs.Post, user *structs.User) {
	tmpl, err := template.ParseFiles("Template/html/post.html")
	if err != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Failed to load post page template", Page: "Home", Path: "/"})
		return
	}
	comments, errLoadComment := database.GetAllComments(post.ID)
	if errLoadComment != nil {
		Errors(w, structs.Error{Code: http.StatusNotFound, Message: "Comments not found", Path: fmt.Sprintf("/post/%d", post.ID)})
		return
	}
	data := struct {
		User     *structs.User
		Post     *structs.Post
		Comments []structs.Comment
	}{
		User:     user,
		Post:     post,
		Comments: comments,
	}
	tmpl.Execute(w, data)
}

func PostComment(w http.ResponseWriter, r *http.Request, id_post int64, user *structs.User, cookie *http.Cookie) {
	if user.Status != "Connected" {
		http.SetCookie(w, &http.Cookie{Name: "session", Value: "", MaxAge: -1})
		Errors(w, structs.Error{Code: http.StatusNotFound, Message: "Please Log in to add Comment", Page: "Home", Path: fmt.Sprintf("/post/%d", id_post)})
		return
	}
	content := strings.TrimSpace(r.FormValue("content"))
	if content == "" {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Check your input", Page: "New-Post", Path: fmt.Sprintf("/post/%d", id_post)})
		return
	}
	newComment := structs.Comment{
		Author:        user.Username,
		Content:       html.EscapeString(content),
		CreatedAt:     database.TimeAgo(time.Now()),
		TotalLikes:    0,
		TotalDislikes: 0,
		UserID:        user.ID,
		PostID:        id_post,
	}
	comment_id, err := database.CreateComment(content, user.ID, id_post)
	if err != nil {
		http.Error(w, "Failed to create comment", http.StatusInternalServerError)
		return
	}
	newComment.ID = comment_id
	cookie.Expires = time.Now().Add(5 * time.Minute)
	http.SetCookie(w, cookie)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newComment)
}
