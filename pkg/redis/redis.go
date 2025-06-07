package redis

import (
	"SnickersShopPet1.0/internal/config"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(cfg *config.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: cfg.Redis.Addr,
	})
}
