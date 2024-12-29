package server

import (
	structs "forum/Data"
	database "forum/Database"
	"net/http"
	"strconv"
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
	user := database.GetUserConnected()
	if user == nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error Adding Like", Page: "Home", Path: "/"})
		return
	}
	if !database.CheckLike(user.UserID, id_post) {
		if err := database.AddLike(user.UserID, id_post); err != nil {
			Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error Adding Like", Page: "Home", Path: "/"})
			return
		}
	} else if err := database.DeleteLike(user.UserID, id_post); err != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error Deleting Like", Page: "Home", Path: "/"})
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
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
	user := database.GetUserConnected()
	if user == nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error Adding Dislike", Page: "Home", Path: "/"})
		return
	}
	if !database.CheckDislike(user.UserID, id_post) {
		if err := database.AddDislike(user.UserID, id_post); err != nil {
			Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error Adding Dislike", Page: "Home", Path: "/"})
			return
		}
	} else if err := database.DeleteDislike(user.UserID, id_post); err != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error Deleting Dislike", Page: "Home", Path: "/"})
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
