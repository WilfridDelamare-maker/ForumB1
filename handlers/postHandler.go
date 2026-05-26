package handlers

import (
	"net/http"
	"strings"
	"forum/fake"
	"strconv"

)

func PostHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/posts/")
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

	if r.Method == http.MethodGet {
		RenderTemplate(w, "post.tmpl", data)
	}
}