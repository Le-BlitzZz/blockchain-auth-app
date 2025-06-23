package server

import (
	"net/http"

	"github.com/Le-BlitzZz/blockchain-auth-app/app/internal/api"
	"github.com/Le-BlitzZz/blockchain-auth-app/app/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

func registerRoutes(router *gin.Engine, conf *config.Config) {
	redis := conf.Redis()

	registerStaticRoutes(router, redis)

	API := router.Group("/api")

	api.ActionStart(API, redis)
	api.ActionAuth(API, redis)
	api.ActionVerify(API, redis)
	api.ActionComplete(API, conf)
}

func registerStaticRoutes(router *gin.Engine, redisClient *redis.Client) {
	router.Static("/static", "./frontend")

	ui := func(c *gin.Context) {
		log.Info("UI triggered")
		sessionId := c.Param("session_id")

		// check if "session:" + sid exists in Redis
		exists, err := redisClient.Exists(c, "session:"+sessionId).Result()
		if err != nil {
			log.Error("redis error:", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if exists == 0 {
			log.Infof("Session %s not found in Redis", sessionId)
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.HTML(http.StatusOK, "sign.html", gin.H{})
	}

	router.GET("/:session_id", ui)
}
