package main

import (
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"

	"learn-web3-go/contracts/erc20"
	"learn-web3-go/contracts/multicall"
)

func main() {
	// 1. 加载环境变量
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error: 找不到 .env 文件")
	}

	// 2. 连接节点
	client, err := ethclient.Dial(os.Getenv("RPC_URL"))
	if err != nil {
		log.Fatal("连接节点失败:", err)
	}

	// 3. 准备地址
	// Multicall3
	mcAddr := common.HexToAddress("0xcA11bde05977b3631167028862bE2a173976CA11")
	// USDT 合约
	usdtAddr := common.HexToAddress(os.Getenv("USDT_CONTRACT_ADDR"))
	// 账户 1
	myAddr := common.HexToAddress(os.Getenv("MY_WALLET_ADDR"))
	// 账户 2
	otherAddr := common.HexToAddress(os.Getenv("TO_WALLET_ADDR"))

	// 准备翻译器 (Parse ABI)
	// 使用 ERC20 的 ABI 来打包/解包数据
	erc20ABI, err := abi.JSON(strings.NewReader(erc20.ERC20MetaData.ABI))
	if err != nil {
		log.Fatal("解析 ERC20 ABI 失败:", err)
	}

	// 写小纸条 (Pack CallData)

	// Task 1: 查我的余额
	callData1, err := erc20ABI.Pack("balanceOf", myAddr)
	if err != nil {
		log.Fatal(err)
	}

	// Task 2: 查总供应量
	callData2, err := erc20ABI.Pack("totalSupply")
	if err != nil {
		log.Fatal(err)
	}

	// Task 3: 查朋友的余额
	callData3, err := erc20ABI.Pack("balanceOf", otherAddr)
	if err != nil {
		log.Fatal(err)
	}

	// 将纸条放入购物车结构体
	// 构建批量执行 3 个任务
	calls := []multicall.Struct0{
		{Target: usdtAddr, CallData: callData1},
		{Target: usdtAddr, CallData: callData2},
		{Target: usdtAddr, CallData: callData3},
	}

	mcInstance, err := multicall.NewMulticallCaller(mcAddr, client)
	if err != nil {
		log.Fatal("实例化 Multicall 失败:", err)
	}

	// 发送请求
	aggregateResult, err := mcInstance.Aggregate(nil, calls)
	if err != nil {
		return
	}
	results := aggregateResult.ReturnData

	// 拆快递 (Unpack Result)

	// --- 1. 解包我的余额 ---
	out1, err := erc20ABI.Unpack("balanceOf", results[0])
	if err != nil {
		log.Fatal(err)
	}
	bal1 := out1[0].(*big.Int)
	fmt.Printf("account1 的余额:   %s Wei\n", bal1.String())

	// --- 2. 解包总供应量 ---
	out2, err := erc20ABI.Unpack("totalSupply", results[1])
	if err != nil {
		log.Fatal(err)
	}
	supply := out2[0].(*big.Int)
	// 简单除以 10^6 显示一点
	humanSupply := new(big.Float).Quo(new(big.Float).SetInt(supply), big.NewFloat(1e6))
	fmt.Printf("总供应量:   %.2f USDT\n", humanSupply)

	// --- 3. 解包朋友余额 ---
	out3, err := erc20ABI.Unpack("balanceOf", results[2])
	if err != nil {
		log.Fatal(err)
	}
	bal2 := out3[0].(*big.Int)
	fmt.Printf("account2 的余额:   %s Wei\n", bal2.String())

}
