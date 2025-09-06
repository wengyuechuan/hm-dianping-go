package service

import (
	"hm-dianping-go/dao"
	"hm-dianping-go/models"
	"hm-dianping-go/utils"
)

// GetShopById 根据ID获取商铺
func GetShopById(id uint) *utils.Result {
	var shop models.Shop
	err := dao.DB.First(&shop, id).Error
	if err != nil {
		return utils.ErrorResult("商铺不存在")
	}

	return utils.SuccessResultWithData(shop)
}

// GetShopList 获取商铺列表
func GetShopList(page, size int) *utils.Result {
	var shops []models.Shop
	var total int64

	offset := (page - 1) * size

	// 获取总数
	dao.DB.Model(&models.Shop{}).Count(&total)

	// 分页查询
	err := dao.DB.Offset(offset).Limit(size).Find(&shops).Error
	if err != nil {
		return utils.ErrorResult("查询失败")
	}

	return utils.SuccessResultWithData(map[string]interface{}{
		"list":  shops,
		"total": total,
		"page":  page,
		"size":  size,
	})
}

// GetShopByType 根据类型获取商铺
func GetShopByType(typeId uint, page, size int) *utils.Result {
	var shops []models.Shop
	var total int64

	offset := (page - 1) * size

	// 获取总数
	dao.DB.Model(&models.Shop{}).Where("type_id = ?", typeId).Count(&total)

	// 分页查询
	err := dao.DB.Where("type_id = ?", typeId).Offset(offset).Limit(size).Find(&shops).Error
	if err != nil {
		return utils.ErrorResult("查询失败")
	}

	return utils.SuccessResultWithData(map[string]interface{}{
		"list":  shops,
		"total": total,
		"page":  page,
		"size":  size,
	})
}

// GetShopByName 根据名称搜索商铺
func GetShopByName(name string, page, size int) *utils.Result {
	var shops []models.Shop
	var total int64

	offset := (page - 1) * size

	query := dao.DB.Model(&models.Shop{})
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	// 获取总数
	query.Count(&total)

	// 分页查询
	err := query.Offset(offset).Limit(size).Find(&shops).Error
	if err != nil {
		return utils.ErrorResult("查询失败")
	}

	return utils.SuccessResultWithData(map[string]interface{}{
		"list":  shops,
		"total": total,
		"page":  page,
		"size":  size,
	})
}