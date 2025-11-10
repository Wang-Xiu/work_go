package ratelimit

import "errors"

var (
	// ErrRateLimitExceeded 限流错误
	ErrRateLimitExceeded = errors.New("rate limit exceeded")

	// ErrInvalidConfig 配置错误
	ErrInvalidConfig = errors.New("invalid rate limit config")

	// ErrRedisUnavailable Redis不可用
	ErrRedisUnavailable = errors.New("redis unavailable")
)
