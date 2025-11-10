/**
 * 搜索建议项
 */
export interface SuggestionItem {
  /** 唯一标识 */
  id: string
  /** 显示文本 */
  text: string
  /** 热门度分数 (0-100) */
  hotScore: number
  /** 分类 */
  category: string
  /** 搜索关键词（用于匹配） */
  keywords?: string[]
  /** 拼音（自动生成） */
  pinyin?: string
  /** 拼音首字母（自动生成） */
  pinyinFirst?: string
  /** 描述 */
  description?: string
  /** 图标 */
  icon?: string
}

/**
 * 匹配类型
 */
export enum MatchType {
  /** 前缀匹配 */
  PREFIX = 'prefix',
  /** 完全包含 */
  CONTAINS = 'contains',
  /** 模糊匹配 */
  FUZZY = 'fuzzy',
  /** 拼音匹配 */
  PINYIN = 'pinyin',
  /** 拼音首字母匹配 */
  PINYIN_FIRST = 'pinyin_first',
}

/**
 * 匹配结果
 */
export interface MatchResult {
  /** 建议项 */
  item: SuggestionItem
  /** 匹配类型 */
  matchType: MatchType
  /** 匹配得分 (0-100) */
  matchScore: number
  /** 综合得分 (0-100) */
  finalScore: number
  /** 高亮范围 */
  highlightRanges?: Array<{ start: number; end: number }>
}

/**
 * 搜索配置
 */
export interface SearchConfig {
  /** 返回结果数量 */
  topN: number
  /** 匹配得分权重 (0-1) */
  matchWeight: number
  /** 热门度权重 (0-1) */
  hotWeight: number
  /** 是否启用拼音搜索 */
  enablePinyin: boolean
  /** 是否启用模糊匹配 */
  enableFuzzy: boolean
  /** 最小匹配得分阈值 */
  minMatchScore: number
  /** 防抖延迟(ms) */
  debounceDelay: number
}

/**
 * 搜索选项
 */
export interface SearchOptions {
  /** 筛选分类 */
  category?: string
  /** 自定义权重 */
  customWeight?: {
    matchWeight?: number
    hotWeight?: number
  }
}

/**
 * 分类统计
 */
export interface CategoryStats {
  /** 分类名称 */
  name: string
  /** 数量 */
  count: number
}

