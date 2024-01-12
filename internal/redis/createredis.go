package redis

import (
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
)

func NewRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

var CacheClient = cache.New(&cache.Options{
	Redis:      NewRedisClient(),
	LocalCache: cache.NewTinyLFU(10, 0),
})
