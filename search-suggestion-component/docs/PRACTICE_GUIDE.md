# 🎯 搜索建议组件 - 练手指导文档

欢迎！这个项目为你准备了一个**完整的搜索建议组件**，其中80%的代码已经实现，**剩余20%的关键功能留给你来完成**。

---

## 📊 项目概览

### 已完成的部分 ✅

1. **项目基础架构**
   - ✅ TypeScript 类型定义
   - ✅ 项目配置（tsconfig, vite, package.json）
   - ✅ 1000条测试数据（10个分类）
   - ✅ 拼音工具类（中文搜索支持）

2. **核心功能（部分实现）**
   - ✅ 前缀匹配算法
   - ✅ 包含匹配算法
   - ✅ 拼音匹配算法
   - ✅ 基础排序逻辑
   - ✅ React组件基础结构

### 你需要完成的部分 🎯

| 模块 | 任务 | 难度 | 预计时间 |
|------|------|------|---------|
| **算法模块** | 实现模糊匹配算法 | ⭐⭐⭐ | 1-2小时 |
| **搜索引擎** | 实现分类筛选功能 | ⭐⭐ | 30分钟 |
| **搜索引擎** | 实现LRU缓存机制 | ⭐⭐⭐ | 1小时 |
| **React组件** | 实现键盘导航（上下键） | ⭐⭐ | 30分钟 |
| **React组件** | 实现关键词高亮 | ⭐⭐ | 30分钟 |

**总计预计时间：3-4小时**

---

## 🗂️ 项目目录结构

```
search-suggestion-component/
├── src/
│   ├── types/
│   │   └── index.ts                 # ✅ 类型定义（已完成）
│   ├── data/
│   │   └── suggestions.ts           # ✅ 1000条测试数据（已完成）
│   ├── utils/
│   │   ├── pinyin.ts               # ✅ 拼音工具（已完成）
│   │   └── cache.ts                # 🎯 LRU缓存（你来完成）
│   ├── core/
│   │   ├── algorithms/
│   │   │   ├── matcher.ts          # 🎯 模糊匹配（你来完成）
│   │   │   └── scorer.ts           # ✅ 评分算法（已完成）
│   │   └── SearchEngine.ts         # 🎯 分类筛选（你来完成）
│   ├── components/
│   │   └── SearchSuggestion/
│   │       ├── index.tsx           # ✅ 主组件（已完成）
│   │       ├── SearchInput.tsx     # 🎯 键盘导航（你来完成）
│   │       ├── SuggestionList.tsx  # ✅ 列表组件（已完成）
│   │       └── SuggestionItem.tsx  # 🎯 高亮显示（你来完成）
│   └── hooks/
│       └── useDebounce.ts          # ✅ 防抖Hook（已完成）
├── examples/
│   └── App.tsx                     # ✅ 示例页面（已完成）
└── docs/
    ├── PRACTICE_GUIDE.md           # 本文档
    └── API.md                      # API文档
```

---

## 🎓 任务详解

### 任务1：实现模糊匹配算法 ⭐⭐⭐

**文件位置**：`src/core/algorithms/matcher.ts`

**任务描述**：
实现一个模糊匹配算法，能够匹配用户输入的关键词与建议项文本，即使有拼写错误或缺少某些字符。

**示例**：
```typescript
// 用户输入：iph
// 应该能匹配：iPhone 15 Pro Max
// 
// 用户输入：macbok
// 应该能匹配：MacBook Pro（允许1个字符错误）
```

**TODO标记位置**：
```typescript
// TODO: 【练手任务1】实现模糊匹配算法
// 提示：
// 1. 可以使用编辑距离（Levenshtein Distance）算法
// 2. 或者实现简单的字符跳跃匹配
// 3. 返回匹配得分 0-100
export function fuzzyMatch(text: string, keyword: string): number {
  // 你的代码在这里
  return 0
}
```

**实现提示**：
```typescript
// 方案1：简单字符跳跃匹配
function simpleF uzzyMatch(text: string, keyword: string): number {
  let textIndex = 0
  let keywordIndex = 0
  let matchCount = 0
  
  while (textIndex < text.length && keywordIndex < keyword.length) {
    if (text[textIndex].toLowerCase() === keyword[keywordIndex].toLowerCase()) {
      matchCount++
      keywordIndex++
    }
    textIndex++
  }
  
  return (matchCount / keyword.length) * 100
}

// 方案2：编辑距离算法（更精确，但更复杂）
// 参考：https://zh.wikipedia.org/wiki/編輯距離
```

