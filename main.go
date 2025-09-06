package main

import (
	"flag"
	"hm-dianping-go/config"
	"hm-dianping-go/dao"
	"hm-dianping-go/models"
	"hm-dianping-go/router"
	"log"
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