package server

import (
	"net/http"

	"github.com/Le-BlitzZz/blockchain-auth-app/app/internal/api"
	"github.com/Le-BlitzZz/blockchain-auth-app/app/internal/cache"
	"github.com/gin-gonic/gin"
)

func registerRoutes(router *gin.Engine) {
	registerStaticRoutes(router)

	API := router.Group("/api")
	api.CreateSession(API)
	api.UpdateSession(API)
	api.DeleteSession(API)
	api.StreamSession(API)
}

func registerStaticRoutes(router *gin.Engine) {
	router.Static("/static", "./frontend")

	ui := func(c *gin.Context) {
		sessionId := c.Param("session_id")

		_, err := cache.GetSession(c.Request.Context(), sessionId)
		if err != nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.HTML(http.StatusOK, "sign.html", gin.H{})
	}

	router.GET("/:session_id", ui)
}
