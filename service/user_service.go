package service

import (
	"hm-dianping-go/dao"
	"hm-dianping-go/models"
	"hm-dianping-go/utils"
)

// UserRegister 用户注册服务
func UserRegister(phone, code, password, nickName string) *utils.Result {
	// TODO: 验证短信验证码
	// 这里暂时跳过验证码验证

	// 检查用户是否已存在
	var existingUser models.User
	err := dao.DB.Where("phone = ?", phone).First(&existingUser).Error
	if err == nil {
		return utils.ErrorResult("手机号已注册")
	}

	// 创建新用户
	user := models.User{
		Phone:    phone,
		Password: utils.HashPassword(password),
		NickName: nickName,
	}

	if nickName == "" {
		user.NickName = "用户" + phone[len(phone)-4:]
	}

	if err := dao.DB.Create(&user).Error; err != nil {
		return utils.ErrorResult("注册失败")
	}

	return utils.SuccessResult("注册成功")
}

// UserLogin 用户登录服务
func UserLogin(phone, password string) *utils.Result {
	var user models.User
	err := dao.DB.Where("phone = ?", phone).First(&user).Error
	if err != nil {
		return utils.ErrorResult("用户不存在")
	}

	if !utils.CheckPassword(password, user.Password) {
		return utils.ErrorResult("密码错误")
	}

	// 生成JWT token
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return utils.ErrorResult("登录失败")
	}

	return utils.SuccessResultWithData(map[string]interface{}{
		"token": token,
		"user": map[string]interface{}{
			"id":       user.ID,
			"phone":    user.Phone,
			"nickName": user.NickName,
			"icon":     user.Icon,
		},
	})
}

// GetUserInfo 获取用户信息服务
func GetUserInfo(userID uint) *utils.Result {
	var user models.User
	err := dao.DB.First(&user, userID).Error
	if err != nil {
		return utils.ErrorResult("用户不存在")
	}

	return utils.SuccessResultWithData(map[string]interface{}{
		"id":       user.ID,
		"phone":    user.Phone,
		"nickName": user.NickName,
		"icon":     user.Icon,
	})
}

// UpdateUserInfo 更新用户信息服务
func UpdateUserInfo(userID uint, nickName, icon string) *utils.Result {
	var user models.User
	err := dao.DB.First(&user, userID).Error
	if err != nil {
		return utils.ErrorResult("用户不存在")
	}

	if nickName != "" {
		user.NickName = nickName
	}
	if icon != "" {
		user.Icon = icon
	}

	if err := dao.DB.Save(&user).Error; err != nil {
		return utils.ErrorResult("更新失败")
	}

	return utils.SuccessResult("更新成功")
}

// SendCode 发送验证码服务
func SendCode(phone string) *utils.Result {
	// TODO: 实现发送短信验证码功能
	// 这里暂时返回成功
	return utils.SuccessResult("验证码发送成功")
}