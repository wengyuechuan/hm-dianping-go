package utils

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisIdWorker struct {
	rdb            *redis.Client
	beginTimestamp int64
	countBits      uint8
}

func NewRedisIdWorker(rdb *redis.Client, countBits uint8) *RedisIdWorker {
	return &RedisIdWorker{
		rdb:            rdb,
		beginTimestamp: 1672531200000, // 2023-01-01 00:00:00 UTC的毫秒时间戳
		countBits:      countBits,
	}
}

func (w *RedisIdWorker) NextId(ctx context.Context, key string) (int64, error) {
	// 这里生成的ID开始的符号位为0，然后拼接时间戳，然后生成序列号，完成一个全局id的生成器

	// 1. 生成时间戳部分
	now := time.Now().UTC().UnixMilli()
	timestamp := now - w.beginTimestamp

	// 2. 生成序列号部分
	date := time.Now().UTC().Format("2006:01:02")
	seq, err := w.rdb.Incr(ctx, "icr:"+key+":"+date).Result()
	if err != nil {
		return 0, err
	}

	return (timestamp << w.countBits) | seq, nil
}
