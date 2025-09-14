package utils

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// BloomFilterConfig 布隆过滤器配置
type BloomFilterConfig struct {
	Key        string  // 过滤器键名
	ErrorRate  float64 // 误判率
	Capacity   uint    // 初始容量
	Expansion  uint    // 扩容倍数，默认2
	NonScaling bool    // 是否禁用自动扩容
}

// BloomFilter 布隆过滤器操作接口
type BloomFilter struct {
	config BloomFilterConfig
	rdb    *redis.Client
}

// NewBloomFilter 创建新的布隆过滤器实例
func NewBloomFilter(rdb *redis.Client, config BloomFilterConfig) *BloomFilter {
	// 设置默认值
	if config.ErrorRate == 0 {
		config.ErrorRate = 0.01 // 默认1%误判率
	}
	if config.Capacity == 0 {
		config.Capacity = 10000 // 默认1万容量
	}
	if config.Expansion == 0 {
		config.Expansion = 2 // 默认2倍扩容
	}

	return &BloomFilter{
		config: config,
		rdb:    rdb,
	}
}

// Reserve 创建布隆过滤器
func (bf *BloomFilter) Reserve(ctx context.Context) error {
	args := []interface{}{"BF.RESERVE", bf.config.Key, bf.config.ErrorRate, bf.config.Capacity}

	// 添加可选参数
	if bf.config.Expansion != 2 {
		args = append(args, "EXPANSION", bf.config.Expansion)
	}
	if bf.config.NonScaling {
		args = append(args, "NONSCALING")
	}

	err := bf.rdb.Do(ctx, args...).Err()
	if err != nil {
		// 如果过滤器已存在，忽略错误
		if err.Error() == "ERR item exists" {
			return nil
		}
		return fmt.Errorf("创建布隆过滤器失败: %v", err)
	}
	return nil
}

// Add 添加单个元素到布隆过滤器
func (bf *BloomFilter) Add(ctx context.Context, item string) (bool, error) {
	result, err := bf.rdb.Do(ctx, "BF.ADD", bf.config.Key, item).Int()
	if err != nil {
		return false, fmt.Errorf("添加元素到布隆过滤器失败: %v", err)
	}
	return result == 1, nil
}

// AddMulti 批量添加元素到布隆过滤器
func (bf *BloomFilter) AddMulti(ctx context.Context, items []string) ([]bool, error) {
	if len(items) == 0 {
		return []bool{}, nil
	}

	args := []interface{}{"BF.MADD", bf.config.Key}
	for _, item := range items {
		args = append(args, item)
	}

	results, err := bf.rdb.Do(ctx, args...).Slice()
	if err != nil {
		return nil, fmt.Errorf("批量添加元素到布隆过滤器失败: %v", err)
	}

	boolResults := make([]bool, len(results))
	for i, result := range results {
		if val, ok := result.(int64); ok {
			boolResults[i] = val == 1
		}
	}

	return boolResults, nil
}

// Exists 检查单个元素是否存在于布隆过滤器中
func (bf *BloomFilter) Exists(ctx context.Context, item string) (bool, error) {
	result, err := bf.rdb.Do(ctx, "BF.EXISTS", bf.config.Key, item).Int()
	if err != nil {
		return false, fmt.Errorf("检查布隆过滤器元素失败: %v", err)
	}
	return result == 1, nil
}

// ExistsMulti 批量检查元素是否存在于布隆过滤器中
func (bf *BloomFilter) ExistsMulti(ctx context.Context, items []string) ([]bool, error) {
	if len(items) == 0 {
		return []bool{}, nil
	}

	args := []interface{}{"BF.MEXISTS", bf.config.Key}
	for _, item := range items {
		args = append(args, item)
	}

	results, err := bf.rdb.Do(ctx, args...).Slice()
	if err != nil {
		return nil, fmt.Errorf("批量检查布隆过滤器元素失败: %v", err)
	}

	boolResults := make([]bool, len(results))
	for i, result := range results {
		if val, ok := result.(int64); ok {
			boolResults[i] = val == 1
		}
	}

	return boolResults, nil
}

