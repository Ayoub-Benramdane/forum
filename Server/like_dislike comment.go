package server

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	structs "forum/Data"
	database "forum/Database"
)

func LikeComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		Errors(w, structs.Error{Code: http.StatusMethodNotAllowed, Message: "Method not allowed", Page: "Home", Path: "/"})
		return
	}
	ids := strings.Split(r.URL.Path[len("/like_comment/"):], "/")
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
		Errors(w, structs.Error{Code: http.StatusBadRequest, Message: "Invalid comment ID", Page: "Post", Path: "/post/" + ids[0]})
		return
	}
	cookie, err := r.Cookie("session")
	if err != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error Adding Like", Page: "Post", Path: "/post/" + ids[0]})
		return
	}
	user := database.GetUserConnected(cookie.Value)
	if user == nil {
		http.SetCookie(w, &http.Cookie{Name: "session", Value: "", MaxAge: -1})
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Please log to Adding Like", Page: "Post", Path: "/post/" + ids[0]})
		return
	} else if !database.CheckLikeComment(user.UserID, id_post, id_comment) {
		if err := database.AddLikeComment(user.UserID, id_post, id_comment); err != nil {
			Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error Adding Like", Page: "Post", Path: "/post/" + ids[0]})
			return
		}
	} else if err := database.DeleteLikeComment(user.UserID, id_post, id_comment); err != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error Deleting Like", Page: "Post", Path: "/post/" + ids[0]})
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
	http.Redirect(w, r, fmt.Sprintf("/post/%d", id_post), http.StatusSeeOther)
}

func DislikeComment(w http.ResponseWriter, r *http.Request) {
	ids := strings.Split(r.URL.Path[len("/dislike_comment/"):], "/")
	if len(ids) != 2 {
		Errors(w, structs.Error{Code: http.StatusBadRequest, Message: "Invalid ID", Page: "Home", Path: "/"})
		return
	}
	id_post, err := strconv.ParseInt(ids[0], 10, 64)
	if err != nil {
		Errors(w, structs.Error{Code: http.StatusBadRequest, Message: "Invalid post ID", Page: "Home", Path: "/"})
		return
	}
	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		Errors(w, structs.Error{Code: http.StatusMethodNotAllowed, Message: "Method not allowed", Page: "Home", Path: "/"})
		return
	}
	id_comment, err := strconv.ParseInt(ids[1], 10, 64)
	if err != nil {
		Errors(w, structs.Error{Code: http.StatusBadRequest, Message: "Invalid comment ID", Page: "Home", Path: "/"})
		return
	}
	cookie, err := r.Cookie("session")
	if err != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error Adding Dislike", Page: "Post", Path: "/post/" + ids[0]})
		return
	}
	user := database.GetUserConnected(cookie.Value)
	if user == nil {
		http.SetCookie(w, &http.Cookie{Name: "session", Value: "", MaxAge: -1})
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Please log to Adding Dislike", Page: "Post", Path: "/post/" + ids[0]})
		return
	} else if !database.CheckDislikeComment(user.UserID, id_post, id_comment) {
		if err := database.AddDislikeComment(user.UserID, id_post, id_comment); err != nil {
			Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error Adding Dislike", Page: "Post", Path: "/post/" + ids[0]})
			return
		}
	} else if err := database.DeleteDislikeComment(user.UserID, id_post, id_comment); err != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error Deleting Dislike", Page: "Post", Path: "/post/" + ids[0]})
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
	http.Redirect(w, r, fmt.Sprintf("/post/%d", id_post), http.StatusSeeOther)
}
