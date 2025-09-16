package models

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Phone     string         `gorm:"uniqueIndex;size:11" json:"phone"`
	Password  string         `gorm:"size:255" json:"-"`
	NickName  string         `gorm:"size:32" json:"nickName"`
	Icon      string         `gorm:"size:255" json:"icon"`
}

func (User) TableName() string {
	return "tb_user"
}
