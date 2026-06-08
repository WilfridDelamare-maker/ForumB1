package database

import (
	"database/sql"
	"forum/models"
)

// Insère un nouveau commentaire en BDD
func CreateComment(content string, authorID, postID int) error {
	// DB.Exec est utilisé à la place de Query car on n'attend aucune ligne en retour
	// Le _ ignore l'objet Result (ID inséré, nb de lignes) dont on n'a pas besoin ici
	_, err := DB.Exec(
		`INSERT INTO comments (content, author_id, post_id) VALUES (?, ?, ?)`,
		content, authorID, postID,
	)
	return err
}

// Retourne tous les commentaires d'un post, avec le nom de l'auteur et les compteurs de likes
// -- fait les commentaires en sql
func GetCommentsByPostID(postID int) []models.Comments {
	rows, err := DB.Query(`
		SELECT c.id, c.post_id, c.content, u.id, u.username, c.created_at,
			COALESCE(SUM(CASE WHEN l.value = 1 THEN 1 ELSE 0 END), 0) AS likes,
			COALESCE(SUM(CASE WHEN l.value = -1 THEN 1 ELSE 0 END), 0) AS dislikes
		FROM comments c

		-- JOIN (non LEFT) : on veut toujours le username, tout commentaire a forcément un auteur
		JOIN users u ON u.id = c.author_id

		-- LEFT JOIN : on garde les commentaires même s'ils n'ont aucun like
		LEFT JOIN likes l ON l.comment_id = c.id

		WHERE c.post_id = ?

		-- GROUP BY nécessaire car on utilise SUM qui regroupe les lignes par commentaire
		GROUP BY c.id

		-- Du plus ancien au plus récent (ordre chronologique de discussion)
		ORDER BY c.created_at ASC
	`, postID)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var comments []models.Comments
	for rows.Next() {
		var c models.Comments
		// L'ordre des & doit correspondre exactement à l'ordre des colonnes dans le SELECT
		err := rows.Scan(&c.ID, &c.PostID, &c.Content, &c.AuthorID, &c.Author,
			&c.CreatedAt, &c.Likes, &c.Dislikes)
		if err == nil {
			comments = append(comments, c)
		}
	}
	return comments
}

// Retourne un commentaire par son ID, et false s'il n'existe pas
// Utilisé dans les handlers de like pour récupérer le PostID (pour la redirection)
func GetCommentByID(id int) (models.Comments, bool) {
	// QueryRow car on attend exactement une seule ligne
	row := DB.QueryRow(`
		SELECT c.id, c.post_id, c.content, u.id, u.username, c.created_at,
			COALESCE(SUM(CASE WHEN l.value = 1 THEN 1 ELSE 0 END), 0) AS likes,
			COALESCE(SUM(CASE WHEN l.value = -1 THEN 1 ELSE 0 END), 0) AS dislikes
		FROM comments c
		JOIN users u ON u.id = c.author_id
		LEFT JOIN likes l ON l.comment_id = c.id
		WHERE c.id = ?
		GROUP BY c.id
	`, id)

	var c models.Comments
	err := row.Scan(&c.ID, &c.PostID, &c.Content, &c.AuthorID, &c.Author,
		&c.CreatedAt, &c.Likes, &c.Dislikes)
	// sql.ErrNoRows = aucun commentaire trouvé avec cet ID
	if err == sql.ErrNoRows || err != nil {
		return models.Comments{}, false
	}
	return c, true
}
