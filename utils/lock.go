package utils

import (
	"context"

	"github.com/go-redis/redis/v8"
)

func TryLock(ctx context.Context, rds *redis.Client, key string) bool {
	// 1. 使用 setnx 模拟互斥锁
	ok, err := rds.SetNX(ctx, key, "1", 0).Result()
	if err != nil {
		return false
	}
	return ok
}

func UnLock(ctx context.Context, rds *redis.Client, key string) {
	// 1. 删除缓存
	rds.Del(ctx, key)
}
