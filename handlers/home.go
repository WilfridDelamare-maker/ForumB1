package handlers

import (
	"net/http"
	"html/template"
	"forum/fake"
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
	// interdire les methodes non necessaires
	if r.Method != http.MethodGet {
		http.Error(w, "Methode interdite", http.StatusMethodNotAllowed)
		return
	}
	data := fake.GetAllPosts()
	RenderTemplate(w, "index.tmpl", data)
}