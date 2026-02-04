package main

import (
	"github.com/ethereum/go-ethereum/common"
	"learn-web3-go/contracts/erc20"
	"learn-web3-go/pkg/chain"
	"learn-web3-go/pkg/chain/model"
	"learn-web3-go/utils"
	"log"
	"math/big"
	"os"
	"time"
)

func main() {
	// 初始化环境
	client := chain.InitClient()

	// 载入用户
	user := model.NewUserFromEnv(client)
	log.Printf("当前用户: %s", user.Address)

	// 构造当前链的交易凭证（Session）
	auth, err := chain.NewAuth(client, user)
	if err != nil {
		log.Fatal("生成凭证失败", err)
		return
	}
	log.Printf("当前 gas 的情况， Tip=%s Wei, MaxFee=%s Wei \n", auth.GasTipCap, auth.GasFeeCap)

	// 连接合约
	usdtAddress := common.HexToAddress(os.Getenv("USDT_CONTRACT_ADDR"))
	toAddress := common.HexToAddress(os.Getenv("TO_WALLET_ADDR"))

	// 是否是有效的地址
	if !utils.IsValidAddress(toAddress.String()) {
		log.Fatal("toAddress 无效的地址")
		return
	}

	amount := big.NewInt(1000000) // 1 USDT

	// 实例化合约
	usdt, err := erc20.NewERC20(usdtAddress, client)
	if err != nil {
		log.Fatal("实例化合约失败", err)
		return
	}

	// 发起交易
	log.Println("正在广播交易中....")
	time.Sleep(5 * time.Second)

	tx, err := usdt.Transfer(auth, toAddress, amount)
	if err != nil {
		log.Fatal("发起交易失败", err)
		return
	}
	log.Println("交易成功，交易 hash:", tx.Hash().Hex())
}
