package server

import (
	structs "forum/Data"
	"html/template"
	"net/http"
)

func Errors(w http.ResponseWriter, err structs.Error) {
	w.WriteHeader(err.Code)
	tmpl, tmplErr := template.ParseFiles("Template/html/errors.html")
	if tmplErr != nil {
		http.Error(w, "Status Internal Server Error 500", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, err)
}
