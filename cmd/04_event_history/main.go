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
	// 加载配置
	godotenv.Load()

	// 连接节点
	client, err := ethclient.Dial(os.Getenv("RPC_URL"))
	if err != nil {
		log.Fatal("err: 连接节点失败", err)
	}

	// 实例化合约
	usdtAddress := common.HexToAddress(os.Getenv("USDT_CONTRACT_ADDR"))
	if usdtAddress.String() == "" {
		log.Fatal("err: 缺少 USDT 主网的合约地址")
	}
	usdt, err := erc20.NewERC20(usdtAddress, client)
	if err != nil {
		log.Fatal("err: 创建合约实例失败", err)
	}

	// 准备查询的范围，不仅需要查询最新的，还需要查询过去发生的
	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal("err: 获取最新的区块头失败", err)
	}

	currentBlack := header.Number.Uint64()

	// 设定查询范围：查询最近的 1000 个区块
	startBlock := uint64(0)
	if currentBlack > 1000 {
		startBlock = currentBlack - 1000
	}

	fmt.Printf("等待搜索区块范围: %d -> %d \n", startBlock, currentBlack)

	// 创建过滤条件
	filterOpts := &bind.FilterOpts{
		Start:   startBlock,
		End:     &currentBlack,
		Context: context.Background(),
	}

	// 设置过滤条件 FilterTransfer(opts, from[], to[])
	// 查询所有“从 account0 发出的”或者“给 account0”的转账
	myWallet := common.HexToAddress(os.Getenv("MY_WALLET_ADDR"))
	if myWallet.String() == "" {
		log.Fatal("err: 缺少账户0地址")
	}

	// 参数 1：过滤参数，参数 2：从来拿来（account0）的地址，参数 3：去往 account2 的地址
	iterator, err := usdt.FilterTransfer(filterOpts, []common.Address{myWallet}, nil)
	if err != nil {
		log.Fatal("err: 创建失败", err)
	}

	defer iterator.Close()

	// 遍历结果
	found := false
	for iterator.Next() {
		found = true
		event := iterator.Event // 获取当前事件对象

		// 输出发票日志
		fmt.Println("------------------------------------------------------------")
		fmt.Printf("发现转账小票\n")
		fmt.Printf("交易哈希: %s\n", event.Raw.TxHash.Hex())
		fmt.Printf("区块高度: %d\n", event.Raw.BlockNumber)
		fmt.Printf("发送方 (From): %s\n", event.From.Hex())
		fmt.Printf("接收方 (To):   %s\n", event.To.Hex())

		// 格式化金额
		fVal := new(big.Float).SetInt(event.Value)
		fVal.Quo(fVal, big.NewFloat(100000)) // 除以精度 6
		fmt.Printf("金额：     %.6f USDT \n", fVal)
	}

	if !found {
		fmt.Println("没有找到任何转账小票")
	}

}
