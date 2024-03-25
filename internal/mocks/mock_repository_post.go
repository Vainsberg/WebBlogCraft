package mocks

import (
	"github.com/Vainsberg/WebBlogCraft/internal/dto"
	"github.com/Vainsberg/WebBlogCraft/internal/response"
	"github.com/stretchr/testify/mock"
)

type MockRepositoryPost struct {
	mock.Mock
}

func (m *MockRepositoryPost) AddContentAndUserId(UsersId int, content string) error {
	args := m.Called(UsersId, content)
	return args.Error(0)
}

func (m *MockRepositoryPost) ContentOutput() (*response.Post, error) {
	args := m.Called()
	return &response.Post{}, args.Error(0)
}

func (m *MockRepositoryPost) CalculatePageOffset(offset int) ([]dto.PostDto, error) {
	args := m.Called(offset)
	return []dto.PostDto{}, args.Error(0)
}

func (m *MockRepositoryPost) CountPosts() (float64, error) {
	args := m.Called()
	return 0, args.Error(0)
}

func (m *MockRepositoryPost) GetLastTenPosts() ([]dto.PostDto, error) {
	args := m.Called()
	return []dto.PostDto{}, args.Error(0)
}
