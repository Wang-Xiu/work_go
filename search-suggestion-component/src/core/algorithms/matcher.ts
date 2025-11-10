import { MatchType } from '../../types'
import { PinyinUtil } from '../../utils/pinyin'

/**
 * 匹配算法模块
 * 负责判断搜索关键词与建议项的匹配程度
 */

/**
 * 前缀匹配
 * 最高优先级，完全匹配开头
 * @example prefixMatch("iPhone", "iph") => true
 */
export function prefixMatch(text: string, keyword: string): boolean {
  return text.toLowerCase().startsWith(keyword.toLowerCase())
}

/**
 * 包含匹配
 * 判断文本是否包含关键词
 * @example containsMatch("MacBook Pro", "book") => true
 */
export function containsMatch(text: string, keyword: string): boolean {
  return text.toLowerCase().includes(keyword.toLowerCase())
}

/**
 * 拼音匹配
 * 支持中文拼音搜索
 * @example pinyinMatch("苹果手机", "pingguo") => true
 */
export function pinyinMatch(text: string, keyword: string): boolean {
  return PinyinUtil.matchPinyin(text, keyword)
}

/**
 * 拼音首字母匹配
 * 支持拼音首字母缩写搜索
 * @example pinyinFirstMatch("苹果手机", "pgsj") => true
 */
export function pinyinFirstMatch(text: string, keyword: string): boolean {
  return PinyinUtil.matchPinyinFirst(text, keyword)
}

// ============================================================
// 🎯 【练手任务1】实现模糊匹配算法
// ============================================================

/**
 * 模糊匹配算法
 * 允许关键词中的字符在文本中不连续出现
 * 例如：输入"iph"可以匹配"iPhone"
 * 
 * @param text 要匹配的文本
 * @param keyword 搜索关键词
 * @returns 匹配得分 0-100（数字越大表示匹配度越高）
 */
