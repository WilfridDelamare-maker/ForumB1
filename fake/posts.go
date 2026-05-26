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
		{
			ID: 3,
			Title: "Sukuna prime ou Gojo",
			Author: "otakudu76",
			Content: "Qui gagne ce combat de fou ?",
			Category: "manga/animé",
		},
	}
}

func GetPostById(id int) (models.Post, bool) {
	for _, post := range GetAllPosts() {
		if id == post.ID {
			return post, true
		}
	}

	return models.Post{}, false
}