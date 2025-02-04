package server

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	structs "forum/Data"
	database "forum/Database"
)

func NewPost(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil {
		Errors(w, structs.Error{Code: http.StatusNotFound, Message: "Page not found", Page: "Home", Path: "/"})
		return
	}
	user := database.GetUserConnected(cookie.Value)
	if user == nil {
		http.SetCookie(w, &http.Cookie{Name: "session", Value: "", MaxAge: -1})
		Errors(w, structs.Error{Code: http.StatusNotFound, Message: "Please Log in to add post", Page: "Home", Path: "/"})
		return
	}
	switch r.Method {
	case http.MethodGet:
		NewPostGet(w, r)
	case http.MethodPost:
		NewPostPost(w, r, cookie, user)
	default:
		Errors(w, structs.Error{Code: http.StatusMethodNotAllowed, Message: "Method not allowed", Page: "Home", Path: "/"})
		return
	}
}

func NewPostGet(w http.ResponseWriter, r *http.Request) {
	tmpl, tmplErr := template.ParseFiles("Template/html/new-post.html")
	if tmplErr != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Failed to load new post page template", Page: "Home", Path: "/"})
		return
	}
	Categories, errLoadPost := database.GetAllCategorys()
	if errLoadPost != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error loading categories", Page: "New-Post", Path: "/new-post"})
		return
	}
	tmpl.Execute(w, Categories)
}

func NewPostPost(w http.ResponseWriter, r *http.Request, cookie *http.Cookie, user *structs.Session) {
	title := strings.TrimSpace(r.FormValue("title"))
	content := strings.TrimSpace(r.FormValue("content"))
	if title == "" || content == "" {
		Errors(w, structs.Error{Code: http.StatusBadRequest, Message: "Title and content cannot be empty", Page: "New-Post", Path: "/new-post"})
		return
	}

	var filePath string
	file, header, err := r.FormFile("file")
	if err != nil && err.Error() != "http: no such file" {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error uploading file", Page: "New-Post", Path: "/new-post"})
		return
	}

	if file != nil {
		defer file.Close()

		slice := strings.Split(header.Filename, ".")
		if len(slice) < 2 || !isValidImageExtension(slice[1]) {
			Errors(w, structs.Error{Code: http.StatusNotAcceptable, Message: "Only image files are allowed", Page: "New-Post", Path: "/new-post"})
			return
		}
		if header.Size > 1024*1024*20 {
			Errors(w, structs.Error{Code: http.StatusNotAcceptable, Message: "File size is too larg", Page: "New-Post", Path: "/new-post"})
		}
		uploadDir := "./Template/uploads"
		err := os.MkdirAll(uploadDir, os.ModePerm)
		if err != nil {
			Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Unable to create upload directory", Page: "New-Post", Path: "/new-post"})
			return
		}

		newFileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), header.Filename)
		filePath = fmt.Sprintf("%s/%s", uploadDir, newFileName)
		outFile, err := os.Create(filePath)
		if err != nil {
			Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Unable to save uploaded file", Page: "New-Post", Path: "/new-post"})
			return
		}
		defer outFile.Close()
		_, err = io.Copy(outFile, file)
		if err != nil {
			Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error saving file", Page: "New-Post", Path: "/new-post"})
			return
		}
	}

	categories := r.Form["category"]
	if err := database.CreatePost(title, content, categories, user.UserID, filePath); err != nil {
		Errors(w, structs.Error{Code: http.StatusInternalServerError, Message: "Error creating post", Page: "New-Post", Path: "/new-post"})
		return
	}

	cookie.Expires = time.Now().Add(5 * time.Minute)
	cookie.Path = "/"
	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func isValidImageExtension(ext string) bool {
	allowedExtensions := []string{"png", "jpeg", "gif", "bmp", "svg", "raw", "tiff", "webp", "jpg"}
	for _, allowed := range allowedExtensions {
		if ext == allowed {
			return true
		}
	}
	return false
}
