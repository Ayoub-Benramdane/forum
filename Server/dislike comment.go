package server

import (
	"fmt"
	database "forum/Database"
	structs "forum/Structs"
	"net/http"
	"strconv"
	"strings"
)

func DislikeComment(w http.ResponseWriter, r *http.Request) {
	ids := strings.Split(r.URL.Path[len("/post/dislike_comment/"):], "/")
	id_post, err := strconv.ParseInt(ids[0], 10, 64)
	if err != nil {
		Errors(w, structs.Error{Code: http.StatusBadRequest, Message: "Invalid post ID"})
		return
	}
	id_comment, err := strconv.ParseInt(ids[1], 10, 64)
	if err != nil {
		Errors(w, structs.Error{Code: http.StatusBadRequest, Message: "Invalid comment ID"})
		return
	}
	user := database.GetUserConnected()
	if !database.CheckDislikeComment(user.UserID, id_post, id_comment) {
		if err := database.AddDislikeComment(user.UserID, id_post, id_comment); err != nil {
			Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error Adding Like"})
			return
		}
	} else if err := database.DeleteDislikeComment(user.UserID, id_post, id_comment); err != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error Deleting Like"})
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/post/%d", id_post), http.StatusSeeOther)
}
