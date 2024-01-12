package service

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Vainsberg/WebBlogCraft/internal/pkg"
	"github.com/Vainsberg/WebBlogCraft/internal/redis"
	"github.com/Vainsberg/WebBlogCraft/internal/repository"
	"github.com/Vainsberg/WebBlogCraft/internal/response"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	Logger             *zap.Logger
	UsersRepository    *repository.RepositoryUsers
	SessionsRepository *repository.RepositorySessions
	PostsRepository    *repository.RepositoryPosts
	ClientRedis        *redis.RedisClient
}

func NewService(logger *zap.Logger, UsersRepository *repository.RepositoryUsers, SessionsRepository *repository.RepositorySessions, PostsRepository *repository.RepositoryPosts, ClientRedis *redis.RedisClient) *Service {
	return &Service{
		Logger:             logger,
		UsersRepository:    UsersRepository,
		SessionsRepository: SessionsRepository,
		PostsRepository:    PostsRepository,
		ClientRedis:        ClientRedis,
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
	s.SessionsRepository.AddSessionCookie(sessionToken, userName, expiresAt)

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

func (s *Service) PublishPostWithSessionUser(sessionToken, content string) {
	searchUserName, err := s.SessionsRepository.SearchUserNameSessionCookie(sessionToken)
	if err != nil {
		s.Logger.Error("SearchUserNameSessionCookie error: ", zap.Error(err))
	}
	err = s.PostsRepository.AddContentAndUserName(searchUserName, content)
	if err != nil {
		s.Logger.Error("AddContentAndUserName error: ", zap.Error(err))
	}
}

func (s *Service) AddContentToPosts(content string, Posts response.StoragePosts) response.StoragePosts {

	postID := pkg.GenerateUserID()
	post := content
	Posts.PostsID = append(Posts.PostsID, postID)
	Posts.Posts = append(Posts.Posts, content)

	if len(Posts.Posts) > 10 {
		lastPostID := s.ClientRedis.SearchLastPostID(Posts)
		s.ClientRedis.DeleteFromCache(redis.CacheClient, lastPostID)
	}

	err := s.ClientRedis.AddToCache(postID, post, 0)
	if err != nil {
		s.Logger.Error("RedisClient.Set error: ", zap.Error(err))
	}
	return Posts
}
