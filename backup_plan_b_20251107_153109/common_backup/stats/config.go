package stats

import (
	"time"
)

// Config 统计模块配置
type Config struct {
	// Redis配置
	Redis RedisConfig `yaml:"redis"`

	// 是否启用统计功能
	Enabled bool `yaml:"enabled"`

	// 是否按路径统计（true=统计每个路径的PV/UV，false=只统计全站PV/UV）
	EnablePathStats bool `yaml:"enable_path_stats"`

	// 数据保留天数（默认90天）
	RetentionDays int `yaml:"retention_days"`

	// 排除的路径（这些路径不参与统计）
	// 例如：健康检查接口、静态资源等
	ExcludePaths []string `yaml:"exclude_paths"`
}

// RedisConfig Redis连接配置
type RedisConfig struct {
	Host     string        `yaml:"host"`
	Port     int           `yaml:"port"`
	Password string        `yaml:"password"` // 支持KMS加密
	DB       int           `yaml:"db"`
	PoolSize int           `yaml:"pool_size"`
	Timeout  time.Duration `yaml:"timeout"`
}

// DailyStats 每日统计数据
type DailyStats struct {
	Date string `json:"date"` // 日期，格式：YYYY-MM-DD
	PV   int64  `json:"pv"`   // 页面浏览量
	UV   int64  `json:"uv"`   // 独立访客数
}

// PathStats 路径统计数据
type PathStats struct {
	Path string `json:"path"` // 路径
	PV   int64  `json:"pv"`   // 该路径的PV
	UV   int64  `json:"uv"`   // 该路径的UV
}

// RangeStats 时间范围统计数据
type RangeStats struct {
	StartDate  string       `json:"start_date"`  // 开始日期
	EndDate    string       `json:"end_date"`    // 结束日期
	TotalPV    int64        `json:"total_pv"`    // 总PV
	TotalUV    int64        `json:"total_uv"`    // 总UV（注意：跨天的UV不能简单相加）
	DailyStats []DailyStats `json:"daily_stats"` // 每日明细
}

// GetRetentionDays 获取数据保留天数（有默认值）
func (c *Config) GetRetentionDays() int {
	if c.RetentionDays <= 0 {
		return 90 // 默认保留90天
	}
	return c.RetentionDays
}

// IsExcludedPath 检查路径是否在排除列表中
func (c *Config) IsExcludedPath(path string) bool {
	for _, excludePath := range c.ExcludePaths {
		if path == excludePath {
			return true
		}
	}
	return false
}
