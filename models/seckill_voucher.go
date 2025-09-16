package models

import "time"

// SeckillVoucher 秒杀优惠券模型
// 与优惠券是一对一关系
type SeckillVoucher struct {
	VoucherID  uint      `gorm:"primaryKey;column:voucher_id" json:"voucherId"` // 关联的优惠券的id
	Stock      int       `gorm:"column:stock;not null" json:"stock"`            // 库存
	CreateTime time.Time `gorm:"column:create_time;not null;default:CURRENT_TIMESTAMP" json:"createTime"` // 创建时间
	BeginTime  time.Time `gorm:"column:begin_time;not null" json:"beginTime"`   // 生效时间
	EndTime    time.Time `gorm:"column:end_time;not null" json:"endTime"`       // 失效时间
	UpdateTime time.Time `gorm:"column:update_time;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updateTime"` // 更新时间
}

// TableName 指定表名
func (SeckillVoucher) TableName() string {
	return "tb_seckill_voucher"
}