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
	} else if post.Status == "blocked" && user.Role != "admin" {
		Errors(w, structs.Error{Code: http.StatusNotFound, Message: "Post blocked", Page: "Home", Path: "/"})
		return
	}
	switch r.Method {
	case http.MethodGet:
		PostGet(w, post, user)
	case http.MethodPost:
		PostComment(w, r, post, user, cookie)
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

func PostComment(w http.ResponseWriter, r *http.Request, post *structs.Post, user *structs.User, cookie *http.Cookie) {
	if user.Status != "Connected" {
		http.SetCookie(w, &http.Cookie{Name: "session", Value: "", MaxAge: -1})
		Errors(w, structs.Error{Code: http.StatusNotFound, Message: "Please Log in to add Comment", Page: "Post", Path: fmt.Sprintf("/post/%d", post.ID)})
		return
	} else if user.Role == "guest" {
		Errors(w, structs.Error{Code: http.StatusNotFound, Message: "Your account is blocked", Page: "Post", Path: fmt.Sprintf("/post/%d", post.ID)})
		return
	}
	content := strings.TrimSpace(r.FormValue("content"))
	if content == "" {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Check your input", Page: "Post", Path: fmt.Sprintf("/post/%d", post.ID)})
		return
	}
	newComment := structs.Comment{
		Author:        user.Username,
		Content:       html.EscapeString(content),
		CreatedAt:     database.TimeAgo(time.Now()),
		TotalLikes:    0,
		TotalDislikes: 0,
		UserID:        user.ID,
		PostID:        post.ID,
	}
	comment_id, err := database.CreateComment(content, user.ID, post.ID)
	if err != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Failed to create comment", Page: "Post", Path: fmt.Sprintf("/post/%d", post.ID)})
		return
	}
	newComment.ID = comment_id
	if post.UserID != user.ID {
		if database.CreateNotification("comment", "post", user.ID, post.ID, post.UserID, -1, post.Title, user.Username) != nil {
			Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Failed to create notification", Page: "Post", Path: fmt.Sprintf("/post/%d", post.ID)})
			return
		}
	}
	cookie.Expires = time.Now().Add(5 * time.Minute)
	cookie.Path = "/"
	http.SetCookie(w, cookie)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newComment)
}
