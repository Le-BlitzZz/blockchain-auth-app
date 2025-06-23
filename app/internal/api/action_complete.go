package api

import (
	"context"
	"math/big"
	"net/http"

	"github.com/Le-BlitzZz/blockchain-auth-app/app/internal/config"
	"github.com/Le-BlitzZz/blockchain-auth-app/app/internal/contracts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type completeRequest struct {
	SessionId string `json:"session_id" binding:"required,uuid"`
}

func ActionComplete(router *gin.RouterGroup, conf *config.Config) {
	router.POST("/complete", func(c *gin.Context) {
		var req completeRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing or invalid session"})
			return
		}

		redisClient := conf.Redis()

		key := "session:" + req.SessionId
		// 1) Fetch session fields
		vals, err := redisClient.HMGet(c, key, "state", "action", "address").Result()
		if err != nil {
			log.Error("redis error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "redis error"})
			return
		}
		state, _ := vals[0].(string)

		switch state {
		case "":
			c.JSON(http.StatusOK, gin.H{"status": "expired"})
			return
		case "started":
			c.JSON(http.StatusOK, gin.H{"status": "pending_wallet"})
			return
		case "auth_sent":
			c.JSON(http.StatusOK, gin.H{"status": "pending_signature"})
			return
		case "verified":
			action, _ := vals[1].(string)
			address, _ := vals[2].(string)
			owner := common.HexToAddress(address)

			switch action {
			case "vip":
				vipPass := conf.VIPPass()
				balance, err := vipPass.BalanceOf(&bind.CallOpts{Context: context.Background()}, owner)
				if err != nil {
					log.Error("chain call failed:", err)
					c.JSON(http.StatusInternalServerError, gin.H{"error": "chain call failed"})
					return
				}

				if balance.Cmp(big.NewInt(0)) == 0 {
					redisClient.Del(c, key)
					c.JSON(http.StatusForbidden, gin.H{"status": "forbidden"})
					return
				}

				redisClient.Del(c, key)

				c.JSON(http.StatusOK, gin.H{"status": "success", "result": "hello winner"})
			case "mint":
				raw := &contracts.VIPPassRaw{Contract: conf.VIPPass()}
				tx, err := raw.Transact(conf.TransactOpts(), "mint", owner)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "mint tx failed"})
					return
				}

				redisClient.Del(c, key)

				c.JSON(http.StatusOK, gin.H{"status": "success", "result": "minted " + tx.Hash().Hex()})
			case "burn":
				vipPass := conf.VIPPass()

				// Check balance first
				balance, err := vipPass.BalanceOf(&bind.CallOpts{Context: context.Background()}, owner)
				if err != nil {
					log.Error("balance check failed:", err)
					c.JSON(http.StatusInternalServerError, gin.H{"error": "balance check failed"})
					return
				}

				// Check if user has exactly 1 token to burn
				if balance.Cmp(big.NewInt(1)) != 0 {
					redisClient.Del(c, key)
					c.JSON(http.StatusBadRequest, gin.H{
						"status": "error",
						"result": "No VIP pass to burn",
					})
					return
				}

				// Proceed with burn
				raw := &contracts.VIPPassRaw{Contract: conf.VIPPass()}
				tx, err := raw.Transact(conf.TransactOpts(), "burn", owner)
				if err != nil {
					log.Error("burn tx failed:", err)
					c.JSON(http.StatusInternalServerError, gin.H{"error": "burn tx failed"})
					return
				}

				redisClient.Del(c, key)
				c.JSON(http.StatusOK, gin.H{"status": "success", "result": "burned " + tx.Hash().Hex()})
			default:
				c.JSON(http.StatusBadRequest, gin.H{"error": "unknown action"})
			}
		}
	})
}
