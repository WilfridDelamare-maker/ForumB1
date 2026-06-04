package handlers

import (
	"forum/fake"
	"forum/models"
	"net/http"
)

func CategoriesHandler(w http.ResponseWriter, r *http.Request) {
	username, isLogged := fake.GetCurrentUser(r)

	data := models.TemplateData {
		Username: username,
		IsLogged: isLogged,
		Categories: fake.GetAllCategories(),
	}

	RenderTemplate(w, "categories.tmpl", data)
}