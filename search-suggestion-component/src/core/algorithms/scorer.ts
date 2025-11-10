import { MatchType } from '../../types'

/**
 * 评分算法模块
 * 负责计算匹配得分和综合排序得分
 */

/**
 * 根据匹配类型计算权重
 * 不同匹配类型的重要程度不同
 */
export function getMatchTypeWeight(matchType: MatchType): number {
  switch (matchType) {
    case MatchType.PREFIX:
      return 1.2  // 前缀匹配最重要，加权20%
    case MatchType.CONTAINS:
      return 1.0  // 包含匹配，标准权重
    case MatchType.PINYIN:
      return 0.9  // 拼音匹配，略降权重
    case MatchType.PINYIN_FIRST:
      return 0.8  // 首字母匹配，进一步降权
    case MatchType.FUZZY:
      return 0.7  // 模糊匹配，权重最低
    default:
      return 1.0
  }
}

/**
 * 计算综合得分
 * @param matchScore 匹配得分 0-100
 * @param hotScore 热度得分 0-100
 * @param matchWeight 匹配权重 0-1
 * @param hotWeight 热度权重 0-1
 * @param matchType 匹配类型
 * @returns 综合得分
 */
export function calculateScore(
  matchScore: number,
  hotScore: number,
  matchWeight: number,
  hotWeight: number,
  matchType: MatchType
): number {
  // 根据匹配类型调整匹配得分
  const typeWeight = getMatchTypeWeight(matchType)
  const adjustedMatchScore = matchScore * typeWeight
  
  // 计算加权得分
  const finalScore = adjustedMatchScore * matchWeight + hotScore * hotWeight
  
  return Math.min(100, Math.max(0, finalScore))  // 限制在0-100之间
}

/**
 * 根据位置计算额外加分
 * 出现在开头的匹配应该得分更高
 */
export function getPositionBonus(position: number, textLength: number): number {
  if (position === 0) return 10  // 开头位置，额外加10分
  if (position < textLength * 0.3) return 5  // 前30%，加5分
  return 0
}

/**
 * 根据长度匹配度计算加分
 * 关键词越接近原文长度，得分越高
 */
export function getLengthMatchBonus(keywordLength: number, textLength: number): number {
  const ratio = keywordLength / textLength
  if (ratio > 0.8) return 10  // 匹配80%以上，加10分
  if (ratio > 0.5) return 5   // 匹配50%以上，加5分
  return 0
}

