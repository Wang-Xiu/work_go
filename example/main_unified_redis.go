package main

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"

	"working-project/common/infrastructure/redis"
	"working-project/common/kms"
	"working-project/common/middleware/ratelimit"
	"working-project/common/middleware/stats"
	"working-project/config"
)

// AppConfig 应用配置结构（对应config/app.yml）
type AppConfig struct {
	Redis      map[string]redis.Config `yaml:"redis"`
	Middleware struct {
		RateLimit ratelimit.Config `yaml:"ratelimit"`
		Stats     stats.Config     `yaml:"stats"`
	} `yaml:"middleware"`
}

func main() {
	ctx := context.Background()

	log.Println("========================================")
	log.Println("应用启动 - 统一Redis管理示例")
	log.Println("========================================")

	// ============================================================
	// 第1步：初始化KMS管理器
	// ============================================================
	log.Println("\n[1/6] 初始化KMS管理器...")

	// 使用MockProvider（仅用于开发测试）
	// 生产环境请替换为真实的KMS Provider（阿里云KMS、腾讯云KMS、AWS KMS等）
	kmsProvider := kms.NewMockProvider()
	kmsManager := kms.NewManager(kmsProvider, "kms://")
	log.Println("✅ KMS管理器初始化成功 (使用MockProvider)")

	// ============================================================
	// 第2步：加载配置（自动KMS解密）
	// ============================================================
	log.Println("\n[2/6] 加载配置文件...")

	configLoader := config.NewLoader(kmsManager)
	var appConfig AppConfig

	// 使用配置加载器，会自动解密 kms:// 开头的配置值
	if err := configLoader.LoadFromFile(ctx, "config/app.yml", &appConfig); err != nil {
		log.Fatalf("❌ 加载配置文件失败: %v", err)
	}
	log.Println("✅ 配置加载成功（已自动解密KMS加密字段）")

	// ============================================================
	// 第3步：初始化Redis Manager（全局单例）
	// ============================================================
	log.Println("\n[3/6] 初始化Redis连接池...")

	// 获取全局Redis Manager
	redisManager := redis.GetGlobalManager()

	// 注册所有Redis连接（这里是关键：统一管理所有Redis连接）
	for name, cfg := range appConfig.Redis {
		log.Printf("  - 注册Redis连接: %s", name)

		// 复制配置以避免指针问题
		redisCfg := cfg
		if err := redisManager.Register(name, &redisCfg); err != nil {
			log.Fatalf("❌ 注册Redis '%s' 失败: %v", name, err)
		}
		log.Printf("    ✅ 连接成功: %s:%d (DB: %d)", cfg.Host, cfg.Port, cfg.DB)
	}

	// 确保程序退出时关闭所有Redis连接
	defer func() {
		log.Println("\n关闭所有Redis连接...")
		if err := redisManager.CloseAll(); err != nil {
			log.Printf("⚠️  关闭Redis连接时出错: %v", err)
		} else {
			log.Println("✅ Redis连接已全部关闭")
		}
	}()

	// ============================================================
	// 第4步：创建限流器（使用统一的Redis Manager）
	// ============================================================
	log.Println("\n[4/6] 初始化限流中间件...")

	limiter, err := ratelimit.NewLimiterFromManager(redisManager, &appConfig.Middleware.RateLimit)
	if err != nil {
		log.Fatalf("❌ 创建限流器失败: %v", err)
	}
	log.Printf("✅ 限流器初始化成功 (Redis: %s)", appConfig.Middleware.RateLimit.RedisName)

	// ============================================================
	// 第5步：创建统计追踪器（使用统一的Redis Manager）
	// ============================================================
	log.Println("\n[5/6] 初始化统计中间件...")

	tracker, err := stats.NewTrackerFromManager(redisManager, &appConfig.Middleware.Stats)
	if err != nil {
		log.Fatalf("❌ 创建统计追踪器失败: %v", err)
	}
	log.Printf("✅ 统计追踪器初始化成功 (Redis: %s)", appConfig.Middleware.Stats.RedisName)

	// ============================================================
	// 第6步：初始化Gin路由并应用中间件
	// ============================================================
	log.Println("\n[6/6] 初始化HTTP服务器...")

	r := gin.Default()

	// 应用中间件（注意顺序：统计在前，限流在后）
	r.Use(ratelimit.Middleware(limiter))
	r.Use(stats.Middleware(tracker))

	// ============================================================
	// 定义业务路由
	// ============================================================

	// API路由
	api := r.Group("/api")
	{
		api.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "OK",
				"time":    time.Now().Format(time.RFC3339),
			})
		})

		api.POST("/login", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "登录成功"})
		})

		api.POST("/register", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "注册成功"})
		})
	}

	// Admin路由 - 统计查询接口
	admin := r.Group("/admin")
	{
		// 查询今日统计
		admin.GET("/stats/today", func(c *gin.Context) {
			today := time.Now().Format("2006-01-02")
			dailyStats, err := tracker.GetDailyStats(ctx, today)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, dailyStats)
		})

		// 查询指定日期统计
		admin.GET("/stats/daily/:date", func(c *gin.Context) {
			date := c.Param("date")
			dailyStats, err := tracker.GetDailyStats(ctx, date)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, dailyStats)
		})

		// 查询时间范围统计
		admin.GET("/stats/range", func(c *gin.Context) {
			startDate := c.Query("start")
			endDate := c.Query("end")

			if startDate == "" || endDate == "" {
				c.JSON(400, gin.H{"error": "缺少start或end参数"})
				return
			}

			rangeStats, err := tracker.GetRangeStats(ctx, startDate, endDate)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, rangeStats)
		})

		// 查询路径统计（如果启用）
		admin.GET("/stats/paths/:date", func(c *gin.Context) {
			date := c.Param("date")
			pathStats, err := tracker.GetPathStats(ctx, date)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, pathStats)
		})
	}

	// 健康检查接口
	r.GET("/health", func(c *gin.Context) {
		// 检查所有Redis连接的健康状态
		healthStatus := make(map[string]string)
		allHealthy := true

		for _, name := range redisManager.List() {
			client, err := redisManager.Get(name)
			if err != nil {
				healthStatus[name] = "error: " + err.Error()
				allHealthy = false
				continue
			}

			if err := client.Ping(ctx).Err(); err != nil {
				healthStatus[name] = "unhealthy: " + err.Error()
				allHealthy = false
			} else {
				healthStatus[name] = "healthy"
			}
		}

		status := "healthy"
		if !allHealthy {
			status = "unhealthy"
		}

		c.JSON(200, gin.H{
			"status": status,
			"redis":  healthStatus,
		})
	})

	// ============================================================
	// 启动HTTP服务器
	// ============================================================
	log.Println("\n========================================")
	log.Println("✅ 初始化完成！")
	log.Println("========================================")
	log.Println("\n🚀 HTTP服务器启动中...")
	log.Println("\n可用接口:")
	log.Println("  - GET  http://localhost:8080/api/test")
	log.Println("  - POST http://localhost:8080/api/login")
	log.Println("  - POST http://localhost:8080/api/register")
	log.Println("  - GET  http://localhost:8080/admin/stats/today")
	log.Println("  - GET  http://localhost:8080/admin/stats/daily/2025-01-07")
	log.Println("  - GET  http://localhost:8080/admin/stats/range?start=2025-01-01&end=2025-01-07")
	log.Println("  - GET  http://localhost:8080/admin/stats/paths/2025-01-07")
	log.Println("  - GET  http://localhost:8080/health")
	log.Println("\n按 Ctrl+C 停止服务器")
	log.Println("========================================")

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("❌ 启动服务器失败: %v", err)
	}
}
