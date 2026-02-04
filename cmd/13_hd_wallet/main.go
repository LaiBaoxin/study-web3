package main

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip39"
	"log"
)

func main() {
	// 生成助记词 (Mnemonic)
	// 128 位随机数 -> 12 个单词
	// 256 位随机数 -> 24 个单词
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		log.Fatal(err)
	}

	// 将随机数查表转换成单词
	mnemonic, _ := bip39.NewMnemonic(entropy)
	// 助记词 (Mnemonic) -> 种子 (Seed)
	seed := bip39.NewSeed(mnemonic, "")
	// 种子 (Seed) -> 私钥 (Private Key)
	// seed 是 64 字节， 私钥只需要 32 字节
	privateKey, err := crypto.ToECDSA(seed[:32])
	if err != nil {
		log.Fatal(err)
	}

	// 私钥 (Private Key) -> 公钥 (Public Key) -> 地址 (Address)
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("err: 推导公钥失败")
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA)

	// 去掉 "0x" 前缀
	privHex := hexutil.Encode(crypto.FromECDSA(privateKey))[2:]

	// 输出
	fmt.Println("最终的钱包详情如下:")
	fmt.Printf("   私钥: %s \n", privHex)
	fmt.Printf("   地址: %s \n", address.Hex())
	fmt.Printf("   助记词: %s \n", mnemonic)
	fmt.Printf("   种子: %x \n", seed)
	fmt.Println("--------------------")

}