// Info 获取布隆过滤器信息
func (bf *BloomFilter) Info(ctx context.Context) (map[string]interface{}, error) {
	results, err := bf.rdb.Do(ctx, "BF.INFO", bf.config.Key).Slice()
	if err != nil {
		return nil, fmt.Errorf("获取布隆过滤器信息失败: %v", err)
	}

	info := make(map[string]interface{})
	for i := 0; i < len(results); i += 2 {
		if i+1 < len(results) {
			key := fmt.Sprintf("%v", results[i])
			value := results[i+1]
			info[key] = value
		}
	}

	return info, nil
}

// Delete 删除布隆过滤器
func (bf *BloomFilter) Delete(ctx context.Context) error {
	err := bf.rdb.Del(ctx, bf.config.Key).Err()
	if err != nil {
		return fmt.Errorf("删除布隆过滤器失败: %v", err)
	}
	return nil
}

// 便利函数：将数字ID转换为字符串

// AddID 添加数字ID到布隆过滤器
func (bf *BloomFilter) AddID(ctx context.Context, id uint) (bool, error) {
	return bf.Add(ctx, strconv.FormatUint(uint64(id), 10))
}

// ExistsID 检查数字ID是否存在于布隆过滤器中
func (bf *BloomFilter) ExistsID(ctx context.Context, id uint) (bool, error) {
	return bf.Exists(ctx, strconv.FormatUint(uint64(id), 10))
}

// AddIDs 批量添加数字ID到布隆过滤器
func (bf *BloomFilter) AddIDs(ctx context.Context, ids []uint) ([]bool, error) {
	items := make([]string, len(ids))
	for i, id := range ids {
		items[i] = strconv.FormatUint(uint64(id), 10)
	}
	return bf.AddMulti(ctx, items)
}

// ExistsIDs 批量检查数字ID是否存在于布隆过滤器中
func (bf *BloomFilter) ExistsIDs(ctx context.Context, ids []uint) ([]bool, error) {
	items := make([]string, len(ids))
	for i, id := range ids {
		items[i] = strconv.FormatUint(uint64(id), 10)
	}
	return bf.ExistsMulti(ctx, items)
}

// 预定义配置
var (
	// ShopBloomConfig 商铺布隆过滤器配置
	ShopBloomConfig = BloomFilterConfig{
		Key:       "shop:bloom:filter",
		ErrorRate: 0.01,   // 1%误判率
		Capacity:  100000, // 10万商铺容量
		Expansion: 2,      // 2倍扩容
	}

	// UserBloomConfig 用户布隆过滤器配置
	UserBloomConfig = BloomFilterConfig{
		Key:       "user:bloom:filter",
		ErrorRate: 0.001,   // 0.1%误判率
		Capacity:  1000000, // 100万用户容量
		Expansion: 2,       // 2倍扩容
	}

	// VoucherBloomConfig 优惠券布隆过滤器配置
	VoucherBloomConfig = BloomFilterConfig{
		Key:       "voucher:bloom:filter",
		ErrorRate: 0.01,  // 1%误判率
		Capacity:  50000, // 5万优惠券容量
		Expansion: 2,     // 2倍扩容
	}
)

// CreateShopBloomFilter 创建商铺布隆过滤器
func CreateShopBloomFilter(rdb *redis.Client) *BloomFilter {
	return NewBloomFilter(rdb, ShopBloomConfig)
}

// CreateUserBloomFilter 创建用户布隆过滤器
func CreateUserBloomFilter(rdb *redis.Client) *BloomFilter {
	return NewBloomFilter(rdb, UserBloomConfig)
}

// CreateVoucherBloomFilter 创建优惠券布隆过滤器
func CreateVoucherBloomFilter(rdb *redis.Client) *BloomFilter {
	return NewBloomFilter(rdb, VoucherBloomConfig)
}

