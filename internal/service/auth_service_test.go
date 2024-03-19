package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func TestAddUserWithHashedPassword(t *testing.T) {
	logger := zaptest.NewLogger(t)
	usersRepo := new(MockRepositoryUsers)

	authService := NewAuthService(logger, usersRepo, nil, nil)

	usersRepo.On("AddPasswordAndUserName", "testuser", "testpassword").Return(nil)

	err := authService.AddUserWithHashedPassword("testuser", "testpassword")

	assert.NoError(t, err)
	usersRepo.AssertCalled(t, "AddPasswordAndUserName", "testuser", "testpassword")
}

func TestDeleteSessionCookie(t *testing.T) {
	logger := zaptest.NewLogger(t)
	authService := NewAuthService(logger, nil, nil, nil)

	cookie := authService.DeleteSessionCookie()
	assert.Equal(t, "session_token", cookie.Name)
	assert.Equal(t, "", cookie.Value)
	assert.Equal(t, "/", cookie.Path)
	assert.Equal(t, -1, cookie.MaxAge)
	assert.True(t, cookie.HttpOnly)
}
