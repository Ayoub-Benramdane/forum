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

func DeletePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		Errors(w, structs.Error{Code: http.StatusMethodNotAllowed, Message: "Method not allowed"})
		return
	}
	id_post, err := strconv.ParseInt(r.URL.Path[len("/post/delete/"):], 10, 64)
	if err != nil {
		Errors(w, structs.Error{Code: http.StatusBadRequest, Message: "Invalid post ID"})
		return
	}
	if database.DeletePostId(id_post) != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error Deleting Post"})
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func EditPost(w http.ResponseWriter, r *http.Request) {
	id_post, err := strconv.ParseInt(r.URL.Path[len("/post/edit/"):], 10, 64)
	if err != nil {
		Errors(w, structs.Error{Code: http.StatusBadRequest, Message: "Invalid post ID"})
		return
	}
	switch r.Method {
	case http.MethodGet:
		EditPostGet(w, r, id_post)
	case http.MethodPost:
		EditPostPost(w, r, id_post)
	default:
		Errors(w, structs.Error{Code: http.StatusMethodNotAllowed, Message: "Method not allowed"})
		return
	}
}

func EditPostGet(w http.ResponseWriter, r *http.Request, id_post int64) {
	post, errLoadPost := database.GetPostByID(id_post)
	if errLoadPost != nil {
		Errors(w, structs.Error{Code: http.StatusNotFound, Message: "Post not found"})
		return
	}
	tmpl, err := template.ParseFiles("Template/editPostComment.html")
	if err != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Failed to load edit post page template"})
		return
	}
	categories, errLoadPost := database.GetAllCategorys()
	if errLoadPost != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error loading categories"})
		return
	}
	data := struct {
		Post       *structs.Post
		Categories []structs.Category
	}{
		Post:       post,
		Categories: categories,
	}
	tmpl.Execute(w, data)
}

func EditPostPost(w http.ResponseWriter, r *http.Request, id_post int64) {
	title := strings.TrimSpace(r.FormValue("title"))
	content := strings.TrimSpace(r.FormValue("content"))
	if title == "" || content == "" {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Check your input"})
		return
	}
	if err := r.ParseForm(); err != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error parsing form"})
		return
	}
	categories := r.Form["category"]
	if errUpdtPost := database.UpdatePost(title, content, categories, id_post); errUpdtPost != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error Updating post"})
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/post/%d", id_post), http.StatusSeeOther)
}
