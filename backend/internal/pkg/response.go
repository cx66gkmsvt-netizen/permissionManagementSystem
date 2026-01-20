package pkg

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

const (
	SUCCESS       = 200
	ERROR         = 500
	UNAUTHORIZED  = 401
	FORBIDDEN     = 403
	NOT_FOUND     = 404
	INVALID_PARAM = 400
)

// OK 成功响应
func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: SUCCESS,
		Msg:  "success",
		Data: data,
	})
}

// OKMsg 成功响应带消息
func OKMsg(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: SUCCESS,
		Msg:  msg,
	})
}

// Fail 失败响应
func Fail(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: ERROR,
		Msg:  msg,
	})
}

// FailCode 失败响应带状态码
func FailCode(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  msg,
	})
}

// Unauthorized 未授权
func Unauthorized(c *gin.Context, msg string) {
	c.JSON(http.StatusUnauthorized, Response{
		Code: UNAUTHORIZED,
		Msg:  msg,
	})
	c.Abort()
}

// Forbidden 禁止访问
func Forbidden(c *gin.Context, msg string) {
	c.JSON(http.StatusForbidden, Response{
		Code: FORBIDDEN,
		Msg:  msg,
	})
	c.Abort()
}
