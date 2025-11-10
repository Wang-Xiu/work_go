package ratelimit

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
)

// setupTestRedis 创建测试用的Redis实例（使用miniredis模拟）
func setupTestRedis(t *testing.T) *redis.Client {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("expected no error on miniredis.Run, got %v", err)
	}

	t.Cleanup(func() {
		mr.Close()
	})

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	return client
}

// TestRuleConfig_Validate 测试规则配置验证
func TestRuleConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		rule    RuleConfig
		wantErr bool
	}{
		{
			name: "valid rule with per second limit",
			rule: RuleConfig{
				Path:           "/api/*",
				LimitPerSecond: 10,
				LimitPerMinute: 0,
			},
			wantErr: false,
		},
		{
			name: "valid rule with per minute limit",
			rule: RuleConfig{
				Path:           "/api/*",
				LimitPerSecond: 0,
				LimitPerMinute: 100,
			},
			wantErr: false,
		},
		{
			name: "valid rule with both limits",
			rule: RuleConfig{
				Path:           "/api/*",
				LimitPerSecond: 10,
				LimitPerMinute: 100,
			},
			wantErr: false,
		},
		{
			name: "invalid rule with empty path",
			rule: RuleConfig{
				Path:           "",
				LimitPerSecond: 10,
			},
			wantErr: true,
		},
		{
			name: "invalid rule with no limits",
			rule: RuleConfig{
				Path:           "/api/*",
				LimitPerSecond: 0,
				LimitPerMinute: 0,
			},
			wantErr: true,
		},
		{
			name: "invalid rule with negative limit",
			rule: RuleConfig{
				Path:           "/api/*",
				LimitPerSecond: -1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.rule.Validate()
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			}
		})
	}
}

// TestRuleConfig_GetBurstSize 测试突发流量大小计算
func TestRuleConfig_GetBurstSize(t *testing.T) {
	tests := []struct {
		name      string
		rule      RuleConfig
		limit     int
		wantBurst int
	}{
		{
			name: "custom burst size",
			rule: RuleConfig{
				BurstSize: 20,
			},
			limit:     10,
			wantBurst: 20,
		},
		{
			name:      "default burst size (2x limit)",
			rule:      RuleConfig{},
			limit:     10,
			wantBurst: 20,
		},
		{
			name:      "zero limit",
			rule:      RuleConfig{},
			limit:     0,
			wantBurst: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			burst := tt.rule.GetBurstSize(tt.limit)
			if cast.ToInt(burst) != cast.ToInt(tt.wantBurst) {
				t.Errorf("expected burst size %d, got %d", tt.wantBurst, burst)
			}
		})
	}
}

// TestLimiter_PathMatch 测试路径匹配
func TestLimiter_PathMatch(t *testing.T) {
	limiter := &Limiter{
		config: &Config{},
	}

	tests := []struct {
		name    string
		pattern string
		path    string
		want    bool
	}{
		{
			name:    "exact match",
			pattern: "/api/users",
			path:    "/api/users",
			want:    true,
		},
		{
			name:    "wildcard match",
			pattern: "/api/*",
			path:    "/api/users",
			want:    true,
		},
		{
			name:    "wildcard no match",
			pattern: "/api/*",
			path:    "/admin/users",
			want:    false,
		},
		{
			name:    "nested wildcard match",
			pattern: "/api/*/profile",
			path:    "/api/users/profile",
			want:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := limiter.pathMatch(tt.pattern, tt.path)
			if cast.ToBool(got) != cast.ToBool(tt.want) {
				t.Errorf("expected %v, got %v", tt.want, got)
			}
		})
	}
}

// TestLimiter_FindMatchingRule 测试查找匹配规则
func TestLimiter_FindMatchingRule(t *testing.T) {
	config := &Config{
		Rules: []RuleConfig{
			{
				Path:           "/api/login",
				LimitPerSecond: 1,
			},
			{
				Path:           "/api/*",
				LimitPerSecond: 10,
			},
			{
				Path:           "/admin/*",
				LimitPerSecond: 20,
			},
		},
	}

	limiter := &Limiter{config: config}

	tests := []struct {
		name         string
		path         string
		wantRuleIdx  int
		wantNotFound bool
	}{
		{
			name:        "exact match has priority",
			path:        "/api/login",
			wantRuleIdx: 0,
		},
		{
			name:        "wildcard match",
			path:        "/api/users",
			wantRuleIdx: 1,
		},
		{
			name:        "admin wildcard match",
			path:        "/admin/dashboard",
			wantRuleIdx: 2,
		},
		{
			name:         "no match",
			path:         "/public/info",
			wantNotFound: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := limiter.findMatchingRule(tt.path)
			if tt.wantNotFound {
				if rule != nil {
					t.Errorf("expected nil rule, got %v", rule)
				}
			} else {
				if rule == nil {
					t.Fatalf("expected non-nil rule, got nil")
				}
				if cast.ToString(rule.Path) != cast.ToString(config.Rules[tt.wantRuleIdx].Path) {
					t.Errorf("expected path %s, got %s", config.Rules[tt.wantRuleIdx].Path, rule.Path)
				}
			}
		})
	}
}

