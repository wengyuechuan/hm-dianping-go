package service

import (
	"context"
	"hm-dianping-go/dao"
	"hm-dianping-go/utils"
)

// GetShopTypeList 获取商铺类型列表
func GetShopTypeList(ctx context.Context) *utils.Result {
	// 从缓存中获取
	shopTypes, err := dao.GetShopTypeListCache(ctx, dao.Redis)
	if err == nil {
		return utils.SuccessResultWithData(shopTypes)
	}

	// 从数据库中获取
	shopTypes, err = dao.GetShopTypeList(ctx, dao.DB)
	if err != nil {
		return utils.ErrorResult("查询失败")
	}
	// 缓存到redis
	err = dao.SetShopTypeListCache(ctx, dao.Redis, shopTypes)
	if err != nil {
		return utils.ErrorResult("查询失败")
	}

	return utils.SuccessResultWithData(shopTypes)
}
