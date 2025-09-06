package service

import (
	"hm-dianping-go/dao"
	"hm-dianping-go/models"
	"hm-dianping-go/utils"
)

// SeckillVoucher 秒杀优惠券
func SeckillVoucher(userId, voucherId uint) *utils.Result {
	// 查询优惠券
	var voucher models.Voucher
	err := dao.DB.First(&voucher, voucherId).Error
	if err != nil {
		return utils.ErrorResult("优惠券不存在")
	}

	// 检查库存
	if voucher.Stock <= 0 {
		return utils.ErrorResult("库存不足")
	}

	// 检查是否已经购买过
	var existingOrder models.VoucherOrder
	err = dao.DB.Where("user_id = ? AND voucher_id = ?", userId, voucherId).First(&existingOrder).Error
	if err == nil {
		return utils.ErrorResult("不能重复购买")
	}

	// 开始事务
	tx := dao.DB.Begin()

	// 扣减库存
	result := tx.Model(&voucher).Where("id = ? AND stock > 0", voucherId).Update("stock", voucher.Stock-1)
	if result.Error != nil || result.RowsAffected == 0 {
		tx.Rollback()
		return utils.ErrorResult("库存不足")
	}

	// 创建订单
	order := models.VoucherOrder{
		UserID:    userId,
		VoucherID: voucherId,
		PayType:   1,
		Status:    1,
	}

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return utils.ErrorResult("创建订单失败")
	}

	tx.Commit()
	return utils.SuccessResultWithData(order.ID)
}