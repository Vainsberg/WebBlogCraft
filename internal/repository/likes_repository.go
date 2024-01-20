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

func (l *RepositoryLikes) RemoveLikeFromPost(userID, postID string) error {
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

func (l *RepositoryLikes) CheckingLikes(userID, postID string) (bool, error) {
	var searchUserId, searchPostId string
	row := l.db.QueryRow("SELECT * FROM likes WHERE Users_id = ? AND Content = ?", userID, postID)
	if err := row.Scan(&searchUserId, &searchPostId); err != sql.ErrNoRows {
		return false, nil
	}
	return true, nil
}

func (l *RepositoryLikes) CountLikes(postID string) (int, error) {
	var count int

	err := l.db.QueryRow("SELECT COUNT(*) FROM likes WHERE Post_id = ?", postID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
