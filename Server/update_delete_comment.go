package server

import (
	"fmt"
	structs "forum/Data"
	database "forum/Database"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

func DeleteComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		Errors(w, structs.Error{Code: http.StatusMethodNotAllowed, Message: "Method not allowed"})
		return
	}
	ids := strings.Split(r.URL.Path[len("/post/delete_comment/"):], "/")
	if len(ids) != 2 {
		Errors(w, structs.Error{Code: http.StatusBadRequest, Message: "Invalid ID"})
		return
	}
	id_post, err := strconv.ParseInt(ids[0], 10, 64)
	if err != nil {
		Errors(w, structs.Error{Code: http.StatusBadRequest, Message: "Invalid post ID"})
		return
	}
	id_comment, err := strconv.ParseInt(ids[1], 10, 64)
	if err != nil {
		Errors(w, structs.Error{Code: http.StatusBadRequest, Message: "Invalid comment ID"})
		return
	}
	user := database.GetUserConnected()
	UserID, errCom := database.GetComment(id_comment)
	if errCom != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error Deleting Comment"})
		return
	}
	if user.UserID == UserID {
		if database.DeleteCommentId(id_post, id_comment) != nil {
			Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error Deleting Comment"})
			return
		}
	}
	http.Redirect(w, r, fmt.Sprintf("/post/%d", id_post), http.StatusSeeOther)
}

func EditComment(w http.ResponseWriter, r *http.Request) {
	ids := strings.Split(r.URL.Path[len("/post/edit_comment/"):], "/")
	if len(ids) != 2 {
		Errors(w, structs.Error{Code: http.StatusBadRequest, Message: "Invalid ID"})
		return
	}
	id_post, err := strconv.ParseInt(ids[0], 10, 64)
	if err != nil {
		Errors(w, structs.Error{Code: http.StatusBadRequest, Message: "Invalid post ID"})
		return
	}
	id_comment, err := strconv.ParseInt(ids[1], 10, 64)
	if err != nil {
		Errors(w, structs.Error{Code: http.StatusBadRequest, Message: "Invalid comment ID"})
		return
	}
	switch r.Method {
	case http.MethodGet:
		EditCommentGet(w, r, id_post, id_comment)
	case http.MethodPost:
		EditCommentPost(w, r, id_post, id_comment)
	default:
		Errors(w, structs.Error{Code: http.StatusMethodNotAllowed, Message: "Method not allowed"})
		return
	}
}

func EditCommentGet(w http.ResponseWriter, r *http.Request, id_post, id_comment int64) {
	comment, errLoadPost := database.GetCommentByID(id_post, id_comment)
	if errLoadPost != nil {
		Errors(w, structs.Error{Code: http.StatusNotFound, Message: "Comment not found"})
		return
	}
	tmpl, err := template.ParseFiles("Template/html/editPostComment.html")
	if err != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Failed to load post page template"})
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

func EditCommentPost(w http.ResponseWriter, r *http.Request, id_post, id_comment int64) {
	content := strings.TrimSpace(r.FormValue("content"))
	if content == "" {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Check your input"})
		return
	} else if errUpdtPost := database.UpdateComment(content, id_comment, id_post); errUpdtPost != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error Updating post"})
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/post/%d", id_post), http.StatusSeeOther)
}
