package redis

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Vainsberg/WebBlogCraft/internal/pkg"
	"github.com/Vainsberg/WebBlogCraft/internal/response"
	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	Client *redis.Client
}

func NewRepositoryRedis(client *redis.Client) *RedisClient {
	return &RedisClient{Client: client}
}

func (r *RedisClient) AddToCache(searchContent []response.PostsRedis) error {
	for _, v := range searchContent {

		jsonContent, err := json.Marshal(v.Content)
		if err != nil {
			log.Fatal(err)
		}

		err = r.Client.Set(context.Background(), pkg.GenerateUserID(), jsonContent, 0).Err()
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *RedisClient) ClearRedisCache() error {
	ctx := context.Background()
	err := r.Client.FlushAll(ctx).Err()
	if err != nil {
		return err
	}
	return nil
}
