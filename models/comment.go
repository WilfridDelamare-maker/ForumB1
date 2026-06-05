package models

type Comments struct {
	ID int
	PostID int
	Author string
	Content string
}