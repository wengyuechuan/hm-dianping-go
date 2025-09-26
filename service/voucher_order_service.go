package service

import (
	"context"
	"fmt"
	"hm-dianping-go/dao"
	"hm-dianping-go/models"
	"hm-dianping-go/utils"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

// SeckillVoucher 秒杀优惠券（使用乐观锁）
// func SeckillVoucher(ctx context.Context, userId, voucherId uint) *utils.Result {
// 	// 1. 检查秒杀券是否存在
// 	seckillVoucher, err := dao.GetSeckillVoucherByID(voucherId)
// 	if err != nil {
// 		log.Printf("查询秒杀券失败: %v", err)
// 		return utils.ErrorResult("秒杀券不存在")
// 	}

// 	// 2. 检查秒杀时间
// 	now := time.Now()
// 	if now.Before(seckillVoucher.BeginTime) {
// 		return utils.ErrorResult("秒杀尚未开始")
// 	}
// 	if now.After(seckillVoucher.EndTime) {
// 		return utils.ErrorResult("秒杀已结束")
// 	}

// 	// 3. 检查库存
// 	if seckillVoucher.Stock <= 0 {
// 		return utils.ErrorResult("库存不足")
// 	}

// 	// 4. 检查用户是否已经购买过该秒杀券（一人一单限制）
// 	exists, err := dao.CheckSeckillVoucherOrderExists(ctx, dao.DB, userId, voucherId)
// 	if err != nil {
// 		log.Printf("检查秒杀券订单是否存在失败: %v", err)
// 		return utils.ErrorResult("系统错误")
// 	}
// 	if exists {
// 		return utils.ErrorResult("不能重复购买")
// 	}

// 	// 5. 使用乐观锁重试机制进行库存扣减和订单创建
// 	const maxRetries = 3
// 	for i := 0; i < maxRetries; i++ {
// 		// 开始事务
// 		tx := dao.DB.Begin()
// 		if tx.Error != nil {
// 			log.Printf("开始事务失败: %v", tx.Error)
// 			return utils.ErrorResult("系统错误")
// 		}

// 		// 扣减库存（乐观锁CAS操作）
// 		err = dao.UpdateSeckillVoucherStock(voucherId, 1)
// 		if err != nil {
// 			tx.Rollback()
// 			if errors.Is(err, gorm.ErrRecordNotFound) {
// 				// 库存不足或并发冲突，重试
// 				if i == maxRetries-1 {
// 					return utils.ErrorResult("库存不足")
// 				}
// 				// 短暂等待后重试
// 				time.Sleep(time.Duration(i+1) * 10 * time.Millisecond)
// 				continue
// 			}
// 			log.Printf("扣减库存失败: %v", err)
// 			return utils.ErrorResult("系统错误")
// 		}

// 		// 6. 创建秒杀券订单（使用DAO层函数）
// 		now = time.Now()
// 		order := &models.VoucherOrder{
// 			UserID:      userId,
// 			VoucherID:   voucherId,
// 			PayType:     1,
// 			Status:      1,
// 			CreateTime:  &now,
// 			VoucherType: 2, // 秒杀券类型
// 		}

// 		err = dao.CreateVoucherOrder(ctx, tx, order) // 这里创建的逻辑，主要依赖于mysql的约束
// 		if err != nil {
// 			tx.Rollback()
// 			log.Printf("创建订单失败: %v", err)
// 			return utils.ErrorResult("创建订单失败")
// 		}

// 		// 7. 提交事务
// 		if err := tx.Commit().Error; err != nil {
// 			tx.Rollback()
// 			log.Printf("提交事务失败: %v", err)
// 			return utils.ErrorResult("系统错误")
// 		}

// 		// 8. 成功，返回订单ID
// 		return utils.SuccessResultWithData(order.ID)
// 	}

// 	// 重试次数用完，返回失败
// 	return utils.ErrorResult("服务繁忙，请稍后重试")
// }

/*
需要注意的是，对于上面的方案，
1. 秒杀券的库存扣减和订单创建是一个原子操作，通过事务确保数据一致性。
2. 普通券的库存扣减和订单创建也是一个原子操作，通过事务确保数据一致性。
3. 秒杀券的唯一索引确保了每个用户只能购买一次秒杀券。
4. 普通券的索引确保了用户可以购买多次普通券。
// 在models/voucher_order.go中
type VoucherOrder struct {
    ID         uint           `gorm:"primarykey" json:"id"`
    CreatedAt  time.Time      `json:"createdAt"`
    UpdatedAt  time.Time      `json:"updatedAt"`
    DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
    UserID     uint           `json:"userId"`
    VoucherID  uint           `json:"voucherId"`
    PayType    int            `json:"payType"`
    Status     int            `json:"status"`
    CreateTime *time.Time     `json:"createTime"`
    PayTime    *time.Time     `json:"payTime"`
    UseTime    *time.Time     `json:"useTime"`
    RefundTime *time.Time     `json:"refundTime"`
    UpdateTime *time.Time     `json:"updateTime"`

    // 新增字段：券类型标识
    VoucherType int `gorm:"index" json:"voucherType"` // 1:普通券 2:秒杀券
}

// 添加复合索引，但不设为唯一
func (VoucherOrder) TableName() string {
    return "tb_voucher_order"
}

-- 只对秒杀券创建唯一约束
CREATE UNIQUE INDEX uk_seckill_user_voucher
ON tb_voucher_order (user_id, voucher_id)
WHERE voucher_type = 2;

-- 普通券只创建普通索引用于查询优化
CREATE INDEX idx_normal_user_voucher
ON tb_voucher_order (user_id, voucher_id, voucher_type);

*/

// SeckillVoucher 秒杀优惠券
// func SeckillVoucher(ctx context.Context, userId, voucherId uint) *utils.Result {
// 	// 从文件当中加载脚本
// 	script, err := os.ReadFile("script/seckill.lua")
// 	if err != nil {
// 		log.Printf("读取秒杀脚本失败: %v", err)
// 		return utils.ErrorResult("系统错误")
// 	}
// 	scriptStr := string(script)

// 	// 1. 执行Lua脚本
// 	result := dao.Redis.Eval(ctx, scriptStr, []string{}, strconv.Itoa(int(voucherId)), strconv.Itoa(int(userId)))
// 	if result.Err() != nil {
// 		log.Printf("执行秒杀脚本失败: %v", result.Err())
// 		return utils.ErrorResult("系统错误")
// 	}

// 	// 2. 判断结果是否为 0，0的时候有资格完成
// 	r, err := result.Int()
// 	if err != nil {
// 		log.Printf("获取秒杀脚本返回值失败: %v", err)
// 		return utils.ErrorResult("系统错误")
// 	}
// 	if r != 0 {
// 		if r == 1 {
// 			return utils.ErrorResult("库存不足")
// 		}
// 		return utils.ErrorResult("不能重复购买")
// 	}

// 	// 3. 有购买资格，将订单信息保存到阻塞队列
// 	err = AddOrderToQueue(userId, voucherId)
// 	if err != nil {
// 		log.Printf("订单入队失败: userId=%d, voucherId=%d, error=%v", userId, voucherId, err)
// 		return utils.ErrorResult("系统繁忙，请稍后重试")
// 	}

// 	// 4. 返回订单ID（这里可以生成一个临时ID或者返回成功信息）
// 	return utils.SuccessResultWithData("秒杀成功，订单处理中...")
// }

// // VoucherOrderInfo 订单信息结构体，用于阻塞队列
// type VoucherOrderInfo struct {
// 	UserID    uint `json:"userId"`
// 	VoucherID uint `json:"voucherId"`
// }

// // 全局阻塞队列和相关变量
// var (
// 	orderQueue  chan VoucherOrderInfo // 订单队列
// 	queueOnce   sync.Once             // 确保队列只初始化一次
// 	workerCount = 5                   // worker数量
// 	queueSize   = 1000                // 队列大小
// )

// // InitOrderQueue 初始化订单队列和worker
// func InitOrderQueue() {
// 	queueOnce.Do(func() {
// 		orderQueue = make(chan VoucherOrderInfo, queueSize)

// 		// 启动多个worker goroutine处理订单
// 		for i := 0; i < workerCount; i++ {
// 			go orderWorker(i)
// 		}

// 		log.Printf("订单队列初始化完成，队列大小: %d, worker数量: %d", queueSize, workerCount)
// 	})
// }

// // orderWorker 订单处理worker
// func orderWorker(workerID int) {
// 	log.Printf("订单处理worker %d 启动", workerID)

// 	for orderInfo := range orderQueue {
// 		err := processOrder(orderInfo)
// 		if err != nil {
// 			log.Printf("Worker %d 处理订单失败: userId=%d, voucherId=%d, error=%v",
// 				workerID, orderInfo.UserID, orderInfo.VoucherID, err)
// 			// 这里可以添加重试逻辑或者将失败的订单放入死信队列
// 		} else {
// 			log.Printf("Worker %d 成功处理订单: userId=%d, voucherId=%d",
// 				workerID, orderInfo.UserID, orderInfo.VoucherID)
// 		}
// 	}
// }

// // processOrder 处理单个订单
// func processOrder(orderInfo VoucherOrderInfo) error {
// 	// 开始事务
// 	tx := dao.DB.Begin()
// 	if tx.Error != nil {
// 		return tx.Error
// 	}

// 	defer func() {
// 		if r := recover(); r != nil {
// 			tx.Rollback()
// 		}
// 	}()

// 	// 创建订单
// 	now := time.Now()
// 	order := &models.VoucherOrder{
// 		UserID:      orderInfo.UserID,
// 		VoucherID:   orderInfo.VoucherID,
// 		PayType:     1,
// 		Status:      1,
// 		CreateTime:  &now,
// 		VoucherType: 2, // 秒杀券类型
// 	}

// 	// 创建订单记录
// 	err := dao.CreateVoucherOrder(context.Background(), tx, order)
// 	if err != nil {
// 		tx.Rollback()
// 		return err
// 	}

// 	// 提交事务
// 	if err := tx.Commit().Error; err != nil {
// 		tx.Rollback()
// 		return err
// 	}

// 	return nil
// }

// // AddOrderToQueue 将订单添加到队列
// func AddOrderToQueue(userID, voucherID uint) error {
// 	orderInfo := VoucherOrderInfo{
// 		UserID:    userID,
// 		VoucherID: voucherID,
// 	}

// 	select {
// 	case orderQueue <- orderInfo:
// 		return nil
// 	default:
// 		return fmt.Errorf("订单队列已满")
// 	}
// }

// SeckillVoucher 秒杀优惠券
func SeckillVoucher(ctx context.Context, userId, voucherId uint) *utils.Result {
	// 从文件当中加载脚本
	script, err := os.ReadFile("script/seckill.lua")
	if err != nil {
		log.Printf("读取秒杀脚本失败: %v", err)
		return utils.ErrorResult("系统错误")
	}
	scriptStr := string(script)

	// 1. 执行Lua脚本
	result := dao.Redis.Eval(ctx, scriptStr, []string{}, strconv.Itoa(int(voucherId)), strconv.Itoa(int(userId)))
	if result.Err() != nil {
		log.Printf("执行秒杀脚本失败: %v", result.Err())
		return utils.ErrorResult("系统错误")
	}

	// 2. 判断结果是否为 0，0的时候有资格完成
	r, err := result.Int()
	if err != nil {
		log.Printf("获取秒杀脚本返回值失败: %v", err)
		return utils.ErrorResult("系统错误")
	}
	if r != 0 {
		if r == 1 {
			return utils.ErrorResult("库存不足")
		}
		return utils.ErrorResult("不能重复购买")
	}

	// 3. 已经加入到消息队列了

	// 4. 返回订单ID（这里可以生成一个临时ID或者返回成功信息）
	return utils.SuccessResultWithData("秒杀成功，订单处理中...")
}

// StreamOrderInfo Redis Stream中的订单信息结构体
type StreamOrderInfo struct {
	UserID    string `json:"userId"`
	VoucherID string `json:"voucherId"`
	OrderID   string `json:"id"`
}

// Stream消费者相关配置
var (
	streamKey     = "stream.orders"     // Stream名称
	groupName     = "order-group"       // 消费者组名称
	consumerCount = 3                   // 消费者数量
	streamOnce    sync.Once             // 确保Stream只初始化一次
	stopChan      = make(chan struct{}) // 停止信号
	wg            sync.WaitGroup        // 等待组，用于优雅关闭
)

// InitStreamConsumer 初始化Redis Stream消费者
func InitStreamConsumer() error {
	var initErr error
	streamOnce.Do(func() {
		ctx := context.Background()

		// 1. 检查Stream是否存在，如果不存在则创建
		exists, err := checkStreamExists(ctx, streamKey)
		if err != nil {
			initErr = fmt.Errorf("检查Stream失败: %v", err)
			return
		}

		if !exists {
			// 创建一个空的Stream（通过添加临时消息然后删除）
			result := dao.Redis.XAdd(ctx, &redis.XAddArgs{
				Stream: streamKey,
				ID:     "*",
				Values: map[string]interface{}{"init": "temp"},
			})
			if result.Err() != nil {
				initErr = fmt.Errorf("创建Stream失败: %v", result.Err())
				return
			}
			// 删除临时消息
			dao.Redis.XDel(ctx, streamKey, result.Val())
		}

		// 2. 创建消费者组（如果不存在）
		err = dao.Redis.XGroupCreateMkStream(ctx, streamKey, groupName, "0").Err()
		if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
			initErr = fmt.Errorf("创建消费者组失败: %v", err)
			return
		}

		// 3. 启动消费者
		for i := 0; i < consumerCount; i++ {
			consumerName := fmt.Sprintf("consumer-%d", i)
			wg.Add(1)
			go streamConsumer(consumerName, i)
		}

		log.Printf("Redis Stream消费者初始化完成，Stream: %s, 消费者组: %s, 消费者数量: %d",
			streamKey, groupName, consumerCount)
	})

	return initErr
}

