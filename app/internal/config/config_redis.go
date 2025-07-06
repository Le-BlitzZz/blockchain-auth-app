package config

import (
	"context"
	"time"

	"github.com/Le-BlitzZz/blockchain-auth-app/app/internal/cache"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

const WalletConnectExpiration = 10 * time.Second
const SignMessageExpiration = 10 * time.Second

func (c *Config) Redis() *redis.Client {
	if c.redis == nil {
		log.Panicf("config: redis not initialized")
	}

	return c.redis
}

func (c *Config) RedisAddr() string {
	if c.options.RedisAddr == "" {
		return "redis:6379"
	}

	return c.options.RedisAddr
}

func (c *Config) connectRedis() error {
	client := redis.NewClient(&redis.Options{
		Addr: c.RedisAddr(),
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return err
	}

	log.Info("redis: opened connection")

	c.redis = client

	return nil
}

func (c *Config) closeRedis() error {
	if c.redis != nil {
		if err := c.redis.Close(); err == nil {
			c.redis = nil
			cache.SetRedis(nil)
		} else {
			return err
		}
	}

	return nil
}
