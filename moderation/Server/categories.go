package server

import (
	"encoding/json"
	structs "forum/Data"
	database "forum/Database"
	"net/http"
	"strconv"
	"strings"
)

func Categories(w http.ResponseWriter, r *http.Request) {
	slc := strings.Split(r.URL.Path[len("/categories/"):], "/")
	var id_category int64
	var err error
	if r.Method != http.MethodPost {
		Errors(w, structs.Error{Code: http.StatusMethodNotAllowed, Message: "Method not allowed", Page: "Admin", Path: "/"})
		return
	}
	cookie, err := r.Cookie("session")
	if err != nil {
		Errors(w, structs.Error{Code: http.StatusNotFound, Message: "Page not Found", Page: "Home", Path: "/"})
		return
	}
	user, err := database.GetUserConnected(cookie.Value)
	if err != nil || user.Role != "admin" {
		if err != nil {
			http.SetCookie(w, &http.Cookie{Name: "session", Value: "", MaxAge: -1})
		}
		Errors(w, structs.Error{Code: http.StatusNotFound, Message: "you cant change role of this user", Page: "Home", Path: "/"})
		return
	}
	if len(slc) == 3 {
		id_category, err = strconv.ParseInt(slc[1], 10, 64)
		if err != nil {
			Errors(w, structs.Error{Code: http.StatusBadRequest, Message: "Invalid category ID", Page: "Admin", Path: "/"})
			return
		}
		if slc[0] == "edit" {
			UpdateCategory(w, r, id_category, slc[2])
		} else if slc[0] == "delete" {
			DeleteCategory(w, r, id_category, slc[2])
		} else {
			Errors(w, structs.Error{Code: http.StatusBadRequest, Message: "Invalid programme category", Page: "Admin", Path: "/"})
			return
		}
	} else if len(slc) == 1 {
		CreateCategory(w, r, slc[0])
		return
	} else {
		Errors(w, structs.Error{Code: http.StatusBadRequest, Message: "Invalid programme category", Page: "Admin", Path: "/"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func CreateCategory(w http.ResponseWriter, r *http.Request, category_name string) {
	category, err := database.CreateCategory(category_name)
	if err != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error create category", Page: "Admin", Path: "/"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

func UpdateCategory(w http.ResponseWriter, r *http.Request, id_category int64, category_name string) {
	if database.UpdateCategory(id_category, category_name) != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error update category", Page: "Admin", Path: "/"})
		return
	}
}

func DeleteCategory(w http.ResponseWriter, r *http.Request, id_category int64, category_name string) {
	if database.DeleteCategory(id_category, category_name) != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error delete category", Page: "Admin", Path: "/"})
		return
	}
}
