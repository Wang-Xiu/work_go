# KMS配置加密 - 使用指南

## 🔐 为什么需要KMS加密配置？

在配置文件中直接存储敏感信息（如Redis密码、数据库密码、API密钥等）是**非常危险**的：

- ❌ 配置文件通常会提交到Git仓库
- ❌ 多人协作时敏感信息会泄露
- ❌ 生产环境的密码容易被窃取

**解决方案**：使用KMS（Key Management Service）加密敏感信息。

---

## 🏗️ 架构设计

```
┌──────────────────────────────────────────────────────┐
│  config/app.yml                                      │
│  redis:                                              │
│    default:                                          │
│      password: "kms://encrypted_value"  ← 加密后的值 │
└──────────────────────────────────────────────────────┘
                      ↓
┌──────────────────────────────────────────────────────┐
│  config.Loader (配置加载器)                          │
│  1. 读取YAML文件                                      │
│  2. 解析为结构体                                      │
│  3. 递归查找 kms:// 开头的字段                        │
│  4. 调用KMS解密                                       │
└──────────────────────────────────────────────────────┘
                      ↓
┌──────────────────────────────────────────────────────┐
│  kms.Manager                                         │
│  - DecryptIfNeeded()                                 │
│  - 调用KMS Provider解密                               │
└──────────────────────────────────────────────────────┘
                      ↓
┌──────────────────────────────────────────────────────┐
│  KMS Provider                                        │
│  - MockProvider (开发测试)                           │
│  - 阿里云KMS (生产环境)                               │
│  - 腾讯云KMS (生产环境)                               │
│  - AWS KMS (生产环境)                                │
└──────────────────────────────────────────────────────┘
```

---

## 📝 使用方法

### 1. 配置文件格式

```yaml
# config/app.yml
redis:
  default:
    host: localhost
    port: 6379
    # 方式1：明文密码（仅限开发环境）
    password: "mypassword"

    # 方式2：KMS加密密码（生产环境推荐）
    # password: "kms://encrypted_value"

    db: 0
```

### 2. 在代码中使用（自动解密）

```go
package main

import (
    "context"
    "working-project/common/kms"
    "working-project/config"
    "working-project/common/infrastructure/redis"
)

func main() {
    ctx := context.Background()

    // 步骤1：创建KMS管理器
    kmsProvider := kms.NewMockProvider()  // 或者真实的KMS Provider
    kmsManager := kms.NewManager(kmsProvider, "kms://")

    // 步骤2：创建配置加载器
    configLoader := config.NewLoader(kmsManager)

    // 步骤3：加载配置（自动解密）
    var appConfig AppConfig
    err := configLoader.LoadFromFile(ctx, "config/app.yml", &appConfig)

    // 此时 appConfig.Redis["default"].Password 已经是解密后的明文密码了
    // 完全透明，业务代码无需关心加密细节
}
```

---

## 🔧 KMS Provider实现

### 开发环境：MockProvider

```go
// 仅用于开发测试，简单的字符串反转模拟加密
kmsProvider := kms.NewMockProvider()
kmsManager := kms.NewManager(kmsProvider, "kms://")

// 加密（生成配置文件中的密文）
encrypted, _ := kmsManager.Encrypt(ctx, "mypassword")
// 结果: "kms://drowssapym" (反转后的字符串)

// 解密（配置加载时自动调用）
decrypted, _ := kmsManager.DecryptIfNeeded(ctx, "kms://drowssapym")
// 结果: "mypassword"
```

### 生产环境：真实KMS Provider

