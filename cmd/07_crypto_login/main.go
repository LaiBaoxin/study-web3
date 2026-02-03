package main

import (
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("err: 找不到 .env 文件")
		return
	}
	// 获取私钥
	privateKey, err := crypto.HexToECDSA(os.Getenv("PRIVATE_KEY"))
	if err != nil {
		log.Fatal("err: 私钥转换错误")
		return
	}

	// 推导公钥地址
	userAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	fmt.Println("user address: ", userAddress.Hex())

	// 构造消息
	message := "welcome to study web3 with golang"

	prefix := fmt.Sprintf("\x19Ethereum Signed Message: \n %d:%s", len(message), message)

	// 计算哈希 (Keccak256)
	hash := crypto.Keccak256Hash([]byte(prefix))
	fmt.Printf("EIP-191 格式的 hash：%s", hash.Hex())

	// 生成签名
	// 使用私钥对这个哈希进行签名，得到 signature 是[]byte 字节数组
	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	signature[64] += 27 // crypto 库生成的是 0 或 1, 需要手动+27

	// 转成 hex 字节显示
	fmt.Printf("生成签名：0x%s\n", hex.EncodeToString(signature))

	// server 收到签名进行验证
	fmt.Println("服务端正在验证签名中，请稍候...")
	time.Sleep(3 * time.Second) // 模拟服务器验签的时间
	if signature[64] >= 27 {
		signature[64] -= 27
	}

	// 从签名中回复公钥（Recover）
	sigPublicKey, err := crypto.SigToPub(hash.Bytes(), signature)
	if err != nil {
		log.Fatal("err: 签名恢复公钥失败")
		return
	}

	// 算出地址
	recoveredAddress := crypto.PubkeyToAddress(*sigPublicKey)
	fmt.Println("后端回复出来的地址：", recoveredAddress.Hex())

	// 验证地址是否一致
	if recoveredAddress == userAddress {
		fmt.Println("签名验证通过")
	} else {
		fmt.Println("签名验证失败")
	}
}
