package utils

import (
	"context"
	"fmt"
	"hm-dianping-go/dao"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware 跨域中间件
func CORSMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
}

// LoggerMiddleware 日志中间件
func LoggerMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	})
}

// JWTMiddleware JWT认证中间件
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.GetHeader("Authorization")
		if authorization == "" {
			ErrorResponse(c, http.StatusUnauthorized, "请先登录")
			c.Abort()
			return
		}

		// 检查Bearer前缀
		if !strings.HasPrefix(authorization, "Bearer ") {
			ErrorResponse(c, http.StatusUnauthorized, "token格式错误")
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authorization, "Bearer ")
		claims, err := ParseToken(token)
		if err != nil {
			ErrorResponse(c, http.StatusUnauthorized, "token无效")
			c.Abort()
			return
		}

		// 将用户ID存储到上下文中
		c.Set("userID", claims.UserID)
		c.Next()
	}
}

// RecoveryMiddleware 恢复中间件
func RecoveryMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		ErrorResponse(c, http.StatusInternalServerError, "服务器内部错误")
	})
}

// UVStatMiddleware UV统计中间件，使用Redis HyperLogLog实现
func UVStatMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 先执行请求
		c.Next()
		
		// 请求处理完成后进行UV统计
		go func() {
			// 获取用户标识，优先使用用户ID，其次使用IP
			var userIdentifier string
			
			// 尝试从JWT中获取用户ID
			if userID, exists := c.Get("userID"); exists {
				userIdentifier = fmt.Sprintf("user:%v", userID)
			} else {
				// 使用客户端IP作为标识
				userIdentifier = fmt.Sprintf("ip:%s", c.ClientIP())
			}
			
			// 获取当前日期作为key的一部分
			today := time.Now().Format("2006-01-02")
			
			// 使用HyperLogLog记录UV
			uvKey := fmt.Sprintf("uv:daily:%s", today)
			
			// 异步记录到Redis，避免影响请求性能
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			
			if err := dao.Redis.PFAdd(ctx, uvKey, userIdentifier).Err(); err != nil {
				// 记录错误但不影响主流程
				fmt.Printf("UV统计记录失败: %v\n", err)
			}
			
			// 设置key的过期时间为7天，避免数据无限增长
			dao.Redis.Expire(ctx, uvKey, 7*24*time.Hour)
		}()
	}
}
