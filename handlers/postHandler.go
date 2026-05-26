package handlers

import (
	"fmt"
	"forum/fake"
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

	RenderTemplate(w, "post.tmpl", data)
}

func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	postID := r.PathValue("id")
	id, err := strconv.Atoi(postID)
	if err != nil {
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