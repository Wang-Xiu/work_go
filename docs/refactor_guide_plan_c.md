# 方案C：完整DDD架构 - 企业级基础设施层设计

## 一、架构理念

### DDD分层架构完整视图

```
┌────────────────────────────────────────────────────────┐
│  表示层 (Presentation Layer)                            │
│  - HTTP Handler (domain/*/handler/)                   │
│  - Proto定义 (domain/*/proto/)                        │
└────────────────────────────────────────────────────────┘
                         ↓
┌────────────────────────────────────────────────────────┐
│  应用层 (Application Layer)                            │
│  - Service编排 (domain/*/service/)                    │
│  - 事件发布 (domain/*/event/)                         │
└────────────────────────────────────────────────────────┘
                         ↓
┌────────────────────────────────────────────────────────┐
│  领域层 (Domain Layer)                                 │
│  - 实体 (domain/*/entity/)                            │
│  - 值对象 (domain/*/valueobject/)                     │
│  - 仓储接口 (domain/*/repository/)                    │
│  - 领域服务 (domain/*/service/)                       │
└────────────────────────────────────────────────────────┘
                         ↓
┌────────────────────────────────────────────────────────┐
│  基础设施层 (Infrastructure Layer) ★ 本次重构重点      │
│  - 持久化 (infrastructure/persistence/)               │
│  - 缓存 (infrastructure/cache/)                       │
│  - 消息队列 (infrastructure/mq/)                      │
│  - HTTP中间件 (infrastructure/middleware/)            │
│  - 配置管理 (infrastructure/config/)                  │
│  - 日志 (infrastructure/logging/)                     │
│  - 监控 (infrastructure/monitoring/)                  │
└────────────────────────────────────────────────────────┘
```

### 核心原则

1. **依赖倒置**：上层依赖接口，基础设施层实现接口
2. **职责单一**：每个子模块只负责一个技术维度
3. **可测试性**：所有基础设施组件可独立测试和Mock
4. **可配置性**：所有基础设施通过配置驱动
5. **可观测性**：内置日志、监控、追踪能力

---

## 二、完整目录结构

