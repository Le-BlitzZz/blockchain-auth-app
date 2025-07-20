package api

import (
	"io"
	"net/http"
	"time"

	"github.com/Le-BlitzZz/blockchain-auth-app/app/internal/cache"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func CreateSession(router *gin.RouterGroup) {
	router.POST("/session", func(c *gin.Context) {
		var body struct {
			Action string `json:"action" binding:"required,oneof=vip mint burn"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		session := cache.NewSession(c, body.Action)

		if err := session.Create(c); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, sessionResponse(session.ID))
	})
}

func UpdateSession(router *gin.RouterGroup) {
	router.PATCH("/session/:id", func(c *gin.Context) {
		sessionId := c.Param("id")

		session, err := cache.GetSession(c, sessionId)
		if err != nil {
			c.JSON(http.StatusBadRequest, sessionError())
			return
		}

		var reqBody struct {
			Status *string `json:"status,omitempty" binding:"omitempty,oneof=started pending_wallet pending_signature"`
			Wallet *string `json:"wallet,omitempty"`
		}
		if err := c.ShouldBindJSON(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		var timeout time.Duration

		if reqBody.Status == nil {
			if reqBody.Wallet == nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "missing status or wallet"})
				return
			}

			session.Wallet = reqBody.Wallet
			timeout = 30 * time.Second
		} else {
			switch *reqBody.Status {
			case cache.SessionStatusPendingWallet:
				timeout = 30 * time.Second
			case cache.SessionStatusPendingSignature:
				timeout = 30 * time.Second
			default:
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
				return
			}

			session.Status = *reqBody.Status
		}

		if err := session.Save(c, timeout); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, sessionResponse(sessionId))
	})
}

func DeleteSession(router *gin.RouterGroup) {
	router.DELETE("/session/:id", func(c *gin.Context) {
		sessionId := c.Param("id")

		session, err := cache.GetSession(c, sessionId)
		if err != nil {
			c.JSON(http.StatusBadRequest, sessionError())
			return
		}

		if err := session.Delete(c); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, sessionResponse(sessionId))
	})
}

func StreamSession(router *gin.RouterGroup) {
	router.GET("/session/:id/stream", func(c *gin.Context) {
		sessionId := c.Param("id")
		session, err := cache.GetSession(c, sessionId)
		if err != nil {
			c.JSON(http.StatusBadRequest, sessionError())
			return
		}

		sessionChan := session.StreamSession(c.Request.Context())

		c.Stream(func(_ io.Writer) bool {
			select {
			case status, ok := <-sessionChan:
				if !ok {
					log.Info("Channel closed, stopping stream")
					return false
				}
				log.Info("Streaming session status:", status)
				c.SSEvent("status", status)
				return true
			case <-c.Request.Context().Done():
				return false
			}
		})
	})
}

func sessionResponse(sessionId string) gin.H {
	return gin.H{"session_id": sessionId}
}

func sessionError() gin.H {
	return gin.H{"error": "missing or invalid session"}
}
