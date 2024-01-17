package redis

import (
	"context"
	"encoding/json"
	"fmt"
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

func (r *RedisClient) AddToCache(searchContent []response.PostsRedis, cachekeys response.PostsIdRedis) error {
	for _, v := range searchContent {
		postId := pkg.GenerateUserID()
		cachekeys.PostId = append(cachekeys.PostId, postId)
		jsonContent, err := json.Marshal(v.Content)
		if err != nil {
			log.Fatal(err)
		}

		err = r.Client.Set(context.Background(), postId, jsonContent, 0).Err()
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

func (r *RedisClient) GetRedisValue(cacheKey response.PostsIdRedis) (response.PostsRedis, error) {
	var postsRedis response.PostsRedis
	ctx := context.Background()
	for _, v := range cacheKey.PostId {
		value, err := r.Client.Get(ctx, v).Result()
		if err == redis.Nil {
			fmt.Println("Key no search")
		} else if err != nil {
			return response.PostsRedis{}, err
		} else {
			postsRedis.Content = append(postsRedis.Content, value)
		}

	}

	return postsRedis, nil
}
