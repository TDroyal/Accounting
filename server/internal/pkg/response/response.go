package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 错误码规范（见 docs/03 §7）：
//   0          成功
//   40000-40099 通用参数/校验错误
//   40100-40199 鉴权错误
//   40400-40499 资源不存在
//   50000-50099 服务器内部错误
const (
	CodeOK            = 0
	CodeErrParam      = 40001 // 参数校验失败
	CodeErrAmount     = 40002 // 金额非法
	CodeErrUnauthorized = 40101 // 未登录 / token 失效
	CodeErrBadCredential = 40102 // 账号或密码错误
	CodeErrForbidden  = 40301 // 无权限
	CodeErrNotFound   = 40401 // 资源不存在
	CodeErrConflict   = 40901 // 资源冲突
	CodeErrRateLimited = 42901 // 请求过于频繁
	CodeErrInternal   = 50000 // 服务器内部错误
)

// code 到 HTTP 状态的映射（docs/03 §6：HTTP 仍按语义使用）
var codeToHTTP = map[int]int{
	CodeOK:             http.StatusOK,
	CodeErrParam:       http.StatusBadRequest,
	CodeErrAmount:      http.StatusBadRequest,
	CodeErrUnauthorized: http.StatusUnauthorized,
	CodeErrBadCredential: http.StatusUnauthorized,
	CodeErrForbidden:   http.StatusForbidden,
	CodeErrNotFound:    http.StatusNotFound,
	CodeErrConflict:    http.StatusConflict,
	CodeErrRateLimited: http.StatusTooManyRequests,
	CodeErrInternal:    http.StatusInternalServerError,
}

// Body 统一响应结构 {code, message, data}
type Body struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// httpStatus 返回错误码对应的 HTTP 状态，默认 500。
func httpStatus(code int) int {
	if s, ok := codeToHTTP[code]; ok {
		return s
	}
	return http.StatusInternalServerError
}

// OK 成功响应，data 可为 nil。
func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Body{Code: CodeOK, Message: "ok", Data: data})
}

// Fail 失败响应，按 code 映射 HTTP 状态。
func Fail(c *gin.Context, code int, msg string) {
	c.JSON(httpStatus(code), Body{Code: code, Message: msg, Data: nil})
}
