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
	id_post, err := strconv.ParseInt(r.URL.Path[len("/delete/"):], 10, 64)
	if err != nil {
		Errors(w, structs.Error{Code: http.StatusBadRequest, Message: "Invalid post ID", Page: "Home", Path: "/"})
		return
	} else if r.Method != http.MethodGet {
		Errors(w, structs.Error{Code: http.StatusMethodNotAllowed, Message: "Method not allowed", Page: "Post", Path: fmt.Sprintf("/post/%d", id_post)})
		return
	}
	cookie, err := r.Cookie("session")
	if err != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error Deleting Post", Page: "Post", Path: fmt.Sprintf("/post/%d", id_post)})
		return
	}
	user := database.GetUserConnected(cookie.Value)
	post, errPost := database.GetPostByID(id_post)
	if errPost != nil {
		Errors(w, structs.Error{Code: http.StatusNotFound, Message: "Post Not Found", Path: fmt.Sprintf("/post/%d", id_post)})
		return
	}
	if user.UserID == post.UserID {
		if database.DeletePostId(id_post) != nil {
			Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error Deleting Post", Page: "Post", Path: fmt.Sprintf("/post/%d", id_post)})
			return
		}
	} else {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "you can't Delete Post", Page: "Post", Path: fmt.Sprintf("/post/%d", id_post)})
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func EditPost(w http.ResponseWriter, r *http.Request) {
	id_post, err := strconv.ParseInt(r.URL.Path[len("/edit/"):], 10, 64)
	if err != nil {
		Errors(w, structs.Error{Code: http.StatusBadRequest, Message: "Invalid post ID", Page: "Home", Path: "/"})
		return
	}
	cookie, err := r.Cookie("session")
	if err != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error Updating Post", Page: "Post", Path: fmt.Sprintf("/post/%d", id_post)})
		return
	}
	user := database.GetUserConnected(cookie.Value)
	post, errPost := database.GetPostByID(id_post)
	if errPost != nil {
		Errors(w, structs.Error{Code: http.StatusNotFound, Message: "Post Not Found", Path: fmt.Sprintf("/post/%d", id_post)})
		return
	}
	if user.UserID != post.UserID {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "you can't Updating Post", Page: "Post", Path: fmt.Sprintf("/post/%d", id_post)})
		return
	}
	switch r.Method {
	case http.MethodGet:
		EditPostGet(w, r, id_post)
	case http.MethodPost:
		EditPostPost(w, r, id_post)
	default:
		Errors(w, structs.Error{Code: http.StatusMethodNotAllowed, Message: "Method not allowed", Page: "Post", Path: fmt.Sprintf("/post/%d", id_post)})
		return
	}
}

func EditPostGet(w http.ResponseWriter, r *http.Request, id_post int64) {
	post, errLoadPost := database.GetPostByID(id_post)
	if errLoadPost != nil {
		Errors(w, structs.Error{Code: http.StatusNotFound, Message: "Post not found", Page: "Home", Path: "/"})
		return
	}
	tmpl, err := template.ParseFiles("Template/html/post&comment-edit.html")
	if err != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Failed to load edit post page template", Page: "Post", Path: fmt.Sprintf("/post/%d", id_post)})
		return
	}
	categories, errLoadPost := database.GetAllCategorys()
	if errLoadPost != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error loading categories", Page: "Post", Path: fmt.Sprintf("/post/%d", id_post)})
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
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Check your input", Page: "Post", Path: fmt.Sprintf("/post/edit/%d", id_post)})
		return
	}
	if err := r.ParseForm(); err != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error parsing form", Page: "Post", Path: fmt.Sprintf("/post/edit/%d", id_post)})
		return
	}
	categories := r.Form["category"]
	if errUpdtPost := database.UpdatePost(title, content, categories, id_post); errUpdtPost != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error Updating post", Page: "Post", Path: fmt.Sprintf("/post/edit/%d", id_post)})
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/post/%d", id_post), http.StatusSeeOther)
}
