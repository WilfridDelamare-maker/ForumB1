package database

import (
	"database/sql"
	"log"
	"os"

	// Le _ importe le package sans l'utiliser directement :
	// son init() interne enregistre le driver SQLite auprès de database/sql
	_ "modernc.org/sqlite"
)

// DB est la variable globale partagée par tous les fichiers du package database
var DB *sql.DB

// Init ouvre la BDD et crée les tables au démarrage du serveur
func Init() {
	// Crée ./data/ si inexistant MkdirAll ne fait rien si le dossier existe déjà
	// 0755 = permissions Unix : propriétaire peut tout, les autres peuvent lire/exécuter
	if err := os.MkdirAll("./data", 0755); err != nil {
		log.Fatal("Impossible de créer le dossier data:", err)
	}

	var err error
	// sql.Open prépare la configuration mais ne crée pas encore la connexion
	DB, err = sql.Open("sqlite", "./data/forum.db")
	if err != nil {
		// log.Fatal affiche l'erreur et arrête le programme immédiatement
		log.Fatal("Impossible d'ouvrir la base de données:", err)
	}

	// DB.Ping() tente réellement de joindre la BDD c'est ici que la connexion est vérifiée
	if err = DB.Ping(); err != nil {
		log.Fatal("Impossible de joindre la base de données:", err)
	}

	createTables()
}

func createTables() {
	// Toutes les requêtes de création sont regroupées dans une slice
	// pour pouvoir les exécuter en boucle plutôt qu'une par une
	queries := []string{
		// IF NOT EXISTS : createTables() est appelée à chaque démarrage,
		// ce if empêche une erreur si les tables existent déjà
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT UNIQUE NOT NULL,
			username TEXT NOT NULL,
			password_hash TEXT NOT NULL,
			provider TEXT NOT NULL DEFAULT 'local',
			provider_id TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,

		`CREATE UNIQUE INDEX IF NOT EXISTS idx_users_provider_provider_id
		ON users(provider, provider_id)`,

		`CREATE TABLE IF NOT EXISTS sessions (
			id TEXT PRIMARY KEY,
			user_id INTEGER NOT NULL,
			expires_at DATETIME NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,

		`CREATE TABLE IF NOT EXISTS categories (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT UNIQUE NOT NULL
		)`,

		`CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			image_path TEXT DEFAULT '',
			author_id INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (author_id) REFERENCES users(id)
		)`,

		// Table de liaison many-to-many : un post peut avoir plusieurs catégories
		// et une catégorie peut contenir plusieurs posts
		`CREATE TABLE IF NOT EXISTS post_categories (
			post_id INTEGER NOT NULL,
			category_id INTEGER NOT NULL,
			PRIMARY KEY (post_id, category_id),
			FOREIGN KEY (post_id) REFERENCES posts(id),
			FOREIGN KEY (category_id) REFERENCES categories(id)
		)`,

		`CREATE TABLE IF NOT EXISTS comments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			content TEXT NOT NULL,
			author_id INTEGER NOT NULL,
			post_id INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (author_id) REFERENCES users(id),
			FOREIGN KEY (post_id) REFERENCES posts(id)
		)`,

		// post_id et comment_id sont nullable (pas de NOT NULL) :
		// un like concerne soit un post, soit un commentaire, jamais les deux
		`CREATE TABLE IF NOT EXISTS likes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			post_id INTEGER,
			comment_id INTEGER,
			value INTEGER NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (post_id) REFERENCES posts(id),
			FOREIGN KEY (comment_id) REFERENCES comments(id)
		)`,

		// INSERT OR IGNORE : insère les catégories par défaut sans planter
		// si elles existent déjà (2ème démarrage du serveur)
		`INSERT OR IGNORE INTO categories (name) VALUES
			('Général'),
			('Jeux vidéos'),
			('Manga/Animé'),
			('Musique'),
			('Sport'),
			('Technologie')`,
	}

	// On exécute chaque requête une par une
	// Le _ ignore l'index (0, 1, 2...) car seul le contenu q nous intéresse
	for _, q := range queries {
		if _, err := DB.Exec(q); err != nil {
			log.Fatal("Erreur création table:", err)
		}
	}
}
