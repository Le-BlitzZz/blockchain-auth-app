package api

import (
	"net/http"

	"github.com/Le-BlitzZz/blockchain-auth-app/app/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type startRequest struct {
	Action string `json:"action" binding:"required,oneof=vip mint burn"`
}

func ActionStart(router *gin.RouterGroup, redisClient *redis.Client) {
	router.POST("/start", func(c *gin.Context) {
		var req startRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "invalid action"})
			return
		}

		sid := uuid.New().String()
		key := "session:" + sid
		entry := map[string]any{
			"state":  "started",
			"action": req.Action,
		}

		if err := redisClient.HSet(c, key, entry).Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create session"})
			return
		}

		if err := redisClient.Expire(c, key, config.WalletConnectExpiration).Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set session expiration"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"session_id": sid})
	})
}
