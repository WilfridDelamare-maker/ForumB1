package handlers

import (
	"forum/database"
	"forum/fake"
	"forum/models"
	"net/http"
	"strconv"
	"strings"
)

// fonction pour diriger vers la page de modification du post
func PostEditHandler(w http.ResponseWriter, r *http.Request) {
	user, isLogged := fake.GetCurrentUserFull(r)
	if !isLogged {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	postID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		NotFoundHandler(w, r)
		return
	}

	post, found := database.GetPostByID(postID)
	if !found {
		NotFoundHandler(w,r)
		return
	}

	if post.AuthorID != user.ID {
		http.Error(w, "Interdit", http.StatusForbidden)
		return
	}

	data := models.TemplateData{
		IsLogged: true,
		Username: user.Username,
		Post: post,
		DarkMode: GetDarkMode(r),
	}

	RenderTemplate(w, "postedit.tmpl", data)
}

// fonction pour modifier le post (envoi du formulaire)
func EditPost(w http.ResponseWriter, r *http.Request) {
	user, isLogged := fake.GetCurrentUserFull(r)
	if !isLogged {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	postID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	title := strings.TrimSpace(r.FormValue("title"))
	content := strings.TrimSpace(r.FormValue("content"))

	if title == "" || content == "" {
		http.Redirect(w, r, "/posts/"+strconv.Itoa(postID)+"/edit?error=1", http.StatusSeeOther)
		return
	}

	err = database.UpdatePost(postID, user.ID, title, content)
	if err != nil {
		http.Error(w, "Erreur modification post", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/posts/"+strconv.Itoa(postID), http.StatusSeeOther)
}

// fonction pour supprimer un post (formulaire post)
func DeletePost(w http.ResponseWriter, r *http.Request) {
	user, isLogged := fake.GetCurrentUserFull(r)
	if !isLogged {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	postID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		NotFoundHandler(w, r)
		return
	}

	err = database.DeletePost(postID, user.ID)
	if err != nil {
		http.Error(w, "Erreur deletion BDD", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}