package handlers

import (
	"net/http"
	"fmt"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "register.tmpl", nil)
}

func PostRegisterHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	username := r.FormValue("username")
	password := r.FormValue("password")

	http.SetCookie(w, &http.Cookie{
		Name: "session_id",
		Value: "session_nbr",
		Path: "/",
		MaxAge: 3600,
		HttpOnly: true,
	})
		
	fmt.Println(email, username, password) //faudra envoyer dans la bdd ces datas...

	http.Redirect(w, r, "/", http.StatusSeeOther) // redirige vers index avec 303
}