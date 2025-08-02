package main

import (
	"os"

	"github.com/Le-BlitzZz/blockchain-auth-app/client/internal/commands"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()

	app.Name = "A client of the blockchain-auth-app"
	app.Usage = "Interact with the blockchain-auth-app"
	app.Commands = commands.BlockchainAuthApp

	if err := app.Run(os.Args); err != nil {
		log.Fatalf("Error running app: %v", err)
	}
}
