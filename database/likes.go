package database

import "database/sql"

// TogglePostLike gère le like/dislike d'un post selon 3 cas :
// - pas encore de like -> INSERT
// - même valeur que ce qu'on clique -> DELETE (toggle off)
// - valeur opposée -> UPDATE (bascule like ->/<- dislike)
func TogglePostLike(userID, postID, value int) error {
	var existing int

	// On cherche si un like existe déjà pour ce user sur ce post
	// AND comment_id IS NULL : évite de confondre avec un like de commentaire
	// qui aurait le même ID numérique que ce post
	err := DB.QueryRow(
		`SELECT value FROM likes WHERE user_id = ? AND post_id = ? AND comment_id IS NULL`,
		userID, postID,
	).Scan(&existing)

	// sql.ErrNoRows = aucun like existant -> premier clic -> on insère
	if err == sql.ErrNoRows {
		_, err = DB.Exec(
			`INSERT INTO likes (user_id, post_id, value) VALUES (?, ?, ?)`,
			userID, postID, value,
		)
		return err
	}
	// Autre erreur BDD inattendue (sécurité)
	if err != nil {
		return err
	}

	// existing == value : l'utilisateur reclique sur le même bouton -> annulation
	if existing == value {
		_, err = DB.Exec(
			`DELETE FROM likes WHERE user_id = ? AND post_id = ? AND comment_id IS NULL`,
			userID, postID,
		)
		return err
	}

	// Valeur différente : l'utilisateur bascule (ex: avait liké, clique dislike)
	_, err = DB.Exec(
		`UPDATE likes SET value = ? WHERE user_id = ? AND post_id = ? AND comment_id IS NULL`,
		value, userID, postID,
	)
	return err
}

// ToggleCommentLike — même logique que TogglePostLike mais pour les commentaires
func ToggleCommentLike(userID, commentID, value int) error {
	var existing int

	// AND post_id IS NULL : évite de confondre avec un like de post
	// qui aurait le même ID numérique que ce commentaire
	err := DB.QueryRow(
		`SELECT value FROM likes WHERE user_id = ? AND comment_id = ? AND post_id IS NULL`,
		userID, commentID,
	).Scan(&existing)

	if err == sql.ErrNoRows {
		_, err = DB.Exec(
			`INSERT INTO likes (user_id, comment_id, value) VALUES (?, ?, ?)`,
			userID, commentID, value,
		)
		return err
	}
	if err != nil {
		return err
	}

	if existing == value {
		_, err = DB.Exec(
			`DELETE FROM likes WHERE user_id = ? AND comment_id = ? AND post_id IS NULL`,
			userID, commentID,
		)
		return err
	}

	_, err = DB.Exec(
		`UPDATE likes SET value = ? WHERE user_id = ? AND comment_id = ? AND post_id IS NULL`,
		value, userID, commentID,
	)
	return err
}
