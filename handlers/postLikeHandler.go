package handlers

import (
	"fmt"
	"forum/fake"
	"net/http"
	"strconv"
)
func PostLikeHandler(w http.ResponseWriter, r *http.Request) {
	PostID := r.PathValue("id")
	id, err := strconv.Atoi(PostID)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	
	post, found := fake.GetPostById(id)
	if !found {
		http.NotFound(w, r)
		return
	}

	fmt.Println("Le post liké est de id:", id, "et de titre:", post.Title)

	http.Redirect(w, r, "/posts/"+PostID, http.StatusSeeOther)

}

func PostDislikeHandler(w http.ResponseWriter, r *http.Request) {
	PostID := r.PathValue("id")
	id, err := strconv.Atoi(PostID)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	post, found := fake.GetPostById(id)
	if !found {
		http.NotFound(w,r)
		return
	}

	fmt.Println("Le post disliké a id:", id, "et de titre:", post.Title)

	http.Redirect(w, r, "/posts/"+PostID, http.StatusSeeOther)
}