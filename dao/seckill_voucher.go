package dao

import (
	"hm-dianping-go/models"

	"gorm.io/gorm"
)

// CreateSeckillVoucher 创建秒杀券
func CreateSeckillVoucher(seckillVoucher *models.SeckillVoucher) error {
	return DB.Create(seckillVoucher).Error
}

// GetSeckillVoucherByID 根据优惠券ID获取秒杀券信息
func GetSeckillVoucherByID(voucherID uint) (*models.SeckillVoucher, error) {
	var seckillVoucher models.SeckillVoucher
	err := DB.Where("voucher_id = ?", voucherID).First(&seckillVoucher).Error
	if err != nil {
		return nil, err
	}
	return &seckillVoucher, nil
}

// UpdateSeckillVoucher 更新秒杀券信息
func UpdateSeckillVoucher(seckillVoucher *models.SeckillVoucher) error {
	return DB.Save(seckillVoucher).Error
}

// DeleteSeckillVoucher 删除秒杀券
func DeleteSeckillVoucher(voucherID uint) error {
	return DB.Where("voucher_id = ?", voucherID).Delete(&models.SeckillVoucher{}).Error
}

// UpdateSeckillVoucherStock 更新秒杀券库存（原子操作）
func UpdateSeckillVoucherStock(voucherID uint, stock int) error {
	result := DB.Model(&models.SeckillVoucher{}).
		Where("voucher_id = ? AND stock >= ?", voucherID, stock).
		Update("stock", gorm.Expr("stock - ?", stock))

	if result.Error != nil {
		return result.Error
	}

	// 检查是否有行被更新，如果没有说明库存不足
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// CheckSeckillVoucherExists 检查秒杀券是否存在
func CheckSeckillVoucherExists(voucherID uint) (bool, error) {
	var count int64
	err := DB.Model(&models.SeckillVoucher{}).Where("voucher_id = ?", voucherID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
