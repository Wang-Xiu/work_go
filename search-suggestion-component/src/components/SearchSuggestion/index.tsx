import React, { useState, useCallback, useMemo } from 'react'
import { SearchEngine } from '../../core/SearchEngine'
import { SuggestionItem, SearchConfig, SearchOptions, MatchResult } from '../../types'
import { SearchInput } from './SearchInput'
import { SuggestionList } from './SuggestionList'
import { useDebounce } from '../../hooks/useDebounce'
import './styles.css'

export interface SearchSuggestionProps {
  /** 所有建议项数据 */
  items: SuggestionItem[]
  
  /** 搜索配置 */
  config?: Partial<SearchConfig>
  
  /** 选中项回调 */
  onSelect?: (item: SuggestionItem) => void
  
  /** 输入变化回调 */
  onChange?: (keyword: string) => void
  
  /** 占位符文本 */
  placeholder?: string
  
  /** 是否显示分类筛选 */
  showCategoryFilter?: boolean
  
  /** 可用的分类列表 */
  categories?: string[]
}

/**
 * 搜索建议组件
 * 提供输入搜索和建议列表功能
 */
export const SearchSuggestion: React.FC<SearchSuggestionProps> = ({
  items,
  config,
  onSelect,
  onChange,
  placeholder,
  showCategoryFilter = true,
  categories,
}) => {
  // ===== 状态管理 =====
  const [keyword, setKeyword] = useState('')
  const [selectedCategory, setSelectedCategory] = useState<string>('')
  const [selectedIndex, setSelectedIndex] = useState(-1)
  const [isOpen, setIsOpen] = useState(false)
  
  // 防抖处理
  const debouncedKeyword = useDebounce(keyword, config?.debounceDelay || 300)
  
  // 创建搜索引擎实例（只创建一次）
  const searchEngine = useMemo(() => {
    return new SearchEngine(items, config)
  }, [items, config])
  
  // 获取分类列表
  const categoryList = useMemo(() => {
    if (categories) return categories
    
    // 从items中提取唯一分类
    const uniqueCategories = Array.from(
      new Set(items.map(item => item.category))
    ).filter(Boolean)
    
    return uniqueCategories
  }, [items, categories])
  
  // 执行搜索
  const searchOptions: SearchOptions = useMemo(() => ({
    category: selectedCategory || undefined,
  }), [selectedCategory])
  
  const results = useMemo(() => {
    return searchEngine.search(debouncedKeyword, searchOptions)
  }, [searchEngine, debouncedKeyword, searchOptions])
  
  // ===== 事件处理 =====
  
  /**
   * 输入变化
   */
  const handleInputChange = useCallback((value: string) => {
    setKeyword(value)
    setSelectedIndex(-1)  // 重置选中项
    setIsOpen(true)       // 显示建议列表
    onChange?.(value)
  }, [onChange])
  
  /**
   * 输入框获得焦点
   */
  const handleFocus = useCallback(() => {
    setIsOpen(true)
  }, [])
  
  /**
   * 输入框失去焦点
   */
  const handleBlur = useCallback(() => {
    // 延迟关闭，以便点击事件能触发
    setTimeout(() => {
      setIsOpen(false)
    }, 200)
  }, [])
  
  /**
   * 选中项变化
   */
  const handleSelectChange = useCallback((index: number) => {
    setSelectedIndex(index)
  }, [])
  
  /**
   * 确认选择
   */
  const handleConfirm = useCallback(() => {
    if (selectedIndex >= 0 && selectedIndex < results.length) {
      const result = results[selectedIndex]
      setKeyword(result.item.text)
      setIsOpen(false)
      setSelectedIndex(-1)
      onSelect?.(result.item)
    }
  }, [selectedIndex, results, onSelect])
  
  /**
   * 关闭建议列表
   */
  const handleClose = useCallback(() => {
    setIsOpen(false)
    setSelectedIndex(-1)
  }, [])
  
  /**
   * 点击选择
   */
  const handleItemSelect = useCallback((result: MatchResult, index: number) => {
    setKeyword(result.item.text)
    setIsOpen(false)
    setSelectedIndex(-1)
    onSelect?.(result.item)
  }, [onSelect])
  
  /**
   * 鼠标悬停
   */
  const handleItemHover = useCallback((index: number) => {
    setSelectedIndex(index)
  }, [])
  
  /**
   * 分类变化
   */
  const handleCategoryChange = useCallback((category: string) => {
    setSelectedCategory(category)
    setSelectedIndex(-1)
  }, [])
  
  return (
    <div className="search-suggestion">
      {/* 分类筛选 */}
      {showCategoryFilter && categoryList.length > 0 && (
        <div className="category-filter">
          <button
            className={`category-button ${selectedCategory === '' ? 'active' : ''}`}
            onClick={() => handleCategoryChange('')}
          >
            全部
          </button>
          {categoryList.map(category => (
            <button
              key={category}
              className={`category-button ${selectedCategory === category ? 'active' : ''}`}
              onClick={() => handleCategoryChange(category)}
            >
              {category}
            </button>
          ))}
        </div>
      )}
      
      {/* 搜索输入框 */}
      <SearchInput
        value={keyword}
        onChange={handleInputChange}
        onFocus={handleFocus}
        onBlur={handleBlur}
        placeholder={placeholder}
        selectedIndex={selectedIndex}
        resultCount={results.length}
        onSelectChange={handleSelectChange}
        onConfirm={handleConfirm}
        onClose={handleClose}
      />
      
      {/* 建议列表 */}
      {isOpen && (
        <SuggestionList
          results={results}
          keyword={keyword}
          selectedIndex={selectedIndex}
          onSelect={handleItemSelect}
          onHover={handleItemHover}
        />
      )}
    </div>
  )
}

export default SearchSuggestion

