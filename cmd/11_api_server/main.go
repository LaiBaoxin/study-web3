package main

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	"learn-web3-go/cmd/11_api_server/request"
	"learn-web3-go/cmd/11_api_server/response"
	"learn-web3-go/contracts/erc20"
	"learn-web3-go/pkg/chain"
	"learn-web3-go/pkg/chain/model"
	"log"
	"math/big"
	"net/http"
	"os"
)

var (
	client    *ethclient.Client
	adminUser *model.User // 服务器的“热钱包”账户
	usdt      *erc20.ERC20
)

func main() {
	// 初始化连接
	client = chain.InitClient()

	// 加载管理员账户(使用 account1)
	adminUser = model.NewUserFromEnv(client)
	log.Printf("热钱包地址:%s", adminUser.Address.Hex())

	// 初始化 USDT
	usdtAddr := common.HexToAddress(os.Getenv("USDT_CONTRACT_ADDR"))
	usdt, err := erc20.NewERC20(usdtAddr, client)
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	// 解决跨域
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // 允许任何来源
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// 注册路由
	// 获取 RPC_URL
	r.GET("/getRpcUrl", func(c *gin.Context) {
		response.Success(c, gin.H{
			"rpcUrl": os.Getenv("RPC_URL"),
		}, "获取成功")
	})

	// 获取钱包余额
	r.GET("/balance", func(c *gin.Context) {
		addressStr := c.Query("address")
		if !common.IsHexAddress(addressStr) {
			response.Fail(c, http.StatusBadRequest, "无效的参数")
			return
		}

		// 调用链上合约
		targetAddr := common.HexToAddress(addressStr)
		bal, err := usdt.BalanceOf(nil, targetAddr)
		if err != nil {
			response.Fail(c, http.StatusInternalServerError, "查询余额失败")
			return
		}

		// 格式化显示余额（Wei -> Human Readable）
		humanBal := new(big.Float).Quo(new(big.Float).SetInt(bal), big.NewFloat(1e6))
		val, _ := humanBal.Float64() // 转成浮点数

		response.Success(c, gin.H{
			"balance": val,
			"address": addressStr,
			"symbol":  "USDT",
		}, "查询成功")
	})

	// 提现/转账
	r.POST("/transfer", func(c *gin.Context) {
		var req request.TransferRequest

		err = c.ShouldBindJSON(&req)
		if err != nil {
			response.Fail(c, http.StatusBadRequest, "无效的参数")
			return
		}
		if req.Amount <= 0 {
			response.Fail(c, http.StatusBadRequest, "交易金额需要大于 0")
			return
		}
		if !common.IsHexAddress(req.ToAddress) {
			response.Fail(c, http.StatusBadRequest, "交易接收的地址无效")
			return
		}

		// 生成交易凭证 auth - 自动计算 gas 费用
		auth, err := chain.NewAuth(client, adminUser)
		if err != nil {
			response.Fail(c, http.StatusInternalServerError, "生成签名失败")
			return
		}

		// 数值转换 (Human -> Wei) USDT 精度 6: amount * 10^6
		amountBig := big.NewInt(int64(req.Amount * 1000000))

		// 发起交易
		toAddress := common.HexToAddress(req.ToAddress)
		// 开始转账
		log.Println("正在广播交易中....")
		tx, err := usdt.Transfer(auth, toAddress, amountBig)
		if err != nil {
			response.Fail(c, http.StatusInternalServerError, "交易广播失败")
			log.Fatal("交易广播失败", err.Error())
			return
		}
		response.Success(c, gin.H{
			"txHash": tx.Hash().Hex(),
		}, "交易已广播，等待上链...")
	})

	r.Run(":8888")

}
