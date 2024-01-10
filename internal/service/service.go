package service

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Vainsberg/WebBlogCraft/internal/repository"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	Logger          *zap.Logger
	UsersRepository *repository.RepositoryUsers
}

func NewService(logger *zap.Logger, UsersRepository *repository.RepositoryUsers) *Service {
	return &Service{
		Logger:          logger,
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

func (s *Service) ParseHtml(textHtml, templateName string) *template.Template {
	tmpl, err := template.New(templateName).Parse(s.HtmlContent(textHtml))
	if err != nil {
		s.Logger.Error("Error parsing HTML content", zap.Error(err))
		return nil
	}
	return tmpl
}

func (s *Service) CreateSessionCookie(userName string) *http.Cookie {
	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(24 * time.Hour)
	s.UsersRepository.AddSessionCookie(sessionToken, userName, expiresAt)

	return &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: expiresAt,
	}
}

func (s *Service) AddUserWithHashedPassword(userName, userPassword string) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userPassword), 8)
	if err != nil {
		s.Logger.Error("Error GenerateFromPassword:", zap.Error(err))
	}
	s.UsersRepository.AddPasswordAndUserName(userName, string(hashedPassword))
}

func (s *Service) SearchPassword(userName string) string {
	serchPassword, err := s.UsersRepository.SearchPasswordAndUserName(userName)
	if err != nil {
		s.Logger.Error("Error SearchPasswordAndUserName:", zap.Error(err))
	}
	return serchPassword
}
