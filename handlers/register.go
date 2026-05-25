package handlers

import (
	"net/http"
	"fmt"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		RenderTemplate(w, "register.tmpl", nil)
		return
	}
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		username := r.FormValue("username")
		password := r.FormValue("password")

		fmt.Println(email, username, password) //faudra envoyer dans la bdd ces datas...

		http.Redirect(w, r, "/", http.StatusSeeOther) // redirige vers index avec 303
		return
	}

	http.Error(w, "Methode interdite", http.StatusMethodNotAllowed)
}