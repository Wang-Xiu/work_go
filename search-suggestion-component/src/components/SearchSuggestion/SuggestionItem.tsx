import React from 'react'
import { MatchResult, MatchType } from '../../types'

interface SuggestionItemProps {
  result: MatchResult
  keyword: string
  isSelected: boolean
  onClick: () => void
  onMouseEnter: () => void
}

/**
 * 建议项组件
 * 显示单个搜索建议
 */
export const SuggestionItem: React.FC<SuggestionItemProps> = ({
  result,
  keyword,
  isSelected,
  onClick,
  onMouseEnter,
}) => {
  const { item, matchType, matchScore, finalScore } = result
  
  // ============================================================
  // 🎯 【练手任务5】实现关键词高亮功能
  // ============================================================
  
  /**
   * 高亮文本中的关键词
   * @param text 要处理的文本
   * @param keyword 关键词
   * @returns 带高亮的JSX元素
   */
  const highlightText = (text: string, keyword: string): JSX.Element => {
    // ============================================================
    // 📝 实现思路：
    // ============================================================
    // 1. 检查keyword是否为空
    // 2. 在text中查找keyword的位置（忽略大小写）
    // 3. 如果找不到，直接返回原文本
    // 4. 如果找到了：
    //    - 分割为三部分：前缀、匹配部分、后缀
    //    - 用<mark>标签包裹匹配部分
    // 5. 返回JSX
    //
    // ============================================================
    // 📌 第1步：检查空值
    // ============================================================
    
    // 👉 你的代码：检查keyword是否为空
    // 
    // 条件：
    //   - !keyword：keyword为null/undefined/空字符串
    //   - keyword.trim() === ''：keyword只包含空格
    //
    // 如果为空，返回：<span>{text}</span>
    //
    // 写在这里：
    if (!keyword || keyword.trim() === '') {
      return <span>{text}</span>
    }
    
    // ============================================================
    // 📌 第2步：查找关键词位置
    // ============================================================
    
    // 👉 你的代码：使用indexOf查找位置
    //
    // 注意：要忽略大小写！
    //   - 先把text和keyword都转为小写
    //   - 再使用indexOf
    //
    // 代码示例：
    //   const index = text.toLowerCase().indexOf(keyword.toLowerCase())
    //
    // 写在这里：
    const index = text.toLowerCase().indexOf(keyword.toLowerCase())
    
    // ============================================================
    // 📌 第3步：检查是否找到
    // ============================================================
    
    // 👉 你的代码：检查index
    //
    // 如果index === -1，说明没找到，返回原文本
    //
    // 代码：
    //   if (index === -1) {
    //     return <span>{text}</span>
    //   }
    //
    // 写在这里：
    if (index === -1) {
      return <span>{text}</span>
    }
    
    // ============================================================
    // 📌 第4步：分割文本
    // ============================================================
    
    // 👉 你的代码：使用slice分割
    //
    // 分割成三部分：
    //   1. before：匹配前的部分
    //      从0到index
    //   2. match：匹配的部分
    //      从index到index+keyword.length
    //   3. after：匹配后的部分
    //      从index+keyword.length到末尾
    //
    // 代码示例：
    //   const before = text.slice(0, index)
    //   const match = text.slice(index, index + keyword.length)
    //   const after = text.slice(index + keyword.length)
    //
    // 写在这里：
    const before = text.slice(0, index)
    const match = text.slice(index, index + keyword.length)
    const after = text.slice(index + keyword.length)
    
    // ============================================================
    // 📌 第5步：返回JSX
    // ============================================================
    
    // 👉 你的代码：拼接JSX
    //
    // 使用<mark>标签包裹匹配部分
    // <mark>是HTML5语义化标签，专门用于高亮
    //
    // 代码：
    //   return (
    //     <span>
    //       {before}
    //       <mark>{match}</mark>
    //       {after}
    //     </span>
    //   )
    //
    // 写在这里：
    return (
      <span>
        {before}
        <mark>{match}</mark>
        {after}
      </span>
    )
    
    // ============================================================
    // ✅ 完成！测试你的实现：
    // ============================================================
    //
    // 测试步骤：
    // 1. 启动应用：npm run dev
    // 2. 输入搜索关键词"手机"
    // 3. 在结果列表中，"手机"两个字应该高亮显示
    // 4. 输入"iphone"
    // 5. "iPhone"中的"iphone"部分应该高亮
    // 6. 输入"pg"（拼音首字母）
    // 7. "苹果"应该高亮（拼音匹配）
    //
    // ============================================================
    // 💡 进阶优化（可选）：
    // ============================================================
    //
    // 1. 支持多个关键词高亮：
    //    例如："iPhone 15 Pro"，高亮"iPhone"和"Pro"
    //    提示：使用正则表达式 /keyword1|keyword2/gi
    //
    // 2. 支持拼音高亮：
    //    例如：输入"pingguo"，高亮"苹果"
    //    提示：需要记录拼音匹配的位置
    //
    // 3. 支持模糊匹配高亮：
    //    例如：输入"iph"，高亮"iPhone"中的"iPh"
    //    提示：需要记录fuzzyMatch返回的匹配位置
    //
    // 4. 使用不同颜色标记不同匹配类型：
    //    - 前缀匹配：绿色
    //    - 包含匹配：黄色
    //    - 拼音匹配：蓝色
    //    - 模糊匹配：灰色
    //    提示：根据matchType动态设置className
    //
    // ============================================================
  }
  
  /**
   * 获取匹配类型的中文显示
   */
  const getMatchTypeLabel = (type: MatchType): string => {
    switch (type) {
      case MatchType.PREFIX:
        return '前缀匹配'
      case MatchType.CONTAINS:
        return '包含匹配'
      case MatchType.PINYIN:
        return '拼音匹配'
      case MatchType.PINYIN_FIRST:
        return '首字母'
      case MatchType.FUZZY:
        return '模糊匹配'
      default:
        return ''
    }
  }
  
  /**
   * 获取匹配类型的样式类名
   */
  const getMatchTypeClass = (type: MatchType): string => {
    return `match-type-${type.toLowerCase().replace('_', '-')}`
  }
  
  return (
    <div
      className={`suggestion-item ${isSelected ? 'selected' : ''}`}
      onClick={onClick}
      onMouseEnter={onMouseEnter}
    >
      <div className="suggestion-main">
        <div className="suggestion-icon">
          {item.icon || '🔍'}
        </div>
        
        <div className="suggestion-content">
          <div className="suggestion-text">
            {highlightText(item.text, keyword)}
          </div>
          
          {item.description && (
            <div className="suggestion-description">
              {item.description}
            </div>
          )}
        </div>
      </div>
      
      <div className="suggestion-meta">
        <span className={`suggestion-category ${getMatchTypeClass(matchType)}`}>
          {item.category}
        </span>
        
        {/* 开发模式下显示调试信息 */}
        {process.env.NODE_ENV === 'development' && (
          <div className="suggestion-debug">
            <span className="match-type-badge" title={getMatchTypeLabel(matchType)}>
              {getMatchTypeLabel(matchType)}
            </span>
            <span className="score-badge" title={`匹配分: ${matchScore} | 综合分: ${finalScore.toFixed(1)}`}>
              {finalScore.toFixed(0)}
            </span>
          </div>
        )}
      </div>
    </div>
  )
}

