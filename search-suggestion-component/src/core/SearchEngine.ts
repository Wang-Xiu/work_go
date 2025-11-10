import { SuggestionItem, MatchResult, SearchConfig, SearchOptions, MatchType } from '../types'
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
    const results: MatchResult[] = []
    
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
    // 如果看到其他分类的结果，说明筛选没有生效
    //
    // ============================================================
    // 💡 扩展思考：
    // ============================================================
    // 
    // 1. 如果要支持多个分类筛选怎么办？
    //    提示：options.categories 数组，使用 includes()
    //
    // 2. 如果分类是层级结构怎么办？
    //    例如：电子产品 > 手机 > 苹果手机
    //    提示：递归检查父分类
    //
    // 3. 如果要统计每个分类的结果数量怎么办？
    //    提示：使用 reduce() 或 Map
    //
    // ============================================================
    
    // 截取TOP N
    const topResults = results.slice(0, this.config.topN)
    
    // 存入缓存
    this.cache.set(cacheKey, topResults)
    
    return topResults
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

