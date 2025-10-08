package dao

import (
	"context"
	"hm-dianping-go/models"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

// CreateBlog 创建博客
func CreateBlog(ctx context.Context, blog *models.Blog) error {
	return DB.WithContext(ctx).Create(blog).Error
}

// GetBlogByID 根据ID获取博客
func GetBlogByID(ctx context.Context, id uint) (*models.Blog, error) {
	var blog models.Blog
	err := DB.WithContext(ctx).First(&blog, id).Error
	if err != nil {
		return nil, err
	}
	return &blog, nil
}

// GetBlogList 获取博客列表（分页）
func GetBlogList(ctx context.Context, offset, limit int) ([]models.Blog, int64, error) {
	var blogs []models.Blog
	var total int64

	// 获取总数
	if err := DB.WithContext(ctx).Model(&models.Blog{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	err := DB.WithContext(ctx).Offset(offset).Limit(limit).Order("created_at desc").Find(&blogs).Error
	return blogs, total, err
}

// GetHotBlogList 获取热门博客列表（按点赞数排序）
func GetHotBlogList(ctx context.Context, offset, limit int) ([]models.Blog, int64, error) {
	var blogs []models.Blog
	var total int64

	// 获取总数
	if err := DB.WithContext(ctx).Model(&models.Blog{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询，按点赞数排序
	err := DB.WithContext(ctx).Offset(offset).Limit(limit).Order("liked desc, created_at desc").Find(&blogs).Error
	return blogs, total, err
}

// GetMyBlogList 获取用户的博客列表
func GetMyBlogList(ctx context.Context, userID uint, offset, limit int) ([]models.Blog, int64, error) {
	var blogs []models.Blog
	var total int64

	// 获取总数
	if err := DB.WithContext(ctx).Model(&models.Blog{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	err := DB.WithContext(ctx).Where("user_id = ?", userID).Offset(offset).Limit(limit).Order("created_at desc").Find(&blogs).Error
	return blogs, total, err
}

// GetBlogLike 检查用户是否已点赞博客
func GetBlogLike(ctx context.Context, userID, blogID uint) (*models.BlogLike, error) {
	var like models.BlogLike
	err := DB.WithContext(ctx).Where("user_id = ? AND blog_id = ?", userID, blogID).First(&like).Error
	if err != nil {
		return nil, err
	}
	return &like, nil
}

// CreateBlogLike 创建博客点赞
func CreateBlogLike(ctx context.Context, like *models.BlogLike) error {
	return DB.WithContext(ctx).Create(like).Error
}

// DeleteBlogLike 删除博客点赞
func DeleteBlogLike(ctx context.Context, like *models.BlogLike) error {
	return DB.WithContext(ctx).Delete(like).Error
}

// IncrementBlogLiked 增加博客点赞数
func IncrementBlogLiked(ctx context.Context, blogID uint) error {
	return DB.WithContext(ctx).Model(&models.Blog{}).Where("id = ?", blogID).UpdateColumn("liked", DB.Raw("liked + 1")).Error
}

// DecrementBlogLiked 减少博客点赞数
func DecrementBlogLiked(ctx context.Context, blogID uint) error {
	return DB.WithContext(ctx).Model(&models.Blog{}).Where("id = ?", blogID).UpdateColumn("liked", DB.Raw("liked - 1")).Error
}

// GetBlogByIDs 根据博客 ID 获取博客详情
func GetBlogByIDs(ctx context.Context, blogIDs []uint) ([]models.Blog, error) {
	var blogs []models.Blog
	err := DB.WithContext(ctx).Where("id IN ?", blogIDs).Find(&blogs).Error
	return blogs, err
}

// ======= redis 相关操作 =========

const (
	// 博客点赞集合的键名格式：blog_like:%d
	blogLikeKey = "blog:liked:"
	feedKey     = "feed:"
)

// // IsLikedMember 检查用户是否已经点赞博客
// func IsLikedMember(ctx context.Context, rds *redis.Client, userID, blogID uint) (bool, error) {
// 	// 使用集合进行判断，判断用户是否在点赞集合中
// 	return rds.SIsMember(ctx, blogLikeKey+strconv.Itoa(int(blogID)), strconv.Itoa(int(userID))).Result()
// }

// // RemoveLikedMember 从博客点赞集合中移除用户
// func RemoveLikedMember(ctx context.Context, rds *redis.Client, userID, blogID uint) error {
// 	// 从集合中移除用户
// 	return rds.SRem(ctx, blogLikeKey+strconv.Itoa(int(blogID)), strconv.Itoa(int(userID))).Err()
// }

// // SaveLikedMember 保存用户到博客点赞集合
// func SaveLikedMember(ctx context.Context, rds *redis.Client, userID, blogID uint) error {
// 	// 保存用户 id 到 redis 集合
// 	return rds.SAdd(ctx, blogLikeKey+strconv.Itoa(int(blogID)), strconv.Itoa(int(userID))).Err()
// }

// 使用 SortedSet 对代码进行改造
func IsLikedMember(ctx context.Context, rds *redis.Client, userID, blogID uint) (bool, error) {
	// 使用 SortedSet 检查用户是否在点赞集合中，使用查分数的方式，判断是否大于 0
	count, err := rds.ZCount(ctx, blogLikeKey+strconv.Itoa(int(blogID)), strconv.Itoa(int(userID)), strconv.Itoa(int(userID))).Result()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func RemoveLikedMember(ctx context.Context, rds *redis.Client, userID, blogID uint) error {
	// 从 SortedSet 中移除用户
	return rds.ZRem(ctx, blogLikeKey+strconv.Itoa(int(blogID)), strconv.Itoa(int(userID))).Err()
}

func SaveLikedMember(ctx context.Context, rds *redis.Client, userID, blogID uint) error {
	// 向 SortedSet 中添加用户，分数设为当前时间戳
	return rds.ZAdd(ctx, blogLikeKey+strconv.Itoa(int(blogID)), &redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: strconv.Itoa(int(userID)),
	}).Err()
}

func GetTopKBloglikedMember(ctx context.Context, rds *redis.Client, blogID uint, k int) ([]string, error) {
	// 从 SortedSet 中获取点赞数最多的 k 个用户
	return rds.ZRevRange(ctx, blogLikeKey+strconv.Itoa(int(blogID)), 0, int64(k-1)).Result()
}

func FeedToUserRedis(ctx context.Context, rds *redis.Client, userID uint, blogID uint) error {
	// 向用户的 feed 队列中插入id，使用 zset 存储，分数设为当前时间戳
	return rds.ZAdd(ctx, feedKey+strconv.Itoa(int(userID)), &redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: strconv.Itoa(int(blogID)),
	}).Err()
}

func GetFeedFromUserRedis(ctx context.Context, rds *redis.Client, userID uint, lastId, offset, count int) ([]uint, int64, int, error) {
	// 从用户的 feed 队列中获取博客 ID
	result, err := rds.ZRevRangeByScoreWithScores(ctx, feedKey+strconv.Itoa(int(userID)), &redis.ZRangeBy{
		Min:    "-inf",
		Max:    strconv.Itoa(int(lastId)),
		Offset: int64(offset),
		Count:  int64(count),
	}).Result()
	if err != nil {
		return nil, 0, 0, err
	}

	if len(result) == 0 {
		return nil, 0, 0, nil
	}

	var blogIds []uint
	minTime := result[len(result)-1].Score
	offset = 0
	for _, z := range result {
		blogId, _ := strconv.Atoi(z.Member.(string))
		blogIds = append(blogIds, uint(blogId))

		if z.Score == minTime {
			offset++
		}
	}

	return blogIds, int64(minTime), int(offset), nil
}
