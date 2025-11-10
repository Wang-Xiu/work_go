package redis

import (
	"fmt"
	"time"
)

// Config Redis配置结构
type Config struct {
	Host         string        `yaml:"host" mapstructure:"host"`
	Port         int           `yaml:"port" mapstructure:"port"`
	Password     string        `yaml:"password" mapstructure:"password"` // 支持KMS加密
	DB           int           `yaml:"db" mapstructure:"db"`
	PoolSize     int           `yaml:"pool_size" mapstructure:"pool_size"`
	MinIdleConns int           `yaml:"min_idle_conns" mapstructure:"min_idle_conns"`
	MaxRetries   int           `yaml:"max_retries" mapstructure:"max_retries"`
	DialTimeout  time.Duration `yaml:"dial_timeout" mapstructure:"dial_timeout"`
	ReadTimeout  time.Duration `yaml:"read_timeout" mapstructure:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout" mapstructure:"write_timeout"`
}

// Validate 验证配置
func (c *Config) Validate() error {
	if c.Host == "" {
		return fmt.Errorf("redis host is required")
	}
	if c.Port <= 0 || c.Port > 65535 {
		return fmt.Errorf("invalid redis port: %d", c.Port)
	}
	if c.PoolSize <= 0 {
		c.PoolSize = 10 // 默认连接池大小
	}
	if c.DialTimeout == 0 {
		c.DialTimeout = 5 * time.Second
	}
	if c.ReadTimeout == 0 {
		c.ReadTimeout = 3 * time.Second
	}
	if c.WriteTimeout == 0 {
		c.WriteTimeout = 3 * time.Second
	}
	return nil
}

// GetAddress 获取Redis地址
func (c *Config) GetAddress() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
