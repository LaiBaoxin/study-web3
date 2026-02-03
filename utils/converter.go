package utils

import (
	"github.com/ethereum/go-ethereum/params"
	"math/big"
	"regexp"
)

// EtherToWei 将 ETH 字符串转换为 Wei
func EtherToWei(eth string) *big.Int {
	// 字符串转换为 big.Float
	amountVal, _, err := big.ParseFloat(eth, 10, 0, big.ToNearestEven)
	if err != nil {
		return big.NewInt(0)
	}

	// 创建 1 Ether 的 Wei 值(10^18)
	etherVal := new(big.Float).SetInt(big.NewInt(params.Ether))

	// 计算 amount * 10 ^ 18
	weiVal := new(big.Float).Mul(amountVal, etherVal)

	// 转换回 big.Int
	result := new(big.Int)
	weiVal.Int(result) // Float 转换为 Int
	return result
}

// WeiToEther 将 Wei 转换为 ETH 字符串
func WeiToEther(wei *big.Int) string {
	// 创建 1 Ether 的 Wei 值(10^18)
	etherVal := new(big.Float).SetInt(big.NewInt(params.Ether))

	// 创建 wei 的 big.Float
	weiVal := new(big.Float).SetInt(wei)

	// 计算 wei / 10^18
	etherVal.Quo(weiVal, etherVal)

	// 转换为字符串
	return etherVal.Text('f', 18)
}

// IsValidAddress 检查地址是否合法
func IsValidAddress(address string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	return re.MatchString(address)
}
