package repository

import (
	"database/sql"
	"fmt"
)

type RepositoryEmail struct {
	db *sql.DB
}

func NewRepositoryEmail(db *sql.DB) *RepositoryEmail {
	return &RepositoryEmail{db: db}
}

func (e *RepositoryEmail) AddEmailAndUserId(UsersId int, content string) error {
	_, err := e.db.Exec(`
        INSERT INTO EmailVerifications (Users_id, Email)
        VALUES (?, ?);
    `, UsersId, content)

	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (e *RepositoryEmail) SearchEmail(UsersId int) (string, error) {
	var email string
	row := e.db.QueryRow(`SELECT email FROM EmailVerifications WHERE Users_id = ?`, UsersId)
	if err := row.Scan(&email); err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		fmt.Println(err)
	}
	return email, nil
}

func (e *RepositoryEmail) UpdateEmailVerificationStatus(email string) error {
	_, err := e.db.Exec("UPDATE EmailVerifications SET is_email_verified = ? WHERE email = ?", true, email)

	if err != nil {
		return err
	}
	return nil
}

func (e *RepositoryEmail) IsEmailVerified(userID int) (bool, error) {
	var isVerified bool
	row := e.db.QueryRow("SELECT is_email_verified FROM EmailVerifications WHERE Users_id = ?", userID)
	if err := row.Scan(&isVerified); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		fmt.Println(err)
	}
	return isVerified, nil
}
