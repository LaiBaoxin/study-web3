package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Success(c *gin.Context, data interface{}, msg string) {
	c.JSON(http.StatusOK, &Response{
		Code:    200,
		Message: msg,
		Data:    data,
	})
}

func Fail(c *gin.Context, httpStatus int, msg string) {
	c.JSON(httpStatus, &Response{
		Code:    httpStatus,
		Message: msg,
		Data:    nil,
	})
}
