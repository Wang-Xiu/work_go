# 🚀 SearchEngine 高级功能使用指南

## 概览

本文档介绍SearchEngine的三个高级扩展功能：
1. **多分类筛选**：同时筛选多个分类
2. **层级分类筛选**：支持父子分类关系
3. **分类统计**：统计每个分类的结果数量

---

## 📚 功能1：多分类筛选

### 应用场景

用户想同时查看多个分类的搜索结果。

**示例**：
- 搜索"充电"，同时显示"电子产品"和"数码配件"的结果
- 搜索"旅游"，同时显示"国内游"和"国外游"的结果

### 使用方法

```typescript
import { SearchEngine } from './core/SearchEngine'
import { suggestions } from './data/suggestions'

const searchEngine = new SearchEngine(suggestions)

// 方式1：筛选单个分类（基础功能）
const results1 = searchEngine.search('手机', {
  category: '电子产品'
})

// 方式2：筛选多个分类（扩展功能）
const results2 = searchEngine.search('充电', {
  categories: ['电子产品', '数码配件', '智能家居']
})

console.log(results2)
// 只显示这三个分类的搜索结果
```

### 类型定义

```typescript
interface SearchOptions {
  // 单个分类筛选
  category?: string
  
  // 多个分类筛选（优先级高于 category）
  categories?: string[]
}
```

### 优先级规则

- 如果同时指定了 `category` 和 `categories`，**categories 优先生效**
- `categories` 为空数组时，等同于不筛选

### React 组件示例

```typescript
import React, { useState } from 'react'
import { SearchSuggestion } from './components/SearchSuggestion'

function App() {
  const [selectedCategories, setSelectedCategories] = useState<string[]>([])
  
  const handleCategoryToggle = (category: string) => {
    setSelectedCategories(prev => 
      prev.includes(category)
        ? prev.filter(c => c !== category)  // 取消选中
        : [...prev, category]                // 添加选中
    )
  }
  
  return (
    <div>
      {/* 分类多选 */}
      <div className="category-filters">
        {categories.map(cat => (
          <label key={cat}>
            <input
              type="checkbox"
              checked={selectedCategories.includes(cat)}
              onChange={() => handleCategoryToggle(cat)}
            />
            {cat}
          </label>
        ))}
      </div>
      
      {/* 搜索组件 */}
      <SearchSuggestion
        items={suggestions}
        searchOptions={{
          categories: selectedCategories
        }}
      />
    </div>
  )
}
```

---

## 🌳 功能2：层级分类筛选

### 应用场景

分类是树形结构，选择父分类时，应该包含所有子分类的结果。

**分类树示例**：
```
电子产品
  ├─ 手机
  │   ├─ 苹果手机
  │   └─ 安卓手机
  └─ 电脑
      ├─ 笔记本
      └─ 台式机

美食餐饮
  ├─ 中餐
  │   ├─ 川菜
  │   └─ 粤菜
  └─ 西餐
      ├─ 意大利菜
      └─ 法国菜
```

**需求**：
- 选择"手机"，应该显示"苹果手机"和"安卓手机"的结果
- 选择"中餐"，应该显示"川菜"和"粤菜"的结果

### 数据准备

首先需要为数据添加 `parentCategory` 字段：

```typescript
const suggestions: SuggestionItem[] = [
  {
    id: '1',
    text: 'iPhone 15 Pro',
    category: '苹果手机',
    parentCategory: '手机',  // 👈 指定父分类
    hotScore: 95
  },
  {
    id: '2',
    text: 'Samsung Galaxy',
    category: '安卓手机',
    parentCategory: '手机',  // 👈 指定父分类
    hotScore: 88
  },
  {
    id: '3',
    text: '小米13',
    category: '安卓手机',
    parentCategory: '手机',  // 👈 指定父分类
    hotScore: 85
  }
]
```

### 使用方法

