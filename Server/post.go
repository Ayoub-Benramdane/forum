package server

import (
	database "forum/Database"
	structs "forum/Structs"
	"html/template"
	"net/http"
	"strconv"
)

func Post(w http.ResponseWriter, r *http.Request) {
	id_post, err := strconv.ParseInt(r.URL.Path[len("/post/"):], 10, 64)
	if err != nil {
		Errors(w, structs.Error{Code: http.StatusBadRequest, Message: "Invalid post ID"})
		return
	}
	user := database.GetUserConnected()
	post, errLoadPost := database.GetPostByID(id_post)
	if errLoadPost != nil {
		Errors(w, structs.Error{Code: http.StatusNotFound, Message: "Post not found"})
		return
	}
	tmpl, err := template.ParseFiles("Template/post.html")
	if err != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error loading post page"})
		return
	}
	if r.Method == http.MethodPost {
		content := r.FormValue("content")
		if errCrePost := database.CreateComment(content, user.ID, id_post); errCrePost != nil {
			Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error Creating Comment"})
			return
		}
	}
	if user == nil {
		user = &structs.Session{ID: 1, Username: "", UserID: 1, Statut: "Login"}
	}
	comments, errLoadComment := database.GetAllComments(id_post, user.Statut)
	if errLoadComment != nil {
		Errors(w, structs.Error{Code: http.StatusNotFound, Message: "Comment not found"})
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
