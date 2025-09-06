package service

import (
	"hm-dianping-go/dao"
	"hm-dianping-go/models"
	"hm-dianping-go/utils"
)

// Follow 关注用户
func Follow(userId, followUserId uint) *utils.Result {
	// 检查是否已经关注
	var existingFollow models.Follow
	err := dao.DB.Where("user_id = ? AND follow_user_id = ?", userId, followUserId).First(&existingFollow).Error
	if err == nil {
		return utils.ErrorResult("已经关注过了")
	}

	// 添加关注
	follow := models.Follow{
		UserID:       userId,
		FollowUserID: followUserId,
	}

	if err := dao.DB.Create(&follow).Error; err != nil {
		return utils.ErrorResult("关注失败")
	}

	return utils.SuccessResult("关注成功")
}

// Unfollow 取消关注
func Unfollow(userId, followUserId uint) *utils.Result {
	var follow models.Follow
	err := dao.DB.Where("user_id = ? AND follow_user_id = ?", userId, followUserId).First(&follow).Error
	if err != nil {
		return utils.ErrorResult("未关注该用户")
	}

	if err := dao.DB.Delete(&follow).Error; err != nil {
		return utils.ErrorResult("取消关注失败")
	}

	return utils.SuccessResult("取消关注成功")
}

// GetCommonFollows 获取共同关注
func GetCommonFollows(userId, targetUserId uint) *utils.Result {
	// 查询共同关注的用户ID
	var commonFollowIds []uint
	err := dao.DB.Table("follows as f1").
		Select("f1.follow_user_id").
		Joins("INNER JOIN follows as f2 ON f1.follow_user_id = f2.follow_user_id").
		Where("f1.user_id = ? AND f2.user_id = ?", userId, targetUserId).
		Pluck("f1.follow_user_id", &commonFollowIds).Error

	if err != nil {
		return utils.ErrorResult("查询失败")
	}

	// 查询用户信息
	var users []models.User
	if len(commonFollowIds) > 0 {
		dao.DB.Where("id IN ?", commonFollowIds).Find(&users)
	}

	return utils.SuccessResultWithData(users)
}