import { SuggestionItem, MatchResult, SearchConfig, SearchOptions, MatchType, CategoryStats } from '../types'
import { match } from './algorithms/matcher'
import { PinyinUtil } from '../utils/pinyin'
import { LRUCache } from '../utils/cache'

/**
 * 搜索引擎类
 * 负责搜索建议项并返回排序后的结果
 */
export class SearchEngine {
  private items: SuggestionItem[]
  private config: SearchConfig
  private cache: LRUCache<string, MatchResult[]>
  
  /**
   * 构造函数
   * @param items 所有建议项数据
   * @param config 搜索配置
   */
  constructor(items: SuggestionItem[], config?: Partial<SearchConfig>) {
    this.items = items
    this.config = {
      topN: config?.topN || 10,
      matchWeight: config?.matchWeight || 0.6,
      hotWeight: config?.hotWeight || 0.4,
      enablePinyin: config?.enablePinyin !== false,
      enableFuzzy: config?.enableFuzzy !== false,
      minMatchScore: config?.minMatchScore || 0,
      debounceDelay: config?.debounceDelay || 300,
    }
    this.cache = new LRUCache<string, MatchResult[]>(100)
    
    // 预处理：为每个item生成拼音
    this.preprocessItems()
  }
  
  /**
   * 预处理：生成拼音索引
   */
  private preprocessItems(): void {
    this.items.forEach(item => {
      if (!item.pinyin) {
        item.pinyin = PinyinUtil.getPinyin(item.text)
      }
      if (!item.pinyinFirst) {
        item.pinyinFirst = PinyinUtil.getFirstLetter(item.text)
      }
    })
  }
  
