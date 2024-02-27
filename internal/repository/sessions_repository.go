package repository

import (
	"database/sql"
	"fmt"
	"time"
)

type RepositorySessions struct {
	db *sql.DB
}

func NewRepositorySessions(db *sql.DB) *RepositorySessions {
	return &RepositorySessions{db: db}
}

func (s *RepositorySessions) AddSessionCookie(session_token string, userName string, time time.Time) error {
	_, err := s.db.Exec(`
        INSERT INTO Sessions (Session_id, Users_id, Expiry)
        SELECT ?, Users.Id, ?
        FROM Users
        WHERE Users.UserName = ?;
    `, session_token, time, userName)

	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (s *RepositorySessions) SearchSessionCookie(session_token string) string {
	var searchSessionToken string
	row := s.db.QueryRow("SELECT Session_id FROM Sessions WHERE Session_id = ?;", session_token)
	if err := row.Scan(&searchSessionToken); err != nil && err != sql.ErrNoRows {
		return ""
	}
	return searchSessionToken
}

func (r *RepositorySessions) DeleteExpiredSessions() error {
	thresholdTime := time.Now().Add(-24 * time.Hour)

	_, err := r.db.Exec("DELETE FROM Sessions WHERE Expiry < ?;", thresholdTime)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (s *RepositorySessions) SearchUsersIdSessionCookie(session_token string) (int, error) {
	var searcUsersId int

	row := s.db.QueryRow(`SELECT Sessions.Users_Id 
	FROM Sessions
	JOIN Users ON Sessions.Users_id = Users_Id
	WHERE Sessions.Session_id = ? `, session_token)
	if err := row.Scan(&searcUsersId); err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	return searcUsersId, nil
}

func (s *RepositorySessions) SearchAccountInSessions(username string) string {
	var sessionID string

	row := s.db.QueryRow(`
	SELECT Sessions.Session_id
	FROM Sessions
	JOIN Users ON Sessions.Users_id = Users_Id
	WHERE Users.UserName = ?`, username)
	if err := row.Scan(&sessionID); err != nil {
		if err == sql.ErrNoRows {
			return ""
		}
		fmt.Println(err)
	}
	return sessionID
}

func (r *RepositorySessions) DeleteSessionCookieAccount(session_token string) error {
	_, err := r.db.Exec(`DELETE FROM Sessions
 	WHERE Session_id = ?`, session_token)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
