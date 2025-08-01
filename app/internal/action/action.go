package action

import (
	"context"
	"errors"
	"math/big"

	"github.com/Le-BlitzZz/blockchain-auth-app/app/internal/app"
	"github.com/Le-BlitzZz/blockchain-auth-app/app/internal/config"
	"github.com/Le-BlitzZz/blockchain-auth-app/app/internal/contracts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

func Result(action string, address string) (string, error) {
	conf := app.Config()
	owner := common.HexToAddress(address)
	switch action {
	case "vip":
		balance, err := fetchBalance(conf, owner, 0)
		if err != nil {
			return "", err
		}
		if balance == 0 {
			return "You dont hold a VIP pass!", nil
		}

		return "Hello Winner!", nil
	case "mint":
		balance, err := fetchBalance(conf, owner, 0)
		if err != nil {
			return "", err
		}

		if balance != 0 {
			return "You already hold a VIP pass!", nil
		}

		raw := &contracts.VIPPassRaw{Contract: conf.VIPPass()}
		tx, err := raw.Transact(conf.TransactOpts(), "mint", owner)
		if err != nil {
			return "", err
		}

		return "Minted: " + tx.Hash().Hex(), nil
	case "burn":
		balance, err := fetchBalance(conf, owner, 1)
		if err != nil {
			return "", err
		}

		if balance != 0 {
			return "No VIP pass to burn!", nil
		}

		raw := &contracts.VIPPassRaw{Contract: conf.VIPPass()}
		tx, err := raw.Transact(conf.TransactOpts(), "burn", owner)
		if err != nil {
			log.Error("burn tx failed:", err)
			return "", err
		}

		return "Burned: " + tx.Hash().Hex(), nil
	default:
		return "", errors.New("unknown action")
	}
}

func fetchBalance(conf *config.Config, owner common.Address, x int64) (int, error) {
	vipPass := conf.VIPPass()
	balance, err := vipPass.BalanceOf(&bind.CallOpts{Context: context.Background()}, owner)
	if err != nil {
		return 0, err
	}

	return balance.Cmp(big.NewInt(x)), nil

}
