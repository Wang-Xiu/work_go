import { useState, useEffect } from 'react'

/**
 * 防抖Hook
 * 用于延迟更新值，避免频繁触发
 * 
 * @param value 要防抖的值
 * @param delay 延迟时间（毫秒）
 * @returns 防抖后的值
 * 
 * @example
 * const [keyword, setKeyword] = useState('')
 * const debouncedKeyword = useDebounce(keyword, 300)
 * 
 * useEffect(() => {
 *   // 只有在用户停止输入300ms后才会触发搜索
 *   search(debouncedKeyword)
 * }, [debouncedKeyword])
 */
export function useDebounce<T>(value: T, delay: number): T {
  const [debouncedValue, setDebouncedValue] = useState<T>(value)

  useEffect(() => {
    // 设置延迟更新
    const handler = setTimeout(() => {
      setDebouncedValue(value)
    }, delay)

    // 清理函数：组件卸载或value/delay变化时清除定时器
    return () => {
      clearTimeout(handler)
    }
  }, [value, delay])

  return debouncedValue
}

