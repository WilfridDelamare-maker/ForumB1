package main

import (
	"fmt"
	"forum/database"
	"net/http"
	"forum/handlers"
)

const port = ":8080"

func main() {
	db, err := database.InitDB()
	if err != nil {
		fmt.Println("Erreur database: ", err)
		return
	}

	defer db.Close()
	fmt.Println("Database créée et fonctionnelle")

	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", handlers.Home)

	mux.HandleFunc("GET /register", handlers.RegisterHandler)
	mux.HandleFunc("POST /register", handlers.PostRegisterHandler)


	mux.HandleFunc("GET /login", handlers.LoginHandler)
	mux.HandleFunc("POST /login", handlers.PostLoginHandler)
	mux.HandleFunc("/posts/", handlers.PostHandler)

	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("Serveur lancé sur (http://localhost" + port + ")")
	
	err = http.ListenAndServe(port, mux)
	if err != nil {
		fmt.Println("erreur serveur:", err)
	}
}