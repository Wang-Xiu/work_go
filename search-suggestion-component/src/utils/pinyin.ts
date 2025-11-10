import { pinyin } from 'pinyin-pro'

/**
 * 拼音工具类
 */
export class PinyinUtil {
  /**
   * 获取全拼音（带音调）
   * @example getPinyin("你好") => "ni hao"
   */
  static getPinyin(text: string): string {
    return pinyin(text, { toneType: 'none', type: 'array' }).join(' ')
  }

  /**
   * 获取拼音首字母
   * @example getFirstLetter("你好") => "nh"
   */
  static getFirstLetter(text: string): string {
    return pinyin(text, { pattern: 'first', type: 'array' }).join('')
  }

  /**
   * 判断是否匹配拼音
   * @param text 原文本
   * @param keyword 搜索关键词
   */
  static matchPinyin(text: string, keyword: string): boolean {
    const pinyinText = this.getPinyin(text).toLowerCase()
    const searchKeyword = keyword.toLowerCase()
    return pinyinText.includes(searchKeyword)
  }

  /**
   * 判断是否匹配拼音首字母
   * @param text 原文本
   * @param keyword 搜索关键词
   */
  static matchPinyinFirst(text: string, keyword: string): boolean {
    const firstLetters = this.getFirstLetter(text).toLowerCase()
    const searchKeyword = keyword.toLowerCase()
    return firstLetters.includes(searchKeyword)
  }
}