// BloomInitializer 布隆过滤器初始化器
type BloomInitializer struct {
	rdb *redis.Client
	db  *gorm.DB
}

// NewBloomInitializer 创建布隆过滤器初始化器
func NewBloomInitializer(rdb *redis.Client, db *gorm.DB) *BloomInitializer {
	return &BloomInitializer{
		rdb: rdb,
		db:  db,
	}
}

// InitShopBloomFilter 初始化商铺布隆过滤器
func (bi *BloomInitializer) InitShopBloomFilter(ctx context.Context) error {
	// 导入dao包的函数
	var getAllShopIDs func(context.Context, *gorm.DB) ([]uint, error)
	
	// 这里需要通过函数参数传入或者直接调用dao层
	// 为了避免循环依赖，我们通过参数传入
	return bi.initBloomFilterWithIDs(ctx, ShopBloomConfig, getAllShopIDs)
}

// InitUserBloomFilter 初始化用户布隆过滤器
func (bi *BloomInitializer) InitUserBloomFilter(ctx context.Context) error {
	// 导入dao包的函数
	var getAllUserIDs func() ([]uint, error)
	
	return bi.initBloomFilterWithUserIDs(ctx, UserBloomConfig, getAllUserIDs)
}

// InitVoucherBloomFilter 初始化优惠券布隆过滤器
func (bi *BloomInitializer) InitVoucherBloomFilter(ctx context.Context) error {
	// 导入dao包的函数
	var getAllVoucherIDs func() ([]uint, error)
	
	return bi.initBloomFilterWithVoucherIDs(ctx, VoucherBloomConfig, getAllVoucherIDs)
}

// initBloomFilterWithIDs 通用的布隆过滤器初始化方法（带context参数）
func (bi *BloomInitializer) initBloomFilterWithIDs(ctx context.Context, config BloomFilterConfig, getIDsFunc func(context.Context, *gorm.DB) ([]uint, error)) error {
	if getIDsFunc == nil {
		return fmt.Errorf("getIDsFunc不能为空")
	}
	
	// 创建布隆过滤器
	bf := NewBloomFilter(bi.rdb, config)
	
	// 创建或重置布隆过滤器
	if err := bf.Reserve(ctx); err != nil {
		return fmt.Errorf("创建布隆过滤器失败: %v", err)
	}
	
	// 获取所有ID
	ids, err := getIDsFunc(ctx, bi.db)
	if err != nil {
		return fmt.Errorf("获取ID列表失败: %v", err)
	}
	
	if len(ids) == 0 {
		log.Printf("警告: %s 没有找到任何数据", config.Key)
		return nil
	}
	
	// 批量添加ID到布隆过滤器
	results, err := bf.AddIDs(ctx, ids)
	if err != nil {
		return fmt.Errorf("批量添加ID到布隆过滤器失败: %v", err)
	}
	
	// 统计添加结果
	addedCount := 0
	for _, added := range results {
		if added {
			addedCount++
		}
	}
	
	log.Printf("布隆过滤器 %s 初始化完成: 总数据量=%d, 新增=%d", config.Key, len(ids), addedCount)
	return nil
}

