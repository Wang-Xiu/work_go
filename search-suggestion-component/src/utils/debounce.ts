/**
 * 防抖函数
 * 在事件被触发n秒后再执行回调，如果在n秒内又被触发，则重新计时
 * 
 * @param func 要执行的函数
 * @param delay 延迟时间（毫秒）
 * @returns 防抖后的函数
 * 
 * @example
 * const debouncedSearch = debounce((keyword) => {
 *   console.log('搜索:', keyword)
 * }, 300)
 * 
 * debouncedSearch('a')  // 不会立即执行
 * debouncedSearch('ab') // 取消上次，重新计时
 * debouncedSearch('abc') // 取消上次，重新计时
 * // 300ms后才会执行最后一次
 */
export function debounce<T extends (...args: any[]) => any>(
  func: T,
  delay: number
): (...args: Parameters<T>) => void {
  let timeoutId: ReturnType<typeof setTimeout> | null = null
  
  return function(...args: Parameters<T>) {
    // 如果已有定时器，清除它
    if (timeoutId) {
      clearTimeout(timeoutId)
    }
    
    // 设置新的定时器
    timeoutId = setTimeout(() => {
      func(...args)
      timeoutId = null
    }, delay)
  }
}

/**
 * 带立即执行的防抖函数
 * 第一次触发立即执行，后续触发需要等待
 * 
 * @param func 要执行的函数
 * @param delay 延迟时间（毫秒）
 * @returns 防抖后的函数
 */
export function debounceImmediate<T extends (...args: any[]) => any>(
  func: T,
  delay: number
): (...args: Parameters<T>) => void {
  let timeoutId: ReturnType<typeof setTimeout> | null = null
  
  return function(...args: Parameters<T>) {
    const shouldCallNow = timeoutId === null
    
    if (timeoutId) {
      clearTimeout(timeoutId)
    }
    
    timeoutId = setTimeout(() => {
      timeoutId = null
    }, delay)
    
    if (shouldCallNow) {
      func(...args)
    }
  }
}

