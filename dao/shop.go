package dao

import (
	"context"
	"hm-dianping-go/models"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

const (
	SHOP_CHACHE = "cache:shop:"
)

func GetShopById(ctx context.Context, db *gorm.DB, rds *redis.Client, shopId uint) (*models.Shop, error) {
	// 从缓存中获取店铺信息
	shop := &models.Shop{}
	result := rds.Get(ctx, SHOP_CHACHE+strconv.Itoa(int(shopId)))
	if result.Err() == nil { // 此时没有redis错误
		err := result.Scan(shop)
		if err == nil { // 此时缓存正确命中
			return shop, nil
		}
	}
	err := db.Where("id = ?", shopId).First(shop).Error
	if err != nil {
		return nil, err
	}
	// 缓存店铺信息
	err = rds.Set(ctx, SHOP_CHACHE+strconv.Itoa(int(shopId)), shop, time.Hour).Err()
	if err != nil {
		return nil, err
	}
	return shop, nil
}
