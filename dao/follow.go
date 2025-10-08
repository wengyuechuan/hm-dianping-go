package dao

import (
	"context"
	"hm-dianping-go/models"
	"strconv"

	"github.com/go-redis/redis/v8"
)

// GetFollowByUserAndTarget 根据用户ID和目标用户ID查询关注关系
func GetFollowByUserAndTarget(ctx context.Context, userId, followUserId uint) (*models.Follow, error) {
	var follow models.Follow
	err := DB.WithContext(ctx).Where("user_id = ? AND follow_user_id = ?", userId, followUserId).First(&follow).Error
	if err != nil {
		return nil, err
	}
	return &follow, nil
}

// CreateFollow 创建关注关系
func CreateFollow(ctx context.Context, follow *models.Follow) error {
	return DB.WithContext(ctx).Create(follow).Error
}

// DeleteFollow 删除关注关系
func DeleteFollow(ctx context.Context, follow *models.Follow) error {
	return DB.WithContext(ctx).Delete(follow).Error
}

// GetFollowers 获取关注用户的人列表
func GetFollowers(ctx context.Context, userId uint) ([]models.User, error) {
	var users []models.User
	err := DB.WithContext(ctx).Table("tb_follow").
		Select("u.*").
		Joins("JOIN tb_user u ON tb_follow.user_id = u.id").
		Where("tb_follow.follow_user_id = ?", userId).
		Find(&users).Error

	return users, err
}

// GetCommonFollowIds 获取共同关注的用户ID列表
func GetCommonFollowIds(ctx context.Context, userId, targetUserId uint) ([]uint, error) {
	var commonFollowIds []uint
	err := DB.WithContext(ctx).Table("follows as f1").
		Select("f1.follow_user_id").
		Joins("INNER JOIN follows as f2 ON f1.follow_user_id = f2.follow_user_id").
		Where("f1.user_id = ? AND f2.user_id = ?", userId, targetUserId).
		Pluck("f1.follow_user_id", &commonFollowIds).Error

	return commonFollowIds, err
}

// GetUsersByIds 根据用户ID列表查询用户信息
func GetUsersByIds(ctx context.Context, userIds []uint) ([]models.User, error) {
	var users []models.User
	if len(userIds) == 0 {
		return users, nil
	}

	err := DB.WithContext(ctx).Where("id IN ?", userIds).Find(&users).Error
	return users, err
}

// IsFollowing 检查是否已关注
func IsFollowing(ctx context.Context, userId, followUserId uint) (bool, error) {
	var count int64
	err := DB.WithContext(ctx).Model(&models.Follow{}).
		Where("user_id = ? AND follow_user_id = ?", userId, followUserId).
		Count(&count).Error

	return count > 0, err
}

// GetFollowingList 获取用户关注的人列表
func GetFollowingList(ctx context.Context, userId uint, limit, offset int) ([]models.User, error) {
	var users []models.User
	err := DB.WithContext(ctx).Table("tb_follow").
		Select("u.*").
		Joins("JOIN tb_user u ON tb_follow.follow_user_id = u.id").
		Where("tb_follow.user_id = ?", userId).
		Limit(limit).
		Offset(offset).
		Find(&users).Error

	return users, err
}

// GetFollowersList 获取关注用户的人列表
func GetFollowersList(ctx context.Context, userId uint, limit, offset int) ([]models.User, error) {
	var users []models.User
	err := DB.WithContext(ctx).Table("tb_follow").
		Select("u.*").
		Joins("JOIN tb_user u ON tb_follow.user_id = u.id").
		Where("tb_follow.follow_user_id = ?", userId).
		Limit(limit).
		Offset(offset).
		Find(&users).Error

	return users, err
}

// GetFollowingCount 获取关注数量
func GetFollowingCount(ctx context.Context, userId uint) (int64, error) {
	var count int64
	err := DB.WithContext(ctx).Model(&models.Follow{}).
		Where("user_id = ?", userId).
		Count(&count).Error

	return count, err
}

// GetFollowersCount 获取粉丝数量
func GetFollowersCount(ctx context.Context, userId uint) (int64, error) {
	var count int64
	err := DB.WithContext(ctx).Model(&models.Follow{}).
		Where("follow_user_id = ?", userId).
		Count(&count).Error

	return count, err
}

// =========== redis 存储关注的信息，存储到一个 set 中
const (
	FollowKeyPrefix = "follow:"
)

func SetFollowing(ctx context.Context, rds *redis.Client, userId, followUserId uint) error {
	return rds.SAdd(ctx, FollowKeyPrefix+strconv.Itoa(int(userId)), strconv.Itoa(int(followUserId))).Err()
}

func RemoveFollowing(ctx context.Context, rds *redis.Client, userId, followUserId uint) error {
	return rds.SRem(ctx, FollowKeyPrefix+strconv.Itoa(int(userId)), strconv.Itoa(int(followUserId))).Err()
}

// 求解两个用户的共同关注
func GetCommonFollows(ctx context.Context, rds *redis.Client, userId, targetUserId uint) ([]uint, error) {
	var commonFollowIds []uint
	err := rds.SInter(ctx, FollowKeyPrefix+strconv.Itoa(int(userId)), FollowKeyPrefix+strconv.Itoa(int(targetUserId))).ScanSlice(&commonFollowIds)
	if err != nil {
		return nil, err
	}
	return commonFollowIds, nil
}
