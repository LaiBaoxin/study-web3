package chain

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"learn-web3-go/pkg/chain/model"
	"log"
	"math/big"
)

func NewAuth(client *ethclient.Client, user *model.User) (*bind.TransactOpts, error) {
	// 获取 ChainID
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatal("获取 ChainID 失败", err)
		return nil, err
	}
	// 生成凭证
	auth, err := bind.NewKeyedTransactorWithChainID(user.PrivateKey, chainID)
	if err != nil {
		log.Fatal("生成凭证失败", err)
		return nil, err
	}

	// 动态获取 Gas 费用，获取小费和基础费

	// 小费
	tip, err := client.SuggestGasTipCap(context.Background())
	if err == nil {
		auth.GasTipCap = tip
	} else {
		// 保底给 1 Gwei
		auth.GasTipCap = big.NewInt(1000000000)
	}

	// 基础费
	header, err := client.HeaderByNumber(context.Background(), nil)
	if err == nil {
		baseFree := header.BaseFee
		auth.GasFeeCap = new(big.Int).Add(
			new(big.Int).Mul(baseFree, big.NewInt(2)),
			auth.GasTipCap,
		)
	}
	// 设置为 0，为后续自动构造 CallMsg
	auth.GasLimit = 0

	return auth, nil
}
