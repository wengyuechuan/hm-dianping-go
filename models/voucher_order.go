package models

import (
	"time"

	"gorm.io/gorm"
)

// VoucherOrder 优惠券订单模型
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
}

func (VoucherOrder) TableName() string {
	return "tb_voucher_order"
}
