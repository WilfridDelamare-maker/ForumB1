package main

import (
	"forum/database"
	"forum/handlers"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

const port = ":8085"

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("Fichier .env non chargé, utilisation des variables d'environnement système")
	}

	database.Init()
	log.Println("Database créée et fonctionnelle")

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

	mux.HandleFunc("GET /posts/{id}/edit", handlers.PostEditHandler)
	mux.HandleFunc("POST /posts/{id}/edit", handlers.EditPost)
	mux.HandleFunc("POST /posts/{id}/delete", handlers.DeletePost)

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