// TestLimiter_Allow 测试限流逻辑
func TestLimiter_Allow(t *testing.T) {
	t.Parallel()

	rdb := setupTestRedis(t)
	config := &Config{
		Enabled: true,
		Rules: []RuleConfig{
			{
				Path:           "/api/test",
				LimitPerSecond: 2,
				BurstSize:      2,
			},
		},
	}

	limiter := NewLimiter(rdb, config)
	ctx := context.Background()

	// 第1次请求：应该允许
	allowed, remaining, _, err := limiter.Allow(ctx, "192.168.1.1", "/api/test")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !cast.ToBool(allowed) {
		t.Errorf("expected allowed to be true, got false")
	}
	if cast.ToInt(remaining) != 1 {
		t.Errorf("expected remaining to be 1, got %d", remaining)
	}

	// 第2次请求：应该允许
	allowed, remaining, _, err = limiter.Allow(ctx, "192.168.1.1", "/api/test")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !cast.ToBool(allowed) {
		t.Errorf("expected allowed to be true, got false")
	}
	if cast.ToInt(remaining) != 0 {
		t.Errorf("expected remaining to be 0, got %d", remaining)
	}

	// 第3次请求：应该被限流
	allowed, remaining, resetTime, err := limiter.Allow(ctx, "192.168.1.1", "/api/test")
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if !errors.Is(err, ErrRateLimitExceeded) {
		t.Errorf("expected ErrRateLimitExceeded, got %v", err)
	}
	if cast.ToBool(allowed) {
		t.Errorf("expected allowed to be false, got true")
	}
	if cast.ToInt(remaining) != 0 {
		t.Errorf("expected remaining to be 0, got %d", remaining)
	}
	if !resetTime.After(time.Now()) {
		t.Errorf("expected resetTime to be in the future")
	}

	// 等待一段时间后，令牌应该恢复
	time.Sleep(1 * time.Second)
	allowed, _, _, err = limiter.Allow(ctx, "192.168.1.1", "/api/test")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !cast.ToBool(allowed) {
		t.Errorf("expected allowed to be true, got false")
	}
}

// TestLimiter_Allow_DifferentKeys 测试不同Key的限流独立性
func TestLimiter_Allow_DifferentKeys(t *testing.T) {
	t.Parallel()

	rdb := setupTestRedis(t)
	config := &Config{
		Enabled: true,
		Rules: []RuleConfig{
			{
				Path:           "/api/*",
				LimitPerSecond: 1,
				BurstSize:      1,
			},
		},
	}

	limiter := NewLimiter(rdb, config)
	ctx := context.Background()

	// IP1的第1次请求：应该允许
	allowed, _, _, err := limiter.Allow(ctx, "192.168.1.1", "/api/test")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !cast.ToBool(allowed) {
		t.Errorf("expected allowed to be true, got false")
	}

	// IP1的第2次请求：应该被限流
	allowed, _, _, err = limiter.Allow(ctx, "192.168.1.1", "/api/test")
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if cast.ToBool(allowed) {
		t.Errorf("expected allowed to be false, got true")
	}

	// IP2的第1次请求：应该允许（独立限流）
	allowed, _, _, err = limiter.Allow(ctx, "192.168.1.2", "/api/test")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !cast.ToBool(allowed) {
		t.Errorf("expected allowed to be true, got false")
	}
}

