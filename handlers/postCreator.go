package handlers

import (
	"bytes"
	"fmt"
	"forum/database"
	"forum/fake"
	"forum/models"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// gere le html création de post, affiche les erreurs si elle existe. 
func PostCreateHandler(w http.ResponseWriter, r *http.Request) {
	username, isLogged := fake.GetCurrentUser(r)
	if !isLogged {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	errMsg := ""
	switch r.URL.Query().Get("error") {
	case "1":
		errMsg = "Titre, Contenu ou Catégories manquant(s)"
	case "2":
		errMsg = "Image trop lourde : max 20 Mo"
	case "3":
		errMsg = "Format d'image invalide (JPEG, PNG ou GIF uniquement)"
	case "4":
		errMsg = "Problème lors de l'enregistrement de l'image"
	case "5":
		errMsg = "Erreur lors de la publication du post"
	}

	data := models.TemplateData{
		Username:   username,
		IsLogged:   isLogged,
		Categories: fake.GetAllCategories(),
		Error:      errMsg,
		DarkMode: GetDarkMode(r),
	}
	RenderTemplate(w, "postcreate.tmpl", data)
}

/* fonction pour créer un post si les conditions sont respectées (pas de champ vide, 
fichier pas trop lourd, fichier de bon format, erreurs de sauvegarde interne) */
func PostCreator(w http.ResponseWriter, r *http.Request) {
	user, isLogged := fake.GetCurrentUserFull(r)
	if !isLogged {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	const maxUploadSize = 20 * 1024 * 1024
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		http.Redirect(w, r, "/posts/create?error=2", http.StatusSeeOther)
		return
	}

	title := strings.TrimSpace(r.FormValue("title"))
	content := strings.TrimSpace(r.FormValue("content"))
	categoryStrs := r.Form["categories"]

	if title == "" || content == "" || len(categoryStrs) == 0 {
		http.Redirect(w, r, "/posts/create?error=1", http.StatusSeeOther)
		return
	}

	var categoryIDs []int
	for _, s := range categoryStrs {
		id, err := strconv.Atoi(s)
		if err == nil {
			categoryIDs = append(categoryIDs, id)
		}
	}

	imagePath := ""
	file, header, err := r.FormFile("image")
	if err == nil {
		defer file.Close()

		if header.Size > maxUploadSize {
			http.Redirect(w, r, "/posts/create?error=2", http.StatusSeeOther)
			return
		}

		buf := make([]byte, 512)
		n, _ := file.Read(buf)
		if !isAllowedImageType(buf[:n]) {
			http.Redirect(w, r, "/posts/create?error=3", http.StatusSeeOther)
			return
		}
		file.Seek(0, io.SeekStart)

		filename := fmt.Sprintf("%d_%s", user.ID, header.Filename)
		dest, err := os.Create("./static/upload/" + filename)
		if err != nil {
			http.Redirect(w, r, "/posts/create?error=4", http.StatusSeeOther)
			return
		}
		defer dest.Close()

		if _, err = io.Copy(dest, file); err != nil {
			http.Redirect(w, r, "/posts/create?error=4", http.StatusSeeOther)
			return
		}
		imagePath = "/static/upload/" + filename
	}

	if err := database.CreatePost(title, content, imagePath, user.ID, categoryIDs); err != nil {
		http.Redirect(w, r, "/posts/create?error=5", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// fonction pour vérifier si le type d'image est accepté 
func isAllowedImageType(buf []byte) bool {
	jpegMagic := []byte{0xFF, 0xD8, 0xFF}
	pngMagic := []byte{0x89, 0x50, 0x4E, 0x47}
	gifMagic := []byte{0x47, 0x49, 0x46}

	return bytes.HasPrefix(buf, jpegMagic) ||
		bytes.HasPrefix(buf, pngMagic) ||
		bytes.HasPrefix(buf, gifMagic)
}