---

### 任务2：实现分类筛选功能 ⭐⭐

**文件位置**：`src/core/SearchEngine.ts`

**任务描述**：
在搜索时，如果用户选择了某个分类，只返回该分类下的结果。

**TODO标记位置**：
```typescript
search(keyword: string, options?: SearchOptions): MatchResult[] {
  // ... 前面的代码已实现
  
  // TODO: 【练手任务2】实现分类筛选
  // 提示：
  // 1. 检查 options.category 是否存在
  // 2. 如果存在，过滤results数组
  // 3. 只保留匹配该分类的结果
  
  // 你的代码在这里
  
  return results.slice(0, this.config.topN)
}
```

**实现提示**：
```typescript
// 如果指定了分类筛选
if (options?.category) {
  results = results.filter(result => 
    result.item.category === options.category
  )
}
```

---

### 任务3：实现LRU缓存机制 ⭐⭐⭐

**文件位置**：`src/utils/cache.ts`

**任务描述**：
实现一个LRU（Least Recently Used）缓存，存储最近的搜索结果，提高性能。

**TODO标记位置**：
```typescript
// TODO: 【练手任务3】实现LRU缓存
// 提示：
// 1. 使用 Map 保存缓存数据（key: 搜索词, value: 结果）
// 2. 记录访问顺序
// 3. 超过maxSize时，删除最久未使用的项
export class LRUCache<K, V> {
  private maxSize: number
  // 你的代码在这里
  
  constructor(maxSize: number = 100) {
    this.maxSize = maxSize
    // 初始化
  }
  
  get(key: K): V | undefined {
    // 实现获取逻辑
  }
  
  set(key: K, value: V): void {
    // 实现设置逻辑
  }
}
```

**实现提示**：
```typescript
// Map在JavaScript中保持插入顺序
// 可以利用这个特性实现LRU

private cache = new Map<K, V>()大min

get(key: K): V | undefined {
  if (!this.cache.has(key)) return undefined
  
  // 重新插入以更新顺序
  const value = this.cache.get(key)!
  this.cache.delete(key)
  this.cache.set(key, value)
  
  return value
}

set(key: K, value: V): void {
  if (this.cache.has(key)) {
    this.cache.delete(key)
  }
  this.cache.set(key, value)
  
  // 如果超过大小限制，删除最旧的项
  if (this.cache.size > this.maxSize) {
    const firstKey = this.cache.keys().next().value
    this.cache.delete(firstKey)
  }
}
```

---

### 任务4：实现键盘导航 ⭐⭐

**文件位置**：`src/components/SearchSuggestion/SearchInput.tsx`

**任务描述**：
支持键盘上下键选择建议项，Enter键确认。

**TODO标记位置**：
```typescript
// TODO: 【练手任务4】实现键盘导航
// 提示：
// 1. 监听 onKeyDown 事件
// 2. 处理 ArrowUp、ArrowDown、Enter、Escape
// 3. 维护 selectedIndex 状态
const handleKeyDown = (e: React.KeyboardEvent) => {
  // 你的代码在这里
}
```

**实现提示**：
```typescript
const [selectedIndex, setSelectedIndex] = useState(-1)

const handleKeyDown = (e: React.KeyboardEvent) => {
  switch (e.key) {
    case 'ArrowDown':
      e.preventDefault()
      setSelectedIndex(prev => 
        Math.min(prev + 1, results.length - 1)
      )
      break
      
    case 'ArrowUp':
      e.preventDefault()
      setSelectedIndex(prev => Math.max(prev - 1, -1))
      break
      
    case 'Enter':
      if (selectedIndex >= 0) {
        handleSelect(results[selectedIndex])
      }
      break
      
    case 'Escape':
      onClose()
      break
  }
}
```

---

### 任务5：实现关键词高亮 ⭐⭐

**文件位置**：`src/components/SearchSuggestion/SuggestionItem.tsx`

**任务描述**：
将匹配的关键词在建议项中高亮显示。

**TODO标记位置**：
```typescript
// TODO: 【练手任务5】实现关键词高亮
// 提示：
// 1. 查找keyword在text中的位置
// 2. 将text分割为 [前缀, 匹配部分, 后缀]
// 3. 匹配部分用<mark>标签包裹
function highlightText(text: string, keyword: string): React.ReactNode {
  // 你的代码在这里
  return text
}
```

