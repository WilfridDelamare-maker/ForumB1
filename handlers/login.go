package handlers

import (
	"fmt"
	"net/http"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		RenderTemplate(w, "login.tmpl", nil)
		return
	}

	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")

		fmt.Println(email, password) //envoyer dans la bdd en vrai :)

		// if (bdd response = ok) {
		http.Redirect(w, r, "/", http.StatusSeeOther) // redirige vers accueil
		return
	}

	http.Error(w, "Erreur: methode interdite", http.StatusMethodNotAllowed)
}