package repository

import (
	"database/sql"
	"fmt"

	"github.com/Vainsberg/WebBlogCraft/internal/response"
)

type RepositoryComments struct {
	db *sql.DB
}

func NewRepositoryComments(db *sql.DB) *RepositoryComments {
	return &RepositoryComments{db: db}
}

func (l *RepositoryComments) AddCommentToPost(comment string, userID, postID int) error {
	_, err := l.db.Exec(`
        INSERT INTO Comments (Comment, Users_id, Posts_id, Comment_at)
        VALUES (?, ?, ?, CURRENT_TIMESTAMP());
    `, comment, userID, postID)

	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (l *RepositoryComments) GetComments(postID string) ([]response.Comment, error) {
	rows, err := l.db.Query(`
	SELECT Comments.Id, Comments.Comment, Users.UserName, COUNT(Comments_like.Id) AS LikesCount
	FROM Comments
	LEFT JOIN Users ON Users.Id = Comments.Users_id
	LEFT JOIN Comments_like ON Comments_like.Comments_id = Comments.Id
	WHERE Comments.Posts_id = ?
	GROUP BY Comments.Id, Users.UserName, Comments.Comment
	ORDER BY Comments.Comment_at DESC;`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []response.Comment
	for rows.Next() {
		var id, comment, UsersName string
		var like int

		err := rows.Scan(&id, &comment, &UsersName, &like)
		if err != nil {
			return nil, err
		}
		comments = append(comments, response.Comment{Comment: comment, UserName: UsersName, CommentId: id, Likes: like})
	}
	return comments, nil
}

func (l *RepositoryComments) CheckingLikesToComment(userID, commentID int) (bool, error) {
	var searchUserId, searchPostId string
	row := l.db.QueryRow("SELECT * FROM Comments_like WHERE Users_id = ? AND Comments_id = ?", userID, commentID)
	if err := row.Scan(&searchUserId, &searchPostId); err == sql.ErrNoRows {
		return false, nil
	}
	return true, nil
}

func (l *RepositoryComments) AddLikesToComments(userID, commentID int) error {
	_, err := l.db.Exec(`
        INSERT INTO Comments_like (Users_id, Comments_id, Comment_like_at)
        VALUES (?, ?, CURRENT_TIMESTAMP());
    `, userID, commentID)

	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (l *RepositoryComments) RemoveLikeFromComments(userID, commentID int) error {
	_, err := l.db.Exec(`
	DELETE Comments_like
	FROM Comments_like
	JOIN Users ON Comments_like.Users_id = Users.Id
	JOIN Comments ON Comments_like.Comments_id = Comments.Id
	WHERE Users.Id = ? AND Comments.Id = ?
`, userID, commentID)
	if err != nil {
		return err
	}
	return nil
}

func (l *RepositoryComments) CountLikesComments(commentsID int) (int, error) {
	var count int

	err := l.db.QueryRow("SELECT COUNT(*) FROM Comments_like WHERE Comments_id = ?", commentsID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (l *RepositoryComments) SearchCommentId(comment string) (int, error) {
	var commentId int

	err := l.db.QueryRow("SELECT Id FROM Comments WHERE Comment = ?", comment).Scan(&commentId)
	if err != nil {
		return 0, err
	}
	return commentId, nil
}
