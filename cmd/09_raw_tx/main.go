package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	"learn-web3-go/contracts/erc20"
	"log"
	"math/big"
	"os"
	"strings"
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

	// 创建链接
	client, err := ethclient.Dial(os.Getenv("RPC_URL"))
	if err != nil {
		log.Fatal(err)
	}
	// account1 的私钥地址
	privateKey, _ := crypto.HexToECDSA(os.Getenv("PRIVATE_KEY"))
	// 推导公钥
	myAddress := crypto.PubkeyToAddress(privateKey.PublicKey)

	usdtAddress := common.HexToAddress(os.Getenv("USDT_ADDRESS"))
	toAddress := common.HexToAddress(os.Getenv("TO_WALLET_ADDR"))

	// 获取交易基础信息（Nonce & ChainID）
	// 获取 ChainID
	chainId, _ := client.ChainID(ctx)

	// 获取 Nonce: PendingNonceAt 会查询链上“我的下一笔交易编号应该是多少”
	nonce, err := client.PendingNonceAt(ctx, myAddress)
	if err != nil {
		log.Fatal("获取 Nonce 失败 ", err)
		return
	}
	log.Printf("Nonce 账号: %d \n", nonce)

	// 估算一下 gas 费: EIP-1559 费率 = BaseFree(基础费) + Tip(小费)

	// 获取建议的小费(Tip / PriorityFree)
	gasTipCap, err := client.SuggestGasTipCap(ctx)
	if err != nil {
		log.Fatal("获取建议的小费失败 ", err)
		return
	}

	// 获取区块基础费(BaseFree)
	head, _ := client.HeaderByNumber(ctx, nil)
	baseFee := head.BaseFee

	// 计算费用的上限（GasFeeCap）
	// 公式: GasFeeCap = BaseFee * 2 + Tip
	gasFeeCap := new(big.Int).Add(
		new(big.Int).Mul(baseFee, big.NewInt(2)),
		gasTipCap,
	)

	fmt.Printf("预估油费: %s , 小费: %s \n", baseFee, gasTipCap)

	// 数据打包（Pack Data）
	parsedABI, _ := abi.JSON(strings.NewReader(erc20.ERC20MetaData.ABI))

	// 转账金额: 1 USDT 10^6
	amount := big.NewInt(1000000)

	// pack 生成二进制
	data, err := parsedABI.Pack("transfer", toAddress, amount)
	if err != nil {
		log.Fatal("pack 失败 ", err)
		return
	}

	// 组装交易结构体

	// 估算 gas limit
	msg := ethereum.CallMsg{
		From: myAddress,
		To:   &usdtAddress,
		Gas:  0,
		Data: data,
	}
	gasLimit, _ := client.EstimateGas(ctx, msg)

	// 创建动态的费率交易(DynamicFeeTx)
	txData := &types.DynamicFeeTx{
		ChainID:   chainId,
		Nonce:     nonce,
		GasTipCap: gasTipCap, // 愿意支付矿工小费
		GasFeeCap: gasFeeCap, // 愿意支付更高的费用
		Gas:       gasLimit,
		To:        &usdtAddress,
		Value:     big.NewInt(0), // 0 ETH
		Data:      data,
	}

	// 包装成 Transaction 对象
	tx := types.NewTx(txData)

	// 进行签名并且广播

	// 私钥进行签名
	signer := types.LatestSignerForChainID(chainId)
	signerTx, err := types.SignTx(tx, signer, privateKey)
	if err != nil {
		log.Fatal("签名失败 ", err)
		return
	}

	// 广播发出
	err = client.SendTransaction(ctx, signerTx)
	if err != nil {
		log.Fatal("广播失败 ", err)
		return
	}

	log.Printf("交易成功，Hash: %s \n", signerTx.Hash().Hex())
}
