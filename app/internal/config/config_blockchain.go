package config

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/Le-BlitzZz/blockchain-auth-app/app/internal/contracts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"
)

func (c *Config) ContractAddr() string {
	if c.options.ContractAddr == "" {
		return "0x5FbDB2315678afecb367f032d93F642f64180aa3"
	}

	return c.options.ContractAddr
}

func (c *Config) DeployerAddr() string {
	if c.options.DeployerAddr == "" {
		return "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
	}

	return c.options.DeployerAddr
}

func (c *Config) DeployerKey() string {
	if c.options.DeployerKey == "" {
		return "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	}

	return c.options.DeployerKey
}

func (c *Config) EthClientUrl() string {
	if c.options.EthClientUrl == "" {
		return "http://hardhat:8545/"
	}

	return c.options.EthClientUrl
}

func (c *Config) ChainID() string {
	if c.options.ChainID == "" {
		return "31337"
	}

	return c.options.ChainID
}

func (c *Config) EthClient() *ethclient.Client { return c.ethClient }

func (c *Config) VIPPass() *contracts.VIPPass { return c.vipPass }

func (c *Config) TransactOpts() *bind.TransactOpts { return c.transactOpts }

func (c *Config) connectEthereum() error {
	rpc := c.EthClientUrl()
	client, err := ethclient.Dial(rpc)
	if err != nil {
		return fmt.Errorf("ethclient dial %s: %w", rpc, err)
	}

	log.Infof("ethereum: connected to %s", rpc)

	c.ethClient = client

	pkBytes, err := hex.DecodeString(c.DeployerKey())
	if err != nil {
		return fmt.Errorf("invalid private key: %w", err)
	}
	pk, err := crypto.ToECDSA(pkBytes)
	if err != nil {
		return fmt.Errorf("parse private key: %w", err)
	}
	chainID, ok := new(big.Int).SetString(c.ChainID(), 10)
	if !ok {
		return fmt.Errorf("invalid chain ID: %s", c.ChainID())
	}

	auth, err := bind.NewKeyedTransactorWithChainID(pk, chainID)
	if err != nil {
		return fmt.Errorf("create transactor: %w", err)
	}

	c.transactOpts = auth

	return nil
}

func (c *Config) closeEthereum() error {
	if c.ethClient != nil {
		c.ethClient.Close()
	}

	return nil
}

func (c *Config) setupContract() error {
	addr := common.HexToAddress(c.ContractAddr())

	vipPass, err := contracts.NewVIPPass(addr, c.ethClient)
	if err != nil {
		return fmt.Errorf("instantiate VIPPass: %w", err)
	}

	c.vipPass = vipPass

	log.Infof("contract: VIPPass at %s", addr.Hex())

	return nil
}