```
/Users/xiu/work/work_go/
├── common/
│   ├── infrastructure/              # 基础设施层（新增完整结构）
│   │   ├── persistence/            # 持久化层
│   │   │   ├── redis/              # Redis管理
│   │   │   │   ├── manager.go      # Redis连接管理器
│   │   │   │   ├── config.go       # Redis配置
│   │   │   │   ├── client.go       # 客户端封装
│   │   │   │   ├── health.go       # 健康检查
│   │   │   │   └── errors.go
│   │   │   ├── mysql/              # MySQL管理
│   │   │   │   ├── manager.go
│   │   │   │   ├── config.go
│   │   │   │   └── ...
│   │   │   └── mongo/              # MongoDB管理
│   │   │
│   │   ├── cache/                  # 缓存层
│   │   │   ├── local/              # 本地缓存（groupcache）
│   │   │   └── distributed/        # 分布式缓存（Redis）
│   │   │
│   │   ├── middleware/             # HTTP中间件
│   │   │   ├── core/               # 中间件核心抽象
│   │   │   │   ├── types.go        # 通用接口定义
│   │   │   │   ├── context.go      # 上下文增强
│   │   │   │   └── errors.go       # 错误处理
│   │   │   │
│   │   │   ├── ratelimit/          # 限流中间件
│   │   │   │   ├── config.go
│   │   │   │   ├── limiter.go
│   │   │   │   ├── middleware.go
│   │   │   │   ├── strategy/       # 限流策略
│   │   │   │   │   ├── token_bucket.go
│   │   │   │   │   ├── leaky_bucket.go
│   │   │   │   │   └── sliding_window.go
│   │   │   │   └── limiter_test.go
│   │   │   │
│   │   │   ├── stats/              # 统计中间件
│   │   │   │   ├── config.go
│   │   │   │   ├── tracker.go
│   │   │   │   ├── middleware.go
│   │   │   │   ├── aggregator/     # 数据聚合
│   │   │   │   │   ├── daily.go
│   │   │   │   │   ├── hourly.go
│   │   │   │   │   └── realtime.go
│   │   │   │   └── stats_test.go
│   │   │   │
│   │   │   ├── auth/               # 认证中间件（未来）
│   │   │   ├── logging/            # 日志中间件（未来）
│   │   │   └── tracing/            # 追踪中间件（未来）
│   │   │
│   │   ├── config/                 # 配置管理
│   │   │   ├── loader.go           # 配置加载器
│   │   │   ├── watcher.go          # 配置热更新
│   │   │   ├── validator.go        # 配置验证
│   │   │   ├── kms/                # KMS加密支持
│   │   │   │   ├── manager.go
│   │   │   │   └── provider.go
│   │   │   └── env.go              # 环境变量管理
│   │   │
│   │   ├── logging/                # 日志系统
│   │   │   ├── logger.go           # 日志接口
│   │   │   ├── zap.go              # Zap实现
│   │   │   ├── context.go          # 上下文日志
│   │   │   └── rotation.go         # 日志轮转
│   │   │
│   │   ├── monitoring/             # 监控系统
│   │   │   ├── metrics/            # 指标采集
│   │   │   │   ├── prometheus.go
│   │   │   │   └── custom.go
│   │   │   ├── tracing/            # 链路追踪
│   │   │   │   └── jaeger.go
│   │   │   └── health/             # 健康检查
│   │   │       ├── checker.go
│   │   │       └── endpoint.go
│   │   │
│   │   └── eventbus/               # 事件总线
│   │       ├── publisher.go
│   │       ├── subscriber.go
│   │       └── manager.go
│   │
│   ├── base/                        # 跨domain常量、枚举
│   │   ├── errors/                 # 通用错误定义
│   │   │   ├── codes.go            # 错误码
│   │   │   └── common.go           # 通用错误
│   │   ├── constants/              # 常量定义
│   │   └── enums/                  # 枚举定义
│   │
│   ├── util/                        # 工具类
│   │   ├── time/                   # 时间工具
│   │   ├── crypto/                 # 加密工具
│   │   ├── http/                   # HTTP工具
│   │   └── validator/              # 验证工具
│   │
│   └── sys.go                       # 全局变量
│
├── config/
│   ├── app.yml                      # 应用主配置
│   ├── infrastructure.yml           # 基础设施配置（新增）
│   ├── middleware.yml               # 中间件配置（新增）
│   ├── logging.yml                  # 日志配置
│   └── monitoring.yml               # 监控配置
│
├── domain/                          # 业务领域（保持不变）
│   ├── user/
│   ├── order/
│   └── ...
│
├── deployment/                      # 部署配置
├── job/                            # 定时任务
└── example/                        # 示例代码
```

---

## 三、核心组件设计

### 3.1 中间件核心抽象

创建 `common/infrastructure/middleware/core/types.go`：

```go
package core

import (
    "context"
    "time"

    "github.com/gin-gonic/gin"
)

// Middleware 中间件接口
type Middleware interface {
    // Name 中间件名称
    Name() string

    // Priority 优先级（数字越小越先执行）
    Priority() int

    // Handle 处理函数
    Handle() gin.HandlerFunc

    // Initialize 初始化
    Initialize() error

    // Shutdown 关闭
    Shutdown(ctx context.Context) error

    // HealthCheck 健康检查
    HealthCheck(ctx context.Context) error
}

// Config 中间件配置接口
type Config interface {
    // Validate 验证配置
    Validate() error

    // IsEnabled 是否启用
    IsEnabled() bool
}

// Metrics 中间件指标接口
type Metrics interface {
    // RecordRequest 记录请求
    RecordRequest(ctx context.Context, duration time.Duration, status int)

    // RecordError 记录错误
    RecordError(ctx context.Context, err error)
}

// BaseMiddleware 中间件基类（提供默认实现）
type BaseMiddleware struct {
    name     string
    priority int
}

func NewBaseMiddleware(name string, priority int) *BaseMiddleware {
    return &BaseMiddleware{
        name:     name,
        priority: priority,
    }
}

func (m *BaseMiddleware) Name() string {
    return m.name
}

func (m *BaseMiddleware) Priority() int {
    return m.priority
}

func (m *BaseMiddleware) Initialize() error {
    return nil
}

func (m *BaseMiddleware) Shutdown(ctx context.Context) error {
    return nil
}

func (m *BaseMiddleware) HealthCheck(ctx context.Context) error {
    return nil
}
```

