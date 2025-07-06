package api

import (
	"io"
	"net/http"
	"time"

	"github.com/Le-BlitzZz/blockchain-auth-app/app/internal/cache"
	"github.com/gin-gonic/gin"
)

type Action struct {
	Name string `json:"name" binding:"required,oneof=vip mint burn"`
}

func CreateSession(router *gin.RouterGroup) {
	router.POST("/session", func(c *gin.Context) {
		var action Action
		if err := c.ShouldBindJSON(&action); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		sess, err := cache.NewSession(c, action.Name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, sessionResponse(sess.ID))
	})
}

func UpdateSession(router *gin.RouterGroup) {
	router.PUT("/session/:id", func(c *gin.Context) {
		var timeout time.Duration

		status := c.Query("status")
		switch status {
		case cache.SessionStatusPendingWallet:
			timeout = cache.WalletConnectExpiration
		case cache.SessionStatusPendingSignature:
			timeout = cache.SignMessageExpiration
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
			return
		}

		sessionId := c.Param("id")
		if _, err := cache.GetSession(c, sessionId); err != nil {
			c.JSON(http.StatusBadRequest, sessionError())
			return
		}

		if err := cache.UpdateSession(c, sessionId, status, timeout); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, sessionResponse(sessionId))
	})
}

func DeleteSession(router *gin.RouterGroup) {
	router.DELETE("/session/:id", func(c *gin.Context) {
		sessionId := c.Param("id")

		if _, err := cache.GetSession(c, sessionId); err != nil {
			c.JSON(http.StatusBadRequest, sessionError())
			return
		}

		if err := cache.DeleteSession(c, sessionId); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, sessionResponse(sessionId))
	})
}

func StreamSessionEvents(router *gin.RouterGroup) {
	router.GET("/session/:id/events", func(c *gin.Context) {
		sessionId := c.Param("id")

		c.SSEvent()
		
		// SSE for streaming session events
		c.Stream(func(w io.Writer) bool {
			// Implement SSE logic here

			return true
		})
	})
}

func sessionResponse(sessionId string) gin.H {
	return gin.H{"session_id": sessionId}
}

func sessionError() gin.H {
	return gin.H{"error": "missing or invalid session"}
}
