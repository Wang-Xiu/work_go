# 重构迁移脚本使用指南

## 目录结构

```
scripts/
├── migrate_plan_a.sh    # 方案A：统一中间件目录
├── migrate_plan_b.sh    # 方案B：抽取Redis层
├── rollback.sh          # 通用回滚脚本
└── README.md           # 本文件
```

## 准备工作

### 1. 赋予执行权限

```bash
cd /Users/xiu/work/work_go/scripts
chmod +x migrate_plan_a.sh
chmod +x migrate_plan_b.sh
chmod +x rollback.sh
```

### 2. 确保Git状态干净

```bash
cd /Users/xiu/work/work_go
git status

# 如果有未提交的更改，先提交或暂存
git add .
git commit -m "保存当前状态"
```

### 3. 创建Git分支（推荐）

```bash
# 为每个方案创建独立分支
git checkout -b refactor/plan-a
# 或
git checkout -b refactor/plan-b
```

---

## 方案选择

### 方案A：最小改动（推荐快速开始）

**适用场景：**
- 快速统一目录结构
- 不想大改现有代码
- 作为方案B的前置步骤

**执行：**
```bash
./scripts/migrate_plan_a.sh
```

**预期变更：**
- `common/stats` → `common/middleware/stats`
- 所有import路径自动更新
- 编译验证自动运行

**耗时：** 约2-5分钟

---

### 方案B：抽取Redis层（推荐生产环境）

**适用场景：**
- 消除Redis连接重复
- 建立基础设施层架构
- 为未来扩展打基础

**前置条件：**
1. 已执行方案A（或stats在middleware下）
2. 已创建以下文件：
   - `common/infrastructure/redis/manager.go`
   - `common/infrastructure/redis/config.go`
   - `common/infrastructure/redis/errors.go`
   - `config/app.yml`

**执行：**
```bash
# 1. 确保infrastructure代码已创建（已由Claude生成）
ls -la common/infrastructure/redis/

# 2. 确保config/app.yml已创建
cat config/app.yml

# 3. 运行迁移脚本
./scripts/migrate_plan_b.sh
```

**后续步骤：**
1. 参考 `docs/refactor_guide_plan_b.md`
2. 更新 `ratelimit` 代码
3. 更新 `stats` ��码
4. 更新 `example/main.go`
5. 运行测试

**耗时：** 约2-4小时（包含代码修改）

---

### 方案C：完整DDD架构（推荐企业级）

**适用场景：**
- 大型项目，需要规范化架构
- 团队协作，需要统一标准
- 长期演进，��要可扩展性

**说明：**
方案C没有自动脚本，因为涉及大量架构设计决策。

**执行步骤：**
1. 阅读 `docs/refactor_guide_plan_c.md`
2. 按阶段（Phase 1-5）逐步实施
3. 每个阶段完成后运行测试
4. 建议分5个Sprint完成

**耗时：** 约4-5周

---

## 回滚操作

如果迁移出现问题，可以快速回滚：

```bash
# 1. 查看可用备份
ls -lt /Users/xiu/work/work_go/backup_*

# 2. 选择要恢复的备份
./scripts/rollback.sh /Users/xiu/work/work_go/backup_20250107_143022

# 3. 验证回滚结果
go build ./...
go test ./...
```

---

## 常见问题

### Q1: 执行方案A后编译失败

**原因：** Import路径替换不完整

**解决：**
```bash
# 手动查找未替换的路径
grep -r "common/stats" . --include="*.go"

# 手动替换
# 将 "working-project/common/stats" 改为 "working-project/common/middleware/stats"
```

### Q2: 方案B脚本提示infrastructure文件不存在

**原因：** 忘记创建基础设施层代码

**解决：**
```bash
# 确认文件已创建
ls -la common/infrastructure/redis/

# 如果没有，参考docs/refactor_guide_plan_b.md创建
```

### Q3: 回滚后仍然编译失败

**原因：** 可能有其他未提交的更改

**解决：**
```bash
# 使用Git恢复到最后一次提交
git reset --hard HEAD

# 或者重新从远程仓库拉取
git fetch origin
git reset --hard origin/main
```

### Q4: 想跳过方案A直接执行方案B

**不推荐，但如果必须：**
```bash
# 手动移动目录
mv common/stats common/middleware/stats

# 手动更新import（macOS）
find . -name "*.go" -type f -exec sed -i '' 's|common/stats|common/middleware/stats|g' {} +

# 然后执行方案B
./scripts/migrate_plan_b.sh
```

---

## 脚本安全特性

### ✅ 自动备份
每次运行都会创建时间戳备份，可随时回滚

### ✅ 编译验证
迁移后自动运行 `go build` 验证

### ✅ 用户确认
关键操作前会要求用户确认

### ✅ 错误停止
遇到错误立即停止，防止破坏代码

---

## 推荐流程

### 开发环境（本地测试）

```bash
# Day 1: 快速验证方案A
git checkout -b refactor/plan-a
./scripts/migrate_plan_a.sh
go test ./...
git commit -am "完成方案A：统一中间件目录"

# Day 2-3: 实施方案B
git checkout -b refactor/plan-b
./scripts/migrate_plan_b.sh
# 按照文档修改代码
go test ./...
git commit -am "完成方案B：建立基础设施层"
```

### 生产环境

```bash
# Week 1: 方案B实施和测试
git checkout -b refactor/production
./scripts/migrate_plan_b.sh
# 完整测试
git push origin refactor/production

# Week 2: Code Review和合并
# 团队Review后合并到main

# Week 3: 灰度发布
# 先发布到staging环境，监控1周

# Week 4: 全量发布
```

---

## 技术支持

遇到问题时的检查清单：

1. ✅ 是否在项目根目录执行？
2. ✅ 是否有文件权限？
3. ✅ Go版本是否>=1.19？
4. ✅ 是否有未提交的Git更改？
5. ✅ 备份目录是否完整？

如果以上都检查过仍有问题，建议：
- 查看详细文档：`docs/refactor_guide_plan_*.md`
- 使用回滚脚本恢复
- 手动逐步执行迁移步骤

---

## 完成后的验证

### 方案A验证

```bash
# 1. 检查目录结构
ls -la common/middleware/stats/

# 2. 检查import路径
grep -r "common/middleware/stats" . --include="*.go" | head -5

# 3. 编译测试
go build ./...
go test ./common/middleware/stats/
```

### 方案B验证

```bash
# 1. 检查基础设施层
ls -la common/infrastructure/redis/

# 2. 检查配置文件
cat config/app.yml

# 3. 运行集成测试
go test ./common/infrastructure/...
go test ./common/middleware/...

# 4. 启动示例程序
go run example/main.go
```

---

祝重构顺利！🚀
