package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	database "forum/Database"
)

func LikePost(w http.ResponseWriter, r *http.Request) {
	idPost, err := strconv.ParseInt(r.URL.Path[len("/like/"):], 10, 64)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("session")
	if err != nil {
		http.Error(w, "Session error", http.StatusUnauthorized)
		return
	}

	user := database.GetUserConnected(cookie.Value)
	if user == nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	if !database.CheckLike(user.UserID, idPost) {
		database.AddLike(user.UserID, idPost)
	} else {
		database.DeleteLike(user.UserID, idPost)
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
	updatedLikes, _ := database.CountLikes(idPost)
	updatedDislikes, _ := database.CountDislikes(idPost)
	response := map[string]interface{}{
		"updatedLikes":    updatedLikes,
		"updatedDislikes": updatedDislikes,
		"isLiked":         database.CheckLike(user.UserID, idPost),
		"isDisliked":      database.CheckDislike(user.UserID, idPost),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func DislikePost(w http.ResponseWriter, r *http.Request) {
	idPost, err := strconv.ParseInt(r.URL.Path[len("/dislike/"):], 10, 64)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("session")
	if err != nil {
		http.Error(w, "Session error", http.StatusUnauthorized)
		return
	}

	user := database.GetUserConnected(cookie.Value)
	if user == nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	if !database.CheckDislike(user.UserID, idPost) {
		database.AddDislike(user.UserID, idPost)
	} else {
		database.DeleteDislike(user.UserID, idPost)
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
	updatedLikes, _ := database.CountLikes(idPost)
	updatedDislikes, _ := database.CountDislikes(idPost)
	response := map[string]interface{}{
		"updatedLikes":    updatedLikes,
		"updatedDislikes": updatedDislikes,
		"isLiked":         database.CheckLike(user.UserID, idPost),
		"isDisliked":      database.CheckDislike(user.UserID, idPost),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
