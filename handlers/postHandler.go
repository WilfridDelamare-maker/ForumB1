package handlers

import (
	"forum/database"
	"forum/fake"
	"forum/models"
	"net/http"
	"strconv"
	"strings"
)

// envoie le template pour voir un post selon l'id de l'url. renvoie une 404 si erreur.
func PostHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		NotFoundHandler(w, r)
		return
	}

	post, ok := fake.GetPostById(id)
	if !ok {
		NotFoundHandler(w, r)
		return
	}

	username, isLogged := fake.GetCurrentUserFull(r)
	comments := fake.GetCommentByPostID(id)

	datas := models.TemplateData{
		Username: username.Username,
		CurrentUserID: username.ID,
		IsLogged: isLogged,
		Post:     post,
		Comments: comments,
	}

	RenderTemplate(w, "post.tmpl", datas)
}

// gere l'ajout de commentaire et les erreurs
func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	user, isLogged := fake.GetCurrentUserFull(r)
	if !isLogged {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	postIDStr := r.PathValue("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		NotFoundHandler(w, r)
		return
	}

	_, found := fake.GetPostById(postID)
	if !found {
		NotFoundHandler(w, r)
		return
	}

	content := strings.TrimSpace(r.FormValue("content"))
	if content == "" {
		http.Redirect(w, r, "/posts/"+postIDStr, http.StatusSeeOther)
		return
	}

	if err := database.CreateComment(content, user.ID, postID); err != nil {
		http.Error(w, "Erreur lors de l'ajout du commentaire", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/posts/"+postIDStr, http.StatusSeeOther)
}
