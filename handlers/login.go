package handlers

import (
	"net/http"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "login.tmpl", nil)
}

func PostLoginHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	// faux cookie crée
	if email == "wilfrid.delamare@ynov.com" && password == "toto" {
		
		http.SetCookie(w, &http.Cookie{
		Name: "session_id",
		Value: "session_nbr",
		Path: "/",
		MaxAge: 3600,
		HttpOnly: true,
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	http.Error(w, "Email ou mot de passe erroné", http.StatusUnauthorized)

}