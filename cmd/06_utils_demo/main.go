package main

import (
	"fmt"
	"learn-web3-go/utils"
	"math/big"
)

func main() {
	fmt.Println(" ------- Utils Test -------")

	// ----- ETH -> Wei -----
	inputEth := "0.05"
	wei := utils.EtherToWei(inputEth)
	fmt.Printf("输入: %s ETH\n", inputEth)
	fmt.Printf("输出: %s Wei\n", wei.String())
	fmt.Println("--------------------")

	// ----- Wei -> ETH ------
	hugeWei, _ := new(big.Int).SetString("12345600000000000000", 10)

	eth := utils.WeiToEther(hugeWei)
	fmt.Printf("输入: %s Wei\n", hugeWei.String())
	fmt.Printf("输出: %s ETH\n", eth) // 应该去除多余的0
	fmt.Println("--------------------")

	// ----- 地址校验 -----
	fmt.Println("--------- use utils -----------")
	addr1 := "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266" // 正确
	addr2 := "0xBadAddress"                               // 错误

	fmt.Printf("地址1 (%s) 是否有效吗? %v\n", addr1, utils.IsValidAddress(addr1))
	fmt.Printf("地址2 (%s) 是否有效吗? %v\n", addr2, utils.IsValidAddress(addr2))

}
