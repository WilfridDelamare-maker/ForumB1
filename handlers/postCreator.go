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

	switch r.URL.Query().Get("error") {
	case "1":
		err = "Titre, Contenu ou Catégories manquant(s)"
	case "2":
		err = "Image trop lourde: max 20 Mo"
	case "3":
		err = "Probleme d'upload ou chargement de l'image"
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

	// refuser l'image si trop volumineuse (>20Mo)
	const maxUploadSize = 20* 1024* 1024
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	err := r.ParseMultipartForm(maxUploadSize)
	if err != nil {
		http.Redirect(w, r, "/posts/create?error=2", http.StatusSeeOther)
		return
	}

	title := strings.TrimSpace(r.FormValue("title"))
	content := strings.TrimSpace(r.FormValue("content"))
	categories := r.Form["categories"]

	if title == "" || content == "" || len(categories) ==0 {
		http.Redirect(w, r, "/posts/create?error=1", http.StatusSeeOther)
		return
	}

	// pour gérer l'upload de fichier image : avoir le path et sauvegarder l'img
	file, header, err := r.FormFile("image")
	if err != nil {
		fmt.Println("aucune image uploadée askip")
	} else {
		defer file.Close()

		if header.Size > maxUploadSize {
			http.Redirect(w, r, "/posts/create?error=2", http.StatusSeeOther)
			return
		}

		imgPath := "./static/upload/" + header.Filename

		destination, err := os.Create(imgPath)
		if err != nil {
			http.Redirect(w, r, "/posts/create?error=3", http.StatusSeeOther)
			return
		}
		defer destination.Close()

		_, err = io.Copy(destination, file)
		if err != nil {
			http.Redirect(w, r, "/posts/create?error=3", http.StatusSeeOther)
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

