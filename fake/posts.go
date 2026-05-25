package fake

import "forum/models"

func GetAllPosts() []models.Post {
	return []models.Post {
		{
			ID: 1,
			Title: "Bienvenue",
			Author: "admin",
			Content: "bonjour bienvenue sur ce forum. test...",
			Category: "général",
		},
		{
			ID: 2,
			Title: "Cherche jeu cool",
			Author: "totodu86",
			Content: "Yo je cherche un jeu cool",
			Category: "jeu vidéo",
		},
	}
}





