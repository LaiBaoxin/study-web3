package model

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"os"
)

type User struct {
	Address    common.Address     // 钱包地址
	PrivateKey *ecdsa.PrivateKey  // 私钥对象
	Auth       *bind.TransactOpts // 发送交易的凭证
}

// NewUserFromEnv 环境变量加载身份
func NewUserFromEnv(client *ethclient.Client) *User {
	// 获取私钥字符串
	privateKeyStr := os.Getenv("PRIVATE_KEY")
	if privateKeyStr == "" {
		log.Fatal("err: .env 文件中缺少私钥")
		return nil
	}

	// 解析私钥
	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		log.Fatal("err: 解析私钥失败,私钥格式有误 ", err)
		return nil
	}

	// 推导地址
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("err: 推导公钥失败")
		return nil
	}

	// 返回 User 对象
	return &User{
		Address:    crypto.PubkeyToAddress(*publicKeyECDSA),
		PrivateKey: privateKey,
	}
}
