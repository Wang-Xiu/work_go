package ratelimit

// RuleConfig 单个路径的限流规则配置
type RuleConfig struct {
	Path           string `yaml:"path"`             // 路径，支持通配符，如 "/api/*", "/admin/user/*"
	LimitPerSecond int    `yaml:"limit_per_second"` // 每秒最大请求数，0表示不限制
	LimitPerMinute int    `yaml:"limit_per_minute"` // 每分钟最大请求数，0表示不限制
	BurstSize      int    `yaml:"burst_size"`       // 突发流量大小（令牌桶容量），默认为limit的2倍
}

// Config 限流器配置
type Config struct {
	// Redis连接名称（引用infrastructure/redis.Manager中的连接）
	RedisName string `yaml:"redis_name"`

	// 限流规则列表
	Rules []RuleConfig `yaml:"rules"`

	// 默认规则（当路径不匹配任何规则时使用）
	DefaultRule *RuleConfig `yaml:"default_rule,omitempty"`

	// 是否启用限流，false则限流器不生效
	Enabled bool `yaml:"enabled"`
}

// GetBurstSize 获取突发流量大小，如果未配置则返回limit的2倍
func (r *RuleConfig) GetBurstSize(limit int) int {
	if r.BurstSize > 0 {
		return r.BurstSize
	}
	// 默认允许2倍的突发流量
	if limit > 0 {
		return limit * 2
	}
	return 0
}

// Validate 验证规则配置是否合法
func (r *RuleConfig) Validate() error {
	if r.Path == "" {
		return ErrInvalidConfig
	}
	if r.LimitPerSecond < 0 || r.LimitPerMinute < 0 {
		return ErrInvalidConfig
	}
	if r.LimitPerSecond == 0 && r.LimitPerMinute == 0 {
		return ErrInvalidConfig
	}
	return nil
}
