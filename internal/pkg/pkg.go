package pkg

import (
	"github.com/Vainsberg/WebBlogCraft/internal/response"
	"github.com/google/uuid"
)

func GenerateUserID() string {
	return uuid.New().String()
}

func AddContentToPosts(content string, pageValiable response.Page) response.Page {
	pageValiable.Posts = append(pageValiable.Posts, content)
	return pageValiable
}

func AddAndRetrieveLastUserName(name string, pageVariables response.Page) string {
	pageVariables.UserName = append(pageVariables.UserName, name)
	userName := pageVariables.UserName[len(pageVariables.UserName)-1]
	return userName
}