```typescript
// 搜索"手机"，同时包含子分类
const results = searchEngine.search('Pro', {
  category: '手机',
  includeSubCategories: true  // 👈 启用子分类包含
})

console.log(results)
// 结果包含：
// - category='手机' 的项
// - parentCategory='手机' 的项（苹果手机、安卓手机）
```

### 多分类 + 层级筛选

```typescript
// 同时筛选多个父分类，并包含它们的子分类
const results = searchEngine.search('智能', {
  categories: ['手机', '智能家居'],
  includeSubCategories: true
})

// 结果包含：
// - 手机及其子分类（苹果手机、安卓手机）
// - 智能家居及其子分类
```

### 类型定义

```typescript
interface SuggestionItem {
  category: string
  parentCategory?: string  // 父分类（可选）
}

interface SearchOptions {
  includeSubCategories?: boolean  // 是否包含子分类
}
```

### 实现原理

```typescript
// 伪代码
if (includeSubCategories) {
  results = results.filter(result => {
    // 检查当前分类
    if (targetCategories.includes(result.item.category)) {
      return true
    }
    
    // 检查父分类
    if (result.item.parentCategory && 
        targetCategories.includes(result.item.parentCategory)) {
      return true
    }
    
    return false
  })
}
```

### UI 设计建议

```typescript
// 分类树形选择器
<TreeSelect
  data={categoryTree}
  onChange={(selected) => {
    setSearchOptions({
      categories: selected,
      includeSubCategories: true
    })
  }}
/>

// 或者面包屑导航
<Breadcrumb>
  <span onClick={() => setCategory('all')}>全部</span>
  <span onClick={() => setCategory('电子产品')}>电子产品</span>
  <span onClick={() => setCategory('手机')}>手机</span>
</Breadcrumb>
```

---

## 📊 功能3：分类统计

### 应用场景

在搜索结果旁边显示每个分类有多少匹配项，帮助用户快速了解数据分布。

**效果示例**：
```
搜索结果 (共 46 项)

分类统计：
- 电子产品 (23)
- 数码配件 (15)
- 智能家居 (8)
```

### 使用方法

#### 方式1：使用 `searchWithStats()` 方法

```typescript
import { SearchEngine } from './core/SearchEngine'

const searchEngine = new SearchEngine(suggestions)

// 返回结果 + 分类统计
const result = searchEngine.searchWithStats('手机')

console.log(result.results)         // MatchResult[] - 搜索结果
console.log(result.categoryStats)   // CategoryStats[] - 分类统计

// 输出示例：
// categoryStats = [
//   { name: "电子产品", count: 15 },
//   { name: "数码配件", count: 8 },
//   { name: "智能家居", count: 3 }
// ]
```

#### 方式2：单独计算统计（内部方法）

```typescript
// calculateCategoryStats 是私有方法
// 如果需要单独使用，可以提取为公共工具函数
function getCategoryStats(results: MatchResult[]): CategoryStats[] {
  const statsMap = new Map<string, number>()
  
  for (const result of results) {
    const category = result.item.category
    statsMap.set(category, (statsMap.get(category) || 0) + 1)
  }
  
  const stats: CategoryStats[] = []
  statsMap.forEach((count, name) => {
    stats.push({ name, count })
  })
  
  return stats.sort((a, b) => b.count - a.count)
}
```

### 类型定义

```typescript
interface CategoryStats {
  name: string   // 分类名称
  count: number  // 结果数量
}

interface SearchResult {
  results: MatchResult[]        // 搜索结果（TOP N）
  categoryStats: CategoryStats[]  // 分类统计（全部匹配项）
}
```

### React 组件示例

