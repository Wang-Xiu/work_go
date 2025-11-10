package main

import (
	"context"
	"fmt"
	"os"

	"working-project/common/kms"
)

// 工具：加密密码，用于生成配置文件中的KMS密文
// 用法: go run tools/encrypt_password.go <明文密码>

func main() {
	if len(os.Args) < 2 {
		fmt.Println("❌ 错误：缺少明文密码参数")
		fmt.Println()
		fmt.Println("用法:")
		fmt.Println("  go run tools/encrypt_password.go <明文密码>")
		fmt.Println()
		fmt.Println("示例:")
		fmt.Println("  go run tools/encrypt_password.go mypassword")
		fmt.Println()
		fmt.Println("生产环境提示:")
		fmt.Println("  当前使用MockProvider（字符串反转），仅用于开发测试")
		fmt.Println("  生产环境请修改代码，使用真实的KMS Provider")
		fmt.Println("  - 阿里云KMS: https://help.aliyun.com/document_detail/28950.html")
		fmt.Println("  - 腾讯云KMS: https://cloud.tencent.com/document/product/573")
		fmt.Println("  - AWS KMS: https://docs.aws.amazon.com/kms/")
		os.Exit(1)
	}

	plaintext := os.Args[1]
	ctx := context.Background()

	// 创建KMS管理器（开发环境使用MockProvider）
	kmsProvider := kms.NewMockProvider()
	kmsManager := kms.NewManager(kmsProvider, "kms://")

	// 加密
	encrypted, err := kmsManager.Encrypt(ctx, plaintext)
	if err != nil {
		fmt.Printf("❌ 加密失败: %v\n", err)
		os.Exit(1)
	}

	// 输出结果
	fmt.Println("========================================")
	fmt.Println("KMS密码加密工具")
	fmt.Println("========================================")
	fmt.Println()
	fmt.Printf("📝 明文密码: %s\n", plaintext)
	fmt.Printf("🔐 加密密文: %s\n", encrypted)
	fmt.Println()
	fmt.Println("📋 将以下内容复制到 config/app.yml:")
	fmt.Println("----------------------------------------")
	fmt.Printf("password: \"%s\"\n", encrypted)
	fmt.Println("----------------------------------------")
	fmt.Println()
	fmt.Println("✅ 加密完成！")
	fmt.Println()
	fmt.Println("⚠️  注意:")
	fmt.Println("  - 当前使用MockProvider（开发环境）")
	fmt.Println("  - 生产环境请使用真实的KMS服务")
	fmt.Println("  - 详细说明请查看: docs/KMS_CONFIG_GUIDE.md")
	fmt.Println()
}
