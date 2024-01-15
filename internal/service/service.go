package service

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/Vainsberg/WebBlogCraft/internal/pkg"
	"github.com/Vainsberg/WebBlogCraft/internal/redis"
	"github.com/Vainsberg/WebBlogCraft/internal/repository"
	"github.com/Vainsberg/WebBlogCraft/internal/response"
	"github.com/go-redis/cache/v8"
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
	Cache              *cache.Cache
	PostsRedis         response.StoragePostsRedis
}

func NewService(logger *zap.Logger,
	UsersRepository *repository.RepositoryUsers,
	SessionsRepository *repository.RepositorySessions,
	PostsRepository *repository.RepositoryPosts,
	ClientRedis *redis.RedisClient,
	cache *cache.Cache,
	PostsRedis response.StoragePostsRedis) *Service {
	return &Service{
		Logger:             logger,
		UsersRepository:    UsersRepository,
		SessionsRepository: SessionsRepository,
		PostsRepository:    PostsRepository,
		ClientRedis:        ClientRedis,
		Cache:              cache,
		PostsRedis:         PostsRedis,
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

func (s *Service) PublishPostWithSessionUser(searchUserName, content string) {
	err := s.PostsRepository.AddContentAndUserName(searchUserName, content)
	if err != nil {
		s.Logger.Error("AddContentAndUserName error: ", zap.Error(err))
	}
}

func (s *Service) AddContentToPosts(content string) {
	postID := pkg.GenerateUserID()
	s.PostsRedis.PostsID = append(s.PostsRedis.PostsID, postID)
	s.PostsRedis.Content = append(s.PostsRedis.Content, content)

	if len(s.PostsRedis.Content) > 10 {
		lastPostID := s.ClientRedis.SearchLastPostID(s.PostsRedis)
		s.ClientRedis.DeleteFromCache(s.Cache, lastPostID)
	}

	err := s.ClientRedis.AddToCache(postID, content, 0)
	if err != nil {
		s.Logger.Error("RedisClient.Set error: ", zap.Error(err))
	}
}

func (s *Service) ParsePageAndCalculateOffset(pageStr string) (int, int) {
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		s.Logger.Error("Invalid page parameter", zap.Error(err))
		return 0, 0
	}
	if err != nil || page < 1 {
		page = 1
	}

	offset := (page - 1) * 10
	return page, offset
}

func (s *Service) CheckUserExistence(userName string) string {
	searchName, err := s.UsersRepository.CheckingPresenceUser(userName)
	if err != nil {
		s.Logger.Error("CheckingPresenceUser error: ", zap.Error(err))
	}
	return searchName
}

func (s *Service) SearchCountPage(page int) response.PageData {
	var PageData response.PageData
	var count float64

	sumPosts, err := s.PostsRepository.CountPosts()
	if err != nil {
		s.Logger.Error("CountPosts error: ", zap.Error(err))
	}
	count = sumPosts / 10.0
	countInt := pkg.FormatInt(count)

	PageData.TotalPages = countInt
	PageData.CurrentPage = page
	return PageData
}

func (s *Service) CreateTemplateData(page, offset int) response.TemplateData {
	countPage := s.SearchCountPage(page)
	offsetPosts := s.PostsRepository.CalculatePageOffset(offset)

	data := response.TemplateData{
		Posts:      offsetPosts,
		Pagination: countPage,
	}

	return data
}