### 3.2 配置管理架构

创建 `config/infrastructure.yml`：

```yaml
# 基础设施配置
infrastructure:
  # 持久化配置
  persistence:
    redis:
      # 默认实例（限流、统计共享）
      default:
        host: localhost
        port: 6379
        password: ""
        db: 0
        pool_size: 20
        min_idle_conns: 5
        max_retries: 3
        dial_timeout: 5s
        read_timeout: 3s
        write_timeout: 3s

      # 缓存专用实例
      cache:
        host: localhost
        port: 6379
        password: ""
        db: 1
        pool_size: 10
        min_idle_conns: 3
        max_retries: 3
        dial_timeout: 5s
        read_timeout: 3s
        write_timeout: 3s

      # 会话专用实例
      session:
        host: localhost
        port: 6379
        password: ""
        db: 2
        pool_size: 10
        min_idle_conns: 3
        max_retries: 3
        dial_timeout: 5s
        read_timeout: 3s
        write_timeout: 3s

    mysql:
      # 主库
      master:
        host: localhost
        port: 3306
        database: app_db
        username: root
        password: "kms://encrypted"
        max_open_conns: 100
        max_idle_conns: 10
        conn_max_lifetime: 3600s

      # 只读从库
      slave:
        host: localhost
        port: 3307
        database: app_db
        username: readonly
        password: "kms://encrypted"
        max_open_conns: 50
        max_idle_conns: 5
        conn_max_lifetime: 3600s

  # 日志配置
  logging:
    level: info  # debug, info, warn, error
    format: json  # json, text
    output: stdout  # stdout, file, both
    file:
      path: /var/log/app/app.log
      max_size: 100  # MB
      max_backups: 10
      max_age: 30  # days
      compress: true

  # 监控配置
  monitoring:
    metrics:
      enabled: true
      port: 9090
      path: /metrics

    tracing:
      enabled: true
      jaeger_endpoint: http://localhost:14268/api/traces
      sample_rate: 0.1  # 10%采样率

    health:
      enabled: true
      path: /health
      checks:
        - name: redis
          timeout: 3s
        - name: mysql
          timeout: 5s
```

创建 `config/middleware.yml`：

```yaml
# 中间件配置
middleware:
  # 限流中间件
  ratelimit:
    enabled: true
    priority: 10  # 优先级（数字越小越先执行）
    redis_name: default
    strategy: token_bucket  # token_bucket, leaky_bucket, sliding_window

    rules:
      - path: /api/*
        limit_per_second: 10
        limit_per_minute: 100
        burst_size: 20

      - path: /api/login
        limit_per_second: 1
        limit_per_minute: 5
        burst_size: 2

      - path: /callback/payment/*
        limit_per_second: 100
        limit_per_minute: 1000
        burst_size: 200

    default_rule:
      path: "*"
      limit_per_second: 5
      limit_per_minute: 50
      burst_size: 10

  # 统计中间件
  stats:
    enabled: true
    priority: 20
    redis_name: default
    enable_path_stats: true
    enable_realtime: true  # 实时统计
    retention_days: 90
    aggregation:
      hourly: true  # 按小时聚合
      daily: true   # 按天聚合
      monthly: true # 按月聚合
    exclude_paths:
      - /health
      - /metrics
      - /favicon.ico
      - /robots.txt

  # 日志中间件
  logging:
    enabled: true
    priority: 5
    log_request: true
    log_response: false
    log_headers: false
    log_body: false
    max_body_size: 1024  # bytes

  # 追踪中间件
  tracing:
    enabled: true
    priority: 1
    sample_rate: 0.1

  # 认证中间件
  auth:
    enabled: true
    priority: 15
    exclude_paths:
      - /api/login
      - /api/register
      - /health
```

