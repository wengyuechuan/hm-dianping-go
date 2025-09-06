package handler

import (
	"hm-dianping-go/service"
	"hm-dianping-go/utils"

	"github.com/gin-gonic/gin"
)

// GetShopTypeList 获取商铺类型列表
func GetShopTypeList(c *gin.Context) {
	result := service.GetShopTypeList()
	utils.Response(c, result)
}