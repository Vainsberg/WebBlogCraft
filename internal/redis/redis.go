package redis

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Vainsberg/WebBlogCraft/internal/dto"
	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	Client *redis.Client
}

func NewRepositoryRedis(client *redis.Client) *RedisClient {
	return &RedisClient{Client: client}
}

func (r *RedisClient) AddToCache(searchContent []dto.PostDto, cachekey string) error {

	jsonContent, err := json.Marshal(searchContent)
	if err != nil {
		log.Fatal(err)
	}

	err = r.Client.Set(context.Background(), cachekey, jsonContent, 0).Err()
	if err != nil {
		return err

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

func (r *RedisClient) GetRedisValue(cacheKey string) ([]dto.PostDto, error) {
	var postRedis []dto.PostDto

	ctx := context.Background()

	value, err := r.Client.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return nil, err
	} else if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(value), &postRedis)
	if err != nil {
		return nil, err
	}

	return postRedis, nil
}
