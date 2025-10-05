package service

import (
	"context"
	"hm-dianping-go/dao"
	"hm-dianping-go/models"
	"hm-dianping-go/utils"
)

// Follow 关注用户
func Follow(ctx context.Context, userId, followUserId uint) *utils.Result {
	// 检查是否已经关注
	_, err := dao.GetFollowByUserAndTarget(ctx, userId, followUserId)
	if err == nil {
		return utils.ErrorResult("已经关注过了")
	}

	// 添加关注
	follow := models.Follow{
		UserID:       userId,
		FollowUserID: followUserId,
	}

	if err := dao.CreateFollow(ctx, &follow); err != nil {
		return utils.ErrorResult("关注失败")
	}

	return utils.SuccessResult("关注成功")
}

// Unfollow 取消关注
func Unfollow(ctx context.Context, userId, followUserId uint) *utils.Result {
	follow, err := dao.GetFollowByUserAndTarget(ctx, userId, followUserId)
	if err != nil {
		return utils.ErrorResult("未关注该用户")
	}

	if err := dao.DeleteFollow(ctx, follow); err != nil {
		return utils.ErrorResult("取消关注失败")
	}

	return utils.SuccessResult("取消关注成功")
}

// GetCommonFollows 获取共同关注
func GetCommonFollows(ctx context.Context, userId, targetUserId uint) *utils.Result {
	// 查询共同关注的用户ID
	commonFollowIds, err := dao.GetCommonFollowIds(ctx, userId, targetUserId)
	if err != nil {
		return utils.ErrorResult("查询失败")
	}

	// 查询用户信息
	users, err := dao.GetUsersByIds(ctx, commonFollowIds)
	if err != nil {
		return utils.ErrorResult("查询用户信息失败")
	}

	return utils.SuccessResultWithData(users)
}