```go
// 示例：阿里云KMS Provider（需要自己实现）
type AliyunKMSProvider struct {
    client *kmsclient.Client
    keyId  string
}

func (p *AliyunKMSProvider) Decrypt(ctx context.Context, ciphertext string) (string, error) {
    // 调用阿里云KMS SDK解密
    request := kms.CreateDecryptRequest()
    request.CiphertextBlob = ciphertext

    response, err := p.client.Decrypt(request)
    if err != nil {
        return "", err
    }

    return response.Plaintext, nil
}

func (p *AliyunKMSProvider) Encrypt(ctx context.Context, plaintext string) (string, error) {
    // 调用阿里云KMS SDK加密
    request := kms.CreateEncryptRequest()
    request.KeyId = p.keyId
    request.Plaintext = plaintext

    response, err := p.client.Encrypt(request)
    if err != nil {
        return "", err
    }

    return response.CiphertextBlob, nil
}

// 使用
func main() {
    // 初始化阿里云KMS客户端
    aliyunProvider := &AliyunKMSProvider{
        client: createAliyunKMSClient(),
        keyId:  "your-kms-key-id",
    }

    kmsManager := kms.NewManager(aliyunProvider, "kms://")
    configLoader := config.NewLoader(kmsManager)

    // 后续使用相同
}
```

---

## 🎯 生成加密密码

### 方法1：使用Go代码生成

```go
package main

import (
    "context"
    "fmt"
    "working-project/common/kms"
)

func main() {
    ctx := context.Background()

    // 创建KMS管理器
    kmsProvider := kms.NewMockProvider()  // 或真实Provider
    kmsManager := kms.NewManager(kmsProvider, "kms://")

    // 加密明文密码
    plaintext := "mypassword"
    encrypted, err := kmsManager.Encrypt(ctx, plaintext)
    if err != nil {
        panic(err)
    }

    fmt.Printf("明文: %s\n", plaintext)
    fmt.Printf("密文: %s\n", encrypted)
    // 输出: kms://drowssapym

    // 将密文复制到配置文件中
    fmt.Println("\n将以下内容复制到config/app.yml:")
    fmt.Printf("password: \"%s\"\n", encrypted)
}
```

### 方法2：使用云厂商CLI工具

```bash
# 阿里云KMS加密
aliyun kms Encrypt --KeyId your-key-id --Plaintext "mypassword"

# 腾讯云KMS加密
tccli kms Encrypt --KeyId your-key-id --Plaintext "mypassword"

# AWS KMS加密
aws kms encrypt --key-id your-key-id --plaintext "mypassword"
```

---

## 🔍 验证加密配置

### 示例配置文件

```yaml
# config/app.yml
redis:
  default:
    host: localhost
    port: 6379
    # MockProvider示例："mypassword" 反转后是 "drowssapym"
    password: "kms://drowssapym"
    db: 0
```

### 运行测试

```bash
# 运行示例程序
go run example/main_unified_redis.go

# 输出应该显示：
[1/6] 初始化KMS管理器...
✅ KMS管理器初始化成功 (使用MockProvider)

[2/6] 加载配置文件...
✅ 配置加载成功（已自动解密KMS加密字段）

[3/6] 初始化Redis连接池...
  - 注册Redis连接: default
    ✅ 连接成功: localhost:6379 (DB: 0)

# 如果能连接成功，说明密码已正确解密
```

---

## 🚨 安全最佳实践

### ✅ 推荐做法

1. **生产环境必须使用KMS**
   ```yaml
   # ✅ 正确
   password: "kms://encrypted_value"
   ```

2. **开发环境可以明文（但不提交）**
   ```yaml
   # ✅ 本地开发环境配置（不提交到Git）
   password: "dev_password_123"
   ```

3. **提交到Git的配置文件使用占位符**
   ```yaml
   # ✅ 提交到Git的模板
   password: "${REDIS_PASSWORD}"  # 或 "kms://please_replace_with_real_encrypted_value"
   ```

4. **敏感配置与代码分离**
   ```bash
   # .gitignore
   config/app.local.yml
   config/secrets.yml
   ```

### ❌ 禁止做法

1. **❌ 明文密码提交到Git**
   ```yaml
   # ❌ 绝对禁止
   password: "prod_redis_password_123"
   ```

2. **❌ 在代码中硬编码密码**
   ```go
   // ❌ 绝对禁止
   redisPassword := "prod_password_123"
   ```

3. **❌ 在日志中打印密码**
   ```go
   // ❌ 绝对禁止
   log.Printf("Redis密码: %s", config.Redis.Password)
   ```

---

## 📊 不同环境的配置策略

### 本地开发环境

