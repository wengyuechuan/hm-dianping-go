package models

import (
	"time"

	"gorm.io/gorm"
)

// Blog 博客模型
type Blog struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	ShopID    uint           `json:"shopId"`
	UserID    uint           `json:"userId"`
	Title     string         `gorm:"size:255" json:"title"`
	Images    string         `gorm:"size:2048" json:"images"`
	Content   string         `gorm:"size:2048" json:"content"`
	Liked     int            `json:"liked"`
	Comments  int            `json:"comments"`
}

func (Blog) TableName() string {
	return "tb_blog"
}
