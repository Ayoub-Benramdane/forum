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
	if r.Method != http.MethodPost {
		Errors(w, structs.Error{Code: http.StatusMethodNotAllowed, Message: "Method not allowed", Page: "Home", Path: "/"})
		return
	}
	idPost, err := strconv.ParseInt(r.URL.Path[len("/like/"):], 10, 64)
	if err != nil {
		Errors(w, structs.Error{Code: http.StatusBadRequest, Message: "Invalid post ID", Page: "Home", Path: "/"})
		return
	}
	cookie, err := r.Cookie("session")
	if err != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error Adding Like", Page: "Post", Path: "/"})
		return
	}
	user, err := database.GetUserConnected(cookie.Value)
	if err != nil {
		http.SetCookie(w, &http.Cookie{Name: "session", Value: "", MaxAge: -1})
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Please log to Adding Like", Page: "Post", Path: "/"})
		return
	} else if !database.CheckLike(user.ID, idPost) {
		if database.AddLike(user.ID, idPost) != nil {
			Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error Adding Like", Page: "Post", Path: "/"})
			return
		}
	} else if database.DeleteLike(user.ID, idPost) != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error Deleting Like", Page: "Post", Path: "/"})
		return
	}
	cookie.Expires = time.Now().Add(5 * time.Minute)
	http.SetCookie(w, cookie)
	updatedLikes, errLikesPost := database.CountLikes(idPost)
	if errLikesPost != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error counting Like", Page: "Post", Path: "/"})
		return
	}
	updatedDislikes, errDislikesPost := database.CountDislikes(idPost)
	if errDislikesPost != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error counting Dislike", Page: "Post", Path: "/"})
		return
	}
	response := map[string]interface{}{
		"updatedLikes":    updatedLikes,
		"updatedDislikes": updatedDislikes,
		"isLiked":         database.CheckLike(user.ID, idPost),
		"isDisliked":      database.CheckDislike(user.ID, idPost),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func DislikePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		Errors(w, structs.Error{Code: http.StatusMethodNotAllowed, Message: "Method not allowed", Page: "Home", Path: "/"})
		return
	}
	idPost, err := strconv.ParseInt(r.URL.Path[len("/dislike/"):], 10, 64)
	if err != nil {
		Errors(w, structs.Error{Code: http.StatusBadRequest, Message: "Invalid post ID", Page: "Home", Path: "/"})
		return
	}
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Error(w, "Session error", http.StatusUnauthorized)
		return
	}
	user, err := database.GetUserConnected(cookie.Value)
	if err != nil {
		http.SetCookie(w, &http.Cookie{Name: "session", Value: "", MaxAge: -1})
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Please log to Adding Dislike", Page: "Post", Path: "/"})
		return
	} else if !database.CheckDislike(user.ID, idPost) {
		if database.AddDislike(user.ID, idPost) != nil {
			Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error Adding Dislike", Page: "Post", Path: "/"})
			return
		}
	} else if database.DeleteDislike(user.ID, idPost) != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error Deleting Dislike", Page: "Post", Path: "/"})
		return
	}
	cookie.Expires = time.Now().Add(5 * time.Minute)
	http.SetCookie(w, cookie)
	updatedLikes, errLikesPost := database.CountLikes(idPost)
	if errLikesPost != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error counting Like", Page: "Post", Path: "/"})
		return
	}
	updatedDislikes, errDislikesPost := database.CountDislikes(idPost)
	if errDislikesPost != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error counting Dislike", Page: "Post", Path: "/"})
		return
	}
	response := map[string]interface{}{
		"updatedLikes":    updatedLikes,
		"updatedDislikes": updatedDislikes,
		"isLiked":         database.CheckLike(user.ID, idPost),
		"isDisliked":      database.CheckDislike(user.ID, idPost),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
