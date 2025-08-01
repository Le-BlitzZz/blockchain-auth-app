package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type verifyRequest struct {
	SessionId string `json:"session_id" binding:"required,uuid"`
	Message   string `json:"message" binding:"required"`
	Signature string `json:"signature" binding:"required"`
}

func ActionVerify(router *gin.RouterGroup, redisClient *redis.Client) {
	router.POST("/verify", func(c *gin.Context) {
		var req verifyRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		key := "session:" + req.SessionId
		vals, err := redisClient.HMGet(c, key, "state", "action", "address", "nonce").Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "redis error"})
			return
		}
		state, _ := vals[0].(string)
		// action, _ := vals[1].(string)
		address, _ := vals[2].(string)
		nonce, _ := vals[3].(string)

		if state != "auth_sent" {
			c.JSON(http.StatusGone, gin.H{"error": "session not in auth_sent state"})
			return
		}

		expectedMessage := fmt.Sprintf(
			"Sign-In With Ethereum\nAddress: %s\nNonce: %s",
			address, nonce,
		)

		if req.Message != expectedMessage {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid message"})
			return
		}

		if !verifyEthereumSignature(address, req.Message, req.Signature) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid signature"})
			return
		}

		if err := redisClient.HSet(c, key, "state", "verified").Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update session"})
			return
		}
		if err := redisClient.Expire(c, key, 60*time.Second).Err(); err != nil {
			c.Error(fmt.Errorf("failed to reset TTL: %w", err))
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
}

func verifyEthereumSignature(address, message, signature string) bool {
	// 1. Decode the hex signature
	sigBytes, err := hexutil.Decode(signature)
	if err != nil {
		fmt.Printf("Failed to decode signature: %v\n", err)
		return false
	}

	// 2. Ethereum signatures are exactly 65 bytes
	if len(sigBytes) != 65 {
		fmt.Printf("Invalid signature length: %d\n", len(sigBytes))
		return false
	}

	// 3. Adjust recovery ID (Ethereum uses 27/28, go-ethereum expects 0/1)
	if sigBytes[64] >= 27 {
		sigBytes[64] -= 27
	}

	// 4. Hash the message using Ethereum's personal sign format
	// This prefixes the message with "\x19Ethereum Signed Message:\n{length}"
	messageHash := accounts.TextHash([]byte(message))

	// 5. Recover the public key from signature + hash
	publicKey, err := crypto.SigToPub(messageHash, sigBytes)
	if err != nil {
		fmt.Printf("Failed to recover public key: %v\n", err)
		return false
	}

	// 6. Derive address from public key
	recoveredAddress := crypto.PubkeyToAddress(*publicKey)

	// 7. Compare with expected address
	expectedAddress := common.HexToAddress(address)

	isValid := recoveredAddress == expectedAddress

	return isValid
}
