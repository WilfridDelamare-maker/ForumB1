package handlers

import (
	"fmt"
	"forum/models"
	"net/http"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	err := ""

	if r.URL.Query().Get("error") == "1" {
		err = "wrong email or username or password"
	}

	data := models.TemplateData{
		Error: err,
	}
	RenderTemplate(w, "register.tmpl", data)
}

func PostRegisterHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	username := r.FormValue("username")
	password := r.FormValue("password")

	if email == "wilfrid.delamare@ynov.com" && len(password) <= 25 {
		http.SetCookie(w, &http.Cookie{
			Name: "session_id",
			Value: "session_nbr",
			Path: "/",
			MaxAge: 3600,
			HttpOnly: true,
		})

	fmt.Println(email, username, password) //faudra envoyer dans la bdd ces datas...

	http.Redirect(w, r, "/", http.StatusSeeOther) // redirige vers index avec 303
	return
	}	

	http.Redirect(w, r, "/register?error=1", http.StatusSeeOther)
}