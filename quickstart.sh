#!/bin/bash

# 快速开始脚本 - 用于快速验证限流器功能

set -e

echo "================================"
echo "限流器快速测试脚本"
echo "================================"
echo ""

# 检查Redis是否运行
echo "[1/4] 检查Redis连接..."
if ! command -v redis-cli &> /dev/null; then
    echo "❌ 未找到redis-cli命令"
    echo "请先安装Redis: brew install redis (macOS)"
    exit 1
fi

# 尝试连接Redis
if ! redis-cli -h localhost -p 6379 ping &> /dev/null; then
    echo "❌ 无法连接到Redis (localhost:6379)"
    echo "请先启动Redis: redis-server"
    exit 1
fi
echo "✅ Redis连接正常"
echo ""

# 安装依赖
echo "[2/4] 安装Go依赖..."
go mod download
go mod tidy
echo "✅ 依赖安装完成"
echo ""

# 运行单元测试
echo "[3/4] 运行单元测试..."
go test -v ./common/middleware/ratelimit/ -run TestLimiter_Allow
if [ $? -eq 0 ]; then
    echo "✅ 单元测试通过"
else
    echo "❌ 单元测试失败"
    exit 1
fi
echo ""

# 启动示例服务器（后台运行）
echo "[4/4] 启动示例服务器..."
go run example/main.go &
SERVER_PID=$!
echo "服务器PID: $SERVER_PID"

# 等待服务器启动
sleep 2

# 测试API
echo ""
echo "================================"
echo "开始压力测试限流功能"
echo "================================"
echo ""

# 测试1：正常请求（应该成功）
echo "测试1: 正常请求（应该成功）"
for i in {1..3}; do
    echo -n "  请求 $i: "
    RESPONSE=$(curl -s -w "\n%{http_code}" http://localhost:8080/api/users 2>/dev/null)
    HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
    if [ "$HTTP_CODE" == "200" ]; then
        echo "✅ 成功 (HTTP $HTTP_CODE)"
    else
        echo "❌ 失败 (HTTP $HTTP_CODE)"
    fi
    sleep 0.1
done
echo ""

# 测试2：超频请求（应该被限流）
echo "测试2: 快速发送20个请求（部分应该被限流 429）"
SUCCESS=0
RATE_LIMITED=0
for i in {1..20}; do
    HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/api/users)
    if [ "$HTTP_CODE" == "200" ]; then
        ((SUCCESS++))
    elif [ "$HTTP_CODE" == "429" ]; then
        ((RATE_LIMITED++))
    fi
done
echo "  ✅ 成功请求: $SUCCESS 个"
echo "  🚫 被限流请求: $RATE_LIMITED 个"
echo ""

# 测试3：登录接口限流（更严格）
echo "测试3: 快速发送10次登录请求（应该大部分被限流）"
SUCCESS=0
RATE_LIMITED=0
for i in {1..10}; do
    HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST http://localhost:8080/api/login)
    if [ "$HTTP_CODE" == "200" ]; then
        ((SUCCESS++))
    elif [ "$HTTP_CODE" == "429" ]; then
        ((RATE_LIMITED++))
    fi
done
echo "  ✅ 成功请求: $SUCCESS 个"
echo "  🚫 被限流请求: $RATE_LIMITED 个"
echo ""

# 停止服务器
echo "================================"
echo "测试完成，正在停止服务器..."
kill $SERVER_PID 2>/dev/null || true
echo "✅ 服务器已停止"
echo ""

echo "================================"
echo "🎉 限流器测试全部完成！"
echo "================================"
echo ""
echo "你可以："
echo "  1. 查看 README.md 了解详细使用文档"
echo "  2. 运行 'make test' 执行完整的单元测试"
echo "  3. 运行 'go run example/main.go' 启动示例服务器"
echo "  4. 运行 'make test-coverage' 查看测试覆盖率"
