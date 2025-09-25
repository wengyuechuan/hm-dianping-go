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
	var req service.AddSeckillVoucherRequest

	// 绑定JSON数据到请求结构体
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	// 调用service层处理业务逻辑
	result := service.AddSeckillVoucher(c.Request.Context(), &req)
	utils.Response(c, result)
}

// GetSeckillVoucher 获取秒杀券详情
func GetSeckillVoucher(c *gin.Context) {
	voucherIdStr := c.Param("id")
	if voucherIdStr == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "优惠券ID不能为空")
		return
	}

	voucherId, err := strconv.ParseUint(voucherIdStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的优惠券ID")
		return
	}

	result := service.GetSeckillVoucher(uint(voucherId))
	utils.Response(c, result)
}
