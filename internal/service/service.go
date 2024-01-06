package service

import (
	"WebBlogCraft/internal/repository"
	"fmt"
	"io/ioutil"
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
