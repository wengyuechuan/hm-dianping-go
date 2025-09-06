package models

import (
	"time"

	"gorm.io/gorm"
)

// BlogLike 博客点赞模型
type BlogLike struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	UserID    uint           `json:"userId"`
	BlogID    uint           `json:"blogId"`
}