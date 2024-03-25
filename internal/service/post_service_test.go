package service

import (
	"testing"

	"github.com/Vainsberg/WebBlogCraft/internal/mocks"
	"github.com/stretchr/testify/assert"
)

func TestSearchCountPage(t *testing.T) {
	postRepo := new(mocks.MockRepositoryPost)

	postService := NewPostService(nil, nil, nil, postRepo, nil, nil, nil, nil, nil, nil)

	postRepo.On("CountPosts").Return(nil)

	_, err := postService.SearchCountPage(10)

	assert.NoError(t, err)
}
