package handler

import (
	"hm-dianping-go/service"
	"hm-dianping-go/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// UserRegister 用户注册
func UserRegister(c *gin.Context) {
	var req struct {
		Phone    string `json:"phone" binding:"required"`
		Code     string `json:"code" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
		NickName string `json:"nickName"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	result := service.UserRegister(req.Phone, req.Code, req.Password, req.NickName)
	utils.Response(c, result)
}

// UserLogin 用户登录
func UserLogin(c *gin.Context) {
	// 这里只能支持验证码登录
	var req struct {
		Phone string `json:"phone" binding:"required"`
		Code  string `json:"code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	// 判断手机号格式是否正确
	if ok := utils.IsPhoneValid(req.Phone); !ok {
		utils.ErrorResponse(c, http.StatusBadRequest, "手机号格式不正确")
		return
	}

	// 判断验证码格式是否正确
	if utils.IsCodeInvalid(req.Code) {
		utils.ErrorResponse(c, http.StatusBadRequest, "验证码格式不正确")
		return
	}

	result := service.UserLogin(req.Phone, req.Code)
	utils.Response(c, result)
}

// GetUserInfo 获取用户信息
func GetUserInfo(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "用户未登录")
		return
	}

	result := service.GetUserInfo(userID.(uint))
	utils.Response(c, result)
}

// UpdateUserInfo 更新用户信息
func UpdateUserInfo(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "用户未登录")
		return
	}

	var req struct {
		NickName string `json:"nickName"`
		Icon     string `json:"icon"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	result := service.UpdateUserInfo(userID.(uint), req.NickName, req.Icon)
	utils.Response(c, result)
}

// UserLogout 用户登出
func UserLogout(c *gin.Context) {
	// TODO: 实现登出逻辑，清除token等
	utils.SuccessResponse(c, "登出成功")
}

// SendCode 发送验证码
func SendCode(c *gin.Context) {
	phone := c.Query("phone")
	if phone == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "手机号不能为空")
		return
	}

	if !utils.IsPhoneValid(phone) {
		utils.ErrorResponse(c, http.StatusBadRequest, "手机号格式不正确")
		return
	}
	result := service.SendCode(phone)
	utils.Response(c, result)
}

// Sign 用户签到
func Sign(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "用户未登录")
		return
	}

	result := service.Sign(c.Request.Context(), userID.(uint))
	utils.Response(c, result)
}

// CheckSign 获取用户签到状态
func CheckSign(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "用户未登录")
		return
	}

	month := c.Query("month")
	if month == "" {
		// 以当前月为准
		month = time.Now().Format("2006-01")
	}

	result := service.CheckSign(c.Request.Context(), userID.(uint), month)
	utils.Response(c, result)
}