// TestLimiter_Allow_PerMinuteLimit 测试每分钟限流
func TestLimiter_Allow_PerMinuteLimit(t *testing.T) {
	t.Parallel()

	rdb := setupTestRedis(t)
	config := &Config{
		Enabled: true,
		Rules: []RuleConfig{
			{
				Path:           "/api/test",
				LimitPerSecond: 0, // 不限制每秒
				LimitPerMinute: 3,
				BurstSize:      3,
			},
		},
	}

	limiter := NewLimiter(rdb, config)
	ctx := context.Background()

	// 快速发送3次请求，都应该允许
	for i := 0; i < 3; i++ {
		allowed, _, _, err := limiter.Allow(ctx, "192.168.1.1", "/api/test")
		if err != nil {
			t.Fatalf("request %d: expected no error, got %v", i+1, err)
		}
		if !cast.ToBool(allowed) {
			t.Errorf("request %d should be allowed", i+1)
		}
	}

	// 第4次请求应该被限流
	allowed, _, _, err := limiter.Allow(ctx, "192.168.1.1", "/api/test")
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if cast.ToBool(allowed) {
		t.Errorf("expected allowed to be false, got true")
	}
}

// TestLimiter_Allow_Disabled 测试禁用限流
func TestLimiter_Allow_Disabled(t *testing.T) {
	t.Parallel()

	rdb := setupTestRedis(t)
	config := &Config{
		Enabled: false, // 禁用限流
		Rules: []RuleConfig{
			{
				Path:           "/api/*",
				LimitPerSecond: 1,
			},
		},
	}

	limiter := NewLimiter(rdb, config)
	ctx := context.Background()

	// 即使配置了限流规则，但由于disabled，所有请求都应该允许
	for i := 0; i < 10; i++ {
		allowed, remaining, _, err := limiter.Allow(ctx, "192.168.1.1", "/api/test")
		if err != nil {
			t.Fatalf("request %d: expected no error, got %v", i+1, err)
		}
		if !cast.ToBool(allowed) {
			t.Errorf("request %d: expected allowed to be true, got false", i+1)
		}
		if cast.ToInt(remaining) != -1 {
			t.Errorf("request %d: expected remaining to be -1, got %d", i+1, remaining)
		}
	}
}

// TestLimiter_Allow_DefaultRule 测试默认规则
func TestLimiter_Allow_DefaultRule(t *testing.T) {
	t.Parallel()

	rdb := setupTestRedis(t)
	defaultRule := RuleConfig{
		Path:           "*",
		LimitPerSecond: 5,
		BurstSize:      5,
	}
	config := &Config{
		Enabled: true,
		Rules: []RuleConfig{
			{
				Path:           "/api/test",
				LimitPerSecond: 1,
			},
		},
		DefaultRule: &defaultRule,
	}

	limiter := NewLimiter(rdb, config)
	ctx := context.Background()

	// 匹配到具体规则的路径
	allowed, _, _, err := limiter.Allow(ctx, "192.168.1.1", "/api/test")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !cast.ToBool(allowed) {
		t.Errorf("expected allowed to be true, got false")
	}

	// 不匹配任何规则，使用默认规则
	for i := 0; i < 5; i++ {
		allowed, _, _, err := limiter.Allow(ctx, "192.168.1.1", "/other/path")
		if err != nil {
			t.Fatalf("request %d: expected no error, got %v", i+1, err)
		}
		if !cast.ToBool(allowed) {
			t.Errorf("request %d should be allowed by default rule", i+1)
		}
	}

	// 超过默认规则限制，应该被限流
	allowed, _, _, err = limiter.Allow(ctx, "192.168.1.1", "/other/path")
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if cast.ToBool(allowed) {
		t.Errorf("expected allowed to be false, got true")
	}
}

// TestLimiter_Allow_NoMatchingRule 测试无匹配规则且无默认规则
func TestLimiter_Allow_NoMatchingRule(t *testing.T) {
	t.Parallel()

	rdb := setupTestRedis(t)
	config := &Config{
		Enabled: true,
		Rules: []RuleConfig{
			{
				Path:           "/api/*",
				LimitPerSecond: 1,
			},
		},
		DefaultRule: nil, // 没有默认规则
	}

	limiter := NewLimiter(rdb, config)
	ctx := context.Background()

	// 不匹配任何规则且没有默认规则，应该放行
	for i := 0; i < 10; i++ {
		allowed, remaining, _, err := limiter.Allow(ctx, "192.168.1.1", "/other/path")
		if err != nil {
			t.Fatalf("request %d: expected no error, got %v", i+1, err)
		}
		if !cast.ToBool(allowed) {
			t.Errorf("request %d: expected allowed to be true, got false", i+1)
		}
		if cast.ToInt(remaining) != -1 {
			t.Errorf("request %d: expected remaining to be -1, got %d", i+1, remaining)
		}
	}
}
