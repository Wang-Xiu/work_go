package kms

import (
	"context"
	"fmt"
	"strings"
)

// KMSProvider KMS提供商接口，支持不同云厂商的KMS服务
type KMSProvider interface {
	// Decrypt 解密密文
	Decrypt(ctx context.Context, ciphertext string) (string, error)
	// Encrypt 加密明文（用于生成配置文件中的密文）
	Encrypt(ctx context.Context, plaintext string) (string, error)
}

// Manager KMS管理器
type Manager struct {
	provider KMSProvider
	prefix   string // 配置文件中KMS密文的前缀，如 "kms://"
}

// NewManager 创建KMS管理器
// prefix: 配置文件中标识KMS密文的前缀，建议使用 "kms://" 或 "encrypted://"
func NewManager(provider KMSProvider, prefix string) *Manager {
	if prefix == "" {
		prefix = "kms://"
	}
	return &Manager{
		provider: provider,
		prefix:   prefix,
	}
}

// DecryptIfNeeded 判断值是否需要解密，如果需要则解密
// 如果值以prefix开头，则认为是KMS加密的密文，执行解密
// 否则直接返回原值
func (m *Manager) DecryptIfNeeded(ctx context.Context, value string) (string, error) {
	if !strings.HasPrefix(value, m.prefix) {
		return value, nil
	}

	// 去掉前缀，获取真正的密文
	ciphertext := strings.TrimPrefix(value, m.prefix)
	if ciphertext == "" {
		return "", fmt.Errorf("empty ciphertext after prefix")
	}

	plaintext, err := m.provider.Decrypt(ctx, ciphertext)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: %w", err)
	}

	return plaintext, nil
}

// Encrypt 加密明文，返回带前缀的密文
func (m *Manager) Encrypt(ctx context.Context, plaintext string) (string, error) {
	ciphertext, err := m.provider.Encrypt(ctx, plaintext)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt: %w", err)
	}
	return m.prefix + ciphertext, nil
}

// MockProvider 模拟KMS提供商（仅用于开发测试环境）
// 生产环境必须使用真实的KMS服务（如阿里云KMS、腾讯云KMS、AWS KMS等）
type MockProvider struct {
	// 简单的base64模拟加密，生产环境禁止使用！
}

func NewMockProvider() *MockProvider {
	return &MockProvider{}
}

func (m *MockProvider) Decrypt(ctx context.Context, ciphertext string) (string, error) {
	// 这里只是简单模拟，实际应该调用KMS SDK
	// 生产环境示例：
	// - 阿里云: https://help.aliyun.com/document_detail/28950.html
	// - 腾讯云: https://cloud.tencent.com/document/product/573
	// - AWS: https://docs.aws.amazon.com/kms/

	// 这里做简单的逆向处理（仅演示）
	decoded := reverseString(ciphertext)
	return decoded, nil
}

func (m *MockProvider) Encrypt(ctx context.Context, plaintext string) (string, error) {
	// 简单的反转字符串模拟加密（仅演示，生产禁用）
	encoded := reverseString(plaintext)
	return encoded, nil
}

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
