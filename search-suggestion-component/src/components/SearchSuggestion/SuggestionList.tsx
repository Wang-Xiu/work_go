import React from 'react'
import { MatchResult } from '../../types'
import { SuggestionItem } from './SuggestionItem'

interface SuggestionListProps {
  results: MatchResult[]
  keyword: string
  selectedIndex: number
  onSelect: (result: MatchResult, index: number) => void
  onHover: (index: number) => void
}

/**
 * 建议列表组件
 * 显示搜索结果列表
 */
export const SuggestionList: React.FC<SuggestionListProps> = ({
  results,
  keyword,
  selectedIndex,
  onSelect,
  onHover,
}) => {
  if (results.length === 0) {
    return (
      <div className="suggestion-list empty">
        <div className="empty-message">
          <span className="empty-icon">🔍</span>
          <p>未找到相关结果</p>
          <p className="empty-hint">试试其他关键词吧</p>
        </div>
      </div>
    )
  }

  return (
    <div className="suggestion-list">
      {results.map((result, index) => (
        <SuggestionItem
          key={result.item.id}
          result={result}
          keyword={keyword}
          isSelected={index === selectedIndex}
          onClick={() => onSelect(result, index)}
          onMouseEnter={() => onHover(index)}
        />
      ))}
    </div>
  )
}

