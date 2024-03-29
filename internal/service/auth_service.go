package service

import (
	"net/http"
	"time"

	"github.com/Vainsberg/WebBlogCraft/internal/repository"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	Logger             *zap.Logger
	UsersRepository    *repository.RepositoryUsers
	SessionsRepository *repository.RepositorySessions
	PostsRepository    *repository.RepositoryPosts
}

func NewAuthService(logger *zap.Logger,
	UsersRepository *repository.RepositoryUsers,
	SessionsRepository *repository.RepositorySessions,
	PostsRepository *repository.RepositoryPosts) *AuthService {
	return &AuthService{
		Logger:             logger,
		UsersRepository:    UsersRepository,
		SessionsRepository: SessionsRepository,
		PostsRepository:    PostsRepository,
	}
}

func (auth *AuthService) CreateSessionCookie(userName string) *http.Cookie {
	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(24 * time.Hour)
	auth.SessionsRepository.AddSessionCookie(sessionToken, userName, expiresAt)

	return &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: expiresAt,
	}
}

func (auth *AuthService) AddUserWithHashedPassword(userName, userPassword string) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userPassword), 8)
	if err != nil {
		auth.Logger.Error("Error GenerateFromPassword:", zap.Error(err))
	}
	auth.UsersRepository.AddPasswordAndUserName(userName, string(hashedPassword))
}

func (auth *AuthService) SearchPassword(userName string) string {
	serchPassword, err := auth.UsersRepository.SearchPasswordAndUserName(userName)
	if err != nil {
		auth.Logger.Error("Error SearchPasswordAndUserName:", zap.Error(err))
	}
	return serchPassword
}

func (auth *AuthService) CheckUserExistence(userName string) string {
	searchName, err := auth.UsersRepository.CheckingPresenceUser(userName)
	if err != nil {
		auth.Logger.Error("CheckingPresenceUser error: ", zap.Error(err))
	}
	return searchName
}

func (auth *AuthService) DeleteSessionCookie() *http.Cookie {
	return &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}
}
