package models

import (
	"time"

	"gorm.io/gorm"
)

// ShopType 商铺类型模型
type ShopType struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Name      string         `gorm:"size:32" json:"name"`
	Icon      string         `gorm:"size:255" json:"icon"`
	Sort      int            `json:"sort"`
}

func (ShopType) TableName() string {
	return "tb_shop_type"
}