**实现提示**：
```typescript
function highlightText(text: string, keyword: string): React.ReactNode {
  if (!keyword) return text
  
  const index = text.toLowerCase().indexOf(keyword.toLowerCase())
  if (index === -1) return text
  
  const before = text.slice(0, index)
  const match = text.slice(index, index + keyword.length)
  const after = text.slice(index + keyword.length)
  
  return (
    <>
      {before}
      <mark className="highlight">{match}</mark>
      {after}
    </>
  )
}
```

---

## 🚀 开始练手

### 第一步：安装依赖

```bash
cd search-suggestion-component
npm install
```

### 第二步：启动开发服务器

```bash
npm run dev
```

浏览器会自动打开 `http://localhost:3000`

### 第三步：按顺序完成任务

建议按照以下顺序完成：

1. ✅ 先完成简单的：**任务2（分类筛选）** 和 **任务5（关键词高亮）**
2. ✅ 再完成中等的：**任务4（键盘导航）**
3. ✅ 最后挑战难的：**任务1（模糊匹配）** 和 **任务3（LRU缓存）**

### 第四步：测试你的实现

```bash
# 运行测试
npm run test

# 在浏览器中测试
# 打开 http://localhost:3000
# 尝试各种搜索：
# - 输入 "iphone" 测试匹配
# - 输入 "sxj" 测试拼音首字母
# - 选择分类后搜索
# - 使用键盘上下键导航
```

---

## 💡 调试技巧

### 1. 使用Console查看数据

```typescript
console.log('搜索关键词：', keyword)
console.log('匹配结果：', results)
console.log('选中索引：', selectedIndex)
```

### 2. React DevTools

安装 React DevTools 浏览器插件，可以查看组件状态和props。

### 3. 断点调试

在Chrome DevTools中设置断点，逐步执行代码。

---

## 🎓 扩展挑战（可选）

如果你完成了所有基础任务，可以尝试这些进阶功能：

1. **历史搜索记录**
   - 保存最近10条搜索
   - 显示在建议列表顶部

2. **热门搜索推荐**
   - 没有输入时，显示热门搜索词
   - 点击后直接搜索

3. **搜索结果分页**
   - 每页显示10条
   - 滚动加载更多

4. **性能优化**
   - 使用Web Worker进行搜索
   - 虚拟滚动优化长列表

5. **搜索分析**
   - 统计搜索关键词频率
   - 记录用户点击行为

---

## 📚 参考资料

### 算法相关
- [编辑距离算法](https://zh.wikipedia.org/wiki/編輯距離)
- [LRU缓存实现](https://leetcode.cn/problems/lru-cache/)
- [字符串匹配算法](https://en.wikipedia.org/wiki/String-searching_algorithm)

### React相关
- [React Hooks文档](https://react.dev/reference/react)
- [事件处理](https://react.dev/learn/responding-to-events)
- [键盘事件](https://developer.mozilla.org/zh-CN/docs/Web/API/KeyboardEvent)

### TypeScript相关
- [TypeScript官方文档](https://www.typescriptlang.org/zh/docs/)
- [泛型](https://www.typescriptlang.org/docs/handbook/2/generics.html)

---

## ❓ 常见问题

### Q1: 找不到TODO标记？
**A**: 在VSCode中，按 `Ctrl+Shift+F` 全局搜索 `TODO:`

### Q2: 拼音库报错？
**A**: 确保安装了依赖：`npm install pinyin-pro`

### Q3: TypeScript类型错误？
**A**: 检查导入路径，确保使用了正确的类型定义

### Q4: 热更新不生效？
**A**: 重启开发服务器：`npm run dev`

---

## 🎉 完成标准

当你完成所有任务后，应该能够：

- ✅ 输入关键词，看到TOP 10匹配结果
- ✅ 支持中文、拼音、拼音首字母搜索
- ✅ 可以按分类筛选
- ✅ 使用键盘上下键导航
- ✅ 匹配的关键词高亮显示
- ✅ 搜索结果有缓存，性能良好

---

## 🤝 获取帮助

如果遇到问题：

1. 查看本文档的实现提示
2. 搜索相关技术文档
3. 在代码中添加console.log调试
4. 寻求帮助（注明具体问题和已尝试的方法）

---

**祝你练手愉快！加油！** 💪

记住：**每个大牛都是从写TODO开始的**。

