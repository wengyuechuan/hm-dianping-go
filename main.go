package main

import (
	"context"
	"flag"
	"hm-dianping-go/config"
	"hm-dianping-go/dao"
	"hm-dianping-go/models"
	"hm-dianping-go/router"
	"hm-dianping-go/service"
	"hm-dianping-go/utils"
	"log"
	"time"
)

func main() {
	// 解析命令行参数
	configPath := flag.String("config", "config/application.yaml", "Path to configuration file")
	flag.Parse()

	// 加载配置
	if err := config.LoadConfigFromFile(*configPath); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化数据库连接
	if err := dao.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 初始化Redis连接
	if err := dao.InitRedis(); err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}

	// 自动迁移数据库表
	if err := dao.DB.AutoMigrate(
		&models.User{},
		&models.Shop{},
		&models.ShopType{},
		&models.Voucher{},
		&models.VoucherOrder{},
		&models.Blog{},
		&models.Follow{},
		&models.BlogLike{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// 初始化布隆过滤器
	if err := initBloomFilters(); err != nil {
		log.Printf("Warning: Failed to initialize bloom filters: %v", err)
		// 布隆过滤器初始化失败不应该阻止服务启动，只记录警告
	}

	// 初始化订单队列和worker
	service.InitOrderQueue()

	// 设置路由
	r := router.SetupRouter()

	// 启动服务器
	port := config.GetConfig().Server.Port
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// initBloomFilters 初始化布隆过滤器
func initBloomFilters() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 创建布隆过滤器初始化器
	initializer := utils.NewBloomInitializer(dao.Redis, dao.DB)

	// 初始化所有布隆过滤器
	err := initializer.InitAllBloomFilters(
		ctx,
		dao.GetAllShopIDs,    // 商铺ID查询函数
		dao.GetAllUserIDs,    // 用户ID查询函数
		dao.GetAllVoucherIDs, // 优惠券ID查询函数
	)

	if err != nil {
		return err
	}

	// 检查布隆过滤器健康状态
	health := initializer.CheckBloomFilterHealth(ctx)
	log.Printf("布隆过滤器健康状态: %+v", health)

	return nil
}