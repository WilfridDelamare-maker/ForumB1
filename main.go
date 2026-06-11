package main

import (
	"fmt"
	"forum/database"
	"forum/handlers"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

const port = ":8080"

func main() {

	// golang ne gere pas le .env nativement. Il faut donc une bibliotheque pour y avoir accès
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Problème de chargement variables .env")
	}

	database.Init()
	fmt.Println("Database créée et fonctionnelle")

	// on crée notre propre mux (request multiplexer), c'est ce qui permet de recevoir les url et d'appeler les bons handlers.
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", handlers.Home)
	mux.HandleFunc("/", handlers.NotFoundHandler)

	mux.HandleFunc("GET /register", handlers.RegisterHandler)
	mux.HandleFunc("POST /register", handlers.PostRegisterHandler)


	mux.HandleFunc("GET /login", handlers.LoginHandler)
	mux.HandleFunc("POST /login", handlers.PostLoginHandler)

	mux.HandleFunc("POST /logout", handlers.LogoutHandler)

	mux.HandleFunc("GET /posts/{id}", handlers.PostHandler)
	mux.HandleFunc("POST /posts/{id}/comments", handlers.CreateCommentHandler)

	mux.HandleFunc("GET /posts/create", handlers.PostCreateHandler)
	mux.HandleFunc("POST /posts/create", handlers.PostCreator)

	mux.HandleFunc("POST /posts/{id}/like", handlers.PostLikeHandler)
	mux.HandleFunc("POST /posts/{id}/dislike", handlers.PostDislikeHandler)

	mux.HandleFunc("POST /comments/{id}/like", handlers.CommentLikeHandler)
	mux.HandleFunc("POST /comments/{id}/dislike", handlers.CommentDislikeHandler)

	mux.HandleFunc("GET /categories", handlers.CategoriesHandler)

	mux.HandleFunc("GET /categories/{id}", handlers.CategoryByIdHandler)

	mux.HandleFunc("GET /random", handlers.RandomPageHandler)

	mux.HandleFunc("GET /auth/github", handlers.GithubLoginHandler)
	mux.HandleFunc("GET /auth/github/callback", handlers.GitHubCallbackHandler)

	mux.HandleFunc("GET /auth/google", handlers.GoogleLoginHandler)
	mux.HandleFunc("GET /auth/google/callback", handlers.GoogleCallbackHandler)

	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Println("Serveur lancé sur http://localhost" + port)
	
	err = http.ListenAndServe(port, mux)
	if err != nil {
		log.Fatal("erreur serveur:", err)
	}
}