package handlers

import (
	"fmt"
	"forum/fake"
	"net/http"
	"strconv"
)

func CommentLikeHandler(w http.ResponseWriter, r *http.Request) {
	CommentID := r.PathValue("id")
	id, err := strconv.Atoi(CommentID)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	comment, found := fake.GetCommentByID(id)
	if !found {
		http.NotFound(w,r)
		return
	}

	fmt.Println("Commentaire liké de id: ", id, " dont le contenu est:", comment.Content)

	http.Redirect(w, r, "/posts/"+strconv.Itoa(comment.PostID), http.StatusSeeOther)
}

func CommentDislikeHandler(w http.ResponseWriter, r *http.Request) {
	CommentID := r.PathValue("id")
	id, err := strconv.Atoi(CommentID)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	comment, found := fake.GetCommentByID(id)
	if !found {
		http.NotFound(w,r)
		return
	}

	fmt.Println("Commentaire disliké de id: ", id, "dont le contenu est:", comment.Content)

	http.Redirect(w, r, "/posts/"+strconv.Itoa(comment.PostID), http.StatusSeeOther)
}