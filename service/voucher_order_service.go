package service

import (
	"context"
	"errors"
	"hm-dianping-go/dao"
	"hm-dianping-go/models"
	"hm-dianping-go/utils"
	"log"
	"time"

	"gorm.io/gorm"
)

// SeckillVoucher 秒杀优惠券（使用乐观锁）
func SeckillVoucher(ctx context.Context, userId, voucherId uint) *utils.Result {
	// 1. 检查秒杀券是否存在
	seckillVoucher, err := dao.GetSeckillVoucherByID(voucherId)
	if err != nil {
		log.Printf("查询秒杀券失败: %v", err)
		return utils.ErrorResult("秒杀券不存在")
	}

	// 2. 检查秒杀时间
	now := time.Now()
	if now.Before(seckillVoucher.BeginTime) {
		return utils.ErrorResult("秒杀尚未开始")
	}
	if now.After(seckillVoucher.EndTime) {
		return utils.ErrorResult("秒杀已结束")
	}

	// 3. 检查库存
	if seckillVoucher.Stock <= 0 {
		return utils.ErrorResult("库存不足")
	}

	// 4. 检查用户是否已经购买过（使用DAO层函数）
	exists, err := dao.CheckVoucherOrderExists(ctx, dao.DB, userId, voucherId)
	if err != nil {
		log.Printf("检查订单是否存在失败: %v", err)
		return utils.ErrorResult("系统错误")
	}
	if exists {
		return utils.ErrorResult("不能重复购买")
	}

	// 5. 使用乐观锁重试机制进行库存扣减和订单创建
	const maxRetries = 3
	for i := 0; i < maxRetries; i++ {
		// 开始事务
		tx := dao.DB.Begin()
		if tx.Error != nil {
			log.Printf("开始事务失败: %v", tx.Error)
			return utils.ErrorResult("系统错误")
		}

		// 扣减库存（乐观锁CAS操作）
		err = dao.UpdateSeckillVoucherStock(voucherId, 1)
		if err != nil {
			tx.Rollback()
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 库存不足或并发冲突，重试
				if i == maxRetries-1 {
					return utils.ErrorResult("库存不足")
				}
				// 短暂等待后重试
				time.Sleep(time.Duration(i+1) * 10 * time.Millisecond)
				continue
			}
			log.Printf("扣减库存失败: %v", err)
			return utils.ErrorResult("系统错误")
		}

		// 6. 创建订单（使用DAO层函数）
		now = time.Now()
		order := &models.VoucherOrder{
			UserID:     userId,
			VoucherID:  voucherId,
			PayType:    1,
			Status:     1,
			CreateTime: &now,
		}

		err = dao.CreateVoucherOrder(ctx, tx, order)
		if err != nil {
			tx.Rollback()
			log.Printf("创建订单失败: %v", err)
			return utils.ErrorResult("创建订单失败")
		}

		// 7. 提交事务
		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			log.Printf("提交事务失败: %v", err)
			return utils.ErrorResult("系统错误")
		}

		// 8. 成功，返回订单ID
		return utils.SuccessResultWithData(order.ID)
	}

	// 重试次数用完，返回失败
	return utils.ErrorResult("服务繁忙，请稍后重试")
}
