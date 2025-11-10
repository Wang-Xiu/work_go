# PV/UV统计组件 - 完整编码思路

## 📋 目录
1. [核心概念](#核心概念)
2. [Redis数据结构设计](#redis数据结构设计)
3. [关键算法实现](#关键算法实现)
4. [实现步骤指导](#实现步骤指导)
5. [性能优化建议](#性能优化建议)
6. [常见问题](#常见问题)

---

## 核心概念

### PV vs UV

```
PV (Page View - 页面浏览量):
  - 计数方式：每次访问 +1
  - 特点：累加计数
  - 实现：Redis INCR

UV (Unique Visitor - 独立访客):
  - 计数方式：同一visitor当天只计数一次
  - 特点：去重计数
  - 实现：Redis HyperLogLog
```

### 为什么用HyperLogLog？

```
场景：统计100万独立访客

方案A：使用Set存储
  - SADD stats:uv:2024-01-01 visitor1
  - SADD stats:uv:2024-01-01 visitor2
  - ...
  - SCARD stats:uv:2024-01-01

  内存占用：
    100万个visitor × 平均32字节 = 32MB

方案B：使用HyperLogLog ✅
  - PFADD stats:uv:2024-01-01 visitor1
  - PFADD stats:uv:2024-01-01 visitor2
  - ...
  - PFCOUNT stats:uv:2024-01-01

  内存占用：
    固定12KB（不管多少visitor）

  精度：
    误差率 < 0.81%（实际可接受）
```

---

## Redis数据结构设计

### Key命名规范

```go
// 全站PV（当天）
stats:pv:2024-01-15                    → 值：12345（INCR）

// 全站UV（当天）
stats:uv:2024-01-15                    → 值：HyperLogLog（PFADD/PFCOUNT）

// 路径PV（当天+路径）
stats:pv:2024-01-15:/api/users         → 值：1234（INCR）

// 路径UV（当天+路径）
stats:uv:2024-01-15:/api/users         → 值：HyperLogLog

// 路径列表（当天所有访问的路径）
stats:paths:2024-01-15                 → 值：Set {"/api/users", "/api/posts"}
```

### 数据流图

```
用户访问 → Track(visitorID, path)
              ↓
         ┌────┴────┐
         │Pipeline │  ← 批量执行Redis命令（减少RTT）
         └────┬────┘
              ↓
    ┌─────────┴─────────┐
    │                   │
  PV计数              UV去重
    ↓                   ↓
 INCR stats:pv      PFADD stats:uv
    ↓                   ↓
  返回新值            返回是否新增
```

---

## 关键算法实现

### 1. Track方法实现思路

```go
func (t *Tracker) Track(ctx context.Context, visitorID, path string) error {
    // 步骤1：检查是否启用、是否排除路径
    if !t.config.Enabled || t.config.IsExcludedPath(path) {
        return nil
    }

    // 步骤2：获取当前日期
    date := time.Now().Format("2006-01-02")

    // 步骤3：构建Redis Key
    pvKey := t.buildPVKey(date, "")     // stats:pv:2024-01-15
    uvKey := t.buildUVKey(date, "")     // stats:uv:2024-01-15

    // 步骤4：使用Pipeline批量执行（关键优化！）
    pipe := t.redis.Pipeline()

    // 4.1 全站PV +1
    pipe.Incr(ctx, pvKey)

    // 4.2 全站UV去重计数
    pipe.PFAdd(ctx, uvKey, visitorID)

    // 4.3 设置过期时间（保留RetentionDays天）
    expire := t.getExpireTime()
    pipe.Expire(ctx, pvKey, expire)
    pipe.Expire(ctx, uvKey, expire)

    // 4.4 如果启用路径统计
    if t.config.EnablePathStats {
        pathPVKey := t.buildPVKey(date, path)
        pathUVKey := t.buildUVKey(date, path)
        pathsKey := t.buildPathsKey(date)

        pipe.Incr(ctx, pathPVKey)
        pipe.PFAdd(ctx, pathUVKey, visitorID)
        pipe.SAdd(ctx, pathsKey, path)  // 记录路径列表
        pipe.Expire(ctx, pathPVKey, expire)
        pipe.Expire(ctx, pathUVKey, expire)
        pipe.Expire(ctx, pathsKey, expire)
    }

    // 步骤5：执行Pipeline
    _, err := pipe.Exec(ctx)
    return err
}
```

**关键点**：
- 使用Pipeline批量执行，一次网络往返完成所有操作
- 设置过期时间，自动清理历史数据
- 路径列表用Set存储，方便后续查询

### 2. GetDailyStats方法实现思路

```go
func (t *Tracker) GetDailyStats(ctx context.Context, date string) (*DailyStats, error) {
    // 步骤1：验证日期格式
    _, err := time.Parse("2006-01-02", date)
    if err != nil {
        return nil, ErrInvalidDate
    }

    // 步骤2：构建Key
    pvKey := t.buildPVKey(date, "")
    uvKey := t.buildUVKey(date, "")

    // 步骤3：查询PV
    pvStr, err := t.redis.Get(ctx, pvKey).Result()
    if err == redis.Nil {
        // Key不存在，返回0
        pvStr = "0"
    } else if err != nil {
        return nil, err
    }

    pv, _ := strconv.ParseInt(pvStr, 10, 64)

    // 步骤4：查询UV（使用PFCOUNT）
    uv, err := t.redis.PFCount(ctx, uvKey).Result()
    if err == redis.Nil {
        uv = 0
    } else if err != nil {
        return nil, err
    }

    // 步骤5：返回结果
    return &DailyStats{
        Date: date,
        PV:   pv,
        UV:   uv,
    }, nil
}
```

### 3. GetRangeStats方法实现思路（难点！）

```go
func (t *Tracker) GetRangeStats(ctx context.Context, startDate, endDate string) (*RangeStats, error) {
    // 步骤1：解析和验证日期
    start, err := time.Parse("2006-01-02", startDate)
    if err != nil {
        return nil, ErrInvalidDate
    }

    end, err := time.Parse("2006-01-02", endDate)
    if err != nil {
        return nil, ErrInvalidDate
    }

    if end.Before(start) {
        return nil, ErrInvalidDate
    }

    // 步骤2：遍历日期范围
    var dailyStats []DailyStats
    var totalPV int64
    var uvKeys []string  // 收集所有UV的Key，用于合并

    current := start
    for !current.After(end) {
        dateStr := current.Format("2006-01-02")

        // 查询每天的数据
        stats, err := t.GetDailyStats(ctx, dateStr)
        if err != nil {
            return nil, err
        }

        dailyStats = append(dailyStats, *stats)
        totalPV += stats.PV

        // 收集UV Key
        uvKey := t.buildUVKey(dateStr, "")
        uvKeys = append(uvKeys, uvKey)

        // 下一天
        current = current.AddDate(0, 0, 1)
    }

    // 步骤3：计算总UV（关键！）
    // 使用PFMERGE合并多个HyperLogLog，然后PFCOUNT

    // 创建临时Key存储合并结果
    tempKey := fmt.Sprintf("stats:uv:temp:%s:%s", startDate, endDate)

    // PFMERGE合并所有UV
    err = t.redis.PFMerge(ctx, tempKey, uvKeys...).Err()
    if err != nil {
        return nil, err
    }

    // PFCOUNT获取总UV
    totalUV, err := t.redis.PFCount(ctx, tempKey).Result()
    if err != nil {
        return nil, err
    }

    // 删除临时Key
    t.redis.Del(ctx, tempKey)

    // 步骤4：返回结果
    return &RangeStats{
        StartDate:  startDate,
        EndDate:    endDate,
        TotalPV:    totalPV,
        TotalUV:    totalUV,  // 去重后的真实UV
        DailyStats: dailyStats,
    }, nil
}
```

**为什么TotalUV不能直接累加？**

```
例子：
  Day1: user1访问 → UV=1
  Day2: user1访问 → UV=1

错误算法：TotalUV = 1 + 1 = 2  ❌
正确算法：TotalUV = 1（去重）  ✅

解决方案：使用PFMERGE合并HyperLogLog
```

### 4. GetPathStats方法实现思路

```go
func (t *Tracker) GetPathStats(ctx context.Context, date string) ([]PathStats, error) {
    if !t.config.EnablePathStats {
        return []PathStats{}, nil
    }

    // 步骤1：获取路径列表
    pathsKey := t.buildPathsKey(date)
    paths, err := t.redis.SMembers(ctx, pathsKey).Result()
    if err == redis.Nil {
        return []PathStats{}, nil
    } else if err != nil {
        return nil, err
    }

    // 步骤2：查询每个路径的PV和UV
    var result []PathStats

    for _, path := range paths {
        pvKey := t.buildPVKey(date, path)
        uvKey := t.buildUVKey(date, path)

        // 查询PV
        pvStr, _ := t.redis.Get(ctx, pvKey).Result()
        pv, _ := strconv.ParseInt(pvStr, 10, 64)

        // 查询UV
        uv, _ := t.redis.PFCount(ctx, uvKey).Result()

        result = append(result, PathStats{
            Path: path,
            PV:   pv,
            UV:   uv,
        })
    }

    return result, nil
}
```

### 5. CleanExpiredData方法实现思路

```go
func (t *Tracker) CleanExpiredData(ctx context.Context) error {
    // 步骤1：计算过期日期
    retentionDays := t.config.GetRetentionDays()
    expireDate := time.Now().AddDate(0, 0, -retentionDays)

    // 步骤2：构建匹配模式
    // 删除所有早于expireDate的Key

    // 例如：RetentionDays=90，今天是2024-01-15
    // 要删除 2023-10-17 之前的所有Key

    // 步骤3：使用SCAN遍历Key（不阻塞Redis）
    var cursor uint64
    var deletedCount int

    for {
        // SCAN匹配 stats:* 的Key
        keys, newCursor, err := t.redis.Scan(ctx, cursor, "stats:*", 100).Result()
        if err != nil {
            return err
        }

        // 遍历Key，提取日期部分
        for _, key := range keys {
            // 从Key中提取日期
            // stats:pv:2023-10-10 → 2023-10-10
            parts := strings.Split(key, ":")
            if len(parts) < 3 {
                continue
            }

            dateStr := parts[2]
            keyDate, err := time.Parse("2006-01-02", dateStr)
            if err != nil {
                continue
            }

            // 如果早于过期日期，删除
            if keyDate.Before(expireDate) {
                t.redis.Del(ctx, key)
                deletedCount++
            }
        }

        cursor = newCursor
        if cursor == 0 {
            break  // 遍历完成
        }
    }

    log.Printf("Cleaned %d expired keys", deletedCount)
    return nil
}
```

---

## 实现步骤指导

### Step 1: 实现基础的Track方法

```go
// 先实现最简单的版本：只统计全站PV和UV
func (t *Tracker) Track(ctx context.Context, visitorID, path string) error {
    date := time.Now().Format("2006-01-02")

    pvKey := fmt.Sprintf("stats:pv:%s", date)
    uvKey := fmt.Sprintf("stats:uv:%s", date)

    pipe := t.redis.Pipeline()
    pipe.Incr(ctx, pvKey)
    pipe.PFAdd(ctx, uvKey, visitorID)
    _, err := pipe.Exec(ctx)

    return err
}
```

**测试这个版本**：
```go
tracker.Track(ctx, "user1", "/api/test")
tracker.Track(ctx, "user1", "/api/test")  // 同一用户再次访问
tracker.Track(ctx, "user2", "/api/test")  // 不同用户访问

// 验证：PV=3, UV=2
```

### Step 2: 添加过期时间

```go
expire := t.getExpireTime()
pipe.Expire(ctx, pvKey, expire)
pipe.Expire(ctx, uvKey, expire)
```

### Step 3: 添加路径统计

```go
if t.config.EnablePathStats {
    // 实现路径级别的PV/UV统计
}
```

### Step 4: 实现查询方法

按顺序实现：
1. GetDailyStats（最简单）
2. GetPathStats（中等）
3. GetRangeStats（最复杂，涉及PFMERGE）

### Step 5: 实现清理方法

```go
// CleanExpiredData
// 建议先手动测试，确认逻辑正确后再接入定时任务
```

---

## 性能优化建议

### 1. 使用Pipeline批量执行

```go
// ❌ 错误：多次网络往返
t.redis.Incr(ctx, pvKey)        // RTT 1
t.redis.PFAdd(ctx, uvKey, id)   // RTT 2
t.redis.Expire(ctx, pvKey, exp) // RTT 3

// ✅ 正确：一次网络往返
pipe := t.redis.Pipeline()
pipe.Incr(ctx, pvKey)
pipe.PFAdd(ctx, uvKey, id)
pipe.Expire(ctx, pvKey, exp)
pipe.Exec(ctx)  // 只有1次RTT
```

### 2. 异步统计（高级）

```go
// 如果统计失败不应该影响业务请求
// 可以考虑异步统计

// 方式A：goroutine异步
go func() {
    tracker.Track(ctx, visitorID, path)
}()

// 方式B：消息队列
// Track → Redis Stream → Consumer处理
```

### 3. 路径统计的内存优化

```go
// 如果路径很多，EnablePathStats会占用大量内存
// 优化方案：只统计热门路径

// 使用Sorted Set记录路径访问频率
pipe.ZIncrBy(ctx, "stats:hot_paths", 1, path)

// 只统计访问量前100的路径
hotPaths := t.redis.ZRevRange(ctx, "stats:hot_paths", 0, 99).Val()
for _, path := range hotPaths {
    // 只统计这些路径
}
```

---

## 常见问题

### Q1: HyperLogLog的精度够用吗？

**A**: 对于PV/UV统计完全够用

```
误差率：< 0.81%

实例：
  真实UV：10000
  HyperLogLog结果：9920 ~ 10080

对业务影响：可忽略
```

### Q2: 如何处理跨天的UV统计？

**A**: 使用PFMERGE

```go
// 统计1月1日到1月7日的总UV
uvKeys := []string{
    "stats:uv:2024-01-01",
    "stats:uv:2024-01-02",
    // ...
    "stats:uv:2024-01-07",
}

// 合并到临时Key
t.redis.PFMerge(ctx, "temp", uvKeys...)

// 获取去重后的总UV
totalUV := t.redis.PFCount(ctx, "temp").Val()
```

### Q3: visitorID应该用什么？

**A**: 根据业务场景选择

```
精确度从高到低：

1. 用户ID（已登录）
   - 最精确
   - 只能统计登录用户

2. 设备指纹
   - IP + User-Agent + 其他特征
   - 较精确，覆盖所有用户

3. IP地址
   - 简单但不够精确
   - NAT网络下会重复计数

4. Cookie/LocalStorage
   - 需要前端配合
   - 可被清除
```

### Q4: 如何统计实时在线人数？

**A**: 这不是PV/UV，需要另外的方案

```go
// 使用Sorted Set + TTL
// Key: online_users
// Score: 当前时间戳
// Member: userID

// 用户活跃时更新时间戳
redis.ZAdd(ctx, "online_users", redis.Z{
    Score:  float64(time.Now().Unix()),
    Member: userID,
})

// 统计最近5分钟活跃的用户数
cutoff := time.Now().Unix() - 300  // 5分钟前
count := redis.ZCount(ctx, "online_users", cutoff, "+inf").Val()
```

---

## 总结

**核心技术点**：
1. PV用INCR计数
2. UV用HyperLogLog去重
3. 跨天UV用PFMERGE合并
4. Pipeline批量执行
5. SCAN遍历清理过期数据

**实现顺序**：
1. Track方法（核心）
2. GetDailyStats（简单查询）
3. GetRangeStats（复杂查询）
4. GetPathStats（可选）
5. CleanExpiredData（清理）

**测试重点**：
1. 同一用户多次访问，UV只计数一次
2. 不同用户访问，UV正确累加
3. 跨天UV统计正确去重
4. 路径统计独立计数
5. 过期数据正确清理

祝你编码顺利！🚀
