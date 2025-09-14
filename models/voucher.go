package models

import (
	"time"

	"gorm.io/gorm"
)

// Voucher 优惠券模型
type Voucher struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	ShopID      uint           `json:"shopId"`
	Title       string         `gorm:"size:255" json:"title"`
	SubTitle    string         `gorm:"size:255" json:"subTitle"`
	Rules       string         `gorm:"size:1024" json:"rules"`
	PayValue    int64          `json:"payValue"`
	ActualValue int64          `json:"actualValue"`
	Type        int            `json:"type"`   // 0-普通券，1-秒杀券
	Status      int            `json:"status"` // 1-上架，2-下架
	Stock       int            `json:"stock"`
	BeginTime   *time.Time     `json:"beginTime"`
	EndTime     *time.Time     `json:"endTime"`
}

func (Voucher) TableName() string {
	return "tb_voucher"
}
