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
	follow, err := dao.GetFollowByUserAndTarget(ctx, userId, followUserId)
	if err == nil {
		// 合并接口，这里取消关注
		if err := dao.DeleteFollow(ctx, follow); err != nil {
			return utils.ErrorResult("取消关注失败")
		}
		// 从 redis 中删除关注关系
		if err := dao.RemoveFollowing(ctx, dao.Redis, userId, followUserId); err != nil {
			return utils.ErrorResult("取消关注失败")
		}
		return utils.SuccessResult("取消关注成功")
	}

	// 添加关注
	follow = &models.Follow{
		UserID:       userId,
		FollowUserID: followUserId,
	}

	if err := dao.CreateFollow(ctx, follow); err != nil {
		return utils.ErrorResult("关注失败")
	}

	// 使用 redis 存储关注的信息
	if err := dao.SetFollowing(ctx, dao.Redis, userId, followUserId); err != nil {
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

	// 从 redis 中删除关注关系
	if err := dao.RemoveFollowing(ctx, dao.Redis, userId, followUserId); err != nil {
		return utils.ErrorResult("取消关注失败")
	}
	if err := dao.DeleteFollow(ctx, follow); err != nil {
		return utils.ErrorResult("取消关注失败")
	}

	return utils.SuccessResult("取消关注成功")
}

// GetCommonFollows 获取共同关注
func GetCommonFollows(ctx context.Context, userId, targetUserId uint) *utils.Result {
	// 使用 redis 存储共同关注的用户ID
	commonFollowIds, err := dao.GetCommonFollows(ctx, dao.Redis, userId, targetUserId)
	if err != nil {
		return utils.ErrorResult("查询失败")
	}

	// 判断是否有共同关注
	if len(commonFollowIds) == 0 {
		return utils.SuccessResultWithData([]models.User{})
	}

	// 查询用户信息
	users, err := dao.GetUsersByIds(ctx, commonFollowIds)
	if err != nil {
		return utils.ErrorResult("查询用户信息失败")
	}

	return utils.SuccessResultWithData(users)
}
