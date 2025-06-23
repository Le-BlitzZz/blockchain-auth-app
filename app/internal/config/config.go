package config

import (
	"github.com/Le-BlitzZz/blockchain-auth-app/app/internal/contracts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

type Config struct {
	options      *Options
	redis        *redis.Client
	ethClient    *ethclient.Client
	vipPass      *contracts.VIPPass
	transactOpts *bind.TransactOpts
}

func NewConfig(ctx *cli.Context) (*Config, error) {
	c := &Config{
		options: NewOptions(ctx),
	}

	if err := c.connectEthereum(); err != nil {
		return nil, err
	}
	if err := c.connectRedis(); err != nil {
		c.shutdownEthereum()
		return nil, err
	}
	if err := c.setupContract(); err != nil {
		c.Shutdown()

		return nil, err
	}

	log.Info("config: successfully initialized")

	return c, nil
}

func (c *Config) Shutdown() {
	c.shutdownEthereum()
	c.shutdownRedis()
}

func (c *Config) shutdownEthereum() {
	c.closeEthereum()
	log.Info("close ethereum connection")
}

func (c *Config) shutdownRedis() {
	if err := c.closeRedis(); err != nil {
		log.Errorf("could not close redis connection: %s", err)
	} else {
		log.Info("close redis connection")
	}
}