package dao

import (
	"context"
	"hm-dianping-go/models"
	"gorm.io/gorm"
)

// CreateVoucherOrder 创建优惠券订单
func CreateVoucherOrder(ctx context.Context, db *gorm.DB, order *models.VoucherOrder) error {
	return db.WithContext(ctx).Create(order).Error
}

// GetVoucherOrderByID 根据订单ID获取订单信息
func GetVoucherOrderByID(ctx context.Context, db *gorm.DB, orderID uint) (*models.VoucherOrder, error) {
	var order models.VoucherOrder
	err := db.WithContext(ctx).First(&order, orderID).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// GetVoucherOrderByUserAndVoucher 根据用户ID和优惠券ID获取订单
func GetVoucherOrderByUserAndVoucher(ctx context.Context, db *gorm.DB, userID, voucherID uint) (*models.VoucherOrder, error) {
	var order models.VoucherOrder
	err := db.WithContext(ctx).Where("user_id = ? AND voucher_id = ?", userID, voucherID).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// CheckVoucherOrderExists 检查用户是否已购买该优惠券
func CheckVoucherOrderExists(ctx context.Context, db *gorm.DB, userID, voucherID uint) (bool, error) {
	var count int64
	err := db.WithContext(ctx).Model(&models.VoucherOrder{}).Where("user_id = ? AND voucher_id = ?", userID, voucherID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// UpdateVoucherOrder 更新订单信息
func UpdateVoucherOrder(ctx context.Context, db *gorm.DB, order *models.VoucherOrder) error {
	return db.WithContext(ctx).Save(order).Error
}

// UpdateVoucherOrderStatus 更新订单状态
func UpdateVoucherOrderStatus(ctx context.Context, db *gorm.DB, orderID uint, status int) error {
	return db.WithContext(ctx).Model(&models.VoucherOrder{}).Where("id = ?", orderID).Update("status", status).Error
}

// DeleteVoucherOrder 删除订单（软删除）
func DeleteVoucherOrder(ctx context.Context, db *gorm.DB, orderID uint) error {
	return db.WithContext(ctx).Delete(&models.VoucherOrder{}, orderID).Error
}

// GetVoucherOrdersByUser 获取用户的所有订单
func GetVoucherOrdersByUser(ctx context.Context, db *gorm.DB, userID uint, page, size int) ([]models.VoucherOrder, error) {
	var orders []models.VoucherOrder
	offset := (page - 1) * size
	err := db.WithContext(ctx).Where("user_id = ?", userID).Offset(offset).Limit(size).Find(&orders).Error
	return orders, err
}

// GetVoucherOrdersByVoucher 获取某优惠券的所有订单
func GetVoucherOrdersByVoucher(ctx context.Context, db *gorm.DB, voucherID uint, page, size int) ([]models.VoucherOrder, error) {
	var orders []models.VoucherOrder
	offset := (page - 1) * size
	err := db.WithContext(ctx).Where("voucher_id = ?", voucherID).Offset(offset).Limit(size).Find(&orders).Error
	return orders, err
}

// CountVoucherOrdersByUser 统计用户订单数量
func CountVoucherOrdersByUser(ctx context.Context, db *gorm.DB, userID uint) (int64, error) {
	var count int64
	err := db.WithContext(ctx).Model(&models.VoucherOrder{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

// CountVoucherOrdersByVoucher 统计某优惠券的订单数量
func CountVoucherOrdersByVoucher(ctx context.Context, db *gorm.DB, voucherID uint) (int64, error) {
	var count int64
	err := db.WithContext(ctx).Model(&models.VoucherOrder{}).Where("voucher_id = ?", voucherID).Count(&count).Error
	return count, err
}