package handler

import (
	"hm-dianping-go/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetDailyUV 获取指定日期的UV统计
func GetDailyUV(c *gin.Context) {
	date := c.Query("date")
	if date == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请提供日期参数",
		})
		return
	}
	
	result := service.GetDailyUV(c.Request.Context(), date)
	if result.Success {
		c.JSON(http.StatusOK, result)
	} else {
		c.JSON(http.StatusBadRequest, result)
	}
}

// GetTodayUV 获取今日UV统计
func GetTodayUV(c *gin.Context) {
	result := service.GetTodayUV(c.Request.Context())
	if result.Success {
		c.JSON(http.StatusOK, result)
	} else {
		c.JSON(http.StatusInternalServerError, result)
	}
}

// GetUVRange 获取指定日期范围的UV统计
func GetUVRange(c *gin.Context) {
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")
	
	if startDate == "" || endDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请提供开始日期和结束日期参数",
		})
		return
	}
	
	result := service.GetUVRange(c.Request.Context(), startDate, endDate)
	if result.Success {
		c.JSON(http.StatusOK, result)
	} else {
		c.JSON(http.StatusBadRequest, result)
	}
}

// GetRecentUV 获取最近N天的UV统计
func GetRecentUV(c *gin.Context) {
	daysStr := c.Query("days")
	if daysStr == "" {
		daysStr = "7" // 默认7天
	}
	
	days, err := strconv.Atoi(daysStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "天数参数格式错误",
		})
		return
	}
	
	result := service.GetRecentUV(c.Request.Context(), days)
	if result.Success {
		c.JSON(http.StatusOK, result)
	} else {
		c.JSON(http.StatusBadRequest, result)
	}
}

// GetUVSummary 获取UV统计摘要
func GetUVSummary(c *gin.Context) {
	result := service.GetUVSummary(c.Request.Context())
	if result.Success {
		c.JSON(http.StatusOK, result)
	} else {
		c.JSON(http.StatusInternalServerError, result)
	}
}