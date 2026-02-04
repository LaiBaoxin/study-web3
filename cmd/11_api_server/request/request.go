package request

type TransferRequest struct {
	ToAddress string  `json:"toAddress" binding:"required"`
	Amount    float64 `json:"amount" binding:"required"`
}
