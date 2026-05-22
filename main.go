package main

import (
	"fmt"
	"net/http"
	"html/template"
)

const port = ":8080"

func renderTemplate(w http.ResponseWriter, tmpl string) {
	t, err := template.ParseFiles("./templates/" + tmpl )
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func Home( w http.ResponseWriter, r *http.Request) {
	// gérer les routes non prévues
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	renderTemplate(w, "index.tmpl")
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		renderTemplate(w, "register.tmpl")
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

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		renderTemplate(w, "login.tmpl")
		return
	}

	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")

		fmt.Println(email, password) //envoyer dans la bdd en vrai :)

		// if (bdd response = ok) {
		http.Redirect(w, r, "/", http.StatusSeeOther) // redirige vers accueil
	}

	http.Error(w, "Erreur: methode interdite", http.StatusMethodNotAllowed)
}

func main() {
	http.HandleFunc("/", Home)
	http.HandleFunc("/register", RegisterHandler)
	http.HandleFunc("/login", LoginHandler)

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("Serveur lancé sur (http://localhost" + port + ")")
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Println("erreur serveur:", err)
	}
}