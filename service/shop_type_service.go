package service

import (
	"hm-dianping-go/dao"
	"hm-dianping-go/models"
	"hm-dianping-go/utils"
)

// GetShopTypeList 获取商铺类型列表
func GetShopTypeList() *utils.Result {
	var shopTypes []models.ShopType
	err := dao.DB.Order("sort asc").Find(&shopTypes).Error
	if err != nil {
		return utils.ErrorResult("查询失败")
	}

	return utils.SuccessResultWithData(shopTypes)
}