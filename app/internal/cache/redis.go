package cache

import (
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

var redisClient *redis.Client

func SetRedis(client *redis.Client) {
	redisClient = client
}

func Redis() *redis.Client {
	if redisClient == nil {
		log.Panic("cache.Redis() called before cache.SetRedis()")
	}
	return redisClient
}
