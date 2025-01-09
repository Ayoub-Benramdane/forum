package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	structs "forum/Data"
	database "forum/Database"
)

func LikePost(w http.ResponseWriter, r *http.Request) {
	id_post, err := strconv.ParseInt(r.URL.Path[len("/like/"):], 10, 64)
	if err != nil {
		Errors(w, structs.Error{Code: http.StatusBadRequest, Message: "Invalid post ID", Page: "Home", Path: "/"})
		return
	}
	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		Errors(w, structs.Error{Code: http.StatusMethodNotAllowed, Message: "Method not allowed", Page: "Home", Path: "/"})
		return
	}
	cookie, err := r.Cookie("session")
	if err != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error Adding Like", Page: "Post", Path: "/"})
		return
	}
	user := database.GetUserConnected(cookie.Value)
	if user == nil {
		http.SetCookie(w, &http.Cookie{Name: "session", Value: "", MaxAge: -1})
		Errors(w, structs.Error{Code: http.StatusNotFound, Message: "Please log to add a Like", Page: "Home", Path: "/"})
		return
	} else if !database.CheckLike(user.UserID, id_post) {
		if err := database.AddLike(user.UserID, id_post); err != nil {
			Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error Adding Like", Page: "Home", Path: "/"})
			return
		}
	} else if err := database.DeleteLike(user.UserID, id_post); err != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error Deleting Like", Page: "Home", Path: "/"})
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
	posts, errLoadPost := database.GetAllPosts()
	if errLoadPost != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error loading posts", Page: "Home", Path: "/"})
		return
	}
	categories, errLoadPost := database.GetAllCategorys()
	if errLoadPost != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error loading categories", Page: "Home", Path: "/"})
		return
	}
	pagination, errPage := Pagination([]string{"All"}, 0)
	if errPage != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error loading pagination", Page: "Home", Path: "/"})
		return
	}
	res, err := json.Marshal(struct {
		User       *structs.Session
		Posts      []structs.Post
		Categories []structs.Category
		Pagination []int64
	}{
		User:       user,
		Posts:      posts,
		Categories: categories,
		Pagination: pagination,
	})
	if err != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error loading API", Page: "Home", Path: "/"})
		return
	}
	w.Write(res)
}

func DislikePost(w http.ResponseWriter, r *http.Request) {
	id_post, err := strconv.ParseInt(r.URL.Path[len("/dislike/"):], 10, 64)
	if err != nil {
		Errors(w, structs.Error{Code: http.StatusBadRequest, Message: "Invalid post ID", Page: "Home", Path: "/"})
		return
	}
	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		Errors(w, structs.Error{Code: http.StatusMethodNotAllowed, Message: "Method not allowed", Page: "Home", Path: "/"})
		return
	}
	cookie, err := r.Cookie("session")
	if err != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error Adding Dislike", Page: "Home", Path: "/"})
		return
	}
	user := database.GetUserConnected(cookie.Value)
	if user == nil {
		http.SetCookie(w, &http.Cookie{Name: "session", Value: "", MaxAge: -1})
		Errors(w, structs.Error{Code: http.StatusNotFound, Message: "Please log to add a Dislike", Page: "Home", Path: "/"})
		return
	} else if !database.CheckDislike(user.UserID, id_post) {
		if err := database.AddDislike(user.UserID, id_post); err != nil {
			Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error Adding Dislike", Page: "Home", Path: "/"})
			return
		}
	} else if err := database.DeleteDislike(user.UserID, id_post); err != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error Deleting Dislike", Page: "Home", Path: "/"})
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
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
