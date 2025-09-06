package handler

import (
	"hm-dianping-go/service"
	"hm-dianping-go/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetVoucherList 获取优惠券列表
func GetVoucherList(c *gin.Context) {
	shopIdStr := c.Param("shopId")
	if shopIdStr == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "商铺ID不能为空")
		return
	}

	shopId, err := strconv.ParseUint(shopIdStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的商铺ID")
		return
	}

	result := service.GetVoucherList(uint(shopId))
	utils.Response(c, result)
}

// AddVoucher 新增普通券
func AddVoucher(c *gin.Context) {
	// TODO: 实现新增普通券功能
	utils.ErrorResponse(c, http.StatusNotImplemented, "功能未完成")
}

// AddSeckillVoucher 新增秒杀券
func AddSeckillVoucher(c *gin.Context) {
	// TODO: 实现新增秒杀券功能
	utils.ErrorResponse(c, http.StatusNotImplemented, "功能未完成")
}