  /**
   * 搜索方法
   * @param keyword 搜索关键词
   * @param options 搜索选项（可选分类筛选）
   * @returns 匹配结果列表（已排序）
   */
  search(keyword: string, options?: SearchOptions): MatchResult[] {
    // 空关键词，返回热门推荐
    if (!keyword || keyword.trim() === '') {
      return this.getHotRecommendations(options)
    }
    
    // 检查缓存
    const cacheKey = this.getCacheKey(keyword, options)
    const cached = this.cache.get(cacheKey)
    if (cached) {
      return cached
    }
    
    // 执行搜索
    let results: MatchResult[] = []
    
    for (const item of this.items) {
      // 执行匹配
      const matchResult = match(item.text, keyword)
      
      if (matchResult.score > this.config.minMatchScore) {
        // 计算综合得分
        const finalScore = this.calculateFinalScore(
          matchResult.score,
          item.hotScore,
          matchResult.matchType
        )
        
        results.push({
          item,
          matchType: matchResult.matchType,
          matchScore: matchResult.score,
          finalScore,
        })
      }
    }
    
    // 按综合得分排序
    results.sort((a, b) => b.finalScore - a.finalScore)
    
    // ============================================================
    // 🎯 【练手任务2】实现分类筛选功能
    // ============================================================
    // 
    // 任务说明：
    //   如果用户选择了某个分类，只返回该分类下的结果
    //   例如：用户选择了"电子产品"，只显示电子产品的搜索结果
    //
    // ============================================================
    // 📝 实现思路：
    // ============================================================
    // 1. 检查 options 参数是否存在
    // 2. 检查 options.category 是否有值
    // 3. 如果有值，使用 Array.filter() 过滤 results 数组
    // 4. 只保留 result.item.category === options.category 的项
    //
    // ============================================================
    // 📌 实现步骤：
    // ============================================================
    
    // 👉 你的代码：实现分类筛选
    //
    // 步骤1：检查是否指定了分类
    //   使用：options?.category（可选链操作符，安全访问）
    //
    // 步骤2：使用filter过滤results
    //   语法：results = results.filter(条件函数)
    //
    // 步骤3：条件函数检查分类是否匹配
    //   条件：result.item.category === options.category
    //
    // 完整代码框架：
    //   if (options?.category) {
    //     results = results.filter(result => 
    //       result.item.category === options.category
    //     )
    //   }
    //
    // 写在这里：
    if (options?.category) {
      results = results.filter(result => 
        result.item.category === options.category
      )
    }
    
    // ============================================================
    // ✅ 完成！测试你的实现：
    // ============================================================
    // 
    // 测试步骤：
    // 1. 启动应用：npm run dev
    // 2. 在分类下拉框选择"电子产品"
    // 3. 输入搜索关键词"手机"
    // 4. 应该只显示"电子产品"分类的结果
    // 5. 切换到"美食餐饮"分类
    // 6. 输入"火锅"
    // 7. 应该只显示"美食餐饮"分类的结果
    //
    // ============================================================
    // 🚀 【扩展功能1】支持多个分类筛选
    // ============================================================
    //
    // 应用场景：
    //   用户想同时查看"电子产品"和"数码配件"的搜索结果
    //   例如：搜索"充电"，可能既有充电宝（电子产品）也有充电线（数码配件）
    //
    // 实现思路：
    //   1. 检查 options.categories 是否存在（数组）
    //   2. 使用 Array.includes() 检查 item.category 是否在数组中
    //   3. 注意：categories 和 category 可能同时存在，categories 优先级更高
    //
    // ============================================================
    
    // 多分类筛选（优先级高于单分类）
    if (options?.categories && options.categories.length > 0) {
      // 使用 filter + includes 筛选多个分类
      // 解释：只保留 category 在 categories 数组中的结果
      results = results.filter(result =>
        options.categories!.includes(result.item.category)
      )
    }
    
    // ============================================================
    // 🚀 【扩展功能2】支持层级分类筛选
    // ============================================================
    //
    // 应用场景：
    //   分类是树形结构，例如：
    //     电子产品
    //       ├─ 手机
    //       │   ├─ 苹果手机
    //       │   └─ 安卓手机
    //       └─ 电脑
    //           ├─ 笔记本
    //           └─ 台式机
    //
    //   用户选择"手机"，应该同时显示"苹果手机"和"安卓手机"的结果
    //
    // 实现方式：
    //   1. 在 SuggestionItem 中添加 parentCategory 字段
    //   2. 检查 item.category 或 item.parentCategory 是否匹配
    //   3. 支持递归查找（如果需要多层级）
    //
    // ============================================================
    
    // 层级分类筛选（includeSubCategories=true 时生效）
    if (options?.includeSubCategories && (options?.category || options?.categories)) {
      // 如果启用了子分类包含，需要检查 parentCategory
      // 这样选择父分类时，子分类的结果也会显示
      
      // 获取筛选条件（单个或多个）
      const targetCategories = options.categories || (options.category ? [options.category] : [])
      
      if (targetCategories.length > 0) {
        results = results.filter(result => {
          const item = result.item
          
          // 方式1：直接匹配当前分类
          if (targetCategories.includes(item.category)) {
            return true
          }
          
          // 方式2：匹配父分类（如果当前item是子分类）
          if (item.parentCategory && targetCategories.includes(item.parentCategory)) {
            return true
          }
          
          // 方式3：递归查找祖先分类（可选，用于多层级）
          // 这里可以实现一个 hasAncestorCategory 方法
          // return this.hasAncestorCategory(item, targetCategories)
          
          return false
        })
      }
    }
    
    // ============================================================
    // 🚀 【扩展功能3】统计每个分类的结果数量
    // ============================================================
    //
    // 应用场景：
    //   在搜索结果旁边显示每个分类有多少个匹配项
    //   例如：
    //     电子产品 (23)
    //     美食餐饮 (15)
    //     旅游景点 (8)
    //
    //   用户可以快速了解各分类的分布情况
    //
    // 实现方式：
    //   使用 reduce() 或 Map 统计每个分类的数量
    //   返回格式：{ category: 分类名, count: 数量 }[]
    //
    // ============================================================
    
    // 如果需要返回分类统计
    // 注意：由于 search() 方法的返回类型是 MatchResult[]
    // 无法直接返回统计信息，请使用 searchWithStats() 方法
    // 该方法会返回 { results, categoryStats } 格式的完整数据
    
    // 截取TOP N
    const topResults = results.slice(0, this.config.topN)
    
    // 存入缓存
    this.cache.set(cacheKey, topResults)
    
    return topResults
  }
  
