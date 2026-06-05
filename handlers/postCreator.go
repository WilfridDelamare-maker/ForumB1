package handlers

import (
	"fmt"
	"forum/fake"
	"forum/models"
	"io"
	"net/http"
	"os"
	"strings"
)

// fonction pour rediriger sur page création de post si user = connecté
func PostCreateHandler(w http.ResponseWriter, r *http.Request) {
	username, isLogged := fake.GetCurrentUser(r)

	err := ""
	if !isLogged {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.URL.Query().Get("error") == "1" {
		err = "Titre, Contenu ou Catégories manquant(s)"
	}

	data := models.TemplateData{
		Username: username,
		IsLogged: isLogged,
		Categories: fake.GetAllCategories(),
		Error: err,
	}

	RenderTemplate(w, "postcreate.tmpl", data)
}

// fonction pour poster le post
func PostCreator(w http.ResponseWriter, r *http.Request) {
	_, isLogged := fake.GetCurrentUser(r)
	if !isLogged {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	
	title := strings.TrimSpace(r.FormValue("title"))
	content := strings.TrimSpace(r.FormValue("content"))
	categories := r.Form["categories"]

	if title == "" || content == "" || len(categories) ==0 {
		http.Redirect(w, r, "/posts/create?error=1", http.StatusSeeOther)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		fmt.Println("aucune image uploadée askip")
	} else {
		defer file.Close()

		fmt.Println("file name:", header.Filename)
		fmt.Println("file size:", header.Size)

		imgPath := "./static/upload/" + header.Filename

		destination, err := os.Create(imgPath)
		if err != nil {
			http.Error(w, "Erreur lors de la sauvegarde image", http.StatusInternalServerError)
			return
		}
		defer destination.Close()

		_, err = io.Copy(destination, file)
		if err != nil {
			http.Error(w, "Erreur de copie de l'image", http.StatusInternalServerError)
			return
		}
		fmt.Println("Image sauvegardée au chemin:", imgPath)
	}


	fmt.Println("title:", title, "content:", content, "categories:")
	for _, category := range categories {
	fmt.Println(category)
}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

