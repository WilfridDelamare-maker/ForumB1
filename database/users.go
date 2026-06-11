package database

import (
	"database/sql"
	"errors"
	"forum/models"
	"strconv"
	"golang.org/x/crypto/bcrypt"
)

var ErrEmailTaken = errors.New("email déjà utilisé")
var ErrUserNotFound = errors.New("utilisateur introuvable")

func CreateUser(email, username, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = DB.Exec(
		`INSERT INTO users (email, username, password_hash) VALUES (?, ?, ?)`,
		email, username, string(hash),
	)
	if err != nil && err.Error() == "UNIQUE constraint failed: users.email" {
		return ErrEmailTaken
	}
	return err
}

func GetUserByEmail(email string) (models.User, error) {
	row := DB.QueryRow(
		`SELECT id, email, username, password_hash FROM users WHERE email = ?`,
		email,
	)

	var u models.User
	err := row.Scan(&u.ID, &u.Email, &u.Username, &u.PasswordHash)
	if err == sql.ErrNoRows {
		return models.User{}, ErrUserNotFound
	}
	return u, err
}

func FindOrCreateGitHubUser(githubID int, username, email string) (models.User, error) {

	var user models.User

	providerID := strconv.Itoa(githubID)

	// Vérifie si le user GitHub existe déjà
	err := DB.QueryRow(`
		SELECT id, email, username, password_hash,
		       created_at, provider, provider_id
		FROM users
		WHERE provider = ? AND provider_id = ?
	`,
		"github",
		providerID,
	).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.Provider,
		&user.ProviderID,
	)

	// User trouvé → login direct
	if err == nil {
		return user, nil
	}

	// Vraie erreur SQL
	if err != sql.ErrNoRows {
		return models.User{}, err
	}

	// GitHub peut cacher l'email
	if email == "" {
		email = username + "@github.local"
	}

	// Création du compte OAuth
	_, err = DB.Exec(`
		INSERT INTO users (
			email,
			username,
			password_hash,
			provider,
			provider_id
		)
		VALUES (?, ?, ?, ?, ?)
	`,
		email,
		username,
		"", // pas de mot de passe local
		"github",
		providerID,
	)

	if err != nil {
		return models.User{}, err
	}

	// Récupère le user créé
	err = DB.QueryRow(`
		SELECT id, email, username, password_hash,
		       created_at, provider, provider_id
		FROM users
		WHERE provider = ? AND provider_id = ?
	`,
		"github",
		providerID,
	).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.Provider,
		&user.ProviderID,
	)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func FindOrCreateGoogleUser(googleUser *models.GoogleUser) (*models.User, error) {
	var user models.User

	err := DB.QueryRow(`
		SELECT id, email, username, password_hash, provider, provider_id, created_at
		FROM users
		WHERE provider = ? AND provider_id = ?
	`, "google", googleUser.ID).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.Provider,
		&user.ProviderID,
		&user.CreatedAt,
	)

	if err == nil {
		return &user, nil
	}

	if err != sql.ErrNoRows {
		return nil, err
	}

	username := googleUser.Name
	if username == "" {
		username = googleUser.Email
	}

	_, err = DB.Exec(`
		INSERT INTO users (email, username, password_hash, provider, provider_id)
		VALUES (?, ?, ?, ?, ?)
	`, googleUser.Email, username, "", "google", googleUser.ID)

	if err != nil {
		return nil, err
	}

	err = DB.QueryRow(`
		SELECT id, email, username, password_hash, provider, provider_id, created_at
		FROM users
		WHERE provider = ? AND provider_id = ?
	`, "google", googleUser.ID).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.Provider,
		&user.ProviderID,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}