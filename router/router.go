package router

import (
	"hm-dianping-go/handler"
	"hm-dianping-go/utils"

	"github.com/gin-gonic/gin"
)

// SetupRouter 设置路由
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 添加中间件
	r.Use(utils.CORSMiddleware())
	r.Use(utils.LoggerMiddleware())
	r.Use(utils.UVStatMiddleware()) // UV统计中间件

	// API路由组
	api := r.Group("/api")
	{
		// 用户相关路由
		userGroup := api.Group("/user")
		{
			userGroup.POST("/code", handler.SendCode)
			userGroup.POST("/register", handler.UserRegister)
			userGroup.POST("/login", handler.UserLogin)
			userGroup.POST("/logout", handler.UserLogout)
			userGroup.GET("/me", utils.JWTMiddleware(), handler.GetUserInfo)
			userGroup.PUT("/update", utils.JWTMiddleware(), handler.UpdateUserInfo)
			userGroup.POST("/sign", utils.JWTMiddleware(), handler.Sign) // 签到
		}

		// 商铺相关路由
		shopGroup := api.Group("/shop")
		{
			shopGroup.GET("/list", handler.GetShopList)
			shopGroup.GET("/:id", handler.GetShopById)
			shopGroup.GET("/of/type", handler.GetShopByType)
			shopGroup.GET("/of/name", handler.GetShopByName)
			shopGroup.POST("", handler.SaveShop)
			shopGroup.PUT("", handler.UpdateShop)
			shopGroup.GET("/:id/nearby", utils.JWTMiddleware(), handler.GetNearbyShops) // 获取某个商铺附近的商铺
		}

		// 商铺类型相关路由
		shopTypeGroup := api.Group("/shop-type")
		{
			shopTypeGroup.GET("/list", handler.GetShopTypeList)
		}

		// 优惠券相关路由
		voucherGroup := api.Group("/voucher")
		{
			voucherGroup.GET("/list/:shopId", handler.GetVoucherList)
			voucherGroup.POST("", handler.AddVoucher)
			voucherGroup.POST("/seckill", handler.AddSeckillVoucher)
			voucherGroup.GET("/seckill/:id", handler.GetSeckillVoucher)
		}

		// 优惠券订单相关路由
		voucherOrderGroup := api.Group("/voucher-order")
		{
			voucherOrderGroup.POST("/seckill/:id", utils.JWTMiddleware(), handler.SeckillVoucher)
		}

		// 博客相关路由
		blogGroup := api.Group("/blog")
		{
			blogGroup.POST("", utils.JWTMiddleware(), handler.CreateBlog)
			blogGroup.PUT("/like/:id", utils.JWTMiddleware(), handler.LikeBlog)
			blogGroup.GET("/hot", handler.GetHotBlogList)
			blogGroup.GET("/of/me", utils.JWTMiddleware(), handler.GetMyBlogList)
			blogGroup.GET("/:id", handler.GetBlogById)
			blogGroup.GET("/of/follow", utils.JWTMiddleware(), handler.GetBlogOfFollow)
		}

		// 关注相关路由
		followGroup := api.Group("/follow")
		{
			followGroup.POST("/:id", utils.JWTMiddleware(), handler.Follow)
			followGroup.DELETE("/:id", utils.JWTMiddleware(), handler.Unfollow)
			followGroup.GET("/common/:id", utils.JWTMiddleware(), handler.GetCommonFollows)
		}

		// 统计相关路由
		statGroup := api.Group("/stat")
		{
			statGroup.GET("/uv/today", handler.GetTodayUV)                    // 获取今日UV
			statGroup.GET("/uv/daily", handler.GetDailyUV)                   // 获取指定日期UV
			statGroup.GET("/uv/range", handler.GetUVRange)                   // 获取日期范围UV
			statGroup.GET("/uv/recent", handler.GetRecentUV)                 // 获取最近N天UV
			statGroup.GET("/uv/summary", handler.GetUVSummary)               // 获取UV统计摘要
		}
	}

	// 健康检查
	r.GET("/health", handler.HealthCheck)

	return r
}
