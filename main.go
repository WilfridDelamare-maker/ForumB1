package main

import (
	"fmt"
	"forum/database"
	"net/http"
	"strings"
	"forum/handlers"
)

const port = ":8080"

func PostHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/posts/")
	fmt.Println("Post numéro: ", id)

	if r.Method == http.MethodGet {
		handlers.RenderTemplate(w, "post.tmpl", nil)
	}
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

	mux.HandleFunc("/", handlers.Home)
	mux.HandleFunc("/register", handlers.RegisterHandler)
	mux.HandleFunc("/login", handlers.LoginHandler)
	mux.HandleFunc("/posts/", PostHandler)

	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("Serveur lancé sur (http://localhost" + port + ")")
	
	err = http.ListenAndServe(port, mux)
	if err != nil {
		fmt.Println("erreur serveur:", err)
	}
}