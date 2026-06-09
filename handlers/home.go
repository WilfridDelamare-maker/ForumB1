package handlers

import (
	"forum/fake"
	"forum/models"
	"html/template"
	"net/http"
	"strings"
)

func RenderTemplate(w http.ResponseWriter, tmpl string, data any) {
	t, err := template.ParseFiles("./templates/" + tmpl)
	if err != nil {
		InternalErrorHandler(w, nil)
		return
	}
	if err = t.Execute(w, data); err != nil {
		InternalErrorHandler(w, nil)
	}
}

func Home( w http.ResponseWriter, r *http.Request) {
	var posts []models.Post
	
	username, isLogged := fake.GetCurrentUser(r)

	research := strings.TrimSpace(r.URL.Query().Get("q"))

	if research != "" {
		posts = fake.SearchPosts(research)
	} else {
	posts = fake.GetAllPosts()
	}

	data := models.TemplateData {
		Username: username,
		Posts: posts,
		IsLogged: isLogged,
		Error: "",
	}
	RenderTemplate(w, "index.tmpl", data)
}