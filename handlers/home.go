package handlers

import (
	"forum/fake"
	"forum/models"
	"html/template"
	"log"
	"net/http"
	"strings"
)

// fonction utile pour créer les templates sans se répéter. gere les erreurs si echec.
func RenderTemplate(w http.ResponseWriter, tmpl string, data any) {
	t, err := template.ParseFiles(
	"./templates/layout.tmpl",
	"./partials/header.tmpl",
	"./templates/" + tmpl)

	if err != nil {
		log.Println("not render tmpl", err)
		InternalErrorHandler(w, nil)
		return
	}
	if err = t.ExecuteTemplate(w, "layout", data); err != nil {
		log.Println("wrong name layout", err)
		InternalErrorHandler(w, nil)
	}
}

/* gere la page home ("/") et uniquement elle. Gere la connexion de l'utilisateur et envoie les datas.
get le index.tmpl */
func Home( w http.ResponseWriter, r *http.Request) {
	var posts []models.Post
	
	username, isLogged := fake.GetCurrentUserFull(r)

	research := strings.TrimSpace(r.URL.Query().Get("q"))

	if research != "" {
		posts = fake.SearchPosts(research)
	} else {
	posts = fake.GetAllPosts()
	}

	data := models.TemplateData {
		Username: username.Username,
		CurrentUserID: username.ID,
		Posts: posts,
		IsLogged: isLogged,
		Error: "",
		DarkMode: GetDarkMode(r),
	}
	RenderTemplate(w, "index.tmpl", data)
}