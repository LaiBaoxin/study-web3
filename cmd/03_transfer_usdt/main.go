package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	"learn-web3-go/contracts/erc20"
	"log"
	"math/big"
	"os"
)

func main() {
	// load 加载配置
	if err := godotenv.Load(); err != nil {
		log.Fatal("err: 加载配置文件 .env 失败")
		return
	}

	// 连接测试链节点
	client, err := ethclient.Dial(os.Getenv("RPC_URL"))
	if err != nil {
		log.Fatal("err: 连接测试链节点失败")
		return
	}

	// 获取私钥
	privateKeyStr := os.Getenv("PRIVATE_KEY")
	if privateKeyStr == "" {
		log.Fatal("err: 缺少私钥")
		return
	}
	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		log.Fatal("err: 私钥格式失败")
		return
	}

	// 获取 ChainID
	chainId, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatal("err: 获取 ChainID 失败")
		return
	}

	// 创建交易发起人
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainId)
	if err != nil {
		log.Fatal("err: 创建交易发起人失败")
		return
	}

	// 实例化 USDT 合约
	usdtAddress := common.HexToAddress(os.Getenv("USDT_CONTRACT_ADDR"))
	if usdtAddress.String() == "" {
		log.Fatal("err: 缺少 USDT 主网的合约地址")
		return
	}
	usdt, err := erc20.NewERC20(usdtAddress, client)
	if err != nil {
		log.Fatal("err: 创建 USDT 合约实例失败")
		return
	}

	// 准备转账的参数(账户2)
	toAddress := common.HexToAddress(os.Getenv("TO_WALLET_ADDR"))
	if toAddress.String() == "" {
		log.Fatal("err: 缺少账户2地址")
		return
	}

	// 转账 10 USDT (精度是 6，所以是  10^6 )
	amount := new(big.Int).Mul(big.NewInt(10), big.NewInt(1000000))
	fmt.Printf("正在准备转账,准备从 %s 转账 10 USDT 给 %s ... \n", auth.From.Hex(), toAddress.Hex())

	// 发起交易
	tx, err := usdt.Transfer(auth, toAddress, amount)
	if err != nil {
		log.Fatal("err: 发起交易失败 ", err)
		return
	}

	// 打印交易的哈希
	fmt.Printf("交易哈希(TX Hash): %s \n", tx.Hash().Hex())
}
