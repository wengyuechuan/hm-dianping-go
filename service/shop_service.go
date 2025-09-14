package service

import (
	"context"
	"hm-dianping-go/dao"
	"hm-dianping-go/models"
	"hm-dianping-go/utils"
	"log"
)

// GetShopById 根据ID获取商铺
func GetShopById(ctx context.Context, id uint) *utils.Result {
	// 1. 布隆过滤器检查，防止缓存击穿
	flag, err := utils.CheckIDExistsWithRedis(ctx, dao.Redis, "shop", id)
	if err != nil {
		log.Fatalf("检查布隆过滤器失败: %v", err)
	}
	if !flag {
		// 布隆过滤器判断商铺不存在，直接返回
		return utils.ErrorResult("商铺不存在")
	}

	// 2. 先从缓存查询
	shop, err := dao.GetShopCacheById(ctx, dao.Redis, id)
	if err == nil && shop != nil {
		// 缓存命中，直接返回
		return utils.SuccessResultWithData(shop)
	}

	// 3. 缓存未命中或缓存查询出错，查询数据库
	shop, err = dao.GetShopById(ctx, dao.DB, id)
	if err != nil {
		// 数据库查询失败
		return utils.ErrorResult("查询失败: " + err.Error())
	}

	// 4. 缓存未命中，设置缓存
	err = dao.SetShopCacheById(ctx, dao.Redis, id, shop)
	if err != nil {
		return utils.ErrorResult("查询失败: " + err.Error())
	}

	// 5. 数据库查询成功，返回结果（DAO层已处理缓存写入）
	return utils.SuccessResultWithData(shop)
}

// UpdateShopById 根据ID更新商铺
func UpdateShopById(ctx context.Context, shop *models.Shop) *utils.Result {

	// 0. 启动事务
	tx := dao.DB.Begin()
	defer func() { // 捕获异常
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. 更新数据库
	err := dao.UpdateShop(ctx, tx, shop)

	// 2. 更新失败
	if err != nil {
		tx.Rollback()
		return utils.ErrorResult("更新失败: " + err.Error())
	}

	// 3. 删除缓存
	err = dao.DelShopCacheById(ctx, dao.Redis, shop.ID)
	if err != nil {
		tx.Rollback()
		return utils.ErrorResult("更新失败: " + err.Error())
	}

	// 4. 提交事务
	tx.Commit()

	// 5. 返回结果
	return utils.SuccessResult("更新成功")
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
