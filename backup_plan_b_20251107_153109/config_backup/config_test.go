package config

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"working-project/common/kms"

	"github.com/spf13/cast"
)

// TestLoader_LoadFromFile 测试从文件加载配置
func TestLoader_LoadFromFile(t *testing.T) {
	// 创建临时配置文件
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test_config.yml")

	yamlContent := `
database:
  host: localhost
  port: 3306
  password: "kms://drowssapym"  # "mypassword"的Mock加密
  username: root

redis:
  host: 127.0.0.1
  port: 6379
  password: "plain_password"
`

	err := os.WriteFile(configPath, []byte(yamlContent), 0644)
	if err != nil {
		t.Fatalf("expected no error on WriteFile, got %v", err)
	}

	// 测试配置结构
	type DatabaseConfig struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Password string `yaml:"password"`
		Username string `yaml:"username"`
	}

	type RedisConfig struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Password string `yaml:"password"`
	}

	type TestConfig struct {
		Database DatabaseConfig `yaml:"database"`
		Redis    RedisConfig    `yaml:"redis"`
	}

	// 创建KMS管理器
	kmsProvider := kms.NewMockProvider()
	kmsManager := kms.NewManager(kmsProvider, "kms://")
	loader := NewLoader(kmsManager)

	// 加载配置
	var config TestConfig
	err = loader.LoadFromFile(context.Background(), configPath, &config)
	if err != nil {
		t.Fatalf("expected no error on LoadFromFile, got %v", err)
	}

	// 验证配置
	if cast.ToString(config.Database.Host) != "localhost" {
		t.Errorf("expected database host 'localhost', got %s", config.Database.Host)
	}
	if cast.ToInt(config.Database.Port) != 3306 {
		t.Errorf("expected database port 3306, got %d", config.Database.Port)
	}
	if cast.ToString(config.Database.Password) != "mypassword" {
		t.Errorf("expected database password 'mypassword', got %s", config.Database.Password)
	}
	if cast.ToString(config.Database.Username) != "root" {
		t.Errorf("expected database username 'root', got %s", config.Database.Username)
	}

	if cast.ToString(config.Redis.Host) != "127.0.0.1" {
		t.Errorf("expected redis host '127.0.0.1', got %s", config.Redis.Host)
	}
	if cast.ToInt(config.Redis.Port) != 6379 {
		t.Errorf("expected redis port 6379, got %d", config.Redis.Port)
	}
	if cast.ToString(config.Redis.Password) != "plain_password" {
		t.Errorf("expected redis password 'plain_password', got %s", config.Redis.Password)
	}
}

// TestLoader_LoadFromFile_NoKMS 测试不使用KMS加载配置
func TestLoader_LoadFromFile_NoKMS(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test_config.yml")

	yamlContent := `
server:
  port: 8080
  host: localhost
`

	err := os.WriteFile(configPath, []byte(yamlContent), 0644)
	if err != nil {
		t.Fatalf("expected no error on WriteFile, got %v", err)
	}

	type ServerConfig struct {
		Port int    `yaml:"port"`
		Host string `yaml:"host"`
	}

	type TestConfig struct {
		Server ServerConfig `yaml:"server"`
	}

	// 不使用KMS
	loader := NewLoader(nil)

	var config TestConfig
	err = loader.LoadFromFile(context.Background(), configPath, &config)
	if err != nil {
		t.Fatalf("expected no error on LoadFromFile, got %v", err)
	}

	if cast.ToInt(config.Server.Port) != 8080 {
		t.Errorf("expected port 8080, got %d", config.Server.Port)
	}
	if cast.ToString(config.Server.Host) != "localhost" {
		t.Errorf("expected host 'localhost', got %s", config.Server.Host)
	}
}

// TestLoader_LoadFromFile_InvalidPath 测试加载不存在的文件
func TestLoader_LoadFromFile_InvalidPath(t *testing.T) {
	loader := NewLoader(nil)

	var config struct{}
	err := loader.LoadFromFile(context.Background(), "/nonexistent/path.yml", &config)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to read config file") {
		t.Errorf("expected error to contain 'failed to read config file', got %v", err)
	}
}

// TestLoader_LoadFromFile_InvalidYAML 测试加载无效的YAML
func TestLoader_LoadFromFile_InvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "invalid.yml")

	invalidYaml := `
server:
  port: 8080
  host: [invalid yaml
`

	err := os.WriteFile(configPath, []byte(invalidYaml), 0644)
	if err != nil {
		t.Fatalf("expected no error on WriteFile, got %v", err)
	}

	loader := NewLoader(nil)

	var config struct{}
	err = loader.LoadFromFile(context.Background(), configPath, &config)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to parse yaml") {
		t.Errorf("expected error to contain 'failed to parse yaml', got %v", err)
	}
}

