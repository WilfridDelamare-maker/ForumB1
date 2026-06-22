package handlers

import (
	"errors"
	"forum/database"
	"forum/models"
	"net/http"
	"strings"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	errMsg := ""
	switch r.URL.Query().Get("error") {
	case "1":
		errMsg = "Tous les champs sont obligatoires"
	case "2":
		errMsg = "Cette adresse email est déjà utilisée"
	case "3":
		errMsg = "Une erreur est survenue, réessayez"
	}

	data := models.TemplateData{
		DarkMode: GetDarkMode(r),
		Error: errMsg,
	}
	RenderTemplate(w, "register.tmpl", data)
}

func PostRegisterHandler(w http.ResponseWriter, r *http.Request) {
	email := strings.TrimSpace(r.FormValue("email"))
	username := strings.TrimSpace(r.FormValue("username"))
	password := r.FormValue("password")

	if email == "" || username == "" || password == "" {
		http.Redirect(w, r, "/register?error=1", http.StatusSeeOther)
		return
	}

	err := database.CreateUser(email, username, password)
	if err != nil {
		if errors.Is(err, database.ErrEmailTaken) {
			http.Redirect(w, r, "/register?error=2", http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, "/register?error=3", http.StatusSeeOther)
		return
	}

	user, err := database.GetUserByEmail(email)
	if err != nil {
		http.Redirect(w, r, "/register?error=3", http.StatusSeeOther)
		return
	}

	sessionID, err := database.CreateSession(user.ID)
	if err != nil {
		http.Redirect(w, r, "/register?error=3", http.StatusSeeOther)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		MaxAge:   86400,
		HttpOnly: true,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
