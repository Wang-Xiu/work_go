package stats

import (
	"context"
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

	return redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
}

// TestTracker_Track 测试基础的Track功能
func TestTracker_Track(t *testing.T) {
	t.Parallel()

	// TODO: 你需要实现测试逻辑
	//
	// 测试场景1: 记录单次访问
	//   1. Track一次
	//   2. 查询PV应该为1
	//   3. 查询UV应该为1
	//
	// 测试场景2: 同一用户多次访问
	//   1. 同一个visitorID Track 3次
	//   2. PV应该为3
	//   3. UV应该为1（去重）
	//
	// 测试场景3: 不同用户访问
	//   1. 3个不同的visitorID各Track一次
	//   2. PV应该为3
	//   3. UV应该为3

	// ========================================
	// 👇 在这里实现你的测试代码
	// ========================================

	rdb := setupTestRedis(t)
	config := &Config{
		Enabled:         true,
		EnablePathStats: false,
		RetentionDays:   90,
	}
	tracker := NewTracker(rdb, config)

	ctx := context.Background()
	today := time.Now().Format("2006-01-02")

	// 示例测试框架：
	// 测试场景1: 单次访问
	// err := tracker.Track(ctx, "user1", "/api/test")
	// assert.NoError(t, err)
	//
	// stats, err := tracker.GetDailyStats(ctx, today)
	// assert.NoError(t, err)
	// assert.Equal(t, int64(1), stats.PV)
	// assert.Equal(t, int64(1), stats.UV)

	_ = tracker
	_ = ctx
	_ = today
}

// TestTracker_GetDailyStats 测试获取每日统计
func TestTracker_GetDailyStats(t *testing.T) {
	t.Parallel()

	// TODO: 你需要实现测试逻辑
	//
	// 测试场景1: 获取有数据的日期
	//   1. 先Track记录一些数据
	//   2. GetDailyStats查询
	//   3. 验证PV和UV正确
	//
	// 测试场景2: 获取没有数据的日期
	//   1. GetDailyStats查询未来的日期
	//   2. 应该返回PV=0, UV=0，不报错
	//
	// 测试场景3: 无效的日期格式
	//   1. GetDailyStats("invalid-date")
	//   2. 应该返回ErrInvalidDate错误

	// ========================================
	// 👇 在这里实现你的测试代码
	// ========================================

	rdb := setupTestRedis(t)
	config := &Config{Enabled: true}
	tracker := NewTracker(rdb, config)
	ctx := context.Background()

	_ = tracker
	_ = ctx
}

// TestTracker_GetRangeStats 测试获取范围统计
func TestTracker_GetRangeStats(t *testing.T) {
	t.Parallel()

	// TODO: 你需要实现测试逻辑
	//
	// 测试场景1: 查询3天的数据
	//   1. 分别在Day1、Day2、Day3记录数据
	//   2. GetRangeStats(Day1, Day3)
	//   3. 验证TotalPV = Day1.PV + Day2.PV + Day3.PV
	//   4. 验证DailyStats数组长度为3
	//
	// 测试场景2: UV的去重
	//   重点测试：同一用户在多天访问，总UV应该正确去重
	//   1. user1在Day1和Day2都访问
	//   2. user2只在Day1访问
	//   3. GetRangeStats(Day1, Day2)
	//   4. TotalUV应该为2，不是3
	//
	// 测试场景3: 日期范围无效
	//   1. startDate > endDate
	//   2. 应该返回错误

	// ========================================
	// 👇 在这里实现你的测试代码
	// ========================================

	rdb := setupTestRedis(t)
	config := &Config{Enabled: true}
	tracker := NewTracker(rdb, config)
	ctx := context.Background()

	_ = tracker
	_ = ctx
}

// TestTracker_PathStats 测试路径统计
func TestTracker_PathStats(t *testing.T) {
	t.Parallel()

	// TODO: 你需要实现测试逻辑
	//
	// 前提: EnablePathStats = true
	//
	// 测试场景1: 统计不同路径的访问
	//   1. Track("/api/users", visitor1)
	//   2. Track("/api/posts", visitor1)
	//   3. Track("/api/users", visitor2)
	//   4. GetPathStats查询
	//   5. 验证:
	//      - /api/users: PV=2, UV=2
	//      - /api/posts: PV=1, UV=1
	//
	// 测试场景2: 禁用路径统计
	//   1. EnablePathStats = false
	//   2. Track记录数据
	//   3. GetPathStats应该返回空数组

	// ========================================
	// 👇 在这里实现你的测试代码
	// ========================================

	rdb := setupTestRedis(t)
	config := &Config{
		Enabled:         true,
		EnablePathStats: true,
	}
	tracker := NewTracker(rdb, config)
	ctx := context.Background()

	_ = tracker
	_ = ctx
}

