package models

type TemplateData struct {
	Username string
	IsLogged bool

	Posts []Post
	Post Post

	Comments []Comments

	Categories []Category

	Error string
}