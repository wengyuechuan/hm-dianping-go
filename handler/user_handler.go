package handler

import (
	"hm-dianping-go/service"
	"hm-dianping-go/utils"
	"net/http"

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
	var req struct {
		Phone    string `json:"phone" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	result := service.UserLogin(req.Phone, req.Password)
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

	result := service.SendCode(phone)
	utils.Response(c, result)
}