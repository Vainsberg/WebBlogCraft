package redis

import (
	"context"
	"time"

	"github.com/Vainsberg/WebBlogCraft/internal/response"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	Client *redis.Client
}

func NewRepositoryRedis(client *redis.Client) *RedisClient {
	return &RedisClient{Client: client}
}

func (redis *RedisClient) SearchLastPostID(posts response.StoragePosts) string {
	postid := posts.PostsID[len(posts.PostsID)-1]
	posts.Posts = posts.Posts[:len(posts.Posts)-1]
	posts.PostsID = posts.PostsID[:len(posts.PostsID)-1]
	return postid
}

func (r *RedisClient) DeleteFromCache(c *cache.Cache, key string) error {
	return r.Client.Del(context.Background(), key).Err()
}

func (r *RedisClient) AddToCache(key string, value string, expiration time.Duration) error {
	err := r.Client.Set(context.Background(), key, value, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}
