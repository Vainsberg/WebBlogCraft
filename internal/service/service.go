package service

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Vainsberg/WebBlogCraft/internal/pkg"
	"github.com/Vainsberg/WebBlogCraft/internal/repository"
	"github.com/Vainsberg/WebBlogCraft/internal/response"
	"go.uber.org/zap"
)

type Service struct {
	Logger          zap.Logger
	UsersRepository *repository.RepositoryUsers
}

func NewService(logger *zap.Logger, UsersRepository *repository.RepositoryUsers) *Service {
	return &Service{
		Logger:          *logger,
		UsersRepository: UsersRepository,
	}
}

func (s *Service) HtmlContent(htmltext string) string {
	htmlContent, err := ioutil.ReadFile(htmltext)
	if err != nil {
		s.Logger.Error("Ошибка чтения HTML-файла", zap.Error(err))
		return ""
	}
	return string(htmlContent)
}

func (s *Service) FetchUserDataByIP(ip string) bool {
	fetchedIP, err := s.UsersRepository.GetIpAdress(ip)
	if err != nil {
		s.Logger.Error("Error GetIpAdress:", zap.Error(err))
	}
	return ip == fetchedIP
}

func (s *Service) CreateUserIDCookie(ip string) *http.Cookie {
	userID := pkg.GenerateUserID()
	s.UsersRepository.AddUsers(userID, ip)
	return &http.Cookie{
		Name:    "userId",
		Value:   userID,
		Expires: time.Now(),
		Path:    "/",
	}
}

func (s *Service) SetNameCookie(name string) http.Cookie {
	return http.Cookie{
		Name:    "userId",
		Value:   name,
		Expires: time.Now(),
		Path:    "/",
	}
}

func (s *Service) SetUserNameAndRepository(name string, pageVariables response.Page, ip string) *http.Cookie {
	userName := pkg.AddAndRetrieveLastUserName(name, pageVariables)
	cookie := s.SetNameCookie(userName)
	s.UsersRepository.GetSetName(ip, userName)
	return &cookie
}
