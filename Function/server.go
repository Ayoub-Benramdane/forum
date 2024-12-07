package server

import (
	st "forum/Structs"
	"html/template"
	"net/http"
)

var log st.User
var err st.Error

func Login(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		err = st.Error{Code: http.StatusNotFound, Message: "Page not found"}
		Errors(w, err)
		return
	} else if r.Method != http.MethodGet {
		err = st.Error{Code: http.StatusMethodNotAllowed, Message: "Method not allowed"}
		Errors(w, err)
		return
	}
	tmpl, tmplErr := template.ParseFiles("Html/login page.html")
	if tmplErr != nil {
		err = st.Error{Code: http.StatusInternalServerError, Message: "Error loading login page"}
		Errors(w, err)
		return
	}
	tmpl.Execute(w, log)
}

func Home(w http.ResponseWriter, r *http.Request) {
	if log.Username == "" {
		log.Username = r.FormValue("username")
		log.Password = r.FormValue("password")
	}
	if r.URL.Path != "/home" {
		err = st.Error{Code: http.StatusNotFound, Message: "Page not found"}
		Errors(w, err)
		return
	} else if r.Method != http.MethodPost {
		err = st.Error{Code: http.StatusMethodNotAllowed, Message: "Method not allowed"}
		Errors(w, err)
		return
	}
	tmpl, tmplErr := template.ParseFiles("Html/home page.html")
	if tmplErr != nil {
		err = st.Error{Code: http.StatusInternalServerError, Message: "Error loading home page"}
		Errors(w, err)
		return
	}
	tmpl.Execute(w, log)
}

func Errors(w http.ResponseWriter, err st.Error) {
	w.WriteHeader(err.Code)
	tmpl, tmplErr := template.ParseFiles("Html/errors.html")
	if tmplErr != nil {
		http.Error(w, "Status Internal Server Error 500", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, err)
}
