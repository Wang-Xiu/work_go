package stats

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
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

// NewTrackerFromConfig 从配置创建追踪器（会自动创建Redis客户端）
func NewTrackerFromConfig(config *Config) (*Tracker, error) {
	if config == nil {
		return nil, ErrInvalidConfig
	}

	// 创建Redis客户端
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port),
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
		PoolSize: config.Redis.PoolSize,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), config.Redis.Timeout)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRedisUnavailable, err)
	}

	return NewTracker(rdb, config), nil
}

// Track 记录一次访问
// 参数:
//   - ctx: 上下文
//   - visitorID: 访客唯一标识（如IP、UserID等）
//   - path: 访问的路径
//
// 功能:
//   - PV计数+1（使用INCR）
//   - UV去重计数（使用Set的SADD，自动去重）
func (t *Tracker) Track(ctx context.Context, visitorID, path string) error {
	// 如果未启用统计，直接返回
	if !t.config.Enabled {
		return nil
	}

	// 如果路径在排除列表中，不统计
	if t.config.IsExcludedPath(path) {
		return nil
	}

	// 获取当前日期（格式：YYYY-MM-DD）
	date := time.Now().Format("2006-01-02")
	expire := t.getExpireTime()

	pipe := t.redis.Pipeline()

	// 1. 全站PV统计（String类型，INCR原子递增）
	pvKey := t.buildPVKey(date, "")
	pipe.Incr(ctx, pvKey)
	pipe.Expire(ctx, pvKey, expire)

	// 2. 全站UV统计（Set类型，SADD自动去重）
	uvKey := t.buildUVKey(date, "")
	pipe.SAdd(ctx, uvKey, visitorID)
	pipe.Expire(ctx, uvKey, expire)

	// 3. 如果启用了路径统计
	if t.config.EnablePathStats {
		// 路径PV
		pathPVKey := t.buildPVKey(date, path)
		pipe.Incr(ctx, pathPVKey)
		pipe.Expire(ctx, pathPVKey, expire)

		// 路径UV
		pathUVKey := t.buildUVKey(date, path)
		pipe.SAdd(ctx, pathUVKey, visitorID)
		pipe.Expire(ctx, pathUVKey, expire)

		// 记录路径到路径列表（用于查询）
		pathsKey := t.buildPathsKey(date)
		pipe.SAdd(ctx, pathsKey, path)
		pipe.Expire(ctx, pathsKey, expire)
	}

	// 执行Pipeline
	_, err := pipe.Exec(ctx)
	return err
}

// GetDailyStats 获取某一天的统计数据
// 参数:
//   - ctx: 上下文
//   - date: 日期，格式：YYYY-MM-DD
//
// 返回:
//   - DailyStats: 包含PV和UV的统计数据
func (t *Tracker) GetDailyStats(ctx context.Context, date string) (*DailyStats, error) {
	// 验证日期格式
	_, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, ErrInvalidDate
	}

	// 查询PV（String类型，GET命令，O(1)）
	pvKey := t.buildPVKey(date, "")
	pvStr, err := t.redis.Get(ctx, pvKey).Result()
	if err == redis.Nil {
		pvStr = "0"
	} else if err != nil {
		return nil, err
	}
	pv, _ := strconv.ParseInt(pvStr, 10, 64)

	// 查询UV（Set类型，SCARD命令，O(1)）
	uvKey := t.buildUVKey(date, "")
	uv, err := t.redis.SCard(ctx, uvKey).Result()
	if err == redis.Nil {
		uv = 0
	} else if err != nil {
		return nil, err
	}

	return &DailyStats{
		Date: date,
		PV:   pv,
		UV:   uv,
	}, nil
}

// GetRangeStats 获取时间范围内的统计数据
// 参数:
//   - ctx: 上下文
//   - startDate: 开始日期，格式：YYYY-MM-DD
//   - endDate: 结束日期，格式：YYYY-MM-DD
//
// 返回:
//   - RangeStats: 包含总计和每日明细
//
// 注意:
//   - TotalPV 是每天PV的累加
//   - TotalUV 是每天UV的累加（会重复计数，如需真实UV请使用GetUniqueUVInRange）
func (t *Tracker) GetRangeStats(ctx context.Context, startDate, endDate string) (*RangeStats, error) {
	// 验证日期格式和范围
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

	var dailyStats []DailyStats
	var totalPV, totalUV int64

	// 遍历日期范围，查询每一天的数据
	for current := start; !current.After(end); current = current.AddDate(0, 0, 1) {
		dateStr := current.Format("2006-01-02")

		// 查询每天的PV和UV
		stats, err := t.GetDailyStats(ctx, dateStr)
		if err != nil {
			return nil, err
		}

		dailyStats = append(dailyStats, *stats)
		totalPV += stats.PV
		totalUV += stats.UV
	}

	return &RangeStats{
		StartDate:  startDate,
		EndDate:    endDate,
		TotalPV:    totalPV,
		TotalUV:    totalUV, // 注意：这是简单累加，会重复计数
		DailyStats: dailyStats,
	}, nil
}

