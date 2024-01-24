package repository

import (
	"database/sql"
	"fmt"
)

type RepositoryComments struct {
	db *sql.DB
}

func NewRepositoryComments(db *sql.DB) *RepositoryPosts {
	return &RepositoryPosts{db: db}
}

func (l *RepositoryLikes) AddCommentToPost(comment string, userID, postID int) error {
	_, err := l.db.Exec(`
        INSERT INTO Likes (Comment, Users_id, Posts_id, Comment_at)
        VALUES (?, ?, ?, CURRENT_TIMESTAMP());
    `, comment, userID, postID)

	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
