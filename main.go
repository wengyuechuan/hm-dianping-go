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
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

	// 初始化订单队列和worker，如果需要让后端自行进行阻塞队列的话，可以使用，现在的优化方案是使用redis的消息队列机制来进行
	// service.InitOrderQueue()

	// 初始化Redis Stream消费者
	if err := service.InitStreamConsumer(); err != nil {
		log.Fatalf("Failed to initialize stream consumer: %v", err)
	}

	// 设置路由
	r := router.SetupRouter()

	// 启动服务器
	port := config.GetConfig().Server.Port
	if port == "" {
		port = "8080"
	}
	
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// 在goroutine中启动服务器
	go func() {
		log.Printf("Server starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// 停止Stream消费者
	service.StopStreamConsumers()

	// 关闭HTTP服务器
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
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
