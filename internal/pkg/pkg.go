package pkg

import (
	"WebBlogCraft/internal/response"
	"fmt"
	"net/url"

	"github.com/google/uuid"
)

func GenerateUserID() string {
	return uuid.New().String()
}

func DecodingContentText(content []byte, pageValiable response.Page) response.Page {
	decodedStringContent, err := url.QueryUnescape(string(content))
	if err != nil {
		fmt.Println("Ошибка при раскодировании строки:", err)
		return response.Page{}
	}
	pageValiable.Posts = append(pageValiable.Posts, decodedStringContent)
	return pageValiable
}

func DecodingName(name []byte, pageValiable response.Page) response.Page {
	decodedStringName, err := url.QueryUnescape(string(name))
	if err != nil {
		fmt.Println("Ошибка при раскодировании строки:", err)
		return response.Page{}
	}
	pageValiable.UserName = append(pageValiable.UserName, decodedStringName)
	return pageValiable
}

func RemovingPreposition(text string) string {
	equalSignCount := 0
	resultText := ""
	for _, v := range text {
		if string(v) == "=" {
			equalSignCount++
			continue
		}
		if equalSignCount >= 1 {
			resultText += string(v)
		}
	}
	return resultText
}
