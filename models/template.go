package models

type TemplateData struct {
	Username string
	Posts []Post
	isLogged bool
	Error string
}