  /**
   * 搜索方法（带分类统计）
   * @param keyword 搜索关键词
   * @param options 搜索选项
   * @returns 搜索结果（包含匹配列表和分类统计）
   */
  searchWithStats(keyword: string, options?: SearchOptions) {
    // ============================================================
    // 🚀 扩展方法：返回结果 + 分类统计
    // ============================================================
    //
    // 与普通 search() 的区别：
    //   - search()：只返回 MatchResult[]
    //   - searchWithStats()：返回 { results: MatchResult[], categoryStats: CategoryStats[] }
    //
    // 使用方式：
    //   const result = searchEngine.searchWithStats("iPhone")
    //   console.log(result.results)         // 搜索结果
    //   console.log(result.categoryStats)   // 分类统计
    //
    // ============================================================
    
    // 执行普通搜索（但不截取TOP N）
    const allResults = this.searchInternal(keyword, options)
    
    // 计算分类统计
    const categoryStats = this.calculateCategoryStats(allResults)
    
    // 截取TOP N
    const topResults = allResults.slice(0, this.config.topN)
    
    return {
      results: topResults,
      categoryStats: categoryStats,
    }
  }
  
  /**
   * 内部搜索方法（不截取TOP N）
   * 供 searchWithStats 使用
   */
  private searchInternal(keyword: string, options?: SearchOptions): MatchResult[] {
    // 空关键词，返回热门推荐
    if (!keyword || keyword.trim() === '') {
      return this.getHotRecommendations(options)
    }
    
    // 执行搜索
    const results: MatchResult[] = []
    
    for (const item of this.items) {
      const matchResult = match(item.text, keyword)
      
      if (matchResult.score > this.config.minMatchScore) {
        const finalScore = this.calculateFinalScore(
          matchResult.score,
          item.hotScore,
          matchResult.matchType
        )
        
        results.push({
          item,
          matchType: matchResult.matchType,
          matchScore: matchResult.score,
          finalScore,
        })
      }
    }
    
    // 按综合得分排序
    results.sort((a, b) => b.finalScore - a.finalScore)
    
    // 应用分类筛选（复用现有逻辑）
    return this.applyCategoryFilters(results, options)
  }
  
  /**
   * 应用分类筛选
   * 提取出来便于复用
   */
  private applyCategoryFilters(results: MatchResult[], options?: SearchOptions): MatchResult[] {
    let filtered = results
    
    // 单分类筛选
    if (options?.category) {
      filtered = filtered.filter(result => 
        result.item.category === options.category
      )
    }
    
    // 多分类筛选（优先级高于单分类）
    if (options?.categories && options.categories.length > 0) {
      filtered = filtered.filter(result =>
        options.categories!.includes(result.item.category)
      )
    }
    
    // 层级分类筛选
    if (options?.includeSubCategories && (options?.category || options?.categories)) {
      const targetCategories = options.categories || (options.category ? [options.category] : [])
      
      if (targetCategories.length > 0) {
        filtered = filtered.filter(result => {
          const item = result.item
          
          if (targetCategories.includes(item.category)) {
            return true
          }
          
          if (item.parentCategory && targetCategories.includes(item.parentCategory)) {
            return true
          }
          
          return false
        })
      }
    }
    
    return filtered
  }
  
