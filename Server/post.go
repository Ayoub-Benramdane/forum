package server

import (
	"encoding/json"
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
		Errors(w, structs.Error{
			Code:    http.StatusMethodNotAllowed,
			Message: "Method not allowed",
			Page:    "Home",
			Path:    "/",
		})
		return
	}

	id_post, err := strconv.ParseInt(strings.TrimPrefix(r.URL.Path, "/post/"), 10, 64)
	if err != nil {
		Errors(w, structs.Error{
			Code:    http.StatusBadRequest,
			Message: "Invalid post ID",
			Page:    "Home",
			Path:    "/",
		})
		return
	}

	cookie, err := r.Cookie("session")
	var user *structs.Session
	if err == nil {
		user = database.GetUserConnected(cookie.Value)
	} else {
		user = &structs.Session{Status: "Disconnected"}
	}

	post, errLoadPost := database.GetPostByID(id_post)
	if errLoadPost != nil {
		Errors(w, structs.Error{
			Code:    http.StatusNotFound,
			Message: "Post not found",
			Page:    "Home",
			Path:    "/",
		})
		return
	}

	if r.Method == http.MethodPost {
		if user == nil || user.Status != "Connected" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		content := strings.TrimSpace(r.FormValue("content"))
		if content == "" {
			http.Error(w, "Content cannot be empty", http.StatusBadRequest)
			return
		}

		newComment := structs.Comment{
			Author:        user.Username,
			Content:       content,
			CreatedAt:     time.Now().Format("2006-01-02 15:04:05"),
			TotalLikes:    0,
			TotalDislikes: 0,
			UserID:        user.UserID,
			PostID:        id_post,
		}
		comment_id, err := database.CreateComment(content, user.UserID, id_post)
		if err != nil {
			http.Error(w, "Failed to create comment", http.StatusInternalServerError)
			return
		}
		newComment.ID = comment_id
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(newComment)
		return
	}

	if r.Method == http.MethodGet {
		comments, errLoadComment := database.GetAllComments(id_post)
		if errLoadComment != nil {
			Errors(w, structs.Error{
				Code:    http.StatusNotFound,
				Message: "Comments not found",
				Path:    fmt.Sprintf("/post/%d", id_post),
			})
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

		tmpl, err := template.ParseFiles("Template/html/post.html")
		if err != nil {
			Errors(w, structs.Error{
				Code:    http.StatusInternalServerError,
				Message: "Failed to load post page template",
				Page:    "Home",
				Path:    "/",
			})
			return
		}

		tmpl.Execute(w, data)
	}
}
