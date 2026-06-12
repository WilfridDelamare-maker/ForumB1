package database

import (
	"log"

	"forum/models"
)

// Requête SELECT de base partagée par les 3 fonctions de lecture
// Chaque fonction y ajoute sa propre clause WHERE / GROUP BY à la fin
const postSelectQuery = `
	SELECT p.id, p.title, p.content, p.image_path, u.id, u.username, p.created_at,
		-- GROUP_CONCAT concatène toutes les catégories du post en une chaîne "Général, Jeux vidéos"
		-- DISTINCT évite les doublons, COALESCE remplace NULL par '' si aucune catégorie
		COALESCE(REPLACE(GROUP_CONCAT(DISTINCT c.name), ',', ', '), '') AS category,
		-- Compte les likes (+1) et dislikes (-1) séparément
		COUNT(DISTINCT CASE WHEN l.value = 1 THEN l.id END) AS likes,
		COUNT(DISTINCT CASE WHEN l.value = -1 THEN l.id END) AS dislikes
	FROM posts p
	-- JOIN classique : tout post a forcément un auteur, on veut son username
	JOIN users u ON u.id = p.author_id
	-- LEFT JOIN : garde les posts même s'ils n'ont pas de catégorie
	LEFT JOIN post_categories pc ON pc.post_id = p.id
	LEFT JOIN categories c ON c.id = pc.category_id
	-- LEFT JOIN : garde les posts même s'ils n'ont aucun like
	LEFT JOIN likes l ON l.post_id = p.id
`

// Helper partagé par GetAllPosts, GetPostByID et GetPostsByCategory
// Accepte une interface{Scan} pour fonctionner avec *sql.Rows (Query) ET *sql.Row (QueryRow)
// Sans cette interface, il faudrait deux fonctions séparées pour les deux types
func scanPost(rows interface{ Scan(...any) error }) (models.Post, error) {
	var p models.Post
	// L'ordre des & doit correspondre exactement à l'ordre des colonnes dans le SELECT
	err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.ImagePath,
		&p.AuthorID, &p.Author, &p.CreatedAt, &p.Category, &p.Likes, &p.Dislikes)
	return p, err
}

// Retourne tous les posts, du plus récent au plus ancien
func GetAllPosts() []models.Post {
	query := postSelectQuery + `GROUP BY p.id ORDER BY p.created_at DESC`
	rows, err := DB.Query(query)
	if err != nil {
		log.Println("Erreur récupération posts:", err)
		return nil
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		p, err := scanPost(rows)
		if err == nil {
			posts = append(posts, p)
		} else {
			log.Println("Erreur scan post:", err)
		}
	}
	if err := rows.Err(); err != nil {
		log.Println("Erreur parcours posts:", err)
	}
	return posts
}

// Retourne un post par son ID, et false s'il n'existe pas
func GetPostByID(id int) (models.Post, bool) {
	query := postSelectQuery + `WHERE p.id = ? GROUP BY p.id`
	row := DB.QueryRow(query, id)
	p, err := scanPost(row)
	if err != nil {
		log.Println("Erreur récupération post:", err)
		return models.Post{}, false
	}
	return p, true
}

// retourne les posts qui contiennent le string en parametre
func SearchPosts(research string) ([]models.Post) {
	query := postSelectQuery + `
		WHERE p.title LIKE ? OR p.content LIKE ?
		GROUP BY p.id
		ORDER BY p.created_at DESC
	`

	pattern := "%" + research + "%"

	rows, err := DB.Query(query, pattern, pattern)
	if err != nil {
		log.Println("Erreur recherche posts:", err)
		return nil
	}
	defer rows.Close()

	var posts []models.Post

	for rows.Next() {
		p, err := scanPost(rows)
		if err == nil {
			posts = append(posts, p)
		} else {
			log.Println("Erreur scan recherche post:", err)
		}
	}

	if err := rows.Err(); err != nil {
		log.Println("Erreur parcours recherche posts:", err)
	}

	return posts
}

// Retourne tous les posts d'une catégorie donnée (par son nom)
func GetPostsByCategory(categoryName string) []models.Post {
	// Sous-requête nécessaire : on ne peut pas filtrer directement sur l'alias c
	// car il est déjà utilisé dans postSelectQuery pour le GROUP_CONCAT
	// On utilise de nouveaux alias pc2 et c2 pour éviter le conflit
	query := postSelectQuery + `
		WHERE p.id IN (
			SELECT pc2.post_id FROM post_categories pc2
			JOIN categories c2 ON c2.id = pc2.category_id
			WHERE c2.name = ?
		)
		GROUP BY p.id ORDER BY p.created_at DESC`

	rows, err := DB.Query(query, categoryName)
	if err != nil {
		log.Println("Erreur récupération posts par catégorie:", err)
		return nil
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		p, err := scanPost(rows)
		if err == nil {
			posts = append(posts, p)
		} else {
			log.Println("Erreur scan post par catégorie:", err)
		}
	}
	if err := rows.Err(); err != nil {
		log.Println("Erreur parcours posts par catégorie:", err)
	}
	return posts
}

// Insère un post et ses catégories dans une transaction
func CreatePost(title, content, imagePath string, authorID int, categoryIDs []int) error {
	// Begin démarre une transaction : toutes les opérations suivantes forment un bloc tout ou rien
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	// defer Rollback : si une erreur survient avant Commit, toutes les opérations sont annulées
	// Si Commit réussit, Rollback ne fait rien — c'est un filet de sécurité
	defer tx.Rollback()

	result, err := tx.Exec(
		`INSERT INTO posts (title, content, image_path, author_id) VALUES (?, ?, ?, ?)`,
		title, content, imagePath, authorID,
	)
	if err != nil {
		return err
	}

	// LastInsertId récupère l'ID généré automatiquement par SQLite pour ce nouveau post
	// On en a besoin pour lier les catégories à ce post dans post_categories
	postID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	// On insère une ligne dans post_categories pour chaque catégorie sélectionnée
	for _, catID := range categoryIDs {
		_, err = tx.Exec(
			`INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)`,
			postID, catID,
		)
		if err != nil {
			return err
		}
	}

	// Commit valide toutes les opérations — sans ça, rien n'est réellement écrit en BDD
	return tx.Commit()
}




// fonctions horribles à modifier si besoin (PL ?) : elle est très bien ta fonction mon piche !
func DeletePost(postID int, userID int) error {
	_, err := DB.Exec(`
		DELETE FROM posts
		WHERE id = ? AND author_id = ?
	`, postID, userID)

	return err
}

func UpdatePost(postID int, userID int, title string, content string) error {
	_, err := DB.Exec(`
		UPDATE posts
		SET title = ?, content = ?
		WHERE id = ? AND author_id = ?
	`, title, content, postID, userID)

	return err
}