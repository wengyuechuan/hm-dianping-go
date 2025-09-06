package handler

import (
	"hm-dianping-go/utils"

	"github.com/gin-gonic/gin"
)

// HealthCheck 健康检查
func HealthCheck(c *gin.Context) {
	utils.SuccessResponse(c, gin.H{
		"status": "ok",
		"message": "服务运行正常",
	})
}