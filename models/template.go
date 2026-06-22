package models

type TemplateData struct {
	Username string
	IsLogged bool

	CurrentUserID int
	Posts []Post
	Post Post

	Comments []Comments

	Categories []Category

	Error string

	DarkMode bool
}