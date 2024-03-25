package service

import (
	"testing"

	"github.com/Vainsberg/WebBlogCraft/internal/mocks"
	"github.com/stretchr/testify/assert"
)

func TestAddUserWithHashedPassword(t *testing.T) {
	usersRepo := new(mocks.MockRepositoryUsers)

	authService := NewAuthService(nil, usersRepo, nil, nil)

	usersRepo.On("AddPasswordAndUserName", "testuser", "testpassword").Return(nil)

	err := authService.AddUserWithHashedPassword("testuser", "testpassword")

	assert.NoError(t, err)
	usersRepo.AssertCalled(t, "AddPasswordAndUserName", "testuser", "testpassword")
}

func TestDeleteSessionCookie(t *testing.T) {
	authService := NewAuthService(nil, nil, nil, nil)

	cookie := authService.DeleteSessionCookie()
	assert.Equal(t, "session_token", cookie.Name)
	assert.Equal(t, "", cookie.Value)
	assert.Equal(t, "/", cookie.Path)
	assert.Equal(t, -1, cookie.MaxAge)
	assert.True(t, cookie.HttpOnly)
}
