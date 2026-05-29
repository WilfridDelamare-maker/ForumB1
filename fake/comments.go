package fake

import (
	"forum/models"
)

func GetCommentByPostID(postID int) []models.Comments {
	comments := []models.Comments{
		{ID: 1, PostID: 1, Author: "Kirl", Content: "Ok bvn sur ce forum !" },
		{ ID: 2, PostID: 2, Author: "Scoobisquit", Content: "Joue à Lol tu me remerciera plus tard, ou pas..."},
		{ID: 3, PostID: 2, Author: "LaguelO", Content: "Joue à fifa surtout ultimate team !" },
	}

	var result []models.Comments

	for _, comment := range comments {
		if comment.PostID == postID {
			result = append(result, comment)
		} 
	}
	return result
}