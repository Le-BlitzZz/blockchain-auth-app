package api

import (
	"net/http"
	"time"

	"github.com/Le-BlitzZz/blockchain-auth-app/app/internal/request"
	"github.com/gin-gonic/gin"
)

type Action struct {
	Name string `json:"name" binding:"required,oneof=vip mint burn"`
}

func CreateRequest(router *gin.RouterGroup) {
	router.POST("/request", func(c *gin.Context) {
		var action Action
		if err := c.ShouldBindJSON(&action); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		req, err := request.NewRequest(c, action.Name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, requestResponse(req.ID))
	})
}

func UpdateRequest(router *gin.RouterGroup) {
	router.PUT("/request/:id", func(c *gin.Context) {
		var timeout time.Duration

		status := c.Query("status")
		switch status {
		case request.RequestStatusPendingWallet:
			timeout = request.WalletConnectExpiration
		case request.RequestStatusPendingSignature:
			timeout = request.SignMessageExpiration
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
			return
		}

		requestId := c.Param("id")
		if _, err := request.GetRequest(c, requestId); err != nil {
			c.JSON(http.StatusBadRequest, requestError())
			return
		}

		if err := request.UpdateRequest(c, requestId, status, timeout); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, requestResponse(requestId))
	})
}

func DeleteRequest(router *gin.RouterGroup) {
	router.DELETE("/request/:id", func(c *gin.Context) {
		requestId := c.Param("id")

		if _, err := request.GetRequest(c, requestId); err != nil {
			c.JSON(http.StatusBadRequest, requestError())
			return
		}

		if err := request.DeleteRequest(c, requestId); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, requestResponse(requestId))
	})
}

func requestResponse(requestId string) gin.H {
	return gin.H{"request_id": requestId}
}

func requestError() gin.H {
	return gin.H{"error": "missing or invalid request"}
}