---

## 四、启动流程设计

创建 `example/main_full.go`：

```go
package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/gin-gonic/gin"

    "working-project/common/infrastructure/config"
    "working-project/common/infrastructure/logging"
    "working-project/common/infrastructure/middleware/core"
    infraredis "working-project/common/infrastructure/persistence/redis"
    "working-project/common/infrastructure/middleware/ratelimit"
    "working-project/common/infrastructure/middleware/stats"
    "working-project/common/infrastructure/monitoring"
)

type Application struct {
    logger       logging.Logger
    redisManager *infraredis.Manager
    middlewares  []core.Middleware
    router       *gin.Engine
}

func main() {
    ctx := context.Background()

    // 1. 初始化配置加载器
    cfgLoader := config.NewLoader()
    appConfig, err := cfgLoader.Load(ctx)
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // 2. 初始化日志系统
    logger := logging.NewLogger(appConfig.Infrastructure.Logging)
    logger.Info("Application starting...")

    // 3. 初始化Redis Manager
    redisManager := infraredis.GetGlobalManager()
    for name, cfg := range appConfig.Infrastructure.Persistence.Redis {
        if err := redisManager.Register(name, &cfg); err != nil {
            logger.Fatal("Failed to register redis", "name", name, "error", err)
        }
        logger.Info("Redis connected", "name", name, "addr", cfg.GetAddress())
    }

    // 4. 初始化监控系统
    monitor := monitoring.NewMonitor(appConfig.Infrastructure.Monitoring)
    if err := monitor.Start(); err != nil {
        logger.Fatal("Failed to start monitoring", "error", err)
    }
    logger.Info("Monitoring started")

    // 5. 初始化所有中间件
    middlewares := make([]core.Middleware, 0)

    // 5.1 限流中间件
    if appConfig.Middleware.RateLimit.Enabled {
        rateLimiter, err := ratelimit.NewMiddleware(redisManager, &appConfig.Middleware.RateLimit)
        if err != nil {
            logger.Fatal("Failed to create rate limiter", "error", err)
        }
        middlewares = append(middlewares, rateLimiter)
        logger.Info("Rate limiter initialized")
    }

    // 5.2 统计中间件
    if appConfig.Middleware.Stats.Enabled {
        statsTracker, err := stats.NewMiddleware(redisManager, &appConfig.Middleware.Stats)
        if err != nil {
            logger.Fatal("Failed to create stats tracker", "error", err)
        }
        middlewares = append(middlewares, statsTracker)
        logger.Info("Stats tracker initialized")
    }

    // 6. 初始化Gin路由
    r := gin.New()

    // 按优先级排序并应用中间件
    core.SortMiddlewaresByPriority(middlewares)
    for _, mw := range middlewares {
        r.Use(mw.Handle())
        logger.Info("Middleware registered", "name", mw.Name(), "priority", mw.Priority())
    }

    // 7. 注册健康检查端点
    r.GET("/health", healthCheckHandler(redisManager, middlewares))

    // 8. 注册业务路由
    registerBusinessRoutes(r)

    // 9. 优雅启动和关闭
    app := &Application{
        logger:       logger,
        redisManager: redisManager,
        middlewares:  middlewares,
        router:       r,
    }

    app.Run()
}

func (app *Application) Run() {
    // 启动HTTP服务器
    srv := &http.Server{
        Addr:    ":8080",
        Handler: app.router,
    }

    go func() {
        app.logger.Info("Server starting", "port", 8080)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            app.logger.Fatal("Server failed", "error", err)
        }
    }()

    // 等待中断信号
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    app.logger.Info("Server shutting down...")

    // 优雅关闭
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    // 关闭HTTP服务器
    if err := srv.Shutdown(ctx); err != nil {
        app.logger.Error("Server forced to shutdown", "error", err)
    }

    // 关闭所有中间件
    for _, mw := range app.middlewares {
        if err := mw.Shutdown(ctx); err != nil {
            app.logger.Error("Middleware shutdown failed", "name", mw.Name(), "error", err)
        }
    }

    // 关闭Redis连接
    if err := app.redisManager.CloseAll(); err != nil {
        app.logger.Error("Redis shutdown failed", "error", err)
    }

    app.logger.Info("Server stopped")
}

func healthCheckHandler(redisManager *infraredis.Manager, middlewares []core.Middleware) gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx := c.Request.Context()
        status := "healthy"
        checks := make(map[string]string)

        // 检查Redis
        for _, name := range redisManager.List() {
            client, _ := redisManager.Get(name)
            if err := client.Ping(ctx).Err(); err != nil {
                checks["redis_"+name] = "unhealthy: " + err.Error()
                status = "unhealthy"
            } else {
                checks["redis_"+name] = "healthy"
            }
        }

        // 检查中间件
        for _, mw := range middlewares {
            if err := mw.HealthCheck(ctx); err != nil {
                checks["middleware_"+mw.Name()] = "unhealthy: " + err.Error()
                status = "unhealthy"
            } else {
                checks["middleware_"+mw.Name()] = "healthy"
            }
        }

        c.JSON(200, gin.H{
            "status": status,
            "checks": checks,
        })
    }
}

func registerBusinessRoutes(r *gin.Engine) {
    api := r.Group("/api")
    {
        api.GET("/test", func(c *gin.Context) {
            c.JSON(200, gin.H{"message": "OK"})
        })
    }

    admin := r.Group("/admin")
    {
        admin.GET("/stats/today", func(c *gin.Context) {
            // 实现统计查询逻辑
            c.JSON(200, gin.H{})
        })
    }
}
```

