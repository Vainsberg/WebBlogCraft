package pkg

import "github.com/google/uuid"

func GenerateUserID() string {
	return uuid.New().String()
}
