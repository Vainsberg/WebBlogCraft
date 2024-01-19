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

func (l *RepositoryLikes) AddLikesToPost(post, user string) error {
	var searchUserId, searchPostId int
	rowName := l.db.QueryRow("SELECT Id FROM users WHERE UserName = ? ", user)
	if err := rowName.Scan(&searchUserId); err != nil && err != sql.ErrNoRows {
		return err
	}
	rowPost := l.db.QueryRow("SELECT Id FROM Users_posts WHERE Content = ? ", post)
	if err := rowPost.Scan(&searchPostId); err != nil && err != sql.ErrNoRows {
		return err
	}

	_, err := l.db.Exec(`
        INSERT INTO Likes (Users_id, Content, DtCreate)
        VALUES (?, ?, CURRENT_TIMESTAMP());
    `, searchUserId, searchPostId)

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
