package pkg

import (
	"math"

	"github.com/google/uuid"
)

func GenerateUserID() string {
	return uuid.New().String()
}

func FormatInt(num float64) int {
	return int(math.Ceil(num))
}
