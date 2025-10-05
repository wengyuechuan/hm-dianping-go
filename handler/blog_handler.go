package handler

import (
	"hm-dianping-go/service"
	"hm-dianping-go/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateBlog 创建博客
func CreateBlog(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "用户未登录")
		return
	}

	var req struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
		Images  string `json:"images"`
		ShopId  uint   `json:"shopId"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	result := service.CreateBlog(c.Request.Context(), userID.(uint), req.Title, req.Content, req.Images, req.ShopId)
	utils.Response(c, result)
}

// LikeBlog 点赞博客
func LikeBlog(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "用户未登录")
		return
	}

	blogIdStr := c.Param("id")
	blogId, err := strconv.ParseUint(blogIdStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的博客ID")
		return
	}

	result := service.LikeBlog(c.Request.Context(), userID.(uint), uint(blogId))
	utils.Response(c, result)
}

// GetBlogList 获取博客列表
func GetBlogList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	result := service.GetBlogList(c.Request.Context(), page, size)
	utils.Response(c, result)
}

// GetBlogById 根据ID获取博客
func GetBlogById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的博客ID")
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "用户未登录")
		return
	}

	result := service.GetBlogById(c.Request.Context(), uint(id), userID.(uint))
	utils.Response(c, result)
}

// GetHotBlogList 获取热门博客列表
func GetHotBlogList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("current", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "用户未登录")
		return
	}

	result := service.GetHotBlogList(c.Request.Context(), page, size, userID.(uint))
	utils.Response(c, result)
}

// GetMyBlogList 获取我的博客列表
func GetMyBlogList(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "用户未登录")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("current", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	result := service.GetMyBlogList(c.Request.Context(), userID.(uint), page, size)
	utils.Response(c, result)
}
