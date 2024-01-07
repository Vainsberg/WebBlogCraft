package service

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Vainsberg/WebBlogCraft/internal/pkg"
	"github.com/Vainsberg/WebBlogCraft/internal/repository"
	"github.com/Vainsberg/WebBlogCraft/internal/response"
)

type Service struct {
	UsersRepository *repository.RepositoryUsers
}

func NewService(UsersRepository *repository.RepositoryUsers) *Service {
	return &Service{
		UsersRepository: UsersRepository,
	}
}

func (s *Service) HtmlContent(htmltext string) string {
	htmlContent, err := ioutil.ReadFile(htmltext)
	if err != nil {
		fmt.Println("Ошибка чтения HTML-файла", err)
		return ""
	}
	return string(htmlContent)
}

func (s *Service) FetchUserDataByIP(ip string) bool {
	fetchedIP, err := s.UsersRepository.GetIpAdress(ip)
	if err != nil {
		fmt.Println("Error:", err)
	}
	return ip == fetchedIP
}

func (s *Service) CreateUserIDCookie() (http.Cookie, string) {
	userID := pkg.GenerateUserID()
	return http.Cookie{
		Name:    "userId",
		Value:   userID,
		Expires: time.Now(),
		Path:    "/",
	}, userID
}
func (s *Service) SetNameCookie(name string) http.Cookie {
	return http.Cookie{
		Name:    "userId",
		Value:   name,
		Expires: time.Now(),
		Path:    "/",
	}
}

func (s *Service) SetUserNameAndRepository(name string, pageVariables response.Page, ip string) {
	userName := pkg.AddAndRetrieveLastUserName(name, pageVariables)
	s.SetNameCookie(userName)
	s.UsersRepository.GetSetName(ip, userName)
}
