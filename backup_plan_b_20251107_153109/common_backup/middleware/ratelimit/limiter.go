package ratelimit

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

// Limiter 分布式限流器
type Limiter struct {
	redis  *redis.Client
	config *Config
}

// NewLimiter 创建限流器
func NewLimiter(redisClient *redis.Client, config *Config) *Limiter {
	return &Limiter{
		redis:  redisClient,
		config: config,
	}
}

// NewLimiterFromConfig 从配置创建限流器（会自动创建Redis客户端）
func NewLimiterFromConfig(config *Config) (*Limiter, error) {
	if config == nil {
		return nil, ErrInvalidConfig
	}

	// 创建Redis客户端
	rdb := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port),
		Password:     config.Redis.Password,
		DB:           config.Redis.DB,
		PoolSize:     config.Redis.PoolSize,
		MinIdleConns: config.Redis.MinIdleConns,
		MaxRetries:   config.Redis.MaxRetries,
		DialTimeout:  config.Redis.DialTimeout,
		ReadTimeout:  config.Redis.ReadTimeout,
		WriteTimeout: config.Redis.WriteTimeout,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRedisUnavailable, err)
	}

	return NewLimiter(rdb, config), nil
}

// Allow 检查请求是否允许通过
// key: 限流key，通常是IP地址
// path: 请求路径
// 返回: allowed(是否允许), remaining(剩余配额), resetTime(配额重置时间), error
func (l *Limiter) Allow(ctx context.Context, key, path string) (bool, int, time.Time, error) {
	if !l.config.Enabled {
		return true, -1, time.Time{}, nil
	}

	// 查找匹配的规则
	rule := l.findMatchingRule(path)
	if rule == nil {
		// 如果没有匹配的规则且没有默认规则，则允许通过
		if l.config.DefaultRule == nil {
			return true, -1, time.Time{}, nil
		}
		rule = l.config.DefaultRule
	}

	// 验证规则
	if err := rule.Validate(); err != nil {
		return false, 0, time.Time{}, err
	}

	// 分别检查每秒和每分钟的限制
	var allowed = true
	var remaining = -1
	var resetTime time.Time

	// 检查每秒限制
	if rule.LimitPerSecond > 0 {
		keyPerSec := fmt.Sprintf("ratelimit:sec:%s:%s", key, path)
		allowedSec, remainingSec, resetSec, err := l.allowWithTokenBucket(
			ctx,
			keyPerSec,
			rule.LimitPerSecond,
			rule.GetBurstSize(rule.LimitPerSecond),
			time.Second,
		)
		if err != nil {
			return false, 0, time.Time{}, err
		}
		if !allowedSec {
			return false, remainingSec, resetSec, ErrRateLimitExceeded
		}
		remaining = remainingSec
		resetTime = resetSec
	}

	// 检查每分钟限制
	if rule.LimitPerMinute > 0 {
		keyPerMin := fmt.Sprintf("ratelimit:min:%s:%s", key, path)
		allowedMin, remainingMin, resetMin, err := l.allowWithTokenBucket(
			ctx,
			keyPerMin,
			rule.LimitPerMinute,
			rule.GetBurstSize(rule.LimitPerMinute),
			time.Minute,
		)
		if err != nil {
			return false, 0, time.Time{}, err
		}
		if !allowedMin {
			return false, remainingMin, resetMin, ErrRateLimitExceeded
		}
		// 如果只配置了每分钟限制，则返回每分钟的剩余配额
		if rule.LimitPerSecond == 0 {
			remaining = remainingMin
			resetTime = resetMin
		}
	}

	return allowed, remaining, resetTime, nil
}

