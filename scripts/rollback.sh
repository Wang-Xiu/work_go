#!/bin/bash
# 回滚脚本：恢复到备份状态
# 用法: ./rollback.sh <backup_dir>

set -e

if [ -z "$1" ]; then
    echo "❌ 错误：请提供备份目录路径"
    echo "用法: ./rollback.sh <backup_dir>"
    echo ""
    echo "可用的备份目录:"
    ls -dt /Users/xiu/work/work_go/backup_* 2>/dev/null || echo "  (未找到备份)"
    exit 1
fi

BACKUP_DIR="$1"
PROJECT_ROOT="/Users/xiu/work/work_go"

if [ ! -d "$BACKUP_DIR" ]; then
    echo "❌ 错误：备份目录不存在: $BACKUP_DIR"
    exit 1
fi

echo "========================================="
echo "回滚到备份状态"
echo "========================================="
echo ""
echo "备份目录: $BACKUP_DIR"
echo "项目目录: $PROJECT_ROOT"
echo ""

# 确认
read -p "⚠️  此操作将覆盖当前代码，是否继续? (yes/no) " -r
echo

if [[ ! $REPLY =~ ^[Yy][Ee][Ss]$ ]]; then
    echo "❌ 回滚已取消"
    exit 1
fi

# 回滚
echo ""
echo "[1/3] 恢复common目录..."
if [ -d "${BACKUP_DIR}/common_backup" ]; then
    rm -rf "${PROJECT_ROOT}/common"
    cp -r "${BACKUP_DIR}/common_backup" "${PROJECT_ROOT}/common"
    echo "✅ common目录已恢复"
elif [ -d "${BACKUP_DIR}/stats_backup" ]; then
    # 方案A的备份
    rm -rf "${PROJECT_ROOT}/common/middleware/stats"
    cp -r "${BACKUP_DIR}/stats_backup" "${PROJECT_ROOT}/common/stats"
    echo "✅ stats目录已恢复"
fi

echo ""
echo "[2/3] 恢复config目录..."
if [ -d "${BACKUP_DIR}/config_backup" ]; then
    rm -rf "${PROJECT_ROOT}/config"
    cp -r "${BACKUP_DIR}/config_backup" "${PROJECT_ROOT}/config"
    echo "✅ config目录已恢复"
else
    echo "⏭️  未找到config备份，跳过"
fi

echo ""
echo "[3/3] 恢复example目录..."
if [ -d "${BACKUP_DIR}/example_backup" ]; then
    rm -rf "${PROJECT_ROOT}/example"
    cp -r "${BACKUP_DIR}/example_backup" "${PROJECT_ROOT}/example"
    echo "✅ example目录已恢复"
else
    echo "⏭️  未找到example备份，跳过"
fi

# 验证编译
echo ""
echo "验证编译..."
cd "$PROJECT_ROOT"

if go build ./...; then
    echo "✅ 编译成功"
else
    echo "❌ 编译失败"
    echo "请手动检查问题"
    exit 1
fi

# 完成
echo ""
echo "========================================="
echo "✅ 回滚完成！"
echo "========================================="
echo ""
echo "已恢复到备份状态: $(basename $BACKUP_DIR)"
echo ""