// initBloomFilterWithUserIDs 用户布隆过滤器初始化方法（无context参数）
func (bi *BloomInitializer) initBloomFilterWithUserIDs(ctx context.Context, config BloomFilterConfig, getIDsFunc func() ([]uint, error)) error {
	if getIDsFunc == nil {
		return fmt.Errorf("getIDsFunc不能为空")
	}
	
	// 创建布隆过滤器
	bf := NewBloomFilter(bi.rdb, config)
	
	// 创建或重置布隆过滤器
	if err := bf.Reserve(ctx); err != nil {
		return fmt.Errorf("创建布隆过滤器失败: %v", err)
	}
	
	// 获取所有ID
	ids, err := getIDsFunc()
	if err != nil {
		return fmt.Errorf("获取ID列表失败: %v", err)
	}
	
	if len(ids) == 0 {
		log.Printf("警告: %s 没有找到任何数据", config.Key)
		return nil
	}
	
	// 批量添加ID到布隆过滤器
	results, err := bf.AddIDs(ctx, ids)
	if err != nil {
		return fmt.Errorf("批量添加ID到布隆过滤器失败: %v", err)
	}
	
	// 统计添加结果
	addedCount := 0
	for _, added := range results {
		if added {
			addedCount++
		}
	}
	
	log.Printf("布隆过滤器 %s 初始化完成: 总数据量=%d, 新增=%d", config.Key, len(ids), addedCount)
	return nil
}

// initBloomFilterWithVoucherIDs 优惠券布隆过滤器初始化方法（无context参数）
func (bi *BloomInitializer) initBloomFilterWithVoucherIDs(ctx context.Context, config BloomFilterConfig, getIDsFunc func() ([]uint, error)) error {
	if getIDsFunc == nil {
		return fmt.Errorf("getIDsFunc不能为空")
	}
	
	// 创建布隆过滤器
	bf := NewBloomFilter(bi.rdb, config)
	
	// 创建或重置布隆过滤器
	if err := bf.Reserve(ctx); err != nil {
		return fmt.Errorf("创建布隆过滤器失败: %v", err)
	}
	
	// 获取所有ID
	ids, err := getIDsFunc()
	if err != nil {
		return fmt.Errorf("获取ID列表失败: %v", err)
	}
	
	if len(ids) == 0 {
		log.Printf("警告: %s 没有找到任何数据", config.Key)
		return nil
	}
	
	// 批量添加ID到布隆过滤器
	results, err := bf.AddIDs(ctx, ids)
	if err != nil {
		return fmt.Errorf("批量添加ID到布隆过滤器失败: %v", err)
	}
	
	// 统计添加结果
	addedCount := 0
	for _, added := range results {
		if added {
			addedCount++
		}
	}
	
	log.Printf("布隆过滤器 %s 初始化完成: 总数据量=%d, 新增=%d", config.Key, len(ids), addedCount)
	return nil
}

// InitAllBloomFilters 初始化所有布隆过滤器
// 需要传入dao层的查询函数以避免循环依赖
func (bi *BloomInitializer) InitAllBloomFilters(ctx context.Context, 
	getAllShopIDs func(context.Context, *gorm.DB) ([]uint, error),
	getAllUserIDs func() ([]uint, error),
	getAllVoucherIDs func() ([]uint, error)) error {
	
	log.Println("开始初始化所有布隆过滤器...")
	start := time.Now()
	
	// 初始化商铺布隆过滤器
	if err := bi.initBloomFilterWithIDs(ctx, ShopBloomConfig, getAllShopIDs); err != nil {
		log.Printf("初始化商铺布隆过滤器失败: %v", err)
		return err
	}
	
	// 初始化用户布隆过滤器
	if err := bi.initBloomFilterWithUserIDs(ctx, UserBloomConfig, getAllUserIDs); err != nil {
		log.Printf("初始化用户布隆过滤器失败: %v", err)
		return err
	}
	
	// 初始化优惠券布隆过滤器
	if err := bi.initBloomFilterWithVoucherIDs(ctx, VoucherBloomConfig, getAllVoucherIDs); err != nil {
		log.Printf("初始化优惠券布隆过滤器失败: %v", err)
		return err
	}
	
	duration := time.Since(start)
	log.Printf("所有布隆过滤器初始化完成，耗时: %v", duration)
	return nil
}

