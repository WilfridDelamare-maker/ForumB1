package fake

import (
	"forum/database"
	"forum/models"
)

func GetAllPosts() []models.Post {
	return database.GetAllPosts()
}

func GetPostById(id int) (models.Post, bool) {
	return database.GetPostByID(id)
}

func GetPostsByCategory(category string) []models.Post {
	return database.GetPostsByCategory(category)
}

func SearchPosts(query string) []models.Post {
	return database.SearchPosts(query)
}
