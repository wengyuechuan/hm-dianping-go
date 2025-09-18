package utils

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/go-redis/redis/v8"
)

// DistributedLock 分布式锁结构体
type DistributedLock struct {
	redisClient *redis.Client
	key         string
	value       string
	ttl         time.Duration
}

// NewDistributedLock 创建新的分布式锁实例
func NewDistributedLock(redisClient *redis.Client, key string, ttl time.Duration) *DistributedLock {
	return &DistributedLock{
		redisClient: redisClient,
		key:         key,
		value:       generateLockValue(),
		ttl:         ttl,
	}
}

// generateLockValue 生成唯一的锁值
func generateLockValue() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// TryLock 尝试获取锁（无TTL）
func TryLock(ctx context.Context, rds *redis.Client, key string) bool {
	// 使用 setnx 模拟互斥锁
	ok, err := rds.SetNX(ctx, key, "1", 0).Result()
	if err != nil {
		return false
	}
	return ok
}

// TryLockWithTTL 尝试获取带TTL的锁（原子性操作）
func TryLockWithTTL(ctx context.Context, rds *redis.Client, key string, ttl time.Duration) (bool, string) {
	// 生成唯一锁值
	lockValue := generateLockValue()
	
	// 使用 SET key value NX EX seconds 确保原子性
	result, err := rds.SetNX(ctx, key, lockValue, ttl).Result()
	if err != nil {
		return false, ""
	}
	
	if result {
		return true, lockValue
	}
	return false, ""
}

// TryLock 尝试获取锁
func (dl *DistributedLock) TryLock(ctx context.Context) bool {
	// 使用 SET key value NX EX seconds 确保原子性
	result, err := dl.redisClient.SetNX(ctx, dl.key, dl.value, dl.ttl).Result()
	if err != nil {
		return false
	}
	return result
}

// UnLock 释放锁
func UnLock(ctx context.Context, rds *redis.Client, key string) {
	// 删除缓存
	rds.Del(ctx, key)
}

// UnLockSafe 安全释放锁（检查锁值）
func UnLockSafe(ctx context.Context, rds *redis.Client, key, value string) bool {
	// Lua脚本确保原子性删除
	luaScript := `
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("del", KEYS[1])
		else
			return 0
		end
	`
	
	result, err := rds.Eval(ctx, luaScript, []string{key}, value).Result()
	if err != nil {
		return false
	}
	
	return result.(int64) == 1
}

// UnLock 释放锁
func (dl *DistributedLock) UnLock(ctx context.Context) bool {
	return UnLockSafe(ctx, dl.redisClient, dl.key, dl.value)
}

// Refresh 刷新锁的TTL
func (dl *DistributedLock) Refresh(ctx context.Context) bool {
	// Lua脚本确保原子性刷新TTL
	luaScript := `
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("expire", KEYS[1], ARGV[2])
		else
			return 0
		end
	`
	
	result, err := dl.redisClient.Eval(ctx, luaScript, []string{dl.key}, dl.value, int(dl.ttl.Seconds())).Result()
	if err != nil {
		return false
	}
	
	return result.(int64) == 1
}
