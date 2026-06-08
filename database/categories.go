package database

import (
	"database/sql"
	"forum/models"
)

// Retourne toutes les catégories de la BDD, triées par ordre alphabétique
func GetAllCategories() []models.Category {
	// DB.Query exécute une requête qui retourne plusieurs lignes
	rows, err := DB.Query(`SELECT id, name FROM categories ORDER BY name`)
	if err != nil {
		return nil
	}
	// Libère la connexion à la fin de la fonction, quoi qu'il arrive
	defer rows.Close()

	var categories []models.Category

	// rows.Next() avance d'une ligne à chaque tour, s'arrête quand il n'y en a plus
	for rows.Next() {
		var c models.Category
		// Scan copie les colonnes de la ligne courante dans les champs de c
		// L'ordre correspond exactement à l'ordre du SELECT : id -> c.ID, name -> c.Name
		if err := rows.Scan(&c.ID, &c.Name); err == nil {
			categories = append(categories, c)
		}
	}
	return categories
}

// Retourne une catégorie par son ID, et false si elle n'existe pas
func GetCategoryByID(id int) (models.Category, bool) {
	// QueryRow (sans s) est utilisé quand on attend au max une seule ligne
	// Le ? est un placeholder remplacé par id de façon sécurisée (anti injection SQL)
	row := DB.QueryRow(`SELECT id, name FROM categories WHERE id = ?`, id)

	var c models.Category
	// sql.ErrNoRows = aucune catégorie trouvée avec cet ID
	if err := row.Scan(&c.ID, &c.Name); err == sql.ErrNoRows || err != nil {
		return models.Category{}, false
	}
	return c, true
}
