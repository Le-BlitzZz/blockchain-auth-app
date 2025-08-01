package api

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Le-BlitzZz/blockchain-auth-app/app/internal/action"
	"github.com/Le-BlitzZz/blockchain-auth-app/app/internal/crypto"

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

		session := cache.NewSession(c.Request.Context(), body.Action)

		if err := session.Create(c.Request.Context()); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, sessionResponse(session.ID))
	})
}

func UpdateSession(router *gin.RouterGroup) {
	router.PATCH("/session/:id", func(c *gin.Context) {
		sessionId := c.Param("id")

		session, err := cache.GetSession(c.Request.Context(), sessionId)
		if err != nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		var reqBody struct {
			Status    *string `json:"status,omitempty" binding:"omitempty,oneof=started pending_wallet declined_signature"`
			Wallet    *string `json:"wallet,omitempty"`
			Signature *string `json:"signature,omitempty"`
		}
		if err := c.ShouldBindJSON(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		var response gin.H
		var timeout time.Duration

		if reqBody.Status == nil {
			if reqBody.Wallet == nil && reqBody.Signature == nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "missing status or wallet"})
				return
			}

			if reqBody.Wallet != nil {
				if reqBody.Signature != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "cannot set both wallet and signature"})
					return
				}

				session.Wallet = reqBody.Wallet
				session.Status = cache.SessionStatusPendingSignature

				nonce, err := crypto.GenerateNonce()
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate nonce"})
					return
				}
				session.Nonce = &nonce

				message := fmt.Sprintf(
					"Sign-In With Ethereum\nAddress: %s\nNonce: %s",
					*session.Wallet, nonce,
				)
				session.Message = &message

				timeout = 30 * time.Second

				response = gin.H{"message": message}
			} else if reqBody.Signature != nil {
				if err := crypto.VerifySignature(*session.Wallet, *session.Message, *reqBody.Signature); err != nil {
					c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
					return
				}

				session.Signature = reqBody.Signature
				session.Status = cache.SessionStatusVerified
				*session.Result, err = action.Result(session.Action, *session.Wallet)
				if err != nil {
					log.Info("Action result failed:", err)
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
			}
		} else {
			switch *reqBody.Status {
			case cache.SessionStatusPendingWallet:
				timeout = 30 * time.Second
			case cache.SessionStatusDeclinedSignature:
				if session.Status != cache.SessionStatusPendingSignature {
					c.JSON(http.StatusBadRequest, gin.H{"error": "cannot decline signature in current state"})
					return
				}
				session.Status = cache.SessionStatusDeclinedSignature
			default:
				timeout = 5 * time.Second
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
				return
			}

			session.Status = *reqBody.Status
			response = sessionResponse(sessionId)
		}

		if err := session.Save(c.Request.Context(), timeout); err != nil {
			log.Error("Failed to save session:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, response)
	})
}

func DeleteSession(router *gin.RouterGroup) {
	router.DELETE("/session/:id", func(c *gin.Context) {
		sessionId := c.Param("id")

		session, err := cache.GetSession(c.Request.Context(), sessionId)
		if err != nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		if session.Status != cache.SessionStatusStarted && session.Status != cache.SessionStatusDeclinedSignature &&
			(session.Status != cache.SessionStatusVerified || session.Result == nil) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "cannot delete session in current state"})
			return
		}

		var response gin.H

		if session.Status == cache.SessionStatusVerified && session.Result != nil {
			response = gin.H{"result": *session.Result}
		}

		if err := session.Delete(c.Request.Context()); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, response)
	})
}

func StreamSession(router *gin.RouterGroup) {
	router.GET("/session/:id/stream", func(c *gin.Context) {
		sessionId := c.Param("id")
		session, err := cache.GetSession(c.Request.Context(), sessionId)
		if err != nil {
			c.AbortWithStatus(http.StatusNotFound)
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
