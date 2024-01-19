package repository

import (
	"database/sql"
	"fmt"
)

type RepositoryLikes struct {
	db *sql.DB
}

func NewRepositoryLikes(db *sql.DB) *RepositoryLikes {
	return &RepositoryLikes{db: db}
}

func (l *RepositoryLikes) AddLikesToPost(postId, userId string) error {
	_, err := l.db.Exec(`
        INSERT INTO Likes (Users_id, Content, DtCreate)
        VALUES (?, ?, CURRENT_TIMESTAMP());
    `, userId, postId)

	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil

}

func (l *RepositoryLikes) RemoveLikeFromPost(userID, postID int) error {
	_, err := l.db.Exec(`
	DELETE likes
	FROM likes
	JOIN users ON likes.user_id = users.id
	JOIN Users_posts ON likes.post_id = Users_posts.id
	WHERE users.id = ? AND Users_posts.id = ?
`, userID, postID)
	if err != nil {
		return err
	}
	return nil
}