```typescript
import React, { useState } from 'react'

function SearchWithStats() {
  const [result, setResult] = useState<SearchResult | null>(null)
  
  const handleSearch = (keyword: string) => {
    const searchResult = searchEngine.searchWithStats(keyword)
    setResult(searchResult)
  }
  
  return (
    <div className="search-page">
      {/* 搜索框 */}
      <SearchInput onSearch={handleSearch} />
      
      <div className="search-content">
        {/* 左侧：分类统计 */}
        <aside className="category-stats">
          <h3>分类分布</h3>
          {result?.categoryStats.map(stat => (
            <div key={stat.name} className="stat-item">
              <span className="category-name">{stat.name}</span>
              <span className="category-count">({stat.count})</span>
            </div>
          ))}
        </aside>
        
        {/* 右侧：搜索结果 */}
        <main className="search-results">
          <h3>搜索结果 (共 {result?.results.length} 项)</h3>
          {result?.results.map(item => (
            <SearchResultItem key={item.item.id} result={item} />
          ))}
        </main>
      </div>
    </div>
  )
}
```

### CSS 样式示例

```css
.search-content {
  display: flex;
  gap: 24px;
}

.category-stats {
  width: 200px;
  padding: 16px;
  background: #f5f5f5;
  border-radius: 8px;
}

.stat-item {
  display: flex;
  justify-content: space-between;
  padding: 8px 0;
  border-bottom: 1px solid #e0e0e0;
  cursor: pointer;
  transition: background 0.2s;
}

.stat-item:hover {
  background: #e8e8e8;
}

.category-count {
  color: #666;
  font-size: 13px;
}
```

### 高级用法：可点击的分类统计

```typescript
function InteractiveCategoryStats() {
  const [selectedCategory, setSelectedCategory] = useState<string | null>(null)
  
  const handleCategoryClick = (categoryName: string) => {
    setSelectedCategory(categoryName)
    // 重新搜索，只显示该分类
    const filtered = searchEngine.search(keyword, {
      category: categoryName
    })
    setResults(filtered)
  }
  
  return (
    <div className="category-stats">
      {categoryStats.map(stat => (
        <div
          key={stat.name}
          className={`stat-item ${selectedCategory === stat.name ? 'active' : ''}`}
          onClick={() => handleCategoryClick(stat.name)}
        >
          <span>{stat.name}</span>
          <span className="count">({stat.count})</span>
        </div>
      ))}
    </div>
  )
}
```

---

## 🔗 组合使用

### 示例1：多分类 + 统计

```typescript
const result = searchEngine.searchWithStats('智能', {
  categories: ['电子产品', '数码配件', '智能家居']
})

console.log(result.categoryStats)
// 只统计这三个分类的数量分布
```

### 示例2：层级分类 + 统计

```typescript
const result = searchEngine.searchWithStats('手机', {
  category: '手机',
  includeSubCategories: true
})

console.log(result.categoryStats)
// 统计：
// - 手机 (5)
// - 苹果手机 (12)
// - 安卓手机 (18)
```

### 示例3：完整功能组合

```typescript
// 搜索"智能"，筛选多个分类，包含子分类，返回统计
const result = searchEngine.searchWithStats('智能', {
  categories: ['电子产品', '智能家居'],
  includeSubCategories: true
})

// 结果：
// - 搜索结果：匹配"智能"的所有项
// - 分类筛选：只显示电子产品和智能家居（及其子分类）
// - 分类统计：显示各分类的数量分布
```

---

## 📈 性能优化建议

### 1. 缓存策略

```typescript
// SearchEngine 已内置 LRU 缓存
// 但分类统计不会被缓存，因为返回格式不同

// 如果需要缓存统计结果，可以自己实现：
const statsCache = new LRUCache<string, CategoryStats[]>(50)

const cachedStats = statsCache.get(keyword)
if (!cachedStats) {
  const result = searchEngine.searchWithStats(keyword)
  statsCache.set(keyword, result.categoryStats)
}
```

### 2. 延迟加载分类统计

