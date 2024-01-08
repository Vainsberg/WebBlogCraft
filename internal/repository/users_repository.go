package repository

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type RepositoryUsers struct {
	db *sql.DB
}

func NewRepositoryUsers(db *sql.DB) *RepositoryUsers {
	return &RepositoryUsers{db: db}
}

func (r *RepositoryUsers) AddUsers(Userid, Userip string) error {
	_, err := r.db.Exec("INSERT INTO users_posts (UserID,UserIP,dt) VALUES(?,?,CURRENT_TIMESTAMP())", Userid, Userip)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (r *RepositoryUsers) AddPasswordAndUserName(userName, password string) error {
	_, err := r.db.Exec("INSERT INTO users (UserName,UserPassword,dt) VALUES(?,?,CURRENT_TIMESTAMP())", userName, password)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (r *RepositoryUsers) SearchPasswordAndUserName(userName string) (string, error) {
	var UserPassword string
	row := r.db.QueryRow("SELECT UserPassword FROM users WHERE UserName = ?", userName)
	if err := row.Scan(&UserPassword); err != nil && err != sql.ErrNoRows {
		return "", err
	}
	return UserPassword, nil
}

func (r *RepositoryUsers) AddSessionCookie(session_token string, userName string, time time.Time) error {
	_, err := r.db.Exec("INSERT INTO sessions (Session_id,UserName,Expiry) VALUES(?,?,?);", session_token, userName, time)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (r *RepositoryUsers) SearchSessionCookie(session_token string) bool {
	var searchSessionToken string
	row := r.db.QueryRow("SELECT Session_id FROM sessions WHERE Session_id = ?;", session_token)
	if err := row.Scan(&searchSessionToken); err != nil && err != sql.ErrNoRows {
		return false
	}
	return true
}

func (r *RepositoryUsers) CheckingTimeforCookie(session_token string) bool {
	var timeCookie time.Time
	row := r.db.QueryRow("SELECT Expiry FROM sessions WHERE Session_id = ?;", session_token)
	if err := row.Scan(&timeCookie); err != nil && err != sql.ErrNoRows {
		return false
	}
	return time.Now().After(timeCookie)
}

func (r *RepositoryUsers) DeleteSessionCookie(session_token string) error {
	_, err := r.db.Exec("DELETE FROM sessions WHERE Session_id = ?;", session_token)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
