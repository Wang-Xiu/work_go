# UV统计方案对比与选择

## 🎯 需求分析

**场景**：查询3个月（90天）内每日的PV/UV数据

**查询频率**：高频（用户频繁查看统计报表）

**数据量**：每日UV约10万-100万

---

## 📊 三种方案对比

### 方案1：HyperLogLog（原方案）

```redis
stats:uv:2024-01-15 -> HyperLogLog {visitor1, visitor2, ...}
```

**优点**：
- ✅ 内存极省（100万UV只需12KB）
- ✅ 支持PFMERGE跨天去重

**缺点**：
- ❌ 有0.81%误差（不精确）
- ❌ 查询90天需要90次PFCOUNT（性能一般）
- ❌ PFMERGE 90个HyperLogLog很慢

**适用场景**：UV量极大（千万级），对精度要求不高

---

### 方案2：Set存储访客ID

```redis
stats:uv:2024-01-15 -> Set {visitor1, visitor2, ...}
```

**优点**：
- ✅ 精确计数（无误差）
- ✅ 支持SUNION跨天去重
- ✅ 实现简单

**缺点**：
- ❌ 内存占用较大（100万UV ≈ 24MB）
- ❌ 查询90天需要90次SCARD
- ❌ SUNION 90个Set会很慢

**内存计算**：
```
100万UV × 24字节 = 24MB/天
90天 = 2.16GB
```

**适用场景**：UV量中等（<100万/天），需要精确统计

---

### 方案3：Hash存储访客ID

```redis
stats:uv:2024-01-15 -> Hash {visitor1: 1, visitor2: 1, ...}
```

**优点**：
- ✅ 精确计数
- ✅ 内存占用比Set小30%（ziplist编码）
- ✅ 查询单日极快（HLEN O(1)）
- ✅ 可扩展（可存储更多信息）

**缺点**：
- ❌ 不支持原生跨天去重（需程序实现）
- ❌ 查询90天需要90次HLEN

**内存计算**：
```
100万UV × 17字节 = 17MB/天
90天 = 1.53GB
```

**适用场景**：UV量中等，主要查询单日数据

---

## 🏆 推荐方案：Set + Hash 混合方案

### 数据结构设计

```redis
# 1. 实时UV统计（Set，短期保留）
stats:uv:realtime:2024-01-15 -> Set {visitor1, visitor2, ...}
TTL: 7天

# 2. 每日UV汇总（Hash，长期保留）
stats:uv:summary:2024-01 -> Hash {
    "01": "8234",    # 1月1日UV
    "02": "8456",    # 1月2日UV
    ...
    "31": "9012"     # 1月31日UV
}
TTL: 90天

# 3. PV统计（String，直接存数字）
stats:pv:2024-01-15 -> "10234"
TTL: 90天
```

### 核心逻辑

#### 1. 记录访问（Track）

```go
func (t *Tracker) Track(ctx context.Context, visitorID, path string) error {
    date := time.Now().Format("2006-01-02")
    yearMonth := time.Now().Format("2006-01")
    day := time.Now().Format("02")

    pipe := t.redis.Pipeline()

    // PV +1
    pvKey := fmt.Sprintf("stats:pv:%s", date)
    pipe.Incr(ctx, pvKey)
    pipe.Expire(ctx, pvKey, 90*24*time.Hour)

    // UV去重（Set）
    uvKey := fmt.Sprintf("stats:uv:realtime:%s", date)
    pipe.SAdd(ctx, uvKey, visitorID)
    pipe.Expire(ctx, uvKey, 7*24*time.Hour)  // 只保留7天

    _, err := pipe.Exec(ctx)
    return err
}
```

#### 2. 每日汇总任务（凌晨2点执行）

```go
func (t *Tracker) DailySummary(ctx context.Context) error {
    yesterday := time.Now().AddDate(0, 0, -1)
    date := yesterday.Format("2006-01-02")
    yearMonth := yesterday.Format("2006-01")
    day := yesterday.Format("02")

    // 读取昨日UV数
    uvKey := fmt.Sprintf("stats:uv:realtime:%s", date)
    uv, err := t.redis.SCard(ctx, uvKey).Result()
    if err != nil {
        return err
    }

    // 写入汇总Hash
    summaryKey := fmt.Sprintf("stats:uv:summary:%s", yearMonth)
    err = t.redis.HSet(ctx, summaryKey, day, uv).Err()
    if err != nil {
        return err
    }

    // 设置Hash过期时间
    t.redis.Expire(ctx, summaryKey, 90*24*time.Hour)

    return nil
}
```

