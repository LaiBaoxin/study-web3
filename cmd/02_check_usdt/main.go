package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	"learn-web3-go/contracts/erc20"
	"log"
	"math/big"
	"os"
)

func main() {
	// Load 加载 .env 文件
	err := godotenv.Load()
	if err != nil {
		log.Fatal("err: 找不到 .env 文件")
		return
	}

	// 读取变量
	rpcURL := os.Getenv("RPC_URL")
	contractAddrStr := os.Getenv("USDT_CONTRACT_ADDR")
	walletAddrStr := os.Getenv("MY_WALLET_ADDR")

	if rpcURL == "" || contractAddrStr == "" || walletAddrStr == "" {
		log.Fatal("错误: .env 文件中缺少必要的配置项 (RPC_URL, USDT_CONTRACT_ADDR, MY_WALLET_ADDR)")
		return
	}

	// 连接 tenderly
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatal(err)
	}
	// USDT 合约主网地址 https://goto.etherscan.com/token/0xdac17f958d2ee523a2206206994597c13d831ec7
	usdtAddr := common.HexToAddress(contractAddrStr)

	// 实例化合约绑定
	usdtInstance, err := erc20.NewERC20(usdtAddr, client)
	if err != nil {
		log.Fatal(err)
	}
	// 查询钱包地址
	myWallet := common.HexToAddress(walletAddrStr)

	// 调用合约内的查询余额方法
	bal, err := usdtInstance.BalanceOf(nil, myWallet)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("原生的余额： %s Wei \n", bal.String())

	// 精度的转换（USDT 6 位小数） 100 USDT = 100,000,000
	fBal := new(big.Float).SetInt(bal)
	fDiv := new(big.Float).SetFloat64(1000000) // 10 ^ 6
	realBal := new(big.Float).Quo(fBal, fDiv)

	fmt.Printf("真实余额：%f USDT\n", realBal)
}
