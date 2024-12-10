package server

import (
	database "forum/Database"
	structs "forum/Structs"
	"net/http"
	"strconv"
)

func DislikePost(w http.ResponseWriter, r *http.Request) {
	id_post, err := strconv.ParseInt(r.URL.Path[len("/dislike/"):], 10, 64)
	if err != nil {
		Errors(w, structs.Error{Code: http.StatusBadRequest, Message: "Invalid post ID"})
		return
	}
	user := database.GetUserConnected()
	if !database.CheckDislike(user.UserID, id_post) {
		if err := database.AddDislike(user.UserID, id_post); err != nil {
			Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error Adding Dislike"})
			return
		}
	} else if err := database.DeleteDislike(user.UserID, id_post); err != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error Deleting Dislike"})
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