// TestLoader_DecryptConfig_NestedStruct 测试嵌套结构体解密
func TestLoader_DecryptConfig_NestedStruct(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "nested.yml")

	yamlContent := `
app:
  name: myapp
  secrets:
    api_key: "kms://yek_ipa"  # "api_key"的Mock加密
    secret_key: "kms://yek_terces"  # "secret_key"的Mock加密
  database:
    password: "kms://drowssapbd"  # "dbpassword"的Mock加密
`

	err := os.WriteFile(configPath, []byte(yamlContent), 0644)
	if err != nil {
		t.Fatalf("expected no error on WriteFile, got %v", err)
	}

	type SecretsConfig struct {
		APIKey    string `yaml:"api_key"`
		SecretKey string `yaml:"secret_key"`
	}

	type DatabaseConfig struct {
		Password string `yaml:"password"`
	}

	type AppConfig struct {
		Name     string         `yaml:"name"`
		Secrets  SecretsConfig  `yaml:"secrets"`
		Database DatabaseConfig `yaml:"database"`
	}

	type TestConfig struct {
		App AppConfig `yaml:"app"`
	}

	kmsProvider := kms.NewMockProvider()
	kmsManager := kms.NewManager(kmsProvider, "kms://")
	loader := NewLoader(kmsManager)

	var config TestConfig
	err = loader.LoadFromFile(context.Background(), configPath, &config)
	if err != nil {
		t.Fatalf("expected no error on LoadFromFile, got %v", err)
	}

	// 验证嵌套结构体的字段都被正确解密
	if cast.ToString(config.App.Name) != "myapp" {
		t.Errorf("expected app name 'myapp', got %s", config.App.Name)
	}
	if cast.ToString(config.App.Secrets.APIKey) != "api_key" {
		t.Errorf("expected api_key 'api_key', got %s", config.App.Secrets.APIKey)
	}
	if cast.ToString(config.App.Secrets.SecretKey) != "secret_key" {
		t.Errorf("expected secret_key 'secret_key', got %s", config.App.Secrets.SecretKey)
	}
	if cast.ToString(config.App.Database.Password) != "dbpassword" {
		t.Errorf("expected database password 'dbpassword', got %s", config.App.Database.Password)
	}
}

// TestLoader_DecryptConfig_SliceAndMap 测试切片和Map中的解密
func TestLoader_DecryptConfig_SliceAndMap(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "collection.yml")

	yamlContent := `
passwords:
  - "kms://1drowssap"  # "password1"
  - "plain_password2"
  - "kms://3drowssap"  # "password3"

credentials:
  admin: "kms://drowssap_nimda"  # "admin_password"
  user: "plain_user_password"
`

	err := os.WriteFile(configPath, []byte(yamlContent), 0644)
	if err != nil {
		t.Fatalf("expected no error on WriteFile, got %v", err)
	}

	type TestConfig struct {
		Passwords   []string          `yaml:"passwords"`
		Credentials map[string]string `yaml:"credentials"`
	}

	kmsProvider := kms.NewMockProvider()
	kmsManager := kms.NewManager(kmsProvider, "kms://")
	loader := NewLoader(kmsManager)

	var config TestConfig
	err = loader.LoadFromFile(context.Background(), configPath, &config)
	if err != nil {
		t.Fatalf("expected no error on LoadFromFile, got %v", err)
	}

	// 验证切片中的字段
	if cast.ToInt(len(config.Passwords)) != 3 {
		t.Errorf("expected 3 passwords, got %d", len(config.Passwords))
	}
	if cast.ToString(config.Passwords[0]) != "password1" {
		t.Errorf("expected password1, got %s", config.Passwords[0])
	}
	if cast.ToString(config.Passwords[1]) != "plain_password2" {
		t.Errorf("expected plain_password2, got %s", config.Passwords[1])
	}
	if cast.ToString(config.Passwords[2]) != "password3" {
		t.Errorf("expected password3, got %s", config.Passwords[2])
	}

	// 验证Map中的字段
	if cast.ToString(config.Credentials["admin"]) != "admin_password" {
		t.Errorf("expected admin_password, got %s", config.Credentials["admin"])
	}
	if cast.ToString(config.Credentials["user"]) != "plain_user_password" {
		t.Errorf("expected plain_user_password, got %s", config.Credentials["user"])
	}
}

// TestLoader_MustLoadFromFile_Success 测试MustLoadFromFile成功场景
func TestLoader_MustLoadFromFile_Success(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test.yml")

	yamlContent := `
port: 8080
`

	err := os.WriteFile(configPath, []byte(yamlContent), 0644)
	if err != nil {
		t.Fatalf("expected no error on WriteFile, got %v", err)
	}

	type TestConfig struct {
		Port int `yaml:"port"`
	}

	loader := NewLoader(nil)

	var config TestConfig
	// 应该不会panic
	didPanic := false
	func() {
		defer func() {
			if r := recover(); r != nil {
				didPanic = true
			}
		}()
		loader.MustLoadFromFile(context.Background(), configPath, &config)
	}()

	if didPanic {
		t.Errorf("expected no panic, but panicked")
	}

	if cast.ToInt(config.Port) != 8080 {
		t.Errorf("expected port 8080, got %d", config.Port)
	}
}

// TestLoader_MustLoadFromFile_Panic 测试MustLoadFromFile失败场景
func TestLoader_MustLoadFromFile_Panic(t *testing.T) {
	loader := NewLoader(nil)

	var config struct{}

	// 应该panic
	didPanic := false
	func() {
		defer func() {
			if r := recover(); r != nil {
				didPanic = true
			}
		}()
		loader.MustLoadFromFile(context.Background(), "/nonexistent/path.yml", &config)
	}()

	if !didPanic {
		t.Errorf("expected panic, but did not panic")
	}
}
