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
		}

		// 关注相关路由
		followGroup := api.Group("/follow")
		{
			followGroup.POST("/:id", utils.JWTMiddleware(), handler.Follow)
			followGroup.DELETE("/:id", utils.JWTMiddleware(), handler.Unfollow)
			followGroup.GET("/common/:id", utils.JWTMiddleware(), handler.GetCommonFollows)
		}
	}

	// 健康检查
	r.GET("/health", handler.HealthCheck)

	return r
}