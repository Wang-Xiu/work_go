package stats

import "errors"

// 统计相关错误定义
var (
	// ErrRedisUnavailable Redis不可用
	ErrRedisUnavailable = errors.New("redis unavailable")

	// ErrInvalidConfig 配置错误
	ErrInvalidConfig = errors.New("invalid stats config")

	// ErrInvalidDate 日期格式错误
	ErrInvalidDate = errors.New("invalid date format")

	// ErrStartDate 日期格式错误
	ErrStartDate = errors.New("the start time cannot be greater than the end time")
)
