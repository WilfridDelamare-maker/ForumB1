package handlers

import (
	"forum/fake"
	"forum/models"
	"net/http"
)

func CategoriesHandler(w http.ResponseWriter, r *http.Request) {
	_, isLogged := fake.GetCurrentUser(r)
	data := models.CategoriesLogged {
		Categories: fake.GetAllCategories(),
		IsLogged: isLogged,
	}

	RenderTemplate(w, "categories.tmpl", data)
}