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

	// 4. 检查用户是否已经购买过该秒杀券（一人一单限制）
	exists, err := dao.CheckSeckillVoucherOrderExists(ctx, dao.DB, userId, voucherId)
	if err != nil {
		log.Printf("检查秒杀券订单是否存在失败: %v", err)
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

		// 6. 创建秒杀券订单（使用DAO层函数）
		now = time.Now()
		order := &models.VoucherOrder{
			UserID:      userId,
			VoucherID:   voucherId,
			PayType:     1,
			Status:      1,
			CreateTime:  &now,
			VoucherType: 2, // 秒杀券类型
		}

		err = dao.CreateVoucherOrder(ctx, tx, order) // 这里创建的逻辑，主要依赖于mysql的约束
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

/*
需要注意的是，对于上面的方案，
1. 秒杀券的库存扣减和订单创建是一个原子操作，通过事务确保数据一致性。
2. 普通券的库存扣减和订单创建也是一个原子操作，通过事务确保数据一致性。
3. 秒杀券的唯一索引确保了每个用户只能购买一次秒杀券。
4. 普通券的索引确保了用户可以购买多次普通券。
// 在models/voucher_order.go中
type VoucherOrder struct {
    ID         uint           `gorm:"primarykey" json:"id"`
    CreatedAt  time.Time      `json:"createdAt"`
    UpdatedAt  time.Time      `json:"updatedAt"`
    DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
    UserID     uint           `json:"userId"`
    VoucherID  uint           `json:"voucherId"`
    PayType    int            `json:"payType"`
    Status     int            `json:"status"`
    CreateTime *time.Time     `json:"createTime"`
    PayTime    *time.Time     `json:"payTime"`
    UseTime    *time.Time     `json:"useTime"`
    RefundTime *time.Time     `json:"refundTime"`
    UpdateTime *time.Time     `json:"updateTime"`

    // 新增字段：券类型标识
    VoucherType int `gorm:"index" json:"voucherType"` // 1:普通券 2:秒杀券
}

// 添加复合索引，但不设为唯一
func (VoucherOrder) TableName() string {
    return "tb_voucher_order"
}

-- 只对秒杀券创建唯一约束
CREATE UNIQUE INDEX uk_seckill_user_voucher
ON tb_voucher_order (user_id, voucher_id)
WHERE voucher_type = 2;

-- 普通券只创建普通索引用于查询优化
CREATE INDEX idx_normal_user_voucher
ON tb_voucher_order (user_id, voucher_id, voucher_type);

*/
