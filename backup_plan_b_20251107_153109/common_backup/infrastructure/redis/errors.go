package redis

import "errors"

var (
	// ErrClientNotFound Redis客户端未找到
	ErrClientNotFound = errors.New("redis client not found")

	// ErrInvalidConfig 无效的Redis配置
	ErrInvalidConfig = errors.New("invalid redis config")

	// ErrConnectionFailed Redis连接失败
	ErrConnectionFailed = errors.New("redis connection failed")
)
