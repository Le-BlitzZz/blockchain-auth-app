package api

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/Le-BlitzZz/blockchain-auth-app/app/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type authRequest struct {
	SessionID string `json:"session_id" binding:"required,uuid"`
	Address   string `json:"address" binding:"required,eth_addr"`
}

func ActionAuth(router *gin.RouterGroup, redisClient *redis.Client) {
	router.POST("/auth", func(c *gin.Context) {
		var req authRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}


		key := "session:" + req.SessionID
		vals, err := redisClient.HMGet(c, key, "state", "action").Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "redis error"})
			return
		}

		state, _ := vals[0].(string)
		action, _ := vals[1].(string)
		if state != "started" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "session not in started state"})
			return
		}

		// 2) Generate a 128-bit nonce (hex)
		b := make([]byte, 16)
		if _, err := rand.Read(b); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate nonce"})
			return
		}
		nonce := hex.EncodeToString(b)

		// 3) Store address, nonce, update state, reset TTL to 60 sec
		if err := redisClient.HSet(c, key, map[string]any{
			"state":   "auth_sent",
			"address": req.Address,
			"nonce":   nonce,
			"action":  action,
		}).Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update session"})
			return
		}
		if err := redisClient.Expire(c, key, config.SignMessageExpiration).Err(); err != nil {
			c.Error(fmt.Errorf("failed to reset TTL: %w", err))
		}

		// 4) Build SIWE message
		// issuedAt := time.Now().UTC().Format(time.RFC3339)
		message := fmt.Sprintf(
			"Sign-In With Ethereum\nAddress: %s\nNonce: %s",
			req.Address, nonce,
		)

		c.JSON(http.StatusOK, gin.H{"message": message})
	})
}
