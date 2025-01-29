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

func DeleteComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		Errors(w, structs.Error{Code: http.StatusMethodNotAllowed, Message: "Method not allowed", Page: "Home", Path: "/"})
		return
	}
	ids := strings.Split(r.URL.Path[len("/delete_comment/"):], "/")
	if len(ids) != 2 {
		Errors(w, structs.Error{Code: http.StatusBadRequest, Message: "Invalid ID", Page: "Home", Path: "/"})
		return
	}
	id_post, err := strconv.ParseInt(ids[0], 10, 64)
	if err != nil {
		Errors(w, structs.Error{Code: http.StatusBadRequest, Message: "Invalid post ID", Page: "Home", Path: "/"})
		return
	}
	id_comment, err := strconv.ParseInt(ids[1], 10, 64)
	if err != nil {
		Errors(w, structs.Error{Code: http.StatusBadRequest, Message: "Invalid comment ID", Page: "Post", Path: fmt.Sprintf("/post/%d", id_post)})
		return
	}
	cookie, err := r.Cookie("session")
	if err != nil {
		Errors(w, structs.Error{Code: http.StatusNotFound, Message: "Page not found", Page: "Home", Path: "/"})
		return
	}
	post, errPost := database.GetPostByID(id_post)
	if errPost != nil {
		Errors(w, structs.Error{Code: http.StatusNotFound, Message: "Post Not Found", Path: fmt.Sprintf("/post/%d", id_post)})
		return
	}
	UserID, errCom := database.GetComment(id_comment)
	if errCom != nil {
		Errors(w, structs.Error{Code: http.StatusNotFound, Message: "Comment Not Found", Path: fmt.Sprintf("/post/%d", id_post)})
		return
	}
	user, err := database.GetUserConnected(cookie.Value)
	if err != nil {
		http.SetCookie(w, &http.Cookie{Name: "session", Value: "", MaxAge: -1})
		Errors(w, structs.Error{Code: http.StatusNotFound, Message: "Page not found", Page: "Post", Path: fmt.Sprintf("/post/%d", id_post)})
		return
	} else if user.ID != UserID && user.ID != post.UserID {
		Errors(w, structs.Error{Code: http.StatusUnauthorized, Message: "you can't Delete Comment", Page: "Post", Path: fmt.Sprintf("/post/%d", id_post)})
		return
	}
	if database.DeleteCommentId(id_post, id_comment) != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error Deleting Comment", Page: "Post", Path: fmt.Sprintf("/post/%d", id_post)})
		return
	}
	if database.DeleteNotification("comment", "post", id_post, -1, user.Username) != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Failed to Delete Notification", Page: "Post", Path: "/post/" + ids[0]})
		return
	}
	w.WriteHeader(http.StatusOK)
	cookie.Expires = time.Now().Add(5 * time.Minute)
	cookie.Path = "/"
	http.SetCookie(w, cookie)
	http.Redirect(w, r, fmt.Sprintf("/post/%d", id_post), http.StatusSeeOther)
}

func EditComment(w http.ResponseWriter, r *http.Request) {
	ids := strings.Split(r.URL.Path[len("/edit_comment/"):], "/")
	if len(ids) != 2 {
		Errors(w, structs.Error{Code: http.StatusBadRequest, Message: "Invalid ID", Page: "Home", Path: "/"})
		return
	}
	id_post, err := strconv.ParseInt(ids[0], 10, 64)
	if err != nil {
		Errors(w, structs.Error{Code: http.StatusBadRequest, Message: "Invalid post ID", Page: "Home", Path: "/"})
		return
	}
	id_comment, err := strconv.ParseInt(ids[1], 10, 64)
	if err != nil {
		Errors(w, structs.Error{Code: http.StatusBadRequest, Message: "Invalid comment ID", Page: "Post", Path: fmt.Sprintf("/post/%d", id_post)})
		return
	}
	UserID, errCom := database.GetComment(id_comment)
	if errCom != nil {
		Errors(w, structs.Error{Code: http.StatusNotFound, Message: "Comment Not Found", Path: fmt.Sprintf("/post/%d", id_post)})
		return
	}
	cookie, err := r.Cookie("session")
	if err != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error Updating Comment", Page: "Post", Path: fmt.Sprintf("/post/%d", id_post)})
		return
	}
	user, err := database.GetUserConnected(cookie.Value)
	if err != nil {
		http.SetCookie(w, &http.Cookie{Name: "session", Value: "", MaxAge: -1})
		Errors(w, structs.Error{Code: http.StatusNotFound, Message: "Page not found", Page: "Post", Path: fmt.Sprintf("/post/%d", id_post)})
		return
	} else if user.ID != UserID {
		Errors(w, structs.Error{Code: http.StatusUnauthorized, Message: "you can't Updating Comment", Page: "Post", Path: fmt.Sprintf("/post/%d", id_post)})
		return
	}
	switch r.Method {
	case http.MethodGet:
		EditCommentGet(w, r, id_post, id_comment)
	case http.MethodPost:
		EditCommentPost(w, r, id_post, id_comment, cookie)
	default:
		Errors(w, structs.Error{Code: http.StatusMethodNotAllowed, Message: "Method not allowed", Page: "Post", Path: fmt.Sprintf("/post/%d", id_post)})
		return
	}
}

func EditCommentGet(w http.ResponseWriter, r *http.Request, id_post, id_comment int64) {
	comment, errLoadPost := database.GetCommentByID(id_post, id_comment)
	if errLoadPost != nil {
		Errors(w, structs.Error{Code: http.StatusNotFound, Message: "Comment not found", Page: "Post", Path: fmt.Sprintf("/post/%d", id_post)})
		return
	}
	tmpl, err := template.ParseFiles("Template/html/post&comment_edit.html")
	if err != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Failed to load post page template", Page: "Post", Path: fmt.Sprintf("/post/%d", id_post)})
		return
	}
	data := struct {
		Post    *structs.Post
		Comment *structs.Comment
	}{
		Post:    nil,
		Comment: comment,
	}
	tmpl.Execute(w, data)
}

func EditCommentPost(w http.ResponseWriter, r *http.Request, id_post, id_comment int64, cookie *http.Cookie) {
	content := strings.TrimSpace(r.FormValue("content"))
	if content == "" {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Check your input", Page: "Post", Path: fmt.Sprintf("/post/%d", id_post)})
		return
	} else if errUpdtPost := database.UpdateComment(content, id_comment, id_post); errUpdtPost != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error Updating post", Page: "Post", Path: fmt.Sprintf("/post/%d", id_post)})
		return
	}
	cookie.Expires = time.Now().Add(5 * time.Minute)
	cookie.Path = "/"
	http.SetCookie(w, cookie)
	http.Redirect(w, r, fmt.Sprintf("/post/%d", id_post), http.StatusSeeOther)
}
