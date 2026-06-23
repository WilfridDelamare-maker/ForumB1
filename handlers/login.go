package handlers

import (
	"errors"
	"forum/database"
	"forum/models"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// dirige vers le html pour se login et affiche l'erreur si erreur de login incorrect
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	errMsg := ""
	if r.URL.Query().Get("error") == "1" {
		errMsg = "Identifiant ou mot de passe incorrect"
	}

	data := models.TemplateData{
		Error: errMsg,
		DarkMode: GetDarkMode(r),
	}
	RenderTemplate(w, "login.tmpl", data)
}

// envoie les données du formulaire, redirige si erreur. crée les cookies de connexion si reussite.
func PostLoginHandler(w http.ResponseWriter, r *http.Request) {
	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")

	user, err := database.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, database.ErrUserNotFound) {
			http.Redirect(w, r, "/login?error=1", http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, "/login?error=1", http.StatusSeeOther)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		http.Redirect(w, r, "/login?error=1", http.StatusSeeOther)
		return
	}

	sessionID, err := database.CreateSession(user.ID)
	if err != nil {
		http.Redirect(w, r, "/login?error=1", http.StatusSeeOther)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		MaxAge:   86400, //24heures
		HttpOnly: true,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// supprime les cookies, supprime la session de la db. 
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == nil {
		database.DeleteSession(cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "session_id",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	// revenir au dark theme par defaut
	http.SetCookie(w, &http.Cookie{
		Name: "theme",
		Value: "",
		Path: "/",
		MaxAge: -1,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}