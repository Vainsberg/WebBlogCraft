package service

import (
	"html/template"
	"io/ioutil"
	"strconv"

	"github.com/Vainsberg/WebBlogCraft/internal/dto"
	"github.com/Vainsberg/WebBlogCraft/internal/pkg"
	"github.com/Vainsberg/WebBlogCraft/internal/rabbitmq"
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
	LikesRepository    *repository.RepositoryLikes
	CommentsRepository *repository.RepositoryComments
	EmailRepository    *repository.RepositoryEmail
	RabbitMQRepository *rabbitmq.RepositoryRabbitMQ
	ClientRedis        *redis.RedisClient
	Cache              *cache.Cache
}

func NewPostService(logger *zap.Logger,
	UsersRepository *repository.RepositoryUsers,
	SessionsRepository *repository.RepositorySessions,
	PostsRepository *repository.RepositoryPosts,
	LikesRepository *repository.RepositoryLikes,
	CommentsRepository *repository.RepositoryComments,
	EmailRepository *repository.RepositoryEmail,
	RabbitMQRepository *rabbitmq.RepositoryRabbitMQ,
	ClientRedis *redis.RedisClient,
	cache *cache.Cache) *PostService {
	return &PostService{
		Logger:             logger,
		UsersRepository:    UsersRepository,
		SessionsRepository: SessionsRepository,
		PostsRepository:    PostsRepository,
		LikesRepository:    LikesRepository,
		CommentsRepository: CommentsRepository,
		EmailRepository:    EmailRepository,
		RabbitMQRepository: RabbitMQRepository,
		ClientRedis:        ClientRedis,
		Cache:              cache,
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

func (post *PostService) PublishPostWithSessionUser(searchUsersId int, content string) {
	err := post.PostsRepository.AddContentAndUserId(searchUsersId, content)
	if err != nil {
		post.Logger.Error("AddContentAndUserName error: ", zap.Error(err))
	}
}

func (post *PostService) AddContentToRedis() []dto.PostDto {
	cachekey := "all_posts"

	searchContentRedis, err := post.ClientRedis.GetRedisValue(cachekey)
	if err == nil {
		return searchContentRedis
	}

	searchContent, err := post.PostsRepository.GetLastTenPosts()
	if err != nil {
		post.Logger.Error("GetLastTenPosts error: ", zap.Error(err))
	}

	err = post.ClientRedis.AddToCache(searchContent, cachekey)
	if err != nil {
		post.Logger.Error("AddToCache error: ", zap.Error(err))
	}

	searchContentRedis, err = post.ClientRedis.GetRedisValue(cachekey)
	if err != nil {
		post.Logger.Error("GetRedisValue error: ", zap.Error(err))
	}

	return searchContentRedis
}

func (post *PostService) SearchCountPage(page int) response.PageData {
	var count float64

	sumPosts, err := post.PostsRepository.CountPosts()
	if err != nil {
		post.Logger.Error("CountPosts error: ", zap.Error(err))
	}
	count = sumPosts / 10.0
	countInt := pkg.FormatInt(count)

	PageList := pkg.CreatePageList(countInt, page)
	return PageList
}

func (post *PostService) GenerateTemplateDataPosts(page, offset int) response.ResponseData {
	countPage := post.SearchCountPage(page)
	offsetPosts, err := post.PostsRepository.CalculatePageOffset(offset)
	if err != nil {
		post.Logger.Error("CalculatePageOffset error: ", zap.Error(err))
	}

	posts, err := post.GetPostsWithComments(offsetPosts)
	if err != nil {
		post.Logger.Error("GetPostsWithComments error: ", zap.Error(err))
	}

	data := response.ResponseData{
		Posts:      posts,
		Pagination: countPage,
	}
	return data
}

func (post *PostService) GenerateTemplateDataPostsRedis(page int) response.ResponseData {
	countPage := post.SearchCountPage(page)
	postsRedis := post.AddContentToRedis()

	posts, err := post.GetPostsWithComments(postsRedis)
	if err != nil {
		post.Logger.Error("GetPostsWithComments error: ", zap.Error(err))
	}

	data := response.ResponseData{
		Posts:      posts,
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

func (post *PostService) ClearRedisCache() {
	err := post.ClientRedis.ClearRedisCache()
	if err != nil {
		post.Logger.Error("ClearRedisCache error: ", zap.Error(err))
	}
}

func (post *PostService) ProcessLikeAction(cookie, postIDStr string) (int, error) {
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		return 0, err
	}
	userID, err := post.SessionsRepository.SearchUsersIdSessionCookie(cookie)
	if err != nil {
		post.Logger.Error("SearchUsersIdSessionCookie error:", zap.Error(err))
	}

	chekingLikeToPost, err := post.LikesRepository.CheckingLikes(userID, postID)
	if err != nil {
		post.Logger.Error("CheckingLikes error:", zap.Error(err))
	}

	if !chekingLikeToPost {
		err := post.LikesRepository.AddLikesToPost(userID, postID)
		if err != nil {
			post.Logger.Error("AddLikesToPost error:", zap.Error(err))
		}
	} else {
		err = post.LikesRepository.RemoveLikeFromPost(userID, postID)
		if err != nil {
			post.Logger.Error("RemoveLikeFromPost error:", zap.Error(err))
		}
	}

	countLikes, err := post.LikesRepository.CountLikes(postID)
	if err != nil {
		post.Logger.Error("CountLikes error:", zap.Error(err))
	}

	if countLikes <= 10 {
		post.ClearRedisCache()
	}
	return countLikes, nil
}

func (post *PostService) AddUserCommentToPostAndSearchUserName(cookie, postIDStr, comment string) (int, string, error) {
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		return 0, "", err
	}

	userID, err := post.SessionsRepository.SearchUsersIdSessionCookie(cookie)
	if err != nil {
		post.Logger.Error("SearchUsersIdSessionCookie error:", zap.Error(err))
	}
	post.CommentsRepository.AddCommentToPost(comment, userID, postID)

	userName, err := post.UsersRepository.SearchUserName(userID)
	if err != nil {
		post.Logger.Error("SearchUserName error:", zap.Error(err))
	}

	commentId, err := post.CommentsRepository.SearchCommentId(comment)
	if err != nil {
		post.Logger.Error("SearchCommentId error:", zap.Error(err))
	}

	return commentId, userName, nil
}

func (post *PostService) GetPostsWithComments(offsetPosts []dto.PostDto) ([]response.Post, error) {
	var posts []response.Post

	for _, p := range offsetPosts {
		Comments, err := post.CommentsRepository.GetComments(p.PostId)
		if err != nil {
			post.Logger.Error("GetComments error: ", zap.Error(err))
		}
		data := response.Post{
			Content:  p.Content,
			PostId:   p.PostId,
			UserName: p.UserName,
			Likes:    p.Likes,
			Comments: Comments,
		}
		posts = append(posts, data)
	}
	return posts, nil
}

func (post *PostService) LikeActionToComment(cookie, commentIDStr string) (int, error) {
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		return 0, err
	}
	userID, err := post.SessionsRepository.SearchUsersIdSessionCookie(cookie)
	if err != nil {
		post.Logger.Error("SearchUsersIdSessionCookie error:", zap.Error(err))
	}

	chekingLikeToComments, err := post.CommentsRepository.CheckingLikesToComment(userID, commentID)
	if err != nil {
		post.Logger.Error("CheckingLikesToComment error:", zap.Error(err))
	}

	if !chekingLikeToComments {
		err := post.CommentsRepository.AddLikesToComments(userID, commentID)
		if err != nil {
			post.Logger.Error("AddLikesToComments error:", zap.Error(err))
		}
	} else {
		err = post.CommentsRepository.RemoveLikeFromComments(userID, commentID)
		if err != nil {
			post.Logger.Error("RemoveLikeFromPost error:", zap.Error(err))
		}
	}
	countLikes, err := post.CommentsRepository.CountLikesComments(commentID)
	if err != nil {
		post.Logger.Error("CountLikes error:", zap.Error(err))
	}
	return countLikes, nil
}

func (post *PostService) AddEmailInDB(cookie, email string) {

	userID, err := post.SessionsRepository.SearchUsersIdSessionCookie(cookie)
	if err != nil {
		post.Logger.Error("SearchUsersIdSessionCookie error:", zap.Error(err))
	}
	post.EmailRepository.AddEmailAndUserId(userID, cookie)
}

func (post *PostService) AddCodeToRedis(code int) dto.EmailCode {
	cachekey := "email"

	searchCode, err := post.ClientRedis.GetRedisCode(cachekey)
	if err == nil {
		return searchCode
	}

	err = post.ClientRedis.AddToCacheCode(code, cachekey)
	if err != nil {
		post.Logger.Error("AddToCache error: ", zap.Error(err))
	}

	searchCode, err = post.ClientRedis.GetRedisCode(cachekey)
	if err != nil {
		post.Logger.Error("GetRedisValue error: ", zap.Error(err))
	}

	return searchCode
}
