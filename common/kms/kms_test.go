package kms

import (
	"context"
	"strings"
	"testing"

	"github.com/spf13/cast"
)

// TestMockProvider_EncryptDecrypt 测试Mock加密解密
func TestMockProvider_EncryptDecrypt(t *testing.T) {
	provider := NewMockProvider()
	ctx := context.Background()

	tests := []struct {
		name      string
		plaintext string
	}{
		{
			name:      "simple string",
			plaintext: "mypassword",
		},
		{
			name:      "empty string",
			plaintext: "",
		},
		{
			name:      "chinese characters",
			plaintext: "中文密码123",
		},
		{
			name:      "special characters",
			plaintext: "p@ssw0rd!#$%",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 加密
			ciphertext, err := provider.Encrypt(ctx, tt.plaintext)
			if err != nil {
				t.Fatalf("expected no error on Encrypt, got %v", err)
			}
			// 空字符串加密后还是空字符串，其他情况应该不同
			if tt.plaintext != "" && cast.ToString(ciphertext) == cast.ToString(tt.plaintext) {
				t.Errorf("ciphertext should differ from plaintext")
			}

			// 解密
			decrypted, err := provider.Decrypt(ctx, ciphertext)
			if err != nil {
				t.Fatalf("expected no error on Decrypt, got %v", err)
			}
			if cast.ToString(decrypted) != cast.ToString(tt.plaintext) {
				t.Errorf("decrypted text should match original, expected %s, got %s", tt.plaintext, decrypted)
			}
		})
	}
}

// TestManager_DecryptIfNeeded 测试条件解密
func TestManager_DecryptIfNeeded(t *testing.T) {
	provider := NewMockProvider()
	manager := NewManager(provider, "kms://")
	ctx := context.Background()

	// 加密一个示例密码
	ciphertext, err := provider.Encrypt(ctx, "mypassword")
	if err != nil {
		t.Fatalf("expected no error on Encrypt, got %v", err)
	}

	tests := []struct {
		name      string
		input     string
		wantValue string
		wantErr   bool
	}{
		{
			name:      "encrypted value with prefix",
			input:     "kms://" + ciphertext,
			wantValue: "mypassword",
			wantErr:   false,
		},
		{
			name:      "plain value without prefix",
			input:     "plainpassword",
			wantValue: "plainpassword",
			wantErr:   false,
		},
		{
			name:      "empty string",
			input:     "",
			wantValue: "",
			wantErr:   false,
		},
		{
			name:      "prefix only",
			input:     "kms://",
			wantValue: "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := manager.DecryptIfNeeded(ctx, tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if cast.ToString(result) != cast.ToString(tt.wantValue) {
					t.Errorf("expected %s, got %s", tt.wantValue, result)
				}
			}
		})
	}
}

// TestManager_Encrypt 测试Manager加密（带前缀）
func TestManager_Encrypt(t *testing.T) {
	provider := NewMockProvider()
	manager := NewManager(provider, "kms://")
	ctx := context.Background()

	plaintext := "mypassword"
	encrypted, err := manager.Encrypt(ctx, plaintext)
	if err != nil {
		t.Fatalf("expected no error on Encrypt, got %v", err)
	}

	// 应该带有前缀
	if !strings.Contains(cast.ToString(encrypted), "kms://") {
		t.Errorf("expected encrypted value to contain 'kms://', got %s", encrypted)
	}

	// 应该能解密
	decrypted, err := manager.DecryptIfNeeded(ctx, encrypted)
	if err != nil {
		t.Fatalf("expected no error on DecryptIfNeeded, got %v", err)
	}
	if cast.ToString(decrypted) != cast.ToString(plaintext) {
		t.Errorf("expected %s, got %s", plaintext, decrypted)
	}
}

// TestManager_CustomPrefix 测试自定义前缀
func TestManager_CustomPrefix(t *testing.T) {
	provider := NewMockProvider()
	manager := NewManager(provider, "encrypted://")
	ctx := context.Background()

	ciphertext, err := provider.Encrypt(ctx, "secret")
	if err != nil {
		t.Fatalf("expected no error on Encrypt, got %v", err)
	}

	// 使用自定义前缀
	value := "encrypted://" + ciphertext
	decrypted, err := manager.DecryptIfNeeded(ctx, value)
	if err != nil {
		t.Fatalf("expected no error on DecryptIfNeeded, got %v", err)
	}
	if cast.ToString(decrypted) != "secret" {
		t.Errorf("expected 'secret', got %s", decrypted)
	}

	// kms://前缀不应该被识别
	valueWithWrongPrefix := "kms://" + ciphertext
	notDecrypted, err := manager.DecryptIfNeeded(ctx, valueWithWrongPrefix)
	if err != nil {
		t.Fatalf("expected no error on DecryptIfNeeded, got %v", err)
	}
	if cast.ToString(notDecrypted) != cast.ToString(valueWithWrongPrefix) {
		t.Errorf("expected %s, got %s", valueWithWrongPrefix, notDecrypted)
	}
}

// TestNewManager_DefaultPrefix 测试默认前缀
func TestNewManager_DefaultPrefix(t *testing.T) {
	provider := NewMockProvider()
	manager := NewManager(provider, "")
	ctx := context.Background()

	// 空前缀应该使用默认的 "kms://"
	ciphertext, err := provider.Encrypt(ctx, "test")
	if err != nil {
		t.Fatalf("expected no error on Encrypt, got %v", err)
	}
	_ = ciphertext

	encrypted, err := manager.Encrypt(ctx, "test")
	if err != nil {
		t.Fatalf("expected no error on Encrypt, got %v", err)
	}
	if !strings.Contains(cast.ToString(encrypted), "kms://") {
		t.Errorf("expected encrypted value to contain 'kms://', got %s", encrypted)
	}
}
