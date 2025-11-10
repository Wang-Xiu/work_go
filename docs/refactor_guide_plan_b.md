# 方案B：抽取Redis层 - 重构指南

## 一、架构变更说明

### 核心改动
1. **新增基础设施层** `common/infrastructure/redis/`
2. **统一Redis连接管理**：通过 `redis.Manager` 管理所有Redis连接
3. **配置统一化**：`config/app.yml` 统一管理Redis和中间件配置
4. **依赖注入**：中间件通过Manager获取Redis��户端，而不是自己创建

### 架构分层
```
应用层 (main.go)
    ↓
中间件层 (middleware/)
    ├── ratelimit.Limiter
    └── stats.Tracker
    ↓
基础设施层 (infrastructure/)
    └── redis.Manager
    ↓
Redis服务器
```

---

## 二、代码重构步骤

### Step 1: 迁移目录结构

```bash
# 1. 移动stats到middleware下
mv common/stats common/middleware/stats

# 2. 创建infrastructure目录（已完成）
# common/infrastructure/redis/ 已创建
```

### Step 2: 修改 ratelimit 配置

修改 `common/middleware/ratelimit/config.go`：

```go
package ratelimit

import "time"

// Config 限流配置
type Config struct {
    Enabled     bool          `yaml:"enabled"`
    RedisName   string        `yaml:"redis_name"`  // 【新增】引用Redis Manager中的连接名
    Rules       []*RuleConfig `yaml:"rules"`
    DefaultRule *RuleConfig   `yaml:"default_rule"`

    // 【删除】Redis RedisConfig - 不再需要单独的Redis配置
}

// RuleConfig 单条限流规则配置
type RuleConfig struct {
    Path            string `yaml:"path"`
    LimitPerSecond  int    `yaml:"limit_per_second"`
    LimitPerMinute  int    `yaml:"limit_per_minute"`
    BurstSize       int    `yaml:"burst_size"`
}

// GetBurstSize 获取突发容量
func (r *RuleConfig) GetBurstSize(limit int) int {
    if r.BurstSize > 0 {
        return r.BurstSize
    }
    return limit * 2 // 默认2倍
}

// Validate 验证规则配置
func (r *RuleConfig) Validate() error {
    if r.LimitPerSecond <= 0 && r.LimitPerMinute <= 0 {
        return ErrInvalidConfig
    }
    return nil
}
```

### Step 3: 修改 ratelimit Limiter

修改 `common/middleware/ratelimit/limiter.go`：

```go
package ratelimit

import (
    "context"
    "fmt"
    "time"

    "github.com/redis/go-redis/v9"
    infraredis "working-project/common/infrastructure/redis"  // 【新增】导入infrastructure
)

// Limiter 分布式限流器
type Limiter struct {
    redis  *redis.Client
    config *Config
}

// NewLimiter 创建限流器（通过依赖注入）
func NewLimiter(redisClient *redis.Client, config *Config) *Limiter {
    return &Limiter{
        redis:  redisClient,
        config: config,
    }
}

// NewLimiterFromManager 【新增】从Redis Manager创建限流器
func NewLimiterFromManager(manager *infraredis.Manager, config *Config) (*Limiter, error) {
    if config == nil {
        return nil, ErrInvalidConfig
    }

    // 从Manager获取Redis客户端
    redisName := config.RedisName
    if redisName == "" {
        redisName = "default"  // 默认使用default连接
    }

    redisClient, err := manager.Get(redisName)
    if err != nil {
        return nil, fmt.Errorf("failed to get redis client '%s': %w", redisName, err)
    }

    return NewLimiter(redisClient, config), nil
}

// NewLimiterFromConfig 【保留兼容性】从配置创建限流器
// 注意：此方法会创建独立的Redis连接，不推荐使用
// 推荐使用 NewLimiterFromManager
func NewLimiterFromConfig(config *Config) (*Limiter, error) {
    // ... 保持原有实现，但标记为 Deprecated
    panic("deprecated: use NewLimiterFromManager instead")
}

// Close 【新增】关闭限流器
// 注意：当使用NewLimiterFromManager创建时，不应该调用Close，
// Redis连接由Manager统一管理
func (l *Limiter) Close() error {
    // 此实现为空，因为Redis连接由Manager管理
    return nil
}

// Allow 检查请求是否允许通过（保持不变）
func (l *Limiter) Allow(ctx context.Context, key, path string) (bool, int, time.Time, error) {
    // ... 保持原有实现
    return true, -1, time.Time{}, nil
}

// ... 其他方法保持不变
```

