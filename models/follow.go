package models

import (
	"time"

	"gorm.io/gorm"
)

// Follow 关注模型
type Follow struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	UserID     uint           `json:"userId"`
	FollowUserID uint         `json:"followUserId"`
}