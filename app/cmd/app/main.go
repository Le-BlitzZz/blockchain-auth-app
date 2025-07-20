package main

import (
	"context"
	"os"
	"os/signal"
	"sync"

	"github.com/Le-BlitzZz/blockchain-auth-app/app/internal/cache"
	"github.com/Le-BlitzZz/blockchain-auth-app/app/internal/config"
	"github.com/Le-BlitzZz/blockchain-auth-app/app/internal/server"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:   "app-service",
		Action: run,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "defaults-yaml",
				Aliases: []string{"y"},
				Value:   "configs/local-docker.yaml",
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Panic(err)
	}
}

func run(ctx *cli.Context) error {
	conf, err := config.NewConfig(ctx)
	if err != nil {
		return err
	}

	cache.SetRedis(conf.Redis())

	cctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		server.Start(cctx, conf)
	}()

	<-cctx.Done()

	log.Info("shutting down...")

	wg.Wait()

	conf.Shutdown()

	log.Info("shutdown complete.")

	return nil
}
