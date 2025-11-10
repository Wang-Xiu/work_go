package redis

import (
	"context"
	"fmt"
	"sync"

	"github.com/redis/go-redis/v9"
)

// Manager Redis连接池管理器
// 负责管理多个命名的Redis客户端实例，支持单例模式
type Manager struct {
	mu      sync.RWMutex
	clients map[string]*redis.Client
}

var (
	globalManager     *Manager
	globalManagerOnce sync.Once
)

// GetGlobalManager 获取全局Manager单例
func GetGlobalManager() *Manager {
	globalManagerOnce.Do(func() {
		globalManager = &Manager{
			clients: make(map[string]*redis.Client),
		}
	})
	return globalManager
}

// NewManager 创建新的Manager实例（用于测试或特殊场景）
func NewManager() *Manager {
	return &Manager{
		clients: make(map[string]*redis.Client),
	}
}

// Register 注册一个命名的Redis客户端
func (m *Manager) Register(name string, cfg *Config) error {
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid redis config for '%s': %w", name, err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:         cfg.GetAddress(),
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		MaxRetries:   cfg.MaxRetries,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), cfg.DialTimeout)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return fmt.Errorf("failed to connect to redis '%s': %w", name, err)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// 如果已存在同名连接，先关闭旧的
	if old, exists := m.clients[name]; exists {
		_ = old.Close()
	}

	m.clients[name] = client
	return nil
}

// Get 获取指定名称的Redis客户端
func (m *Manager) Get(name string) (*redis.Client, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	client, exists := m.clients[name]
	if !exists {
		return nil, fmt.Errorf("redis client '%s' not found", name)
	}
	return client, nil
}

// MustGet 获取Redis客户���，不存在则panic（适合启动阶段）
func (m *Manager) MustGet(name string) *redis.Client {
	client, err := m.Get(name)
	if err != nil {
		panic(err)
	}
	return client
}

// Close 关闭指定名称的Redis客户端
func (m *Manager) Close(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	client, exists := m.clients[name]
	if !exists {
		return fmt.Errorf("redis client '%s' not found", name)
	}

	err := client.Close()
	delete(m.clients, name)
	return err
}

// CloseAll 关闭所有Redis客户端
func (m *Manager) CloseAll() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var errs []error
	for name, client := range m.clients {
		if err := client.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close redis '%s': %w", name, err))
		}
	}

	m.clients = make(map[string]*redis.Client)

	if len(errs) > 0 {
		return fmt.Errorf("errors closing redis clients: %v", errs)
	}
	return nil
}

// List 列出所有已注册的Redis客户端名称
func (m *Manager) List() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	names := make([]string, 0, len(m.clients))
	for name := range m.clients {
		names = append(names, name)
	}
	return names
}
