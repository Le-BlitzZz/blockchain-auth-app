package config

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

type Options struct {
	HttpHost     string
	HttpPort     int
	RedisAddr    string
	ContractAddr string `yaml:"ContractAddr"`
	DeployerAddr string `yaml:"DeployerAddr"`
	DeployerKey  string `yaml:"DeployerKey"`
	EthClientUrl string `yaml:"EthClientUrl"`
	ChainID      string `yaml:"ChainID"`
}

func NewOptions(ctx *cli.Context) *Options {
	o := &Options{}

	defaultsYaml := ctx.String("defaults-yaml")
	if defaultsYaml == "" {
		log.Tracef("config: defaults file was not specified")
	}

	return o
}

func (o *Options) Load(fileName string) error {
	if fileName == "" {
		return nil
	}

	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return fmt.Errorf("%s not found", fileName)
	}

	yamlConfig, err := os.ReadFile(fileName)

	if err != nil {
		return err
	}

	return yaml.Unmarshal(yamlConfig, o)
}
