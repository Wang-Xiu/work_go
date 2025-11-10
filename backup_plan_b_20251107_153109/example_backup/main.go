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

func main() {
	// 1. 初始化KMS管理器（用于解密配置文件中的敏感信息）
	// 生产环境应该使用真实的KMS Provider（阿里云KMS、腾讯云KMS等）
	// 这里使用Mock Provider仅作演示
	kmsProvider := kms.NewMockProvider()
	kmsManager := kms.NewManager(kmsProvider, "kms://")

	// 2. 加载限流配置（会自动解密KMS加密的字段）
	configLoader := config.NewLoader(kmsManager)
	var rateLimitConfig ratelimit.Config
	if err := configLoader.LoadFromFile(
		context.Background(),
		"config/ratelimit.yml",
		&rateLimitConfig,
	); err != nil {
		log.Fatalf("Failed to load ratelimit config: %v", err)
	}

	// 打印解密后的Redis密码（仅调试，生产环境禁止打印敏感信息）
	fmt.Printf("Redis password decrypted: %s\n", rateLimitConfig.Redis.Password)

	// 3. 创建限流器
	limiter, err := ratelimit.NewLimiterFromConfig(&rateLimitConfig)
	if err != nil {
		log.Fatalf("Failed to create rate limiter: %v", err)
	}
	defer limiter.Close()

	// 4. 创建Gin路由
	r := gin.Default()

	// 5. 应用限流中间件（全局限流，按IP）
	r.Use(ratelimit.Middleware(limiter))

	// 6. 定义业务接口
	// API接口
	r.GET("/api/users", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "获取用户列表成功",
			"users":   []string{"user1", "user2", "user3"},
		})
	})

	r.POST("/api/login", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "登录成功",
			"token":   "mock_token_123456",
		})
	})

	r.POST("/api/register", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "注册成功",
		})
	})

	// 管理后台接口
	r.GET("/admin/stats", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "统计数据",
			"total":   1000,
		})
	})

	// 回调接口
	r.POST("/callback/payment/alipay", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "支付回调处理成功",
		})
	})

	// 7. 启动服务
	fmt.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
