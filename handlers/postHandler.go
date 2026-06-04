package handlers

import (
	"fmt"
	"forum/fake"
	"forum/models"
	"net/http"
	"strconv"
	"strings"
)

func PostHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	data, ok := fake.GetPostById(id)
	if !ok {
		http.NotFound(w, r)
		return
	}

	username, isLogged := fake.GetCurrentUser(r)

	comments := fake.GetCommentByPostID(id)

	datas := models.TemplateData {
		Username: username,
		IsLogged: isLogged,
		Post: data,
		Comments: comments,
	}
	
	RenderTemplate(w, "post.tmpl", datas)
}

func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	postID := r.PathValue("id")
	id, err := strconv.Atoi(postID)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	_, found := fake.GetPostById(id)
	if !found {
		http.NotFound(w, r)
		return
	} 

	content := strings.TrimSpace(r.FormValue("content"))
	if content == "" {
		http.Error(w, "Erreur contenu vide", http.StatusBadRequest)
		return
	}

	fmt.Println(content, id)

	http.Redirect(w, r, "/posts/" + postID, http.StatusSeeOther)

}