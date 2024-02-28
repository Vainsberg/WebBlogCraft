package pkg

import (
	"math"
	"math/rand"
	"time"

	"github.com/Vainsberg/WebBlogCraft/internal/response"
	"github.com/google/uuid"
)

func GenerateUserID() string {
	return uuid.New().String()
}

func FormatInt(num float64) int {
	return int(math.Ceil(num))
}

func CreatePageList(countInt, page int) response.PageData {
	var Pagelist response.PageData
	for i := 1; i <= countInt; i++ {
		Pagelist.TotalPages = append(Pagelist.TotalPages, i)
	}
	Pagelist.CurrentPage = page
	return Pagelist
}

func GenerateSixDigitCode() int {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	return random.Intn(900000) + 100000
}
