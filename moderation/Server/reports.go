package server

import (
	"fmt"
	structs "forum/Data"
	database "forum/Database"
	"net/http"
	"strconv"
	"strings"
)

func Report(w http.ResponseWriter, r *http.Request) {
	id_post, err := strconv.ParseInt(r.URL.Path[len("/report/"):], 10, 64)
	if err != nil {
		Errors(w, structs.Error{Code: http.StatusBadRequest, Message: "Invalid post ID", Page: "Home", Path: "/"})
		return
	} else if r.Method != http.MethodPost && r.Method != http.MethodGet {
		Errors(w, structs.Error{Code: http.StatusMethodNotAllowed, Message: "Method not allowed", Page: "Post", Path: fmt.Sprintf("/post/edit/%d", id_post)})
		return
	}
	cookie, err := r.Cookie("session")
	if err != nil {
		Errors(w, structs.Error{Code: http.StatusNotFound, Message: "you cant report this post", Page: "Home", Path: "/"})
		return
	}
	user, err := database.GetUserConnected(cookie.Value)
	if err != nil || user.Role != "admin" && user.Role != "moderator" {
		if err != nil {
			http.SetCookie(w, &http.Cookie{Name: "session", Value: "", MaxAge: -1})
		}
		Errors(w, structs.Error{Code: http.StatusNotFound, Message: "you cant report this post", Page: "Home", Path: "/"})
		return
	} else if _, errPost := database.GetPostByID(id_post); errPost != nil {
		Errors(w, structs.Error{Code: http.StatusNotFound, Message: "Post Not Found", Path: fmt.Sprintf("/post/%d", id_post)})
		return
	}
	switch r.Method {
	case http.MethodPost:
		ReportPost(w, r, id_post, user.ID)
	default:
		Errors(w, structs.Error{Code: http.StatusMethodNotAllowed, Message: "Method not allowed", Page: "Post", Path: fmt.Sprintf("/post/%d", id_post)})
		return
	}
}

func ReportPost(w http.ResponseWriter, r *http.Request, id_post, user_id int64) {
	description := strings.TrimSpace(r.FormValue("description"))
	if description == "" {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Check your input", Page: "Post", Path: fmt.Sprintf("/post/edit/%d", id_post)})
		return
	} else if database.InsertReport(id_post, user_id, description) != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error adding report", Page: "Post", Path: fmt.Sprintf("/post/edit/%d", id_post)})
		return
	} else if database.UpdatePost("", "", []string{""}, id_post, "reported") != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error adding report", Page: "Post", Path: fmt.Sprintf("/post/edit/%d", id_post)})
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/post/%d", id_post), http.StatusSeeOther)
}
