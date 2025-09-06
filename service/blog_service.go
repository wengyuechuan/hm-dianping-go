package service

import (
	"hm-dianping-go/dao"
	"hm-dianping-go/models"
	"hm-dianping-go/utils"
)

// CreateBlog 创建博客
func CreateBlog(userId uint, title, content, images string, shopId uint) *utils.Result {
	blog := models.Blog{
		UserID:  userId,
		Title:   title,
		Content: content,
		Images:  images,
		ShopID:  shopId,
	}

	if err := dao.DB.Create(&blog).Error; err != nil {
		return utils.ErrorResult("创建失败")
	}

	return utils.SuccessResultWithData(blog.ID)
}

// LikeBlog 点赞博客
func LikeBlog(userId, blogId uint) *utils.Result {
	// 检查是否已经点赞
	var existingLike models.BlogLike
	err := dao.DB.Where("user_id = ? AND blog_id = ?", userId, blogId).First(&existingLike).Error
	if err == nil {
		// 已经点赞，取消点赞
		dao.DB.Delete(&existingLike)
		dao.DB.Model(&models.Blog{}).Where("id = ?", blogId).UpdateColumn("liked", dao.DB.Raw("liked - 1"))
		return utils.SuccessResult("取消点赞成功")
	}

	// 添加点赞
	like := models.BlogLike{
		UserID: userId,
		BlogID: blogId,
	}

	if err := dao.DB.Create(&like).Error; err != nil {
		return utils.ErrorResult("点赞失败")
	}

	// 更新博客点赞数
	dao.DB.Model(&models.Blog{}).Where("id = ?", blogId).UpdateColumn("liked", dao.DB.Raw("liked + 1"))

	return utils.SuccessResult("点赞成功")
}

// GetBlogList 获取博客列表
func GetBlogList(page, size int) *utils.Result {
	var blogs []models.Blog
	var total int64

	offset := (page - 1) * size

	// 获取总数
	dao.DB.Model(&models.Blog{}).Count(&total)

	// 分页查询
	err := dao.DB.Offset(offset).Limit(size).Order("created_at desc").Find(&blogs).Error
	if err != nil {
		return utils.ErrorResult("查询失败")
	}

	return utils.SuccessResultWithData(map[string]interface{}{
		"list":  blogs,
		"total": total,
		"page":  page,
		"size":  size,
	})
}

// GetBlogById 根据ID获取博客
func GetBlogById(id uint) *utils.Result {
	var blog models.Blog
	err := dao.DB.First(&blog, id).Error
	if err != nil {
		return utils.ErrorResult("博客不存在")
	}

	return utils.SuccessResultWithData(blog)
}

// GetHotBlogList 获取热门博客列表
func GetHotBlogList(page, size int) *utils.Result {
	var blogs []models.Blog
	var total int64

	offset := (page - 1) * size

	// 获取总数
	dao.DB.Model(&models.Blog{}).Count(&total)

	// 分页查询，按点赞数排序
	err := dao.DB.Offset(offset).Limit(size).Order("liked desc, created_at desc").Find(&blogs).Error
	if err != nil {
		return utils.ErrorResult("查询失败")
	}

	return utils.SuccessResultWithData(map[string]interface{}{
		"list":  blogs,
		"total": total,
		"page":  page,
		"size":  size,
	})
}

// GetMyBlogList 获取我的博客列表
func GetMyBlogList(userId uint, page, size int) *utils.Result {
	var blogs []models.Blog
	var total int64

	offset := (page - 1) * size

	// 获取总数
	dao.DB.Model(&models.Blog{}).Where("user_id = ?", userId).Count(&total)

	// 分页查询
	err := dao.DB.Where("user_id = ?", userId).Offset(offset).Limit(size).Order("created_at desc").Find(&blogs).Error
	if err != nil {
		return utils.ErrorResult("查询失败")
	}

	return utils.SuccessResultWithData(map[string]interface{}{
		"list":  blogs,
		"total": total,
		"page":  page,
		"size":  size,
	})
}