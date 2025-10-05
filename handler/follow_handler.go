package handler

import (
	"hm-dianping-go/service"
	"hm-dianping-go/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Follow 关注用户
func Follow(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "用户未登录")
		return
	}

	followIdStr := c.Param("id")
	followId, err := strconv.ParseUint(followIdStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	result := service.Follow(c.Request.Context(), userID.(uint), uint(followId))
	utils.Response(c, result)
}

// Unfollow 取消关注
func Unfollow(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "用户未登录")
		return
	}

	followIdStr := c.Param("id")
	followId, err := strconv.ParseUint(followIdStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	result := service.Unfollow(c.Request.Context(), userID.(uint), uint(followId))
	utils.Response(c, result)
}

// GetCommonFollows 获取共同关注
func GetCommonFollows(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "用户未登录")
		return
	}

	targetIdStr := c.Param("id")
	targetId, err := strconv.ParseUint(targetIdStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	result := service.GetCommonFollows(c.Request.Context(), userID.(uint), uint(targetId))
	utils.Response(c, result)
}
