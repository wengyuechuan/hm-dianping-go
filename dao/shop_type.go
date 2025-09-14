package dao

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"

	"hm-dianping-go/models"
)

func GetShopTypeList(ctx context.Context, db *gorm.DB) ([]*models.ShopType, error) {
	var shopTypes []*models.ShopType
	err := db.WithContext(ctx).Model(&models.ShopType{}).Order("sort").Find(&shopTypes).Error
	if err != nil {
		return nil, err
	}
	return shopTypes, nil
}

// ===========缓存相关=============

const (
	ShopTypeCache = "cache:shop_type"
)

func GetShopTypeListCache(ctx context.Context, rds *redis.Client) ([]*models.ShopType, error) {
	key := ShopTypeCache
	result := rds.Get(ctx, key)
	if result.Err() != nil {
		return nil, result.Err()
	}

	var shopTypes []*models.ShopType
	if err := result.Scan(&shopTypes); err != nil {
		return nil, err
	}
	return shopTypes, nil
}

func SetShopTypeListCache(ctx context.Context, rds *redis.Client, shopTypes []*models.ShopType) error {
	// 设置一小时的过期时间
	err := rds.Set(ctx, ShopTypeCache, shopTypes, time.Hour).Err()
	if err != nil {
		return err
	}
	return nil
}
