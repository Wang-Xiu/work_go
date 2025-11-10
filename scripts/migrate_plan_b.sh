#!/bin/bash
# 方案B迁移脚本：抽取Redis层 - 建立基础设施层
# 用法: ./migrate_plan_b.sh

set -e

PROJECT_ROOT="/Users/xiu/work/work_go"
BACKUP_DIR="${PROJECT_ROOT}/backup_plan_b_$(date +%Y%m%d_%H%M%S)"

echo "========================================="
echo "方案B：抽取Redis层迁移"
echo "========================================="

# 1. 创建备份
echo ""
echo "[1/7] 创建备份..."
mkdir -p "$BACKUP_DIR"
cp -r "${PROJECT_ROOT}/common" "${BACKUP_DIR}/common_backup"
cp -r "${PROJECT_ROOT}/config" "${BACKUP_DIR}/config_backup"
[ -d "${PROJECT_ROOT}/example" ] && cp -r "${PROJECT_ROOT}/example" "${BACKUP_DIR}/example_backup"
echo "✅ 备份已创建: $BACKUP_DIR"

# 2. 执行方案A（统一目录结构）
echo ""
echo "[2/7] 先执行方案A：统一目录结构..."
if [ -d "${PROJECT_ROOT}/common/middleware/stats" ]; then
    echo "✅ stats已在middleware目录下，跳过"
else
    mv "${PROJECT_ROOT}/common/stats" "${PROJECT_ROOT}/common/middleware/stats"
    echo "✅ 已移动: common/stats -> common/middleware/stats"

    # 更新import路径
    if [[ "$OSTYPE" == "darwin"* ]]; then
        find "${PROJECT_ROOT}" -name "*.go" -type f -exec sed -i '' 's|"working-project/common/stats"|"working-project/common/middleware/stats"|g' {} +
    else
        find "${PROJECT_ROOT}" -name "*.go" -type f -exec sed -i 's|"working-project/common/stats"|"working-project/common/middleware/stats"|g' {} +
    fi
    echo "✅ Import路径已更新"
fi

# 3. 创建infrastructure目录结构
echo ""
echo "[3/7] 创建infrastructure目录结构..."
mkdir -p "${PROJECT_ROOT}/common/infrastructure/redis"

# 检查文件是否已存在
if [ -f "${PROJECT_ROOT}/common/infrastructure/redis/manager.go" ]; then
    echo "✅ infrastructure/redis 目录已存在，跳过创建"
else
    echo "⚠️  infrastructure/redis 代码文件不存在"
    echo "请确保以下文件已创建:"
    echo "  - common/infrastructure/redis/manager.go"
    echo "  - common/infrastructure/redis/config.go"
    echo "  - common/infrastructure/redis/errors.go"
    echo ""
    read -p "文件是否已准备好? (y/n) " -n 1 -r
    echo

    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "❌ 请先创建必要的文件后再运行此脚本"
        exit 1
    fi
fi

# 4. 创建config/app.yml
echo ""
echo "[4/7] 创建统一配置文件..."
if [ ! -f "${PROJECT_ROOT}/config/app.yml" ]; then
    echo "⚠️  config/app.yml 不存在"
    echo "请确保config/app.yml已创建"
    read -p "配置文件是否已准备好? (y/n) " -n 1 -r
    echo

    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "❌ 请先创建config/app.yml"
        exit 1
    fi
else
    echo "✅ config/app.yml 已存在"
fi

# 5. 验证编译
echo ""
echo "[5/7] 验证基础设施层编译..."
cd "$PROJECT_ROOT"

if go build ./common/infrastructure/...; then
    echo "✅ 基础设施层编译成功"
else
    echo "❌ 基础设施层编译失败"
    echo "请检查以下文件:"
    echo "  - common/infrastructure/redis/manager.go"
    echo "  - common/infrastructure/redis/config.go"
    echo "  - common/infrastructure/redis/errors.go"
    exit 1
fi

# 6. 验证整体编译
echo ""
echo "[6/7] 验证整体编译..."
if go build ./...; then
    echo "✅ 整体编译成功"
else
    echo "⚠️  编译失败，可能需要更新中间件代码以使用新的Redis Manager"
    echo ""
    echo "请参考文档更新以下文件:"
    echo "  - common/middleware/ratelimit/config.go"
    echo "  - common/middleware/ratelimit/limiter.go"
    echo "  - common/middleware/stats/config.go"
    echo "  - common/middleware/stats/stats.go"
    echo ""
    echo "详细说明请查看: docs/refactor_guide_plan_b.md"
fi

# 7. 运行测试
echo ""
echo "[7/7] 运行测试..."
if go test ./common/infrastructure/...; then
    echo "✅ 基础设施层测试通过"
else
    echo "⚠️  基础设施层测试失败"
fi

# 完成
echo ""
echo "========================================="
echo "✅ 方案B基础设施已创建！"
echo "========================================="
echo ""
echo "变更摘要:"
echo "  ✅ common/stats → common/middleware/stats"
echo "  ✅ 创建 common/infrastructure/redis/"
echo "  ✅ 创建 config/app.yml"
echo ""
echo "⚠️  下一步行动:"
echo "  1. 参考 docs/refactor_guide_plan_b.md"
echo "  2. 更新 ratelimit 和 stats 的代码"
echo "  3. 更新 example/main.go 使用新的初始化方式"
echo "  4. 运行完整测试: go test ./..."
echo ""
echo "备份位置: $BACKUP_DIR"
echo ""
