package main

import (
	"fmt"
	"net/http"
	"html/template"
	"forum/fake"
	"forum/database"
)

const port = ":8080"

func renderTemplate(w http.ResponseWriter, tmpl string, data any) {
	t, err := template.ParseFiles("./templates/" + tmpl )
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, data)
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
	data := fake.GetAllPosts()
	renderTemplate(w, "index.tmpl", data)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		renderTemplate(w, "register.tmpl", nil)
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
		renderTemplate(w, "login.tmpl", nil)
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

func main() {
	db, err := database.InitDB()
	if err != nil {
		fmt.Println("Erreur database: ", err)
		return
	}

	defer db.Close()
	fmt.Println("Database créée et fonctionnelle")

	mux := http.NewServeMux()

	mux.HandleFunc("/", Home)
	mux.HandleFunc("/register", RegisterHandler)
	mux.HandleFunc("/login", LoginHandler)

	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("Serveur lancé sur (http://localhost" + port + ")")
	
	err = http.ListenAndServe(port, mux)
	if err != nil {
		fmt.Println("erreur serveur:", err)
	}
}