// TestTracker_ExcludePaths 测试排除路径
func TestTracker_ExcludePaths(t *testing.T) {
	t.Parallel()

	// TODO: 你需要实现测试逻辑
	//
	// 测试场景: 排除特定路径不统计
	//   1. ExcludePaths = ["/health", "/metrics"]
	//   2. Track("/health", visitor1)
	//   3. Track("/api/users", visitor1)
	//   4. GetDailyStats查询
	//   5. 验证: PV=1, UV=1（只统计了/api/users）

	// ========================================
	// 👇 在这里实现你的测试代码
	// ========================================

	rdb := setupTestRedis(t)
	config := &Config{
		Enabled:      true,
		ExcludePaths: []string{"/health", "/metrics"},
	}
	tracker := NewTracker(rdb, config)
	ctx := context.Background()

	_ = tracker
	_ = ctx
}

// TestTracker_CleanExpiredData 测试清理过期数据
func TestTracker_CleanExpiredData(t *testing.T) {
	t.Parallel()

	// TODO: 你需要实现测试逻辑
	//
	// 测试场景: 清理超过RetentionDays的数据
	//   1. RetentionDays = 7
	//   2. 手动创建8天前的Key
	//   3. 调用CleanExpiredData()
	//   4. 验证8天前的Key被删除
	//   5. 验证7天内的Key保留

	// ========================================
	// 👇 在这里实现你的测试代码
	// ========================================

	rdb := setupTestRedis(t)
	config := &Config{
		Enabled:       true,
		RetentionDays: 7,
	}
	tracker := NewTracker(rdb, config)
	ctx := context.Background()

	_ = tracker
	_ = ctx
}

// TestConfig_IsExcludedPath 测试路径排除逻辑
func TestConfig_IsExcludedPath(t *testing.T) {
	config := &Config{
		ExcludePaths: []string{"/health", "/metrics", "/favicon.ico"},
	}

	tests := []struct {
		path     string
		excluded bool
	}{
		{"/health", true},
		{"/metrics", true},
		{"/favicon.ico", true},
		{"/api/users", false},
		{"/api/posts", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := config.IsExcludedPath(tt.path)
			if cast.ToBool(result) != cast.ToBool(tt.excluded) {
				t.Errorf("expected %v, got %v", tt.excluded, result)
			}
		})
	}
}

// TestConfig_GetRetentionDays 测试获取保留天数
func TestConfig_GetRetentionDays(t *testing.T) {
	tests := []struct {
		name          string
		retentionDays int
		want          int
	}{
		{
			name:          "custom retention days",
			retentionDays: 30,
			want:          30,
		},
		{
			name:          "default retention days",
			retentionDays: 0,
			want:          90,
		},
		{
			name:          "negative value uses default",
			retentionDays: -1,
			want:          90,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{RetentionDays: tt.retentionDays}
			got := config.GetRetentionDays()
			if cast.ToInt(got) != cast.ToInt(tt.want) {
				t.Errorf("expected %d, got %d", tt.want, got)
			}
		})
	}
}

// BenchmarkTracker_Track 性能测试：Track操作
func BenchmarkTracker_Track(b *testing.B) {
	// TODO: 你需要实现性能测试
	//
	// 测试目标: Track操作的性能
	//   - 每次操作耗时应该 < 5ms
	//   - 并发场景下不应该有数据竞争
	//
	// 方法:
	//   b.RunParallel(func(pb *testing.PB) {
	//       for pb.Next() {
	//           tracker.Track(ctx, visitorID, path)
	//       }
	//   })

	// ========================================
	// 👇 在这里实现你的性能测试代码
	// ========================================

	mr, _ := miniredis.Run()
	defer mr.Close()

	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	config := &Config{Enabled: true}
	tracker := NewTracker(rdb, config)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = tracker.Track(ctx, "user1", "/api/test")
	}
}
