package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/joho/godotenv"
	"learn-web3-go/contracts/erc20"
	"learn-web3-go/pkg/chain"
	"log"
	"math/big"
	"os"
)

func main() {
	var (
		err error
		ctx = context.Background()
	)

	err = godotenv.Load()
	if err != nil {
		log.Fatal("err: 找不到 .env 文件")
		return
	}
	wsUrl := os.Getenv("WS_URL")
	if wsUrl == "" {
		log.Fatal("err: 找不到 WS_URL")
		return
	}

	// 创建连接
	client := chain.InitWSClient()
	log.Println("ws 节点链接成功...")

	// 过滤条件,只关心 USDT 的合约事件
	contractAddr := common.HexToAddress(os.Getenv("USDT_CONTRACT_ADDR"))

	// 只看 Transfer 事件
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddr},
	}

	// 创建 channel 接收日志
	logs := make(chan types.Log)

	// 发起订阅
	sub, err := client.SubscribeFilterLogs(ctx, query, logs)
	if err != nil {
		log.Fatal("订阅失败:", err)
	}

	// ABI 解析器: 利用生成的合约绑定来解析日志
	usdtFilter, _ := erc20.NewERC20Filterer(contractAddr, client)
	fmt.Println("------- Start Up Lister USDT Event -------")
	fmt.Printf("Lister Address: %s", contractAddr.Hex())

	for {
		select {
		case err = <-sub.Err():
			log.Fatal("订阅中断:", err)

		case vLog := <-logs:
			// 接收到的 transfer 日志进行解析
			event, err := usdtFilter.ParseTransfer(vLog)
			if err != nil {
				log.Fatal("err: 解析失败")
				continue
			}

			// 进行格式化展示输出
			fmt.Println("-------------- 格式化内容如下 -----------")
			fmt.Printf("   区块高度: %d\n", vLog.BlockNumber)
			fmt.Printf("   交易哈希: %s\n", vLog.TxHash.Hex())
			fmt.Printf("   从 (From): %s\n", event.From.Hex())
			fmt.Printf("   到 (To)  : %s\n", event.To.Hex())

			// 格式化金额
			amount := new(big.Float).Quo(new(big.Float).SetInt(event.Value), big.NewFloat(1e6))
			fmt.Printf("   金额     : %.2f USDT\n", amount)
			fmt.Println("------------------------------------------------")

		}
	}
}