// checkStreamExists 检查Stream是否存在
func checkStreamExists(ctx context.Context, streamKey string) (bool, error) {
	result := dao.Redis.Exists(ctx, streamKey)
	if result.Err() != nil {
		return false, result.Err()
	}
	return result.Val() > 0, nil
}

// streamConsumer Stream消费者worker
func streamConsumer(consumerName string, workerID int) {
	defer wg.Done()

	log.Printf("Stream消费者 %s (Worker %d) 启动", consumerName, workerID)

	ctx := context.Background()

	for {
		select {
		case <-stopChan:
			log.Printf("Stream消费者 %s (Worker %d) 收到停止信号，正在退出", consumerName, workerID)
			return
		default:
			// 从Stream中读取消息
			messages, err := readStreamMessages(ctx, consumerName)
			if err != nil {
				log.Printf("消费者 %s 读取消息失败: %v", consumerName, err)
				time.Sleep(time.Second * 2) // 出错时等待2秒再重试
				continue
			}

			// 处理每条消息
			for _, msg := range messages {
				err := processStreamMessage(ctx, msg, consumerName)
				if err != nil {
					log.Printf("消费者 %s 处理消息失败: msgID=%s, error=%v",
						consumerName, msg.ID, err)
					// 这里可以添加重试逻辑或将失败消息放入死信队列
				} else {
					log.Printf("消费者 %s 成功处理消息: msgID=%s", consumerName, msg.ID)
					// 确认消息已处理
					dao.Redis.XAck(ctx, streamKey, groupName, msg.ID)
				}
			}

			// 如果没有消息，短暂休眠
			if len(messages) == 0 {
				time.Sleep(time.Millisecond * 100)
			}
		}
	}
}

