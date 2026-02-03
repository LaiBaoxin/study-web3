package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {

	// load
	err := godotenv.Load()
	if err != nil {
		log.Fatal("err: 找不到 .env 文件")
		return
	}

	rpcURL := os.Getenv("RPC_URL")
	if rpcURL == "" {
		log.Fatal("err: 缺少 rpcURL 的值")
	}

	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatal(err)
	}

	// 获取 ChainID
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("连接成功，ChainID: %s\n", chainID.String())

	// 查询区块链高度
	blockNumber, _ := client.BlockNumber(context.Background())
	fmt.Printf("当前区块高度: %d\n", blockNumber)
}
