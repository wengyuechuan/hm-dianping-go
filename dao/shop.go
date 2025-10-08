package dao

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"hm-dianping-go/models"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

func GetShopById(ctx context.Context, db *gorm.DB, shopId uint) (*models.Shop, error) {
	// 从缓存中获取店铺信息
	shop := &models.Shop{}
	err := db.Where("id = ?", shopId).First(shop).Error
	if err != nil {
		return nil, err
	}

	return shop, nil
}

// GetAllShopIDs 获取所有商铺ID
func GetAllShopIDs(ctx context.Context, db *gorm.DB) ([]uint, error) {
	var ids []uint
	err := db.WithContext(ctx).Model(&models.Shop{}).Pluck("id", &ids).Error
	if err != nil {
		return nil, err
	}
	return ids, nil
}

func UpdateShop(ctx context.Context, db *gorm.DB, shop *models.Shop) error {
	err := db.Model(&models.Shop{}).Where("id = ?", shop.ID).Updates(shop).Error
	if err != nil {
		return err
	}
	return nil
}

/* ================缓存相关================ */

const (
	ShopCache         = "cache:shop:description:"
	ShopLocationCache = "cache:shop:location:"
)

func GetShopCacheById(ctx context.Context, rds *redis.Client, shopId uint) (*models.Shop, error) {
	key := ShopCache + strconv.Itoa(int(shopId))
	result := rds.Get(ctx, key)

	// 1. 先判断 Redis 是否返回错误
	if result.Err() != nil {
		// 区分"缓存未命中"和"其他错误"
		if errors.Is(result.Err(), redis.Nil) {
			return nil, nil // 缓存未命中：返回 nil, nil（或自定义一个"未命中"错误）
		}
		// 其他错误（如连接失败）：返回 nil + 具体错误
		return nil, fmt.Errorf("redis query failed: %w", result.Err())
	}

	// 2. Redis 键存在，获取JSON字符串
	jsonStr, err := result.Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get cache result: %w", err)
	}

	// 3. JSON反序列化
	shop := &models.Shop{}
	if err := json.Unmarshal([]byte(jsonStr), shop); err != nil {
		// 缓存数据损坏：返回 nil + 反序列化错误
		return nil, fmt.Errorf("cache data unmarshal failed: %w", err)
	}

	// 4. 反序列化成功：返回有效 shop 对象
	return shop, nil
}

func SetShopCacheById(ctx context.Context, rds *redis.Client, shopId uint, shop *models.Shop) error {
	// JSON序列化
	jsonData, err := json.Marshal(shop)
	if err != nil {
		return fmt.Errorf("failed to marshal shop to json: %w", err)
	}

	// 存储到Redis
	err = rds.Set(ctx, ShopCache+strconv.Itoa(int(shopId)), jsonData, time.Hour).Err()
	if err != nil {
		return fmt.Errorf("failed to set cache: %w", err)
	}
	return nil
}

func DelShopCacheById(ctx context.Context, rds *redis.Client, shopId uint) error {
	err := rds.Del(ctx, ShopCache+strconv.Itoa(int(shopId))).Err()
	if err != nil {
		return err
	}
	return nil
}

// LoadShopData 加载店铺地理位置数据到缓存，按照类型进行存到不同key当中
func LoadShopData(ctx context.Context, db *gorm.DB, rds *redis.Client) error {
	// 1. 查询所有的店铺
	var shops []models.Shop
	err := db.WithContext(ctx).Model(&models.Shop{}).Find(&shops).Error
	if err != nil {
		return fmt.Errorf("failed to query shops: %w", err)
	}

	// 2. 遍历店铺，根据类型进行缓存
	for _, shop := range shops {
		// 2.1 使用 GEOADD 存储店铺位置信息
		err = rds.GeoAdd(ctx, ShopLocationCache+strconv.Itoa(int(shop.TypeID)), &redis.GeoLocation{
			Name:      strconv.Itoa(int(shop.ID)),
			Latitude:  shop.Y,
			Longitude: shop.X,
		}).Err()

		if err != nil {
			return fmt.Errorf("failed to set geo cache: %w", err)
		}
	}

	return nil
}

// GetNearbyShops 获取某个店铺的附近某个距离的所有点
func GetNearbyShops(ctx context.Context, rds *redis.Client, shop *models.Shop, radius float64, unit string, count int) ([]uint, error) {
	key := ShopLocationCache + strconv.Itoa(int(shop.TypeID))
	result, err := rds.GeoSearch(ctx, key, &redis.GeoSearchQuery{
		Latitude:   shop.Y,
		Longitude:  shop.X,
		Radius:     radius,
		RadiusUnit: unit,
		Count:      count,
	}).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get geo cache: %w", err)
	}
	// 2. 解析结果，提取店铺ID
	var shopIds []uint
	for _, loc := range result {
		id, _ := strconv.Atoi(loc)
		shopIds = append(shopIds, uint(id))
	}
	return shopIds, nil
}
