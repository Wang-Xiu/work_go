package stats

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

// Middleware 统计中间件
// 自动统计每个请求的PV和UV
//
// 使用方式:
//
//	tracker := stats.NewTracker(redisClient, config)
//	r.Use(stats.Middleware(tracker))
func Middleware(tracker *Tracker) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 如果未启用统计，直接放行
		if !tracker.config.Enabled {
			c.Next()
			return
		}

		// TODO: 你需要实现中间件逻辑
		//
		// 步骤1: 获取访客唯一标识（visitorID）
		//   方案A: 使用IP地址
		//     - 从X-Forwarded-For或RemoteAddr获取
		//     - 简单但不够精确（同一NAT下的多个用户会被认为是一个）
		//
		//   方案B: 使用用户ID（如果已登录）
		//     - 从Token、Session或Cookie中提取用户ID
		//     - 精确但只能统计登录用户
		//
		//   方案C: 使用设备指纹
		//     - 结合IP、User-Agent、Accept-Language等生成唯一标识
		//     - 较为精确
		//
		//   推荐: 优先使用用户ID，未登录则使用IP
		//
		// 步骤2: 获取访问路径
		//   - c.Request.URL.Path
		//   - 注意：是否需要包含query参数？
		//     例如: /api/users 还是 /api/users?page=1
		//     建议：只用Path，不包含query参数
		//
		// 步骤3: 调用Track方法记录
		//   - tracker.Track(c.Request.Context(), visitorID, path)
		//   - 注意错误处理：统计失败不应影响业务请求
		//
		// 步骤4: 继续处理请求
		//   - c.Next()

		// ========================================
		// 👇 在这里实现你的代码
		//没uid的接口统一使用ip做标识，否则使用用户id做标识(暂定请求头Authorization为用户标识)
		var visitorID string
		urlPath := c.Request.URL.Path
		if urlPath == "/api/login" {
			visitorID = getClientIP(c)
		} else {
			visitorID = c.GetHeader("Authorization")
		}

		err := tracker.Track(c, visitorID, urlPath)
		if err != nil {
			fmt.Printf("[Stats] Track failed: %v\n", err)
		}

		// ========================================

		c.Next()
	}
}

// ========================================
// 辅助函数
// ========================================

// getClientIP 获取客户端IP
// 从X-Forwarded-For或RemoteAddr获取
func getClientIP(c *gin.Context) string {
	// 1. 尝试从X-Forwarded-For获取
	xff := c.GetHeader("X-Forwarded-For")
	if xff != "" {
		// 取第一个IP
		// 注意：这里应该配合可信代理验证使用（参考ratelimit模块）
		return xff
	}

	// 2. 尝试从X-Real-IP获取
	xRealIP := c.GetHeader("X-Real-IP")
	if xRealIP != "" {
		return xRealIP
	}

	// 3. 使用RemoteAddr
	return c.ClientIP()
}
