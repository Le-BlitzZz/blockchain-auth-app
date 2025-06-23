package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Le-BlitzZz/blockchain-auth-app/app/internal/config"
	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
)

func Start(ctx context.Context, conf *config.Config) {
	router := gin.Default()

	router.LoadHTMLGlob("assets/templates/*")

	registerRoutes(router, conf)

	server := &http.Server{
		Handler: router,
		Addr:    fmt.Sprintf("%s:%d", conf.HttpHost(), conf.HttpPort()),
	}

	go startHttp(server)

	<-ctx.Done()

	if err := server.Close(); err != nil {
		log.Errorf("server: shutdown failed (%s)", err)
	}
}

func startHttp(s *http.Server) {
	if err := s.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			log.Println("server: shutdown complete")
		} else {
			log.Errorf("server: %s", err)
		}
	}
}
