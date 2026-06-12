package handlers

import (
	"forum/fake"
	"math/rand"
	"net/http"
	"strconv"
)

// fonction qui redirige vers un post aléatoire
func RandomPageHandler(w http.ResponseWriter, r *http.Request) {
	postList := fake.GetAllPosts()

	if len(postList) == 0 {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	randomNbr := rand.Intn(len(postList))
	randomPost := postList[randomNbr]

	http.Redirect(w, r, "/posts/"+strconv.Itoa(randomPost.ID), http.StatusSeeOther)
}