### Step 4: 修改 stats 配置

修改 `common/middleware/stats/config.go`：

```go
package stats

import "time"

// Config 统计配置
type Config struct {
    Enabled         bool     `yaml:"enabled"`
    RedisName       string   `yaml:"redis_name"`      // 【新增】引用Redis Manager中的连接名
    EnablePathStats bool     `yaml:"enable_path_stats"`
    RetentionDays   int      `yaml:"retention_days"`
    ExcludePaths    []string `yaml:"exclude_paths"`

    // 【删除】Redis RedisConfig - 不再需要单独的Redis配置
}

// IsExcludedPath 检查路径是否被排除
func (c *Config) IsExcludedPath(path string) bool {
    for _, p := range c.ExcludePaths {
        if p == path {
            return true
        }
    }
    return false
}
```

### Step 5: 修改 stats Tracker

修改 `common/middleware/stats/stats.go`：

```go
package stats

import (
    "context"
    "fmt"
    "time"

    "github.com/redis/go-redis/v9"
    infraredis "working-project/common/infrastructure/redis"  // 【新增】
)

// Tracker 统计追踪器
type Tracker struct {
    redis  *redis.Client
    config *Config
}

// NewTracker 创建统计追踪器
func NewTracker(redisClient *redis.Client, config *Config) *Tracker {
    return &Tracker{
        redis:  redisClient,
        config: config,
    }
}

// NewTrackerFromManager 【新增】从Redis Manager创建追踪器
func NewTrackerFromManager(manager *infraredis.Manager, config *Config) (*Tracker, error) {
    if config == nil {
        return nil, ErrInvalidConfig
    }

    // 从Manager获取Redis客户端
    redisName := config.RedisName
    if redisName == "" {
        redisName = "default"
    }

    redisClient, err := manager.Get(redisName)
    if err != nil {
        return nil, fmt.Errorf("failed to get redis client '%s': %w", redisName, err)
    }

    return NewTracker(redisClient, config), nil
}

// NewTrackerFromConfig 【废弃】从配置创建追踪器
func NewTrackerFromConfig(config *Config) (*Tracker, error) {
    panic("deprecated: use NewTrackerFromManager instead")
}

// Close 关闭追踪器
func (t *Tracker) Close() error {
    // Redis连接由Manager管理，这里为空实现
    return nil
}

// Track 记录一次访问（保持不变）
func (t *Tracker) Track(ctx context.Context, visitorID, path string) error {
    // ... 保持原有实现
    return nil
}

// ... 其他方法保持不变
```

---

## 三、应用启动代码

修改 `example/main.go`：

