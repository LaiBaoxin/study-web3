package chain

import (
	"context"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	"log"
	"math/big"
	"os"
)

// InitClient 初始化并返回一个 EthClient
func InitClient() *ethclient.Client {
	// 加载配置文件
	err := godotenv.Load()
	if err != nil {
		log.Fatal("err: 加载配置文件 .env 失败")
		return nil
	}
	rpcUrl := os.Getenv("RPC_URL")
	// 判断是否存在一个连接地址
	if rpcUrl == "" {
		log.Fatal("err: RPC_URL 为空")
		return nil
	}
	// 连接节点
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		log.Fatal("err: 节点连接失败", err)
		return nil
	}

	// 连接测试是否成功，获取 ChainID
	cid, _ := client.ChainID(context.Background())
	log.Println("连接成功，ChainID:", cid)
	return client
}

// InitWSClient 创建一个 WebSocket 客户端
func InitWSClient() *ethclient.Client {
	// 加载配置文件
	err := godotenv.Load()
	if err != nil {
		log.Fatal("err: 加载配置文件 .env 失败")
		return nil
	}
	wsUrl := os.Getenv("WS_URL")
	// 判断是否存在一个连接地址
	if wsUrl == "" {
		log.Fatal("err: WS_URL 为空")
		return nil
	}
	// 创建一个 WebSocket 客户端
	client, err := ethclient.Dial(wsUrl)
	if err != nil {
		log.Fatal("err: ws节点连接失败", err)
		return nil
	}
	return client
}

// GetChainID 获取当前链的 ChainID
func GetChainID(ctx context.Context, client *ethclient.Client) *big.Int {
	cId, err := client.ChainID(ctx)
	if err != nil {
		log.Fatal("err: 获取 ChainID 失败", err)
		return nil
	}
	return cId
}