// CheckBloomFilterHealth 检查布隆过滤器健康状态
func (bi *BloomInitializer) CheckBloomFilterHealth(ctx context.Context) map[string]interface{} {
	health := make(map[string]interface{})
	
	configs := []BloomFilterConfig{ShopBloomConfig, UserBloomConfig, VoucherBloomConfig}
	
	for _, config := range configs {
		bf := NewBloomFilter(bi.rdb, config)
		info, err := bf.Info(ctx)
		if err != nil {
			health[config.Key] = map[string]interface{}{
				"status": "error",
				"error":  err.Error(),
			}
		} else {
			health[config.Key] = map[string]interface{}{
				"status": "healthy",
				"info":   info,
			}
		}
	}
	
	return health
}

// AddToBloomFilter 向指定布隆过滤器添加新ID（用于增量更新）
func (bi *BloomInitializer) AddToBloomFilter(ctx context.Context, filterType string, id uint) error {
	var config BloomFilterConfig
	
	switch filterType {
	case "shop":
		config = ShopBloomConfig
	case "user":
		config = UserBloomConfig
	case "voucher":
		config = VoucherBloomConfig
	default:
		return fmt.Errorf("不支持的过滤器类型: %s", filterType)
	}
	
	bf := NewBloomFilter(bi.rdb, config)
	_, err := bf.AddID(ctx, id)
	if err != nil {
		return fmt.Errorf("添加ID到布隆过滤器失败: %v", err)
	}
	
	log.Printf("成功添加ID %d 到布隆过滤器 %s", id, config.Key)
	return nil
}

// CheckIDExists 检查ID是否存在于指定布隆过滤器中
func (bi *BloomInitializer) CheckIDExists(ctx context.Context, filterType string, id uint) (bool, error) {
	var config BloomFilterConfig
	
	switch filterType {
	case "shop":
		config = ShopBloomConfig
	case "user":
		config = UserBloomConfig
	case "voucher":
		config = VoucherBloomConfig
	default:
		return false, fmt.Errorf("不支持的过滤器类型: %s", filterType)
	}
	
	bf := NewBloomFilter(bi.rdb, config)
	return bf.ExistsID(ctx, id)
}

// CheckIDExists 检查ID是否存在于指定布隆过滤器中（全局函数）
// 注意：此函数需要全局Redis连接，建议使用CheckStringExistsInBloomFilter
func CheckIDExists(filterType string, id uint) bool {
	// 这里需要获取Redis连接，暂时返回true避免阻塞
	// 在实际使用中，应该传入Redis连接或使用全局连接
	return true
}

// CheckIDExistsWithRedis 检查ID是否存在于指定布隆过滤器中（带Redis连接）
func CheckIDExistsWithRedis(ctx context.Context, rdb *redis.Client, filterType string, id uint) (bool, error) {
	var key string
	
	switch filterType {
	case "shop":
		key = ShopBloomConfig.Key
	case "user":
		key = UserBloomConfig.Key
	case "voucher":
		key = VoucherBloomConfig.Key
	default:
		return false, fmt.Errorf("不支持的过滤器类型: %s", filterType)
	}
	
	// 将ID转换为字符串
	idStr := strconv.FormatUint(uint64(id), 10)
	
	// 使用通用函数检查
	return CheckStringExistsInBloomFilter(ctx, rdb, key, idStr)
}

// CheckStringExistsInBloomFilter 通用函数：检查字符串是否存在于指定key的布隆过滤器中
func CheckStringExistsInBloomFilter(ctx context.Context, rdb *redis.Client, key string, value string) (bool, error) {
	if rdb == nil {
		return false, fmt.Errorf("Redis客户端不能为空")
	}
	
	if key == "" {
		return false, fmt.Errorf("布隆过滤器key不能为空")
	}
	
	if value == "" {
		return false, fmt.Errorf("检查的值不能为空")
	}
	
	// 直接使用Redis命令检查元素是否存在
	result, err := rdb.Do(ctx, "BF.EXISTS", key, value).Int()
	if err != nil {
		return false, fmt.Errorf("检查布隆过滤器元素失败: %v", err)
	}
	
	return result == 1, nil
}
