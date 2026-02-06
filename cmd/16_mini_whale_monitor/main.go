package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	ens "github.com/wealdtech/go-ens/v3"
	"log"
	"math/big"
	"os"

	"learn-web3-go/contracts/erc20"
	"learn-web3-go/pkg/chain"
)

// WhaleThreshold USDT 精度是 6，所以是 10000 * 10^6
var WhaleThreshold = new(big.Int).Mul(big.NewInt(10000), big.NewInt(1e6))

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("加载 .env 文件失败")
		return
	}

	// 连接 WebSocket
	client := chain.InitWSClient()
	fmt.Println("监听器启动... ")

	// 准备过滤条件
	usdtAddr := common.HexToAddress(os.Getenv("USDT_CONTRACT_ADDR"))
	query := ethereum.FilterQuery{
		Addresses: []common.Address{usdtAddr},
	}

	// 订阅事件
	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatal("订阅失败:", err)
	}

	// 初始化 ABI 解析器
	usdtFilter, _ := erc20.NewERC20Filterer(usdtAddr, client)

	// 开始循环监听
	for {
		select {
		case err := <-sub.Err():
			log.Fatal("连接断开:", err)

		case vLog := <-logs:
			// 解析日志
			event, err := usdtFilter.ParseTransfer(vLog)
			if err != nil {
				continue
			}

			// 筛选大额交易
			// 如果 event.Value < 10000_000000，就跳过
			if event.Value.Cmp(WhaleThreshold) < 0 {
				continue
			}

			// 格式化金额
			amountFloat, _ := new(big.Float).Quo(new(big.Float).SetInt(event.Value), big.NewFloat(1e6)).Float64()

			// 尝试 ENS 反向解析
			fromName := getEnsName(client, event.From)
			toName := getEnsName(client, event.To)

			// 触发报警
			fmt.Printf("金额: %.2f USDT\n", amountFloat)
			fmt.Printf("发送方: %s (%s)\n", fromName, event.From.Hex())
			fmt.Printf("接收方: %s (%s)\n", toName, event.To.Hex())
			fmt.Printf("TxHash: %s\n", vLog.TxHash.Hex())
			fmt.Println("--------------------------------------------------")
		}
	}
}

// getEnsName 尝试获取 ENS 域名
func getEnsName(client *ethclient.Client, addr common.Address) string {

	name, err := ens.ReverseResolve(client, addr)
	if err != nil {
		return "未知"
	}
	return name
}
