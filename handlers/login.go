package handlers

import (
	"forum/models"
	"net/http"
	"time"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	errorMessage := ""

	if r.URL.Query().Get("error") == "1" {
		errorMessage = "Identifiant ou mot de passe incorrect"
	}

	data := models.TemplateData{
		Error: errorMessage,
	}

	RenderTemplate(w, "login.tmpl", data)
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
	time.Sleep(2* time.Second)
	http.Redirect(w, r, "/login?error=1", http.StatusSeeOther)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name: "session_id",
		Value: "",
		Path: "/",
		MaxAge: -1,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}