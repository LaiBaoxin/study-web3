package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	"learn-web3-go/contracts/erc20"
	"log"
	"math/big"
	"os"
)

func main() {
	godotenv.Load()

	// 连接 tenderly 的 socket地址
	client, err := ethclient.Dial(os.Getenv("WS_URL"))
	if err != nil {
		log.Fatal(" socket 连接失败", err)
		return
	}
	fmt.Println("tenderly socket 连接成功...")

	// 实例化合约
	usdtAddress := common.HexToAddress(os.Getenv("USDT_CONTRACT_ADDR"))
	usdt, err := erc20.NewERC20(usdtAddress, client)
	if err != nil {
		log.Fatal("err: 创建 USDT 合约实例失败")
		return
	}

	// 创建一个 channel 通道接受事件
	logs := make(chan *erc20.ERC20Transfer)

	// 开始订阅
	// 参数2: logs 通道
	// 参数3/4: 过滤 From/To，这里填 nil 代表监听所有人
	sub, err := usdt.WatchTransfer(&bind.WatchOpts{Context: context.Background()}, logs, nil, nil)
	if err != nil {
		log.Fatal("err: 订阅失败", err)
		return
	}

	fmt.Println("正在监听 Sepolia 链上的所有  USDT 转账记录")

	// 循环监听,接收事件
	for {
		select {
		case err := <-sub.Err():
			log.Fatal("err: 监听订阅失败", err)

		case vLog := <-logs:
			fmt.Println("\n 捕捉到一笔新的转账")
			fmt.Printf("Tx:   %s \n", vLog.Raw.TxHash.Hex())
			fmt.Printf("From: %s \n", vLog.From.Hex())
			fmt.Printf("To:   %s \n", vLog.To.Hex())

			// 处理交易金额
			val := new(big.Float).SetInt(vLog.Value)
			val.Quo(val, big.NewFloat(1000000))
			fmt.Printf("金额:  %.6f USDT", val)

		}
	}

}
