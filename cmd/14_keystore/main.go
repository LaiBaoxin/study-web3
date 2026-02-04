package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
	"os"
)

func main() {
	// 定义 keystore 的存储路径
	ksDir := "./tmp/keystore"

	// 确保存在这个目录
	if err := os.MkdirAll(ksDir, 0755); err != nil {
		log.Fatal(err)
	}

	// 初始化 Geth 的 Keystore 管理器，[参数1: 路径, 参数2: 加强密度(StandardScryptN, LightScryptN) ]
	ks := keystore.NewKeyStore(ksDir, keystore.StandardScryptN, keystore.StandardScryptP)

	// 创建新的账户
	myPassward := "lbx123456@#$"

	// 1. 生成随机私钥
	// 2. 使用密码加密私钥
	// 3. 在 ksDir 下生成 UTC-xxx-xxx的 json 文件
	account, err := ks.NewAccount(myPassward)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("创建账户成功,新账户地址: %s  , \n 文件存储位置: %s \n", account.Address, account.URL.Path)

	// 读取并且解密keystore加密后的账户文件
	jsonBytes, err := os.ReadFile(account.URL.Path)
	if err != nil {
		log.Fatal(err)
	}
	// 解密 DecryptKey
	key, err := keystore.DecryptKey(jsonBytes, myPassward)
	if err != nil {
		log.Fatal("密码错误或者文件已经损坏了", err)
	}

	// 拿到私钥进行输出
	fmt.Printf("    解密成功!\n")
	fmt.Printf("    还原出的地址: %s\n", key.Address.Hex())
	privateKeyBytes := crypto.FromECDSA(key.PrivateKey)
	fmt.Printf("    还原出的私钥: %s\n", hexutil.Encode(privateKeyBytes))
}