```yaml
# config/app.local.yml (不提交到Git)
redis:
  default:
    host: localhost
    port: 6379
    password: ""  # 本地Redis无密码
```

### 测试环境

```yaml
# config/app.test.yml
redis:
  default:
    host: test-redis.internal
    port: 6379
    password: "kms://test_encrypted_password"  # 测试环境KMS加密
```

### 生产环境

```yaml
# config/app.prod.yml
redis:
  default:
    host: prod-redis.internal
    port: 6379
    password: "kms://prod_encrypted_password"  # 生产环境KMS加密
```

### 代码中根据环境加载

```go
env := os.Getenv("APP_ENV")  // development, test, production
if env == "" {
    env = "development"
}

configFile := fmt.Sprintf("config/app.%s.yml", env)
err := configLoader.LoadFromFile(ctx, configFile, &appConfig)
```

---

## 🎓 原理说明

### config.Loader的解密流程

1. **读取YAML文件**
   ```go
   data, _ := os.ReadFile("config/app.yml")
   yaml.Unmarshal(data, &config)
   ```

2. **递归遍历配置结构体**
   - 遍历所有struct字段
   - 遍历所有slice/array元素
   - 遍历所有map值

3. **检查字符串字段**
   ```go
   if strings.HasPrefix(value, "kms://") {
       decrypted, _ := kmsManager.DecryptIfNeeded(ctx, value)
       // 替换为解密后的明文
   }
   ```

4. **返回解密后的配置**
   - 业务代码无需关心加密细节
   - 完全透明的解密过程

---

## 📦 完整示例

```go
package main

import (
    "context"
    "log"

    "working-project/common/infrastructure/redis"
    "working-project/common/kms"
    "working-project/common/middleware/ratelimit"
    "working-project/common/middleware/stats"
    "working-project/config"
)

type AppConfig struct {
    Redis      map[string]redis.Config `yaml:"redis"`
    Middleware struct {
        RateLimit ratelimit.Config `yaml:"ratelimit"`
        Stats     stats.Config     `yaml:"stats"`
    } `yaml:"middleware"`
}

func main() {
    ctx := context.Background()

    // 1. 初始化KMS（根据环境选择Provider）
    var kmsProvider kms.KMSProvider
    if isProduction() {
        kmsProvider = createProductionKMSProvider()
    } else {
        kmsProvider = kms.NewMockProvider()
    }
    kmsManager := kms.NewManager(kmsProvider, "kms://")

    // 2. 加载配置（自动解密）
    configLoader := config.NewLoader(kmsManager)
    var appConfig AppConfig
    configLoader.LoadFromFile(ctx, "config/app.yml", &appConfig)

    // 3. 使用配置（密码已解密）
    redisManager := redis.GetGlobalManager()
    for name, cfg := range appConfig.Redis {
        redisCfg := cfg
        redisManager.Register(name, &redisCfg)
    }

    // 4. 创建中间件（共享Redis连接）
    limiter, _ := ratelimit.NewLimiterFromManager(redisManager, &appConfig.Middleware.RateLimit)
    tracker, _ := stats.NewTrackerFromManager(redisManager, &appConfig.Middleware.Stats)

    // ... 启动应用
}
```

---

## ✅ 总结

### 核心优势

- ✅ **安全**：敏感信息加密存储
- ✅ **透明**：业务代码无需关心加密细节
- ✅ **灵活**：支持多种KMS Provider
- ✅ **自动**：配置加载时自动解密

### 使用流程

```
1. 创建KMS Provider (开发用Mock，生产用真实KMS)
   ↓
2. 创建KMS Manager
   ↓
3. 创建Config Loader (注入KMS Manager)
   ↓
4. 加载配置文件 (自动解密 kms:// 开头的字段)
   ↓
5. 使用解密后的配置 (完全透明)
```

### 关键要点

1. **配置文件中**：使用 `kms://encrypted_value` 格式
2. **代码中**：使用 `config.Loader` 加载配置
3. **生产环境**：必须使用真实的KMS Provider
4. **安全规范**：永远不要将明文密码提交到Git

---

**现在你的配置文件是安全的了！** 🔐
