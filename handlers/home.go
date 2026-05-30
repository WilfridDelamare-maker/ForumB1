package handlers

import (
	"net/http"
	"html/template"
	"forum/fake"
	"forum/models"
)

func RenderTemplate(w http.ResponseWriter, tmpl string, data any) {
	t, err := template.ParseFiles("./templates/" + tmpl )
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func Home( w http.ResponseWriter, r *http.Request) {
	posts := fake.GetAllPosts()

	IsLogged := true

	data := models.TemplateData {
		Username: "Boss",
		Posts: posts,
		IsLogged: IsLogged,
		Error: "",
	}
	RenderTemplate(w, "index.tmpl", data)
}