  /**
   * 计算分类统计
   * @param results 搜索结果
   * @returns 分类统计数组
   */
  private calculateCategoryStats(results: MatchResult[]): CategoryStats[] {
    // ============================================================
    // 📝 实现思路：统计每个分类的数量
    // ============================================================
    //
    // 方法1：使用 reduce() + 对象
    //   const stats = results.reduce((acc, result) => {
    //     const category = result.item.category
    //     acc[category] = (acc[category] || 0) + 1
    //     return acc
    //   }, {} as Record<string, number>)
    //
    // 方法2：使用 Map（更推荐，性能更好）
    //   const map = new Map<string, number>()
    //   results.forEach(result => {
    //     const category = result.item.category
    //     map.set(category, (map.get(category) || 0) + 1)
    //   })
    //
    // 方法3：使用 for...of（最简单直观）
    //   for (const result of results) {
    //     // 统计逻辑
    //   }
    //
    // ============================================================
    
    // 使用 Map 统计（性能最优）
    const statsMap = new Map<string, number>()
    
    // 遍历所有结果，统计每个分类的数量
    for (const result of results) {
      const category = result.item.category
      
      // 获取当前分类的计数，如果不存在则为0
      const currentCount = statsMap.get(category) || 0
      
      // 更新计数（+1）
      statsMap.set(category, currentCount + 1)
    }
    
    // 将 Map 转换为数组格式
    const categoryStats: CategoryStats[] = []
    
    statsMap.forEach((count, name) => {
      categoryStats.push({
        name,   // 分类名称
        count,  // 数量
      })
    })
    
    // 按数量降序排序（数量多的在前）
    categoryStats.sort((a, b) => b.count - a.count)
    
    return categoryStats
    
    // ============================================================
    // ✅ 完成！使用示例：
    // ============================================================
    //
    // 基本用法：
    //   const result = searchEngine.searchWithStats("手机")
    //   console.log(result.categoryStats)
    //   // 输出：
    //   // [
    //   //   { name: "电子产品", count: 15 },
    //   //   { name: "数码配件", count: 8 },
    //   //   { name: "智能家居", count: 3 }
    //   // ]
    //
    // 在UI中显示：
    //   <div className="category-stats">
    //     {result.categoryStats.map(stat => (
    //       <div key={stat.name} className="stat-item">
    //         <span>{stat.name}</span>
    //         <span className="count">({stat.count})</span>
    //       </div>
    //     ))}
    //   </div>
    //
    // 实际应用场景：
    //   1. 搜索结果页面的侧边栏分类筛选
    //   2. 显示各分类的结果数量
    //   3. 帮助用户快速定位感兴趣的分类
    //   4. 数据分析：了解搜索结果的分布
    //
    // ============================================================
  }
  
  /**
   * 获取热门推荐（keyword为空时）
   */
  private getHotRecommendations(options?: SearchOptions): MatchResult[] {
    let items = [...this.items]
    
    // 分类筛选
    if (options?.category) {
      items = items.filter(item => item.category === options.category)
    }
    
    // 按热度排序
    items.sort((a, b) => b.hotScore - a.hotScore)
    
    // 转换为MatchResult格式
    return items.slice(0, this.config.topN).map(item => ({
      item,
      matchType: MatchType.PREFIX,
      matchScore: 0,
      finalScore: item.hotScore,
    }))
  }
  
  /**
   * 计算综合得分
   * 综合得分 = 匹配得分 × 匹配权重 + 热门度 × 热门权重
   */
  private calculateFinalScore(
    matchScore: number,
    hotScore: number,
    matchType: MatchType
  ): number {
    const { matchWeight, hotWeight } = this.config
    
    // 根据匹配类型调整权重
    let adjustedMatchScore = matchScore
    switch (matchType) {
      case MatchType.PREFIX:
        adjustedMatchScore *= 1.2  // 前缀匹配加权
        break
      case MatchType.CONTAINS:
        adjustedMatchScore *= 1.0
        break
      case MatchType.PINYIN:
        adjustedMatchScore *= 0.9
        break
      case MatchType.PINYIN_FIRST:
        adjustedMatchScore *= 0.8
        break
      case MatchType.FUZZY:
        adjustedMatchScore *= 0.7
        break
    }
    
    return adjustedMatchScore * matchWeight + hotScore * hotWeight
  }
  
  /**
   * 生成缓存key
   */
  private getCacheKey(keyword: string, options?: SearchOptions): string {
    const category = options?.category || 'all'
    return `${keyword}:${category}`
  }
  
  /**
   * 清除缓存
   */
  clearCache(): void {
    this.cache.clear()
  }
  
  /**
   * 更新配置
   */
  updateConfig(config: Partial<SearchConfig>): void {
    this.config = { ...this.config, ...config }
    this.clearCache()  // 配置变更，清除缓存
  }
  
  /**
   * 添加建议项
   */
  addItem(item: SuggestionItem): void {
    // 生成拼音
    if (!item.pinyin) {
      item.pinyin = PinyinUtil.getPinyin(item.text)
    }
    if (!item.pinyinFirst) {
      item.pinyinFirst = PinyinUtil.getFirstLetter(item.text)
    }
    this.items.push(item)
    this.clearCache()
  }
  
  /**
   * 批量添加建议项
   */
  addItems(items: SuggestionItem[]): void {
    items.forEach(item => this.addItem(item))
  }
  
  /**
   * 更新热度分数
   */
  updateHotScore(itemId: string, newScore: number): void {
    const item = this.items.find(i => i.id === itemId)
    if (item) {
      item.hotScore = newScore
      this.clearCache()
    }
  }
}

