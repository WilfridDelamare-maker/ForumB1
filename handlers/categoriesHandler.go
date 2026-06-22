package handlers

import (
	"forum/fake"
	"forum/models"
	"net/http"
	"strconv"
)

// handler pour avoir les catégories et afficher le template html 
func CategoriesHandler(w http.ResponseWriter, r *http.Request) {
	username, isLogged := fake.GetCurrentUser(r)

	data := models.TemplateData {
		Username: username,
		IsLogged: isLogged,
		Categories: fake.GetAllCategories(),
		DarkMode: GetDarkMode(r),
	}

	RenderTemplate(w, "categories.tmpl", data)
}

// handler pour avoir la catégorie de l'url {id} et afficher tous les posts de la catégorie
func CategoryByIdHandler(w http.ResponseWriter, r *http.Request) {
	IdString := r.PathValue("id")

	id, err := strconv.Atoi(IdString)
	if err != nil {
		http.NotFound(w,r)
		return
	}

	category, found := fake.GetCategoryById(id)
	if !found {
		http.NotFound(w, r)
		return
	}

	posts := fake.GetPostsByCategory(category.Name)

	username, isLogged := fake.GetCurrentUser(r)

	data := models.TemplateData{
		Username: username,
		IsLogged: isLogged,
		DarkMode: GetDarkMode(r),
		Posts: posts,
	}

	RenderTemplate(w, "postByCategory.tmpl", data)
}