export function fuzzyMatch(text: string, keyword: string): number {
  // ============================================================
  // 📝 实现步骤说明：
  // ============================================================
  // 
  // 这个函数的目标是：判断keyword中的每个字符是否能在text中
  // 按顺序找到（可以跳过中间的字符）
  //
  // 例如：
  //   text = "iPhone 15 Pro"
  //   keyword = "iph"
  //   
  //   在text中查找：
  //   i -> 找到（位置0）
  //   p -> 找到（位置1）
  //   h -> 找到（位置2）
  //   
  //   所以匹配成功，返回 100（完全匹配）
  //
  // 又如：
  //   text = "MacBook Pro"
  //   keyword = "mbp"
  //   
  //   在text中查找：
  //   m -> 找到（位置0）
  //   b -> 找到（位置3）
  //   p -> 找到（位置8）
  //   
  //   所以匹配成功，返回 100
  //
  // ============================================================
  // 📌 第1步：参数预处理
  // ============================================================
  // 提示：将text和keyword都转为小写，方便比较
  // 代码示例：
  //   const lowerText = text.toLowerCase()
  //   const lowerKeyword = keyword.toLowerCase()
  
  // 👉 你的代码：在这里转换为小写
  const lowerText = text.toLowerCase()
  const lowerKeyword = keyword.toLowerCase()
  
  // ============================================================
  // 📌 第2步：边界条件检查
  // ============================================================
  // 提示：
  // - 如果keyword为空，返回0（没有输入，不匹配）
  // - 如果text为空，返回0（文本为空，无法匹配）
  // - 如果keyword比text还长，返回0（不可能匹配）
  
  // 👉 你的代码：在这里检查边界条件
  if (!lowerKeyword || !lowerText) return 0
  if (lowerKeyword.length > lowerText.length) return 0
  
  // ============================================================
  // 📌 第3步：使用双指针算法进行匹配
  // ============================================================
  // 解释：什么是双指针？
  //   - textIndex：指向text的当前位置
  //   - keywordIndex：指向keyword的当前位置
  //   - 从左到右扫描，如果字符匹配就移动keyword指针
  //
  // 算法流程：
  //   1. textIndex从0开始遍历text的每个字符
  //   2. 如果text[textIndex] == keyword[keywordIndex]
  //      -> 找到匹配，keywordIndex++，matchCount++
  //   3. 继续遍历，直到keyword全部匹配完或text遍历完
  //   4. 最后计算匹配率
  
  // 👉 你的代码：定义变量
  // 提示：需要定义三个变量
  //   - textIndex: number = 0（text的索引）
  //   - keywordIndex: number = 0（keyword的索引）
  //   - matchCount: number = 0（已匹配的字符数）
  
  // 写在这里：
  let textIndex = 0
  let keywordIndex = 0
  let matchCount = 0
  
  // ============================================================
  // 📌 第4步：编写while循环进行匹配
  // ============================================================
  // 👉 你的代码：完成while循环
  // 
  // 循环条件：
  //   while (textIndex < lowerText.length && keywordIndex < lowerKeyword.length)
  //
  // 循环体内的逻辑：
  //   1. 比较 lowerText[textIndex] 和 lowerKeyword[keywordIndex]
  //   2. 如果相等：
  //      - matchCount++（匹配数+1）
  //      - keywordIndex++（keyword指针后移）
  //   3. 无论是否匹配，textIndex都要++（text指针后移）
  //
  // 代码框架：
  //   while (循环条件) {
  //     if (字符相等) {
  //       matchCount++
  //       keywordIndex++
  //     }
  //     textIndex++
  //   }
  
  // 写在这里：
  while (textIndex < lowerText.length && keywordIndex < lowerKeyword.length) {
    if (lowerText[textIndex] === lowerKeyword[keywordIndex]) {
      matchCount++
      keywordIndex++
    }
    textIndex++
  }
  
  // ============================================================
  // 📌 第5步：计算匹配得分
  // ============================================================
  // 解释：
  //   - 如果matchCount < keyword.length，说明没有完全匹配，返回0
  //   - 如果完全匹配，计算得分：
  //     得分 = (匹配字符数 / 关键词长度) × 100
  //
  // 例如：
  //   keyword = "iph"（长度3）
  //   matchCount = 3（全部匹配）
  //   得分 = (3 / 3) × 100 = 100
  
  // 👉 你的代码：计算并返回得分
  // 
  // 步骤1：判断是否完全匹配
  //   if (matchCount < lowerKeyword.length) {
  //     return 0  // 没有完全匹配
  //   }
  //
  // 步骤2：计算得分
  //   return (matchCount / lowerKeyword.length) * 100
  
  // 写在这里：
  if (matchCount < lowerKeyword.length) {
    return 0  // 没有完全匹配，说明有字符没找到
  }
  
  // 计算匹配得分（0-100）
  return (matchCount / lowerKeyword.length) * 100
  
  // ============================================================
  // ✅ 完成！测试你的实现：
  // ============================================================
  // 在浏览器console中测试：
  //   fuzzyMatch("iPhone 15 Pro", "iph")  // 应该返回 100
  //   fuzzyMatch("MacBook Pro", "mbp")    // 应该返回 100
  //   fuzzyMatch("iPad Air", "iph")       // 应该返回 0（没有'h'）
  //   fuzzyMatch("AirPods", "apo")        // 应该返回 100
  //
  // ============================================================
  // 💡 进阶优化（可选）：
  // ============================================================
  // 如果你完成了基础版本，可以考虑这些优化：
  //
  // 1. 考虑匹配密度：
  //    字符越紧密得分越高
  //    例如："iPhone" vs "iPnone"，第一个应该得分更高
  //
  // 2. 考虑匹配位置：
  //    在开头匹配的得分更高
  //    例如："iPhone Pro" 搜索 "iph" 应该比 "Prp ihone" 得分高
  //
  // 3. 使用编辑距离算法（Levenshtein Distance）：
  //    允许1-2个字符的拼写错误
  //    例如："iphone" 也能匹配 "ifone"（f→ph）
  //
  // 这些优化可以在完成基础版本后再实现！
  // ============================================================
}

/**
 * 综合匹配
 * 按优先级尝试各种匹配方式
 * @returns { matchType, score }
 */
export function match(text: string, keyword: string): { matchType: MatchType; score: number } {
  if (!keyword) {
    return { matchType: MatchType.PREFIX, score: 0 }
  }

  // 按优先级顺序尝试匹配
  if (prefixMatch(text, keyword)) {
    return { matchType: MatchType.PREFIX, score: 100 }
  }

  if (containsMatch(text, keyword)) {
    return { matchType: MatchType.CONTAINS, score: 80 }
  }

  if (pinyinMatch(text, keyword)) {
    return { matchType: MatchType.PINYIN, score: 70 }
  }

  if (pinyinFirstMatch(text, keyword)) {
    return { matchType: MatchType.PINYIN_FIRST, score: 60 }
  }

  const fuzzyScore = fuzzyMatch(text, keyword)
  if (fuzzyScore > 0) {
    return { matchType: MatchType.FUZZY, score: 40 }
  }

  return { matchType: MatchType.FUZZY, score: 0 }
}

