package service

import (
	"html/template"
	"io/ioutil"
	"strconv"

	"github.com/Vainsberg/WebBlogCraft/internal/pkg"
	"github.com/Vainsberg/WebBlogCraft/internal/redis"
	"github.com/Vainsberg/WebBlogCraft/internal/repository"
	"github.com/Vainsberg/WebBlogCraft/internal/response"
	"github.com/go-redis/cache/v8"
	"go.uber.org/zap"
)

type PostService struct {
	Logger             *zap.Logger
	UsersRepository    *repository.RepositoryUsers
	SessionsRepository *repository.RepositorySessions
	PostsRepository    *repository.RepositoryPosts
	ClientRedis        *redis.RedisClient
	Cache              *cache.Cache
	PostsRedis         []response.PostsRedis
}

func NewPostService(logger *zap.Logger,
	UsersRepository *repository.RepositoryUsers,
	SessionsRepository *repository.RepositorySessions,
	PostsRepository *repository.RepositoryPosts,
	ClientRedis *redis.RedisClient,
	cache *cache.Cache,
	PostsRedis []response.PostsRedis) *PostService {
	return &PostService{
		Logger:             logger,
		UsersRepository:    UsersRepository,
		SessionsRepository: SessionsRepository,
		PostsRepository:    PostsRepository,
		ClientRedis:        ClientRedis,
		Cache:              cache,
		PostsRedis:         PostsRedis,
	}
}

func (post *PostService) HtmlContent(htmltext string) string {
	htmlContent, err := ioutil.ReadFile(htmltext)
	if err != nil {
		post.Logger.Error("Ошибка чтения HTML-файла", zap.Error(err))
		return ""
	}
	return string(htmlContent)
}

func (post *PostService) ParseHtml(textHtml, templateName string) *template.Template {
	tmpl, err := template.New(templateName).Parse(post.HtmlContent(textHtml))
	if err != nil {
		post.Logger.Error("Error parsing HTML content", zap.Error(err))
		return nil
	}
	return tmpl
}

func (post *PostService) PublishPostWithSessionUser(searchUsersId, content string) {
	err := post.PostsRepository.AddContentAndUserId(searchUsersId, content)
	if err != nil {
		post.Logger.Error("AddContentAndUserName error: ", zap.Error(err))
	}
}

func (post *PostService) AddContentToPosts() {
	err := post.ClientRedis.ClearRedisCache()
	if err != nil {
		post.Logger.Error("ClearRedisCache error: ", zap.Error(err))
	}

	searchContent, err := post.PostsRepository.GetLastTenPosts()
	if err != nil {
		post.Logger.Error("GetLastTenPosts error: ", zap.Error(err))
	}

	err = post.ClientRedis.AddToCache(searchContent)
	if err != nil {
		post.Logger.Error("AddToCache error: ", zap.Error(err))
	}

}

func (post *PostService) SearchCountPage(page int) response.PageData {
	var PageData response.PageData
	var count float64

	sumPosts, err := post.PostsRepository.CountPosts()
	if err != nil {
		post.Logger.Error("CountPosts error: ", zap.Error(err))
	}
	count = sumPosts / 10.0
	countInt := pkg.FormatInt(count)

	PageData.TotalPages = countInt
	PageData.CurrentPage = page
	return PageData
}

func (post *PostService) CreateTemplateData(page, offset int) response.TemplateData {
	countPage := post.SearchCountPage(page)
	offsetPosts := post.PostsRepository.CalculatePageOffset(offset)

	data := response.TemplateData{
		Posts:      offsetPosts,
		Pagination: countPage,
	}
	return data
}

func (post *PostService) ParsePageAndCalculateOffset(pageStr string) (int, int) {
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		post.Logger.Error("Invalid page parameter", zap.Error(err))
		return 0, 0
	}
	if err != nil || page < 1 {
		page = 1
	}

	offset := (page - 1) * 10
	return page, offset
}
