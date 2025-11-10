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
		// ========================================

		// 示例代码框架（需要你完善）：
		/*
			// 获取访客ID
			visitorID := getVisitorID(c)

			// 获取路径
			path := c.Request.URL.Path

			// 记录统计（异步或同步？）
			if err := tracker.Track(c.Request.Context(), visitorID, path); err != nil {
				// 统计失败只记录日志，不影响业务
				fmt.Printf("[Stats] Track failed: %v\n", err)
			}
		*/

		c.Next()
	}
}

// ========================================
// 辅助函数
// ========================================

// getVisitorID 获取访客唯一标识
// 这是一个示例实现，你可以根据实际需求修改
func getVisitorID(c *gin.Context) string {
	// TODO: 实现访客ID获取逻辑
	//
	// 方案1: 优先使用用户ID（如果已登录）
	//   if userID := c.GetString("user_id"); userID != "" {
	//       return "user:" + userID
	//   }
	//
	// 方案2: 使用IP地址
	//   return getClientIP(c)
	//
	// 方案3: 使用设备指纹
	//   return generateFingerprint(c)

	// ========================================
	// 👇 在这里实现你的代码
	// ========================================

	// 临时实现：使用IP地址
	return getClientIP(c)
}

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

// generateFingerprint 生成设备指纹（高级功能）
// 结合多个因素生成唯一标识，降低重复计数
//
// 注意: 这只是示例，生产环境需要更复杂的算法
func generateFingerprint(c *gin.Context) string {
	// TODO: 实现设备指纹生成
	//
	// 可以考虑的因素:
	//   - IP地址
	//   - User-Agent
	//   - Accept-Language
	//   - Accept-Encoding
	//   - Screen Resolution (需要前端配合)
	//   - Timezone
	//
	// 方法:
	//   1. 拼接这些字段
	//   2. 计算哈希值（如MD5、SHA256）
	//   3. 作为唯一标识

	// ========================================
	// 👇 在这里实现你的代码
	// ========================================

	ip := getClientIP(c)
	userAgent := c.GetHeader("User-Agent")
	acceptLang := c.GetHeader("Accept-Language")

	// 简单拼接作为指纹
	fingerprint := fmt.Sprintf("%s|%s|%s", ip, userAgent, acceptLang)

	// TODO: 这里应该计算哈希值，避免指纹太长
	// import "crypto/md5"
	// hash := md5.Sum([]byte(fingerprint))
	// return hex.EncodeToString(hash[:])

	return fingerprint
}

// MiddlewareWithCustomID 自定义访客ID的中间件
// 允许调用方自定义如何获取访客ID
//
// 使用方式:
//
//	r.Use(stats.MiddlewareWithCustomID(tracker, func(c *gin.Context) string {
//	    // 自定义逻辑
//	    if userID := c.GetString("user_id"); userID != "" {
//	        return userID
//	    }
//	    return c.ClientIP()
//	}))
func MiddlewareWithCustomID(tracker *Tracker, idFunc func(*gin.Context) string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !tracker.config.Enabled {
			c.Next()
			return
		}

		// TODO: 实现自定义ID的中间件逻辑
		//
		// 步骤1: 调用idFunc获取访客ID
		// 步骤2: 获取路径
		// 步骤3: 调用Track记录
		// 步骤4: 继续处理请求

		// ========================================
		// 👇 在这里实现你的代码
		// ========================================

		c.Next()
	}
}
