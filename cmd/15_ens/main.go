package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	"github.com/wealdtech/go-ens/v3"
	"log"
	"os"
)

func main() {
	_ = godotenv.Load()

	rpcUrl := os.Getenv("RPC_MAINNET_URL")
	if rpcUrl == "" {
		log.Fatal("RPC_MAINNET_URL not found")
		return
	}

	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		log.Fatal(err)
	}

	// 正向解析 Name -> Address
	domain := "vitalik.eth"
	fmt.Printf("正在查询域名:%s  \n", domain)

	address, err := ens.Resolve(client, domain)
	if err != nil {
		log.Fatal(err)
		return
	} else {
		fmt.Printf("域名解析成功,域名:%s,地址:%s \n", domain, address)
	}

	// 反向解析 Address -> Name
	targetAddr := common.HexToAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045") // V神公用地址
	fmt.Printf("正在反查地址: %s \n", targetAddr.Hex())

	reverseName, err := ens.ReverseResolve(client, targetAddr)
	if err != nil {
		log.Fatal(err)
		return
	} else {
		fmt.Printf("地址反查成功,地址:%s \n ,域名:%s \n", targetAddr.Hex(), reverseName)
	}
}
