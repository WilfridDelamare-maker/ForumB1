package fake

import "forum/models"

func GetAllCategories() []models.Category { 
	var allCategories []models.Category = []models.Category{
		{
			ID: 1,
			Name: "Général",
		},
		{
			ID: 2,
			Name: "Jeux vidéos",
		},
		{
			ID: 3,
			Name: "Manga/Animé",
		},
		{
			ID: 4,
			Name: "Actualités",
		},
	}
	return allCategories
}