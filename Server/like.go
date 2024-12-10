package server

import (
	database "forum/Database"
	structs "forum/Structs"
	"net/http"
	"strconv"
)

func LikePost(w http.ResponseWriter, r *http.Request) {
	id_post, err := strconv.ParseInt(r.URL.Path[len("/like/"):], 10, 64)
	if err != nil {
		Errors(w, structs.Error{Code: http.StatusBadRequest, Message: "Invalid post ID"})
		return
	}
	user := database.GetUserConnected()
	if !database.CheckLike(user.UserID, id_post) {
		if err := database.AddLike(user.UserID, id_post); err != nil {
			Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error Adding Like"})
			return
		}
	} else if err := database.DeleteLike(user.UserID, id_post); err != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error Deleting Like"})
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
