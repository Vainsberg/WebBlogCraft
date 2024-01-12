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
	_, err := s.db.Exec("INSERT INTO sessions (Session_id,UserName,Expiry) VALUES(?,?,?);", session_token, userName, time)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (s *RepositorySessions) SearchSessionCookie(session_token string) string {
	var searchSessionToken string
	row := s.db.QueryRow("SELECT Session_id FROM sessions WHERE Session_id = ?;", session_token)
	if err := row.Scan(&searchSessionToken); err != nil && err != sql.ErrNoRows {
		return ""
	}
	return searchSessionToken
}

func (s *RepositorySessions) CheckingTimeforCookie(session_token string) bool {
	var timeCookie time.Time
	row := s.db.QueryRow("SELECT Expiry FROM sessions WHERE Session_id = ?;", session_token)
	if err := row.Scan(&timeCookie); err != nil && err != sql.ErrNoRows {
		return false
	}
	return !time.Now().After(timeCookie)
}

func (r *RepositorySessions) DeleteSessionCookie(session_token string) error {
	_, err := r.db.Exec("DELETE FROM sessions WHERE Session_id = ?;", session_token)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (s *RepositorySessions) SearchUserNameSessionCookie(session_token string) (string, error) {
	var searchUserName string
	row := s.db.QueryRow("SELECT UserName FROM sessions WHERE Session_id = ?;", session_token)
	if err := row.Scan(&searchUserName); err != nil && err != sql.ErrNoRows {
		return "", err
	}
	return searchUserName, nil
}

func (s *RepositorySessions) SearchAccountInSessions(username string) string {
	var searchUserName string
	row := s.db.QueryRow("SELECT UserName FROM sessions WHERE UserName = ?;", username)
	if err := row.Scan(&searchUserName); err != nil && err != sql.ErrNoRows {
		return ""
	}
	return searchUserName
}

func (r *RepositorySessions) DeleteSessionCookieAccount(userName string) error {
	_, err := r.db.Exec("DELETE FROM sessions WHERE UserName = ?;", userName)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
