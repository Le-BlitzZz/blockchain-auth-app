package app

import (
	"log"

	"github.com/Le-BlitzZz/blockchain-auth-app/app/internal/config"
)

var conf *config.Config

func SetConfig(c *config.Config) {
	if c == nil {
		log.Panic("app.SetConfig() called with nil config")
	}

	conf = c
}

func Config() *config.Config {
	if conf == nil {
		log.Panic("app.Config() called before config was set")
	}

	return conf
}
