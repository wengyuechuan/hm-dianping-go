package service

import (
	"context"
	"fmt"
	"time"

	"hm-dianping-go/dao"
	"hm-dianping-go/models"
	"hm-dianping-go/utils"
)

// UserRegister 用户注册服务
func UserRegister(phone, code, password, nickName string) *utils.Result {
	// 校验手机号格式
	if utils.IsPhoneInvalid(phone) {
		return utils.ErrorResult("手机号格式不正确")
	}

	// 校验验证码格式
	if utils.IsCodeInvalid(code) {
		return utils.ErrorResult("验证码格式不正确")
	}

	// 校验密码格式
	if utils.IsPasswordInvalid(password) {
		return utils.ErrorResult("密码格式不正确，需要4-32位字母、数字或下划线")
	}

	// TODO: 验证短信验证码
	// 这里暂时跳过验证码验证

	// 检查用户是否已存在
	exists, err := dao.CheckUserExistsByPhone(phone)
	if err != nil {
		return utils.ErrorResult("系统错误，请稍后重试")
	}
	if exists {
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

	if err := dao.CreateUser(&user); err != nil {
		return utils.ErrorResult("注册失败")
	}

	return utils.SuccessResult("注册成功")
}

// UserLogin 用户登录服务
func UserLogin(phone, code string) *utils.Result {
	// 从Redis获取验证码进行验证
	storedCode, err := dao.GetLoginCode(phone)
	if err != nil {
		return utils.ErrorResult("验证码已过期或不存在，请重新获取")
	}

	// 验证验证码是否正确
	if storedCode != code {
		return utils.ErrorResult("验证码错误")
	}

	// 验证成功后删除验证码（防止重复使用）
	_ = dao.DeleteLoginCode(phone)

	// 根据手机号查询用户
	user, err := dao.GetUserByPhone(phone)
	if err != nil {
		// 用户不存在，自动注册
		newUser := models.User{
			Phone:    phone,
			NickName: "用户" + phone[7:], // 使用手机号后4位作为昵称
		}
		if err = dao.CreateUser(&newUser); err != nil {
			return utils.ErrorResult("登录失败")
		}
		user = &newUser
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
	user, err := dao.GetUserByID(userID)
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
	user, err := dao.GetUserByID(userID)
	if err != nil {
		return utils.ErrorResult("用户不存在")
	}

	if nickName != "" {
		user.NickName = nickName
	}
	if icon != "" {
		user.Icon = icon
	}

	if err := dao.UpdateUser(user); err != nil {
		return utils.ErrorResult("更新失败")
	}

	return utils.SuccessResult("更新成功")
}

// SendCode 发送验证码服务
func SendCode(phone string) *utils.Result {
	// 检查是否已存在未过期的验证码
	exists, err := dao.CheckLoginCodeExists(phone)
	if err != nil {
		return utils.ErrorResult("系统错误，请稍后重试")
	}
	// 防止频繁请求，进行避免
	if exists {
		// 获取剩余时间
		ttl, _ := dao.GetLoginCodeTTL(phone)
		if ttl > 0 {
			return utils.ErrorResult(fmt.Sprintf("验证码已发送，请%d秒后重试", int(ttl.Seconds())))
		}
	}

	// 生成6位随机验证码
	code := utils.GenerateRandomCode(6)

	// 将验证码存储到Redis，设置5分钟过期
	err = dao.SetLoginCode(phone, code, 0) // 0表示使用默认过期时间
	if err != nil {
		return utils.ErrorResult("验证码发送失败，请稍后重试")
	}

	// TODO: 实现发送短信验证码功能
	// 这里可以集成短信服务商API，如阿里云短信、腾讯云短信等
	// 暂时在日志中输出验证码（仅用于开发测试）
	fmt.Printf("[开发模式] 手机号 %s 的验证码是: %s\n", phone, code)

	return utils.SuccessResult("验证码发送成功")
}

func Sign(ctx context.Context, userID uint) *utils.Result {
	// 检查用户是否已签到
	// 这里使用 redis 的 bitMap 来实现

	// 1. 获取本月的日期
	date := time.Now().Format("200601")

	// 2. 获取今天是本月的第几天
	day := time.Now().Day()

	// 2. 直接签到
	if err := dao.SignUser(ctx, dao.Redis, userID, date, day); err != nil {
		return utils.ErrorResult("签到失败")
	}
	return utils.SuccessResult("签到成功")
}

func CheckSign(ctx context.Context, userID uint, month string) *utils.Result {
	// 检查某个月到某一天的连续签到次数
	// 获取到当前时间的第几天
	day := time.Now().Day()

	result, err := dao.CheckSign(ctx, dao.Redis, userID, month, day)
	if err != nil {
		return utils.ErrorResult("查询签到状态失败")
	}

	count := 0
	for {
		if (result & 1) == 0 {
			break
		} else {
			count++
		}
		result >>= 1
	}

	// 3. 检查签到状态
	return utils.SuccessResultWithData(count)
}
