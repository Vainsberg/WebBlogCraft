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

func (l *RepositoryLikes) AddLikesToPost(userID, postID int) error {
	_, err := l.db.Exec(`
        INSERT INTO Likes (Users_id, Posts_id, Liked_at)
        VALUES (?, ?, CURRENT_TIMESTAMP());
    `, userID, postID)

	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil

}

func (l *RepositoryLikes) RemoveLikeFromPost(userID, postID int) error {
	_, err := l.db.Exec(`
	DELETE Likes
	FROM Likes
	JOIN users ON Likes.Users_id = Users.Id
	JOIN Users_posts ON Likes.Posts_id = Users_posts.id
	WHERE Users.Id = ? AND Users_posts.Id = ?
`, userID, postID)
	if err != nil {
		return err
	}
	return nil
}

func (l *RepositoryLikes) CheckingLikes(userID, postID int) (bool, error) {
	var searchUserId, searchPostId string
	row := l.db.QueryRow("SELECT * FROM Likes WHERE Users_id = ? AND Posts_id = ?", userID, postID)
	if err := row.Scan(&searchUserId, &searchPostId); err == sql.ErrNoRows {
		return false, nil
	}
	return true, nil
}

func (l *RepositoryLikes) CountLikes(postID int) (int, error) {
	var count int

	err := l.db.QueryRow("SELECT COUNT(*) FROM Likes WHERE Posts_id = ?", postID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