// GetUniqueUVInRange 获取时间范围内的真实UV（去重后）
// 参数:
//   - ctx: 上下文
//   - startDate: 开始日期
//   - endDate: 结束日期
//
// 返回:
//   - 去重后的真实独立访客数
//
// 注意:
//   - 使用SUNION合并多天的Set并自动去重
//   - 数据量大时较慢，建议查询天数不超过30天
//   - 示例：用户A在1号和2号都访问，只计数1次
func (t *Tracker) GetUniqueUVInRange(ctx context.Context, startDate, endDate string) (int64, error) {
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return 0, ErrInvalidDate
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return 0, ErrInvalidDate
	}

	// 收集所有UV Set的Key
	var keys []string
	for current := start; !current.After(end); current = current.AddDate(0, 0, 1) {
		dateStr := current.Format("2006-01-02")
		keys = append(keys, t.buildUVKey(dateStr, ""))
	}

	// 使用SUNION合并并去重
	members, err := t.redis.SUnion(ctx, keys...).Result()
	if err != nil && err != redis.Nil {
		return 0, err
	}

	return int64(len(members)), nil
}

// GetPathStats 获取某天各路径的统计数据
// 参数:
//   - ctx: 上下文
//   - date: 日期，格式：YYYY-MM-DD
//
// 返回:
//   - []PathStats: 所有路径的统计数据列表
//
// 注意:
//   - 只有当EnablePathStats=true时才有数据
func (t *Tracker) GetPathStats(ctx context.Context, date string) ([]PathStats, error) {
	if !t.config.EnablePathStats {
		return []PathStats{}, nil
	}

	// 验证日期格式
	_, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, ErrInvalidDate
	}

	// 获取路径列表
	pathsKey := t.buildPathsKey(date)
	paths, err := t.redis.SMembers(ctx, pathsKey).Result()
	if err == redis.Nil {
		return []PathStats{}, nil
	} else if err != nil {
		return nil, err
	}

	// 查询每个路径的PV和UV
	var result []PathStats

	for _, path := range paths {
		// 查询路径PV
		pvKey := t.buildPVKey(date, path)
		pvStr, _ := t.redis.Get(ctx, pvKey).Result()
		pv, _ := strconv.ParseInt(pvStr, 10, 64)

		// 查询路径UV
		uvKey := t.buildUVKey(date, path)
		uv, _ := t.redis.SCard(ctx, uvKey).Result()

		result = append(result, PathStats{
			Path: path,
			PV:   pv,
			UV:   uv,
		})
	}

	return result, nil
}

// CleanExpiredData 清理过期数据
// 该方法应该由定时任务调用（如每天凌晨执行）
//
// 功能:
//   - 删除超过RetentionDays天的统计数据
//   - 释放Redis内存
//
// 注意:
//   - 由于设置了TTL，数据会自动过期
//   - 此方法主要用于手动清理或修正过期时间
func (t *Tracker) CleanExpiredData(ctx context.Context) error {
	retentionDays := t.config.GetRetentionDays()
	expireDate := time.Now().AddDate(0, 0, -retentionDays)

	var cursor uint64
	var deletedCount int

	// 使用SCAN遍历所有stats:*的Key
	for {
		keys, newCursor, err := t.redis.Scan(ctx, cursor, "stats:*", 100).Result()
		if err != nil {
			return err
		}

		// 检查每个Key是否过期
		for _, key := range keys {
			// 从Key中提取日期
			// stats:pv:2024-01-15 -> 2024-01-15
			// stats:uv:2024-01-15 -> 2024-01-15
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
			break
		}
	}

	return nil
}

// Close 关闭追踪器，释放Redis连接
func (t *Tracker) Close() error {
	if t.redis != nil {
		return t.redis.Close()
	}
	return nil
}

// ========================================
// 辅助函数（私有方法）
// ========================================

// buildPVKey 构建PV统计的Redis Key
// 格式: stats:pv:YYYY-MM-DD 或 stats:pv:YYYY-MM-DD:/api/users
func (t *Tracker) buildPVKey(date, path string) string {
	if path == "" {
		return fmt.Sprintf("stats:pv:%s", date)
	}
	return fmt.Sprintf("stats:pv:%s:%s", date, path)
}

// buildUVKey 构建UV统计的Redis Key
// 格式: stats:uv:YYYY-MM-DD 或 stats:uv:YYYY-MM-DD:/api/users
func (t *Tracker) buildUVKey(date, path string) string {
	if path == "" {
		return fmt.Sprintf("stats:uv:%s", date)
	}
	return fmt.Sprintf("stats:uv:%s:%s", date, path)
}

// buildPathsKey 构建路径列表的Redis Key
// 格式: stats:paths:YYYY-MM-DD
func (t *Tracker) buildPathsKey(date string) string {
	return fmt.Sprintf("stats:paths:%s", date)
}

// getExpireTime 获取Key的过期时间（秒）
func (t *Tracker) getExpireTime() time.Duration {
	days := t.config.GetRetentionDays()
	return time.Duration(days) * 24 * time.Hour
}
