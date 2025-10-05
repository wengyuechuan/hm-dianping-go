package service

import (
	"context"
	"hm-dianping-go/dao"
	"hm-dianping-go/models"
	"hm-dianping-go/utils"
	"strconv"

	"gorm.io/gorm"
)

// CreateBlog 创建博客
func CreateBlog(ctx context.Context, userId uint, title, content, images string, shopId uint) *utils.Result {
	blog := models.Blog{
		UserID:  userId,
		Title:   title,
		Content: content,
		Images:  images,
		ShopID:  shopId,
	}

	if err := dao.CreateBlog(ctx, &blog); err != nil {
		return utils.ErrorResult("创建失败")
	}

	return utils.SuccessResultWithData(blog.ID)
}

// LikeBlog 点赞博客
func LikeBlog(ctx context.Context, userId, blogId uint) *utils.Result {
	// 检查是否已经点赞
	liked, err := dao.IsLikedMember(ctx, dao.Redis, userId, blogId)
	if err == nil {
		if liked {
			// 已经点赞，取消点赞
			if err := dao.RemoveLikedMember(ctx, dao.Redis, userId, blogId); err != nil {
				return utils.ErrorResult("取消点赞失败")
			}
			if err := dao.DecrementBlogLiked(ctx, blogId); err != nil {
				return utils.ErrorResult("更新点赞数失败")
			}
			return utils.SuccessResult("取消点赞成功")
		}
		if err := dao.IncrementBlogLiked(ctx, blogId); err != nil {
			return utils.ErrorResult("更新点赞数失败")
		}
		// 保存用户 id 到 redis 集合
		if err := dao.SaveLikedMember(ctx, dao.Redis, userId, blogId); err != nil {
			return utils.ErrorResult("保存点赞失败")
		}
		return utils.SuccessResult("点赞成功")
	}
	return utils.ErrorResult("点赞失败")
}

// GetBlogList 获取博客列表
func GetBlogList(ctx context.Context, page, size int) *utils.Result {
	offset := (page - 1) * size

	blogs, total, err := dao.GetBlogList(ctx, offset, size)
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
func GetBlogById(ctx context.Context, id uint, userId uint) *utils.Result {
	blog, err := dao.GetBlogByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.ErrorResult("博客不存在")
		}
		return utils.ErrorResult("查询失败")
	}

	// 检查是否点赞
	if err := isBlogLiked(ctx, blog, userId); err != nil {
		return utils.ErrorResult("检查点赞状态失败")
	}

	return utils.SuccessResultWithData(blog)
}

// GetHotBlogList 获取热门博客列表
func GetHotBlogList(ctx context.Context, page, size int, userId uint) *utils.Result {
	offset := (page - 1) * size

	blogs, total, err := dao.GetHotBlogList(ctx, offset, size)
	if err != nil {
		return utils.ErrorResult("查询失败")
	}

	// 检查是否点赞
	for _, blog := range blogs {
		if err := isBlogLiked(ctx, &blog, userId); err != nil {
			return utils.ErrorResult("检查点赞状态失败")
		}
	}

	return utils.SuccessResultWithData(map[string]interface{}{
		"list":  blogs,
		"total": total,
		"page":  page,
		"size":  size,
	})
}

// GetMyBlogList 获取我的博客列表
func GetMyBlogList(ctx context.Context, userId uint, page, size int) *utils.Result {
	offset := (page - 1) * size

	blogs, total, err := dao.GetMyBlogList(ctx, userId, offset, size)
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

func GetBlogLikes(ctx context.Context, blogId uint) *utils.Result {
	// 从 SortedSet 中获取点赞数最多的 k 个用户
	ids, err := dao.GetTopKBloglikedMember(ctx, dao.Redis, blogId, 5)
	if err != nil {
		return utils.ErrorResult("查询失败")
	}

	// 转换用户ID为uint类型
	userIds := make([]uint, 0, len(ids))
	for _, id := range ids {
		idUint, err := strconv.ParseUint(id, 10, 32)
		if err != nil {
			return utils.ErrorResult("转换用户ID失败")
		}
		userIds = append(userIds, uint(idUint))
	}

	// 根据用户 ID 获取用户信息
	users, err := dao.GetUserByIDs(userIds)
	if err != nil {
		return utils.ErrorResult("查询失败")
	}

	return utils.SuccessResultWithData(users)
}

func isBlogLiked(ctx context.Context, blog *models.Blog, userId uint) error {
	liked, err := dao.IsLikedMember(ctx, dao.Redis, userId, blog.ID)

	if liked {
		blog.IsLiked = true
	} else {
		blog.IsLiked = false
	}

	return err
}
