package dao

import (
	"hm-dianping-go/models"
)

// GetUserByPhone 根据手机号查询用户
func GetUserByPhone(phone string) (*models.User, error) {
	var user models.User
	err := DB.Where("phone = ?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByID 根据用户ID查询用户
func GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	err := DB.First(&user, userID).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateUser 创建新用户
func CreateUser(user *models.User) error {
	return DB.Create(user).Error
}

// UpdateUser 更新用户信息
func UpdateUser(user *models.User) error {
	return DB.Save(user).Error
}

// CheckUserExistsByPhone 检查手机号是否已注册
func CheckUserExistsByPhone(phone string) (bool, error) {
	var count int64
	err := DB.Model(&models.User{}).Where("phone = ?", phone).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}