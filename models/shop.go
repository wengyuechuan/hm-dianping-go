package models

import (
	"time"

	"gorm.io/gorm"
)

// Shop 商铺模型
type Shop struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Name        string         `gorm:"size:128" json:"name"`
	TypeID      uint           `json:"typeId"`
	Images      string         `gorm:"size:1024" json:"images"`
	Area        string         `gorm:"size:128" json:"area"`
	Address     string         `gorm:"size:255" json:"address"`
	X           float64        `json:"x"`
	Y           float64        `json:"y"`
	AvgPrice    int            `json:"avgPrice"`
	Sold        int            `json:"sold"`
	Comments    int            `json:"comments"`
	Score       int            `json:"score"`
	OpenHours   string         `gorm:"size:32" json:"openHours"`
}