package service

import (
	"net/http"
	"time"

	"github.com/Vainsberg/WebBlogCraft/internal/repository"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthUserRepository interface {
	AddPasswordAndUserName(userName, password string) error
	SearchPasswordAndUserName(userName string) (string, error)
	CheckingPresenceUser(username string) (string, error)
	SearchUserName(UserID int) (string, error)
}

type AuthService struct {
	Logger             *zap.Logger
	UsersRepository    AuthUserRepository
	SessionsRepository *repository.RepositorySessions
	PostsRepository    *repository.RepositoryPosts
}

func NewAuthService(logger *zap.Logger,
	AuthUserRepository AuthUserRepository,
	SessionsRepository *repository.RepositorySessions,
	PostsRepository *repository.RepositoryPosts) *AuthService {
	return &AuthService{
		Logger:             logger,
		UsersRepository:    AuthUserRepository,
		SessionsRepository: SessionsRepository,
		PostsRepository:    PostsRepository,
	}
}

func (auth *AuthService) CreateSessionCookie(userName string) *http.Cookie {
	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(24 * time.Hour)
	err := auth.SessionsRepository.AddSessionCookie(sessionToken, userName, expiresAt)
	if err != nil {
		auth.Logger.Error("Error AddSessionCookie:", zap.Error(err))
	}

	return &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: expiresAt,
	}
}

func (auth *AuthService) AddUserWithHashedPassword(userName, userPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userPassword), 8)
	if err != nil {
		auth.Logger.Error("Error GenerateFromPassword:", zap.Error(err))
	}

	err = auth.UsersRepository.AddPasswordAndUserName(userName, string(hashedPassword))
	if err != nil {
		auth.Logger.Error("Error AddPasswordAndUserName:", zap.Error(err))
	}
	return nil
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