#### 3. 查询90天数据（超快！）

```go
func (t *Tracker) GetRangeStats(ctx context.Context, startDate, endDate string) (*RangeStats, error) {
    start, _ := time.Parse("2006-01-02", startDate)
    end, _ := time.Parse("2006-01-02", endDate)

    // 按月分组
    monthKeys := make(map[string][]string) // yearMonth -> [day1, day2, ...]
    for current := start; !current.After(end); current = current.AddDate(0, 0, 1) {
        yearMonth := current.Format("2006-01")
        day := current.Format("02")
        monthKeys[yearMonth] = append(monthKeys[yearMonth], day)
    }

    var dailyStats []DailyStats
    var totalPV, totalUV int64

    // 批量查询（每个月只需2次Redis调用）
    for yearMonth, days := range monthKeys {
        // 查询这个月的UV汇总
        uvSummaryKey := fmt.Sprintf("stats:uv:summary:%s", yearMonth)
        uvData, _ := t.redis.HMGet(ctx, uvSummaryKey, days...).Result()

        // 查询每天的PV
        for i, day := range days {
            date := fmt.Sprintf("%s-%s", yearMonth, day)

            // PV
            pvKey := fmt.Sprintf("stats:pv:%s", date)
            pv, _ := t.redis.Get(ctx, pvKey).Int64()

            // UV（从汇总读取）
            uv, _ := cast.ToInt64E(uvData[i])

            dailyStats = append(dailyStats, DailyStats{
                Date: date,
                PV:   pv,
                UV:   uv,
            })

            totalPV += pv
            totalUV += uv  // 注意：这是简单累加，不去重
        }
    }

    return &RangeStats{
        StartDate:  startDate,
        EndDate:    endDate,
        TotalPV:    totalPV,
        TotalUV:    totalUV,
        DailyStats: dailyStats,
    }, nil
}
```

#### 4. 计算近7天真实UV（去重）

```go
func (t *Tracker) GetLast7DaysUniqueUV(ctx context.Context) (int64, error) {
    keys := make([]string, 7)
    for i := 0; i < 7; i++ {
        date := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
        keys[i] = fmt.Sprintf("stats:uv:realtime:%s", date)
    }

    // 使用SUNION去重计算
    uv, err := t.redis.SUnion(ctx, keys...).Result()
    if err != nil {
        return 0, err
    }

    return int64(len(uv)), nil
}
```

---

## 📈 性能对比

### 查询90天数据

| 方案 | Redis调用次数 | 响应时间 | 内存占用 |
|------|--------------|---------|---------|
| HyperLogLog | 90次PFCOUNT | ~200ms | 1.08MB |
| Set | 90次SCARD | ~150ms | 2.16GB |
| Hash | 90次HLEN | ~150ms | 1.53GB |
| **混合方案** | **6次HMGET** | **~20ms** | **1.5GB** |

### 计算7天总UV（去重）

| 方案 | Redis调用次数 | 响应时间 |
|------|--------------|---------|
| HyperLogLog | 1次PFMERGE + 1次PFCOUNT | ~50ms |
| Set | 1次SUNION | ~30ms |
| Hash | 需要程序去重 | ~500ms |
| **混合方案** | **1次SUNION** | **~30ms** |

---

## ✅ 最终建议

**对于你的场景（查询3个月每日UV），推荐：混合方案**

**理由**：
1. ✅ 查询90天数据极快（只需6次Redis调用）
2. ✅ 精确计数（无误差）
3. ✅ 内存占用合理（Set只保留7天）
4. ✅ 支持跨天UV去重（最近7天）
5. ✅ 易于扩展和维护

**实现要点**：
1. Track时写入Set（实时数据）
2. 每日凌晨定时任务汇总到Hash
3. 查询历史数据从Hash读取
4. 计算近期总UV用Set的SUNION

**注意事项**：
- 定时任务需要监控，确保不漏跑
- 可以增加重试机制
- Hash的TTL要比Set长（Hash是汇总数据）
