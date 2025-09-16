package handler

import (
	"hm-dianping-go/service"
	"hm-dianping-go/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// SeckillVoucher 秒杀优惠券
func SeckillVoucher(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "用户未登录")
		return
	}

	voucherIdStr := c.Param("id")
	voucherId, err := strconv.ParseUint(voucherIdStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的优惠券ID")
		return
	}

	ctx := c.Request.Context()
	result := service.SeckillVoucher(ctx, userID.(uint), uint(voucherId))
	utils.Response(c, result)
}
