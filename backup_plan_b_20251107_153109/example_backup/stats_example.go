package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"working-project/common/kms"
	"working-project/common/middleware/stats"
	"working-project/config"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. 初始化KMS管理器
	kmsProvider := kms.NewMockProvider()
	kmsManager := kms.NewManager(kmsProvider, "kms://")

	// 2. 加载统计配置
	configLoader := config.NewLoader(kmsManager)
	var statsConfig stats.Config
	if err := configLoader.LoadFromFile(
		context.Background(),
		"config/stats.yml",
		&statsConfig,
	); err != nil {
		log.Fatalf("Failed to load stats config: %v", err)
	}

	// 3. 创建统计追踪器
	tracker, err := stats.NewTrackerFromConfig(&statsConfig)
	if err != nil {
		log.Fatalf("Failed to create tracker: %v", err)
	}
	defer tracker.Close()

	// 4. 创建Gin路由
	r := gin.Default()

	// 5. 应用统计中间件（自动统计所有请求）
	r.Use(stats.Middleware(tracker))

	// 6. 定义业务接口
	r.GET("/api/users", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"users": []string{"Alice", "Bob", "Charlie"},
		})
	})

	r.GET("/api/posts", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"posts": []string{"Post 1", "Post 2"},
		})
	})

	// 7. 提供统计查询接口
	// 获取今日统计
	r.GET("/admin/stats/today", func(c *gin.Context) {
		today := time.Now().Format("2006-01-02")
		dailyStats, err := tracker.GetDailyStats(c.Request.Context(), today)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, dailyStats)
	})

	// 获取指定日期的统计
	r.GET("/admin/stats/daily/:date", func(c *gin.Context) {
		date := c.Param("date")
		dailyStats, err := tracker.GetDailyStats(c.Request.Context(), date)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, dailyStats)
	})

	// 获取时间范围统计
	r.GET("/admin/stats/range", func(c *gin.Context) {
		startDate := c.Query("start")
		endDate := c.Query("end")

		if startDate == "" || endDate == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "start and end date required"})
			return
		}

		rangeStats, err := tracker.GetRangeStats(c.Request.Context(), startDate, endDate)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, rangeStats)
	})

	// 获取路径统计（如果启用了EnablePathStats）
	r.GET("/admin/stats/paths/:date", func(c *gin.Context) {
		date := c.Param("date")
		pathStats, err := tracker.GetPathStats(c.Request.Context(), date)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, pathStats)
	})

	// 8. 启动定时任务：每天凌晨清理过期数据
	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			if err := tracker.CleanExpiredData(ctx); err != nil {
				log.Printf("Failed to clean expired data: %v", err)
			}
			cancel()
		}
	}()

	// 9. 启动服务
	log.Println("Server starting on :8080")
	log.Println("统计接口:")
	log.Println("  - GET /admin/stats/today")
	log.Println("  - GET /admin/stats/daily/:date")
	log.Println("  - GET /admin/stats/range?start=2024-01-01&end=2024-01-31")
	log.Println("  - GET /admin/stats/paths/:date")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