// readStreamMessages 从Stream中读取消息
func readStreamMessages(ctx context.Context, consumerName string) ([]redis.XMessage, error) {
	// 首先尝试读取pending消息（之前未确认的消息）
	pendingResult := dao.Redis.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    groupName,
		Consumer: consumerName,
		Streams:  []string{streamKey, "0"}, // "0"表示读取pending消息
		Count:    10,
		Block:    0, // 不阻塞
	})

	if pendingResult.Err() == nil && len(pendingResult.Val()) > 0 && len(pendingResult.Val()[0].Messages) > 0 {
		return pendingResult.Val()[0].Messages, nil
	}

	// 如果没有pending消息，读取新消息
	result := dao.Redis.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    groupName,
		Consumer: consumerName,
		Streams:  []string{streamKey, ">"}, // ">"表示读取新消息
		Count:    10,
		Block:    time.Second * 1, // 阻塞1秒
	})

	if result.Err() != nil {
		return nil, result.Err()
	}

	if len(result.Val()) > 0 && len(result.Val()[0].Messages) > 0 {
		return result.Val()[0].Messages, nil
	}

	return []redis.XMessage{}, nil
}

// processStreamMessage 处理单条Stream消息
func processStreamMessage(ctx context.Context, msg redis.XMessage, consumerName string) error {
	// 解析消息内容
	orderInfo, err := parseOrderMessage(msg)
	if err != nil {
		return fmt.Errorf("解析消息失败: %v", err)
	}

	// 转换字符串ID为uint
	userID, err := strconv.ParseUint(orderInfo.UserID, 10, 32)
	if err != nil {
		return fmt.Errorf("解析用户ID失败: %v", err)
	}

	voucherID, err := strconv.ParseUint(orderInfo.VoucherID, 10, 32)
	if err != nil {
		return fmt.Errorf("解析优惠券ID失败: %v", err)
	}

	// 处理订单
	return processStreamOrder(ctx, uint(userID), uint(voucherID), orderInfo.OrderID)
}