// allowWithTokenBucket 使用令牌桶算法检查是否允许请求
// 使用Redis + Lua脚本实现分布式令牌桶
func (l *Limiter) allowWithTokenBucket(
	ctx context.Context,
	key string,
	rate int, // 令牌产生速率（每period产生rate个令牌）
	burst int, // 桶容量
	period time.Duration, // 时间周期
) (bool, int, time.Time, error) {
	now := time.Now()

	// Lua脚本：原子性地实现令牌桶算法
	// KEYS[1]: Redis key
	// ARGV[1]: 桶容量(burst)
	// ARGV[2]: 令牌产生速率(rate)
	// ARGV[3]: 时间周期(秒)
	// ARGV[4]: 当前时间戳(秒)
	// ARGV[5]: 过期时间(秒)
	script := redis.NewScript(`
		local key = KEYS[1]
		local burst = tonumber(ARGV[1])
		local rate = tonumber(ARGV[2])
		local period = tonumber(ARGV[3])
		local now = tonumber(ARGV[4])
		local expire_time = tonumber(ARGV[5])

		-- 获取当前令牌数和上次更新时间
		local token_info = redis.call('HMGET', key, 'tokens', 'last_time')
		local tokens = tonumber(token_info[1])
		local last_time = tonumber(token_info[2])

		-- 如果是首次访问，初始化令牌桶
		if tokens == nil then
			tokens = burst
			last_time = now
		end

		-- 计算应该补充的令牌数
		local elapsed = now - last_time
		local new_tokens = tokens + (elapsed / period) * rate

		-- 令牌数不能超过桶容量
		if new_tokens > burst then
			new_tokens = burst
		end

		-- 尝试消费1个令牌
		local allowed = 0
		if new_tokens >= 1 then
			new_tokens = new_tokens - 1
			allowed = 1
		end

		-- 更新Redis中的数据
		redis.call('HMSET', key, 'tokens', new_tokens, 'last_time', now)
		redis.call('EXPIRE', key, expire_time)

		-- 返回: 是否允许(1/0), 剩余令牌数, 下次重置时间
		local reset_after = 0
		if new_tokens < burst then
			reset_after = math.ceil((burst - new_tokens) / rate * period)
		end

		return {allowed, math.floor(new_tokens), reset_after}
	`)

	// 执行Lua脚本
	expireTime := int(period.Seconds() * 2) // 过期时间设置为周期的2倍
	result, err := script.Run(
		ctx,
		l.redis,
		[]string{key},
		burst,
		rate,
		period.Seconds(),
		now.Unix(),
		expireTime,
	).Result()

	if err != nil {
		return false, 0, time.Time{}, fmt.Errorf("%w: %v", ErrRedisUnavailable, err)
	}

	// 解析结果
	resultSlice, ok := result.([]interface{})
	if !ok || len(resultSlice) != 3 {
		return false, 0, time.Time{}, fmt.Errorf("invalid lua script result")
	}

	allowed := resultSlice[0].(int64) == 1
	remaining := int(resultSlice[1].(int64))
	resetAfter := int(resultSlice[2].(int64))
	resetTime := now.Add(time.Duration(resetAfter) * time.Second)

	return allowed, remaining, resetTime, nil
}

// findMatchingRule 查找匹配的限流规则
// 支持通配符匹配，如 "/api/*", "/admin/user/*"
func (l *Limiter) findMatchingRule(path string) *RuleConfig {
	for i := range l.config.Rules {
		rule := &l.config.Rules[i]
		if l.pathMatch(rule.Path, path) {
			return rule
		}
	}
	return nil
}

// pathMatch 路径匹配，支持通配符
// pattern: 规则路径，支持 * 通配符
// path: 实际请求路径
func (l *Limiter) pathMatch(pattern, path string) bool {
	// 精确匹配
	if pattern == path {
		return true
	}

	// 通配符匹配
	if strings.Contains(pattern, "*") {
		matched, _ := filepath.Match(pattern, path)
		return matched
	}

	return false
}

// Close 关闭限流器，释放Redis连接
func (l *Limiter) Close() error {
	if l.redis != nil {
		return l.redis.Close()
	}
	return nil
}
