package models

type User struct {
	ID           int
	Email        string
	Username     string
	PasswordHash string
	CreatedAt    string
	Provider 	 string
	ProviderID   string
}