```go
package main

import (
    "context"
    "log"

    "github.com/gin-gonic/gin"
    "gopkg.in/yaml.v3"

    infraredis "working-project/common/infrastructure/redis"
    "working-project/common/middleware/ratelimit"
    "working-project/common/middleware/stats"
)

// AppConfig 应用配置结构（对应app.yml）
type AppConfig struct {
    Redis      map[string]infraredis.Config `yaml:"redis"`
    Middleware struct {
        RateLimit ratelimit.Config `yaml:"ratelimit"`
        Stats     stats.Config     `yaml:"stats"`
    } `yaml:"middleware"`
}

func main() {
    ctx := context.Background()

    // 1. 加载配置
    var appConfig AppConfig
    if err := loadConfig("config/app.yml", &appConfig); err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // 2. 初始化Redis Manager（全局单例）
    redisManager := infraredis.GetGlobalManager()

    // 注册所有Redis连接
    for name, cfg := range appConfig.Redis {
        if err := redisManager.Register(name, &cfg); err != nil {
            log.Fatalf("Failed to register redis '%s': %v", name, err)
        }
        log.Printf("✅ Redis '%s' connected: %s:%d", name, cfg.Host, cfg.Port)
    }
    defer redisManager.CloseAll()

    // 3. 创建限流器
    limiter, err := ratelimit.NewLimiterFromManager(redisManager, &appConfig.Middleware.RateLimit)
    if err != nil {
        log.Fatalf("Failed to create limiter: %v", err)
    }
    log.Println("✅ Rate limiter initialized")

    // 4. 创建统计追踪器
    tracker, err := stats.NewTrackerFromManager(redisManager, &appConfig.Middleware.Stats)
    if err != nil {
        log.Fatalf("Failed to create tracker: %v", err)
    }
    log.Println("✅ Stats tracker initialized")

    // 5. 初始化Gin
    r := gin.Default()

    // 应用中间件
    r.Use(stats.Middleware(tracker))
    r.Use(ratelimit.Middleware(limiter))

    // 定义路由
    r.GET("/api/test", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "OK"})
    })

    r.GET("/admin/stats/today", func(c *gin.Context) {
        dailyStats, err := tracker.GetDailyStats(ctx, time.Now().Format("2006-01-02"))
        if err != nil {
            c.JSON(500, gin.H{"error": err.Error()})
            return
        }
        c.JSON(200, dailyStats)
    })

    // 启动服务
    log.Println("🚀 Server starting on :8080")
    if err := r.Run(":8080"); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}

func loadConfig(path string, config interface{}) error {
    // 实现配置加载逻辑（支持KMS解密等）
    // 这里简化处理
    return yaml.Unmarshal([]byte{}, config)
}
```

---

## 四、迁移检查清单

- [ ] 创建 `common/infrastructure/redis/` 目录
- [ ] 实现 `redis.Manager`、`redis.Config`、`redis.Errors`
- [ ] 移动 `common/stats` 到 `common/middleware/stats`
- [ ] 创建 `config/app.yml` 统一配置文件
- [ ] 修改 `ratelimit/config.go` 删除RedisConfig，添加RedisName
- [ ] 修改 `ratelimit/limiter.go` 添加NewLimiterFromManager方法
- [ ] 修改 `stats/config.go` 删除RedisConfig，添加RedisName
- [ ] 修改 `stats/stats.go` 添加NewTrackerFromManager方法
- [ ] 更新 `example/main.go` 使用新的初始化方式
- [ ] 批量替换import路径：`common/stats` → `common/middleware/stats`
- [ ] 运行单元测试确保功能正常
- [ ] 更新文档和注释

---

## 五、优势总结

### ✅ 解决的问题
1. **消除Redis连接重复**：全应用共享一个或多个命名的Redis连接池
2. **配置统一管理**：所有Redis配置在app.yml中统一定义
3. **清晰的分层架构**：基础设施层 → 中间件层 → 应用层
4. **易于测试**：可以轻松Mock redis.Manager进行单元测试
5. **易于扩展**：未来新增功能可直接复用Redis Manager

### ✅ 性能提升
- **减少连接数**：从N个独立连接 → 1个共享连接池
- **降低内存占用**：共享连接池减少资源消耗
- **提高连接复用率**：所有模块共享连接

### ✅ 维护性提升
- **配置集中化**：只需维护一份Redis配置
- **代码复用**：所有模块使用统一的Redis管理器
- **职责清晰**：基础设施层专注连接管理，中间件层专注业务逻辑
