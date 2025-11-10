package ratelimit

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Middleware 限流中间件
func Middleware(limiter *Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取客户端IP
		ip := getClientIP(c)
		if ip == "" {
			// 如果无法获取IP，为了安全起见，拒绝请求
			c.JSON(http.StatusForbidden, gin.H{
				"error": "unable to identify client",
			})
			c.Abort()
			return
		}

		// 获取请求路径
		path := c.Request.URL.Path

		// 检查是否允许通过
		allowed, remaining, resetTime, err := limiter.Allow(c.Request.Context(), ip, path)

		// 设置响应头，告知客户端限流状态
		if remaining >= 0 {
			c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
			c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", resetTime.Unix()))
		}

		// 处理错误
		if err != nil {
			if errors.Is(err, ErrRateLimitExceeded) {
				// 被限流，返回429
				c.Header("Retry-After", fmt.Sprintf("%d", resetTime.Unix()))
				c.JSON(http.StatusTooManyRequests, gin.H{
					"error":       "rate limit exceeded",
					"message":     "请求过于频繁，请稍后再试",
					"retry_after": resetTime.Unix(),
				})
				c.Abort()
				return
			}

			// 其他错误（如Redis连接失败）
			// 为了服务可用性，这里选择放行请求，但记录错误日志
			// 生产环境建议接入日志系统
			fmt.Printf("[RateLimit] Error: %v, IP: %s, Path: %s\n", err, ip, path)
			c.Next()
			return
		}

		if !allowed {
			// 被限流，返回429
			c.Header("Retry-After", fmt.Sprintf("%d", resetTime.Unix()))
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "rate limit exceeded",
				"message":     "请求过于频繁，请稍后再试",
				"retry_after": resetTime.Unix(),
			})
			c.Abort()
			return
		}

		// 允许通过
		c.Next()
	}
}

// getClientIP 获取客户端真实IP
// 优先从X-Forwarded-For, X-Real-IP等Header获取，最后才使用RemoteAddr
func getClientIP(c *gin.Context) string {
	// 1. 尝试从X-Forwarded-For获取（可能被伪造，需要配合可信代理列表使用）
	xForwardedFor := c.GetHeader("X-Forwarded-For")
	if xForwardedFor != "" {
		// X-Forwarded-For格式: client, proxy1, proxy2
		// 取第一个IP
		ips := strings.Split(xForwardedFor, ",")
		if len(ips) > 0 {
			ip := strings.TrimSpace(ips[0])
			if ip != "" {
				return ip
			}
		}
	}

	// 2. 尝试从X-Real-IP获取
	xRealIP := c.GetHeader("X-Real-IP")
	if xRealIP != "" {
		return strings.TrimSpace(xRealIP)
	}

	// 3. 使用RemoteAddr
	// RemoteAddr格式: IP:Port
	remoteAddr := c.Request.RemoteAddr
	if idx := strings.LastIndex(remoteAddr, ":"); idx != -1 {
		return remoteAddr[:idx]
	}

	return remoteAddr
}

// MiddlewareWithKeyFunc 自定义限流Key的中间件
// keyFunc: 自定义Key生成函数，例如可以根据用户ID、Token等维度限流
func MiddlewareWithKeyFunc(limiter *Limiter, keyFunc func(*gin.Context) string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取自定义Key
		key := keyFunc(c)
		if key == "" {
			// 如果无法获取Key，拒绝请求
			c.JSON(http.StatusForbidden, gin.H{
				"error": "unable to identify client",
			})
			c.Abort()
			return
		}

		// 获取请求路径
		path := c.Request.URL.Path

		// 检查是否允许通过
		allowed, remaining, resetTime, err := limiter.Allow(c.Request.Context(), key, path)

		// 设置响应头
		if remaining >= 0 {
			c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
			c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", resetTime.Unix()))
		}

		// 处理错误
		if err != nil {
			if errors.Is(err, ErrRateLimitExceeded) {
				c.Header("Retry-After", fmt.Sprintf("%d", resetTime.Unix()))
				c.JSON(http.StatusTooManyRequests, gin.H{
					"error":       "rate limit exceeded",
					"message":     "请求过于频繁，请稍后再试",
					"retry_after": resetTime.Unix(),
				})
				c.Abort()
				return
			}

			fmt.Printf("[RateLimit] Error: %v, Key: %s, Path: %s\n", err, key, path)
			c.Next()
			return
		}

		if !allowed {
			c.Header("Retry-After", fmt.Sprintf("%d", resetTime.Unix()))
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "rate limit exceeded",
				"message":     "请求过于频繁，请稍后再试",
				"retry_after": resetTime.Unix(),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
