package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"working-project/common/kms"
	"working-project/common/middleware/ratelimit"
	"working-project/config"

	"github.com/gin-gonic/gin"
)

// 高级用法示例：按用户ID限流
func main() {
	// 初始化KMS和配置（与basic示例相同）
	kmsProvider := kms.NewMockProvider()
	kmsManager := kms.NewManager(kmsProvider, "kms://")

	configLoader := config.NewLoader(kmsManager)
	var rateLimitConfig ratelimit.Config
	if err := configLoader.LoadFromFile(
		context.Background(),
		"config/ratelimit.yml",
		&rateLimitConfig,
	); err != nil {
		log.Fatalf("Failed to load ratelimit config: %v", err)
	}

	limiter, err := ratelimit.NewLimiterFromConfig(&rateLimitConfig)
	if err != nil {
		log.Fatalf("Failed to create rate limiter: %v", err)
	}
	defer limiter.Close()

	r := gin.Default()

	// 方式1: 全局按IP限流
	r.Use(ratelimit.Middleware(limiter))

	// 方式2: 特定路由组按用户ID限流
	apiGroup := r.Group("/api")
	{
		// 对需要登录的接口，使用自定义Key函数按用户ID限流
		apiGroup.Use(ratelimit.MiddlewareWithKeyFunc(limiter, func(c *gin.Context) string {
			// 从请求Header或Token中提取用户ID
			userID := c.GetHeader("X-User-ID")
			if userID == "" {
				// 如果没有用户ID，回退到IP限流
				return c.ClientIP()
			}
			// 使用 "user:" 前缀区分用户ID和IP
			return "user:" + userID
		}))

		apiGroup.GET("/profile", func(c *gin.Context) {
			userID := c.GetHeader("X-User-ID")
			c.JSON(http.StatusOK, gin.H{
				"message": "获取用户信息成功",
				"user_id": userID,
			})
		})

		apiGroup.POST("/update", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "更新成功",
			})
		})
	}

	// 方式3: 针对单个路由使用不同的限流策略
	// 例如：上传接口需要更严格的限流
	uploadGroup := r.Group("/upload")
	{
		// 可以为这个路由组单独创建限流器实例，使用不同的配置
		uploadGroup.POST("/image", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "上传成功",
			})
		})
	}

	// 方式4: 组合多个维度限流（IP + 用户ID）
	adminGroup := r.Group("/admin")
	{
		// 先按IP限流（防止单个IP暴力请求）
		adminGroup.Use(ratelimit.Middleware(limiter))

		// 再按用户ID限流（防止单个用户滥用）
		adminGroup.Use(ratelimit.MiddlewareWithKeyFunc(limiter, func(c *gin.Context) string {
			token := c.GetHeader("Authorization")
			if token == "" {
				return ""
			}
			// 实际项目中应该解析JWT token获取用户ID
			return "admin:" + token
		}))

		adminGroup.GET("/dashboard", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "管理后台数据",
			})
		})
	}

	fmt.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