---

## 五、迁移步骤

### Phase 1: 基础设施层搭建（第1-2周）
1. 创建完整的infrastructure目录结构
2. 实现Redis Manager和配置管理
3. 建立日志系统
4. 建立监控系统

### Phase 2: 中间件重构（第3周）
1. 移动stats到middleware下
2. 实现middleware/core核心抽象
3. 重构ratelimit适配新架构
4. 重构stats适配新架构

### Phase 3: 配置统一化（第4周）
1. 创建infrastructure.yml和middleware.yml
2. 重构配置加载逻辑
3. 实现配置热更新
4. 实现KMS解密

### Phase 4: 测试和文档（第5周）
1. 编写单元测试
2. 编写集成测试
3. 性能测试和优化
4. 完善文档

---

## 六、收益评估

### 架构层面
- ✅ **清晰的分层架构**：符合DDD最佳实践
- ✅ **高度模块化**：每个子系统独立可测试
- ✅ **易于扩展**：新功能遵循统一模式
- ✅ **技术债务清零**：消除所有历史包袱

### 技术层面
- ✅ **资源优化**：共享连接池，减少资源消耗
- ✅ **性能提升**：统一管理，优化连接复用
- ✅ **可观测性**：内置监控、日志、追踪
- ✅ **高可用性**：健康检查、优雅关闭

### 团队层面
- ✅ **降低学习成本**：统一的编码规范
- ✅ **提升开发效率**：可复用的基础设施
- ✅ **减少Bug**：统一的错误处理和验证
- ✅ **便于Code Review**：清晰的代码组织

---

这是最完整的企业级架构方案，适合中长期演进和大型团队协作。
