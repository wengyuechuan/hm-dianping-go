package handler

import (
	"hm-dianping-go/service"
	"hm-dianping-go/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetShopById 根据ID获取商铺
func GetShopById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的商铺ID")
		return
	}

	result := service.GetShopById(uint(id))
	utils.Response(c, result)
}

// GetShopList 获取商铺列表
func GetShopList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	result := service.GetShopList(page, size)
	utils.Response(c, result)
}

// GetShopByType 根据类型获取商铺
func GetShopByType(c *gin.Context) {
	typeIdStr := c.Query("typeId")
	if typeIdStr == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "类型ID不能为空")
		return
	}

	typeId, err := strconv.ParseUint(typeIdStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的类型ID")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("current", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	result := service.GetShopByType(uint(typeId), page, size)
	utils.Response(c, result)
}

// GetShopByName 根据名称搜索商铺
func GetShopByName(c *gin.Context) {
	name := c.Query("name")
	page, _ := strconv.Atoi(c.DefaultQuery("current", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	result := service.GetShopByName(name, page, size)
	utils.Response(c, result)
}

// SaveShop 新增商铺
func SaveShop(c *gin.Context) {
	// TODO: 实现新增商铺功能
	utils.ErrorResponse(c, http.StatusNotImplemented, "功能未完成")
}

// UpdateShop 更新商铺信息
func UpdateShop(c *gin.Context) {
	// TODO: 实现更新商铺功能
	utils.ErrorResponse(c, http.StatusNotImplemented, "功能未完成")
}