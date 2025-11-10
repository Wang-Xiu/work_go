#!/bin/bash
# 方案A迁移脚本：统一中间件目录结构
# 用法: ./migrate_plan_a.sh

set -e  # 遇到错误立即退出

PROJECT_ROOT="/Users/xiu/work/work_go"
BACKUP_DIR="${PROJECT_ROOT}/backup_$(date +%Y%m%d_%H%M%S)"

echo "========================================="
echo "方案A：统一中间件目录结构迁移"
echo "========================================="

# 1. 创建备份
echo ""
echo "[1/5] 创建备份..."
mkdir -p "$BACKUP_DIR"
cp -r "${PROJECT_ROOT}/common/stats" "${BACKUP_DIR}/stats_backup"
echo "✅ 备份已创建: $BACKUP_DIR"

# 2. 移动目录
echo ""
echo "[2/5] 移动stats目录到middleware下..."
if [ -d "${PROJECT_ROOT}/common/middleware/stats" ]; then
    echo "⚠️  目标目录已存在，跳过移动"
else
    mv "${PROJECT_ROOT}/common/stats" "${PROJECT_ROOT}/common/middleware/stats"
    echo "�� 目录已移动: common/stats -> common/middleware/stats"
fi

# 3. 更新import路径
echo ""
echo "[3/5] 更新import路径..."
echo "查找所有引用common/stats的Go文件..."

# 查找所有go文件中的import
AFFECTED_FILES=$(grep -r "common/stats" "${PROJECT_ROOT}" --include="*.go" -l || true)

if [ -z "$AFFECTED_FILES" ]; then
    echo "✅ 未发现需要更新的文件"
else
    echo "发现以下文件需要更新:"
    echo "$AFFECTED_FILES"

    echo ""
    read -p "是否继续更新这些文件? (y/n) " -n 1 -r
    echo

    if [[ $REPLY =~ ^[Yy]$ ]]; then
        # 批量替换（macOS和Linux兼容）
        if [[ "$OSTYPE" == "darwin"* ]]; then
            # macOS
            find "${PROJECT_ROOT}" -name "*.go" -type f -exec sed -i '' 's|"working-project/common/stats"|"working-project/common/middleware/stats"|g' {} +
        else
            # Linux
            find "${PROJECT_ROOT}" -name "*.go" -type f -exec sed -i 's|"working-project/common/stats"|"working-project/common/middleware/stats"|g' {} +
        fi
        echo "✅ Import路径已更新"
    else
        echo "❌ 用户取消更新"
        exit 1
    fi
fi

# 4. 验证编译
echo ""
echo "[4/5] 验证编译..."
cd "$PROJECT_ROOT"

if go build ./...; then
    echo "✅ 编译成功"
else
    echo "❌ 编译失败，开始回滚..."

    # 回滚
    rm -rf "${PROJECT_ROOT}/common/middleware/stats"
    mv "${BACKUP_DIR}/stats_backup" "${PROJECT_ROOT}/common/stats"

    # 恢复import路径
    if [[ "$OSTYPE" == "darwin"* ]]; then
        find "${PROJECT_ROOT}" -name "*.go" -type f -exec sed -i '' 's|"working-project/common/middleware/stats"|"working-project/common/stats"|g' {} +
    else
        find "${PROJECT_ROOT}" -name "*.go" -type f -exec sed -i 's|"working-project/common/middleware/stats"|"working-project/common/stats"|g' {} +
    fi

    echo "✅ 已回滚到原始状态"
    exit 1
fi

# 5. 运行测试
echo ""
echo "[5/5] 运行测试..."
if go test ./...; then
    echo "✅ 测试通过"
else
    echo "⚠️  测试失败，但目录结构已更新"
    echo "请手动检查测试失败原因"
fi

# 完成
echo ""
echo "========================================="
echo "✅ 方案A迁移完成！"
echo "========================================="
echo ""
echo "变更摘要:"
echo "  - common/stats → common/middleware/stats"
echo "  - 所有import路径已更新"
echo "  - 编译验证: 通过"
echo ""
echo "备份位置: $BACKUP_DIR"
echo "如需回滚，请运行: ./rollback_plan_a.sh $BACKUP_DIR"
echo ""