// parseOrderMessage 解析订单消息
func parseOrderMessage(msg redis.XMessage) (*StreamOrderInfo, error) {
	orderInfo := &StreamOrderInfo{}

	// 从消息中提取字段
	if userID, ok := msg.Values["userId"].(string); ok {
		orderInfo.UserID = userID
	} else {
		return nil, fmt.Errorf("消息中缺少userId字段")
	}

	if voucherID, ok := msg.Values["voucherId"].(string); ok {
		orderInfo.VoucherID = voucherID
	} else {
		return nil, fmt.Errorf("消息中缺少voucherId字段")
	}

	if orderID, ok := msg.Values["id"].(string); ok {
		orderInfo.OrderID = orderID
	} else {
		return nil, fmt.Errorf("消息中缺少id字段")
	}

	return orderInfo, nil
}

// processStreamOrder 处理Stream中的订单
func processStreamOrder(ctx context.Context, userID, voucherID uint, orderID string) error {
	// 开始数据库事务
	tx := dao.DB.Begin()
	if tx.Error != nil {
		return fmt.Errorf("开始事务失败: %v", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("订单处理发生panic: %v", r)
		}
	}()

	// 创建订单
	now := time.Now()
	order := &models.VoucherOrder{
		UserID:      userID,
		VoucherID:   voucherID,
		PayType:     1,
		Status:      1,
		CreateTime:  &now,
		VoucherType: 2, // 秒杀券类型
	}

	// 如果Lua脚本提供了订单ID，可以使用它
	if orderID != "" {
		// 这里可以根据需要设置订单ID或其他字段
		// 注意：GORM的ID字段通常是自增的，需要根据实际情况处理
	}

	// 创建订单记录
	err := dao.CreateVoucherOrder(ctx, tx, order)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("创建订单失败: %v", err)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("提交事务失败: %v", err)
	}

	log.Printf("成功创建订单: userID=%d, voucherID=%d, orderID=%d",
		userID, voucherID, order.ID)

	return nil
}

// StopStreamConsumers 停止所有Stream消费者（用于优雅关闭）
func StopStreamConsumers() {
	log.Println("正在停止Stream消费者...")
	close(stopChan)
	wg.Wait()
	log.Println("所有Stream消费者已停止")
}

// GetStreamInfo 获取Stream状态信息（用于监控）
func GetStreamInfo() (map[string]interface{}, error) {
	ctx := context.Background()

	// 获取Stream基本信息
	streamInfo, err := dao.Redis.XInfoStream(ctx, streamKey).Result()
	if err != nil {
		return nil, fmt.Errorf("获取Stream信息失败: %v", err)
	}

	// 获取消费者组信息
	groupInfo, err := dao.Redis.XInfoGroups(ctx, streamKey).Result()
	if err != nil {
		return nil, fmt.Errorf("获取消费者组信息失败: %v", err)
	}

	// 获取消费者信息
	consumerInfo, err := dao.Redis.XInfoConsumers(ctx, streamKey, groupName).Result()
	if err != nil {
		return nil, fmt.Errorf("获取消费者信息失败: %v", err)
	}

	return map[string]interface{}{
		"stream":    streamInfo,
		"groups":    groupInfo,
		"consumers": consumerInfo,
	}, nil
}