```typescript
// 只在用户点击"查看分类分布"时才计算
const [showStats, setShowStats] = useState(false)
const [stats, setStats] = useState<CategoryStats[]>([])

const handleShowStats = () => {
  if (!showStats) {
    const result = searchEngine.searchWithStats(keyword)
    setStats(result.categoryStats)
  }
  setShowStats(!showStats)
}
```

### 3. 虚拟滚动（大量分类）

```typescript
// 如果有100+个分类，使用虚拟滚动
import { FixedSizeList } from 'react-window'

<FixedSizeList
  height={400}
  itemCount={categoryStats.length}
  itemSize={35}
>
  {({ index, style }) => (
    <div style={style} className="stat-item">
      <span>{categoryStats[index].name}</span>
      <span>({categoryStats[index].count})</span>
    </div>
  )}
</FixedSizeList>
```

---

## 🎯 实战案例

### 电商搜索页面

```typescript
function EcommerceSearch() {
  const [keyword, setKeyword] = useState('')
  const [result, setResult] = useState<SearchResult | null>(null)
  const [priceRange, setPriceRange] = useState<[number, number]>([0, 10000])
  
  const handleSearch = () => {
    const searchResult = searchEngine.searchWithStats(keyword, {
      categories: ['电子产品', '数码配件', '智能穿戴'],
      includeSubCategories: true
    })
    
    // 再根据价格筛选
    const filtered = searchResult.results.filter(r => 
      r.item.price >= priceRange[0] && r.item.price <= priceRange[1]
    )
    
    setResult({
      results: filtered,
      categoryStats: searchResult.categoryStats
    })
  }
  
  return (
    <div className="ecommerce-search">
      {/* 搜索区域 */}
      <div className="search-bar">
        <input 
          value={keyword}
          onChange={(e) => setKeyword(e.target.value)}
          placeholder="搜索商品..."
        />
        <button onClick={handleSearch}>搜索</button>
      </div>
      
      {/* 筛选和结果 */}
      <div className="search-content">
        {/* 左侧筛选 */}
        <aside className="filters">
          <div className="filter-group">
            <h4>分类 ({result?.categoryStats.length})</h4>
            {result?.categoryStats.map(stat => (
              <label key={stat.name}>
                <input type="checkbox" />
                {stat.name} ({stat.count})
              </label>
            ))}
          </div>
          
          <div className="filter-group">
            <h4>价格范围</h4>
            <Slider 
              range
              value={priceRange}
              onChange={setPriceRange}
              min={0}
              max={10000}
            />
          </div>
        </aside>
        
        {/* 右侧结果 */}
        <main className="results">
          <div className="results-header">
            <span>共 {result?.results.length} 个结果</span>
            <select>
              <option>综合排序</option>
              <option>价格从低到高</option>
              <option>价格从高到低</option>
            </select>
          </div>
          
          <div className="product-grid">
            {result?.results.map(item => (
              <ProductCard key={item.item.id} product={item.item} />
            ))}
          </div>
        </main>
      </div>
    </div>
  )
}
```

---

## 📚 总结

### 功能对比

| 功能 | 单分类 | 多分类 | 层级分类 | 分类统计 |
|------|--------|--------|----------|----------|
| **使用场景** | 基础筛选 | 复合筛选 | 树形筛选 | 数据分析 |
| **实现难度** | ⭐ | ⭐⭐ | ⭐⭐⭐ | ⭐⭐ |
| **性能影响** | 很小 | 小 | 中等 | 小 |
| **用户体验** | 基础 | 良好 | 优秀 | 优秀 |

### 最佳实践

1. **基础搜索**：使用 `search()` + `category`
2. **高级搜索**：使用 `search()` + `categories` + `includeSubCategories`
3. **数据分析**：使用 `searchWithStats()` 获取完整信息
4. **性能优化**：合理使用缓存，延迟加载统计数据
5. **用户体验**：提供清晰的分类导航和结果数量提示

---

**开始使用这些高级功能，打造更强大的搜索体验！** 🚀

