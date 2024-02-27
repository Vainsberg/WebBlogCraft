package redis

import (
	config "github.com/Vainsberg/WebBlogCraft/internal/config"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
)

func NewRedisClient(cfg *config.Сonfigurations) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPass,
		DB:       cfg.RedisDb,
	})
}
func CreateRedisCache(cfg *config.Сonfigurations) *cache.Cache {
	return cache.New(&cache.Options{
		Redis:      NewRedisClient(cfg),
		LocalCache: cache.NewTinyLFU(10, 0),
	})
}
