package dao

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

// Redis验证码相关常量
const (
	// LoginCodePrefix 登录验证码Redis key前缀
	LoginCodePrefix = "dianping:user:login:phone:"
	// DefaultCodeExpiration 默认验证码过期时间（5分钟）
	DefaultCodeExpiration = 5 * time.Minute
)

// SetLoginCode 设置登录验证码到Redis
// phone: 手机号
// code: 验证码
// expiration: 过期时间，如果为0则使用默认5分钟
func SetLoginCode(phone, code string, expiration time.Duration) error {
	if Redis == nil {
		return fmt.Errorf("Redis client not initialized")
	}

	if expiration == 0 {
		expiration = DefaultCodeExpiration
	}

	key := LoginCodePrefix + phone
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := Redis.Set(ctx, key, code, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set login code for phone %s: %w", phone, err)
	}

	log.Printf("Login code set for phone: %s, expiration: %v", phone, expiration)
	return nil
}

// GetLoginCode 从Redis获取登录验证码
// phone: 手机号
// 返回验证码和错误信息
func GetLoginCode(phone string) (string, error) {
	if Redis == nil {
		return "", fmt.Errorf("Redis client not initialized")
	}

	key := LoginCodePrefix + phone
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	code, err := Redis.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("login code not found or expired for phone: %s", phone)
	}
	if err != nil {
		return "", fmt.Errorf("failed to get login code for phone %s: %w", phone, err)
	}

	return code, nil
}

// DeleteLoginCode 删除Redis中的登录验证码
// phone: 手机号
func DeleteLoginCode(phone string) error {
	if Redis == nil {
		return fmt.Errorf("Redis client not initialized")
	}

	key := LoginCodePrefix + phone
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := Redis.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete login code for phone %s: %w", phone, err)
	}

	log.Printf("Login code deleted for phone: %s", phone)
	return nil
}

// CheckLoginCodeExists 检查登录验证码是否存在
// phone: 手机号
// 返回是否存在和错误信息
func CheckLoginCodeExists(phone string) (bool, error) {
	if Redis == nil {
		return false, fmt.Errorf("Redis client not initialized")
	}

	key := LoginCodePrefix + phone
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	exists, err := Redis.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check login code existence for phone %s: %w", phone, err)
	}

	return exists > 0, nil
}

// GetLoginCodeTTL 获取登录验证码的剩余过期时间
// phone: 手机号
// 返回剩余时间和错误信息
func GetLoginCodeTTL(phone string) (time.Duration, error) {
	if Redis == nil {
		return 0, fmt.Errorf("Redis client not initialized")
	}

	key := LoginCodePrefix + phone
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ttl, err := Redis.TTL(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get TTL for login code of phone %s: %w", phone, err)
	}

	return ttl, nil
}