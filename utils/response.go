package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Result 统一响应结构
type Result struct {
	Success bool        `json:"success"`
	ErrorMsg string     `json:"errorMsg,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Total   int64       `json:"total,omitempty"`
}

// SuccessResult 成功响应
func SuccessResult(message string) *Result {
	return &Result{
		Success: true,
		Data:    message,
	}
}

// SuccessResultWithData 成功响应带数据
func SuccessResultWithData(data interface{}) *Result {
	return &Result{
		Success: true,
		Data:    data,
	}
}

// ErrorResult 错误响应
func ErrorResult(errorMsg string) *Result {
	return &Result{
		Success:  false,
		ErrorMsg: errorMsg,
	}
}

// Response 统一响应处理
func Response(c *gin.Context, result *Result) {
	if result.Success {
		c.JSON(http.StatusOK, result)
	} else {
		c.JSON(http.StatusOK, result)
	}
}

// SuccessResponse 成功响应
func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Result{
		Success: true,
		Data:    data,
	})
}

// ErrorResponse 错误响应
func ErrorResponse(c *gin.Context, code int, message string) {
	c.JSON(code, Result{
		Success:  false,
		ErrorMsg: message,
	})
}

// PageResult 分页响应
func PageResult(c *gin.Context, data interface{}, total int64, page, size int) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"list":  data,
			"total": total,
			"page":  page,
			"size":  size,
		},
	})
}