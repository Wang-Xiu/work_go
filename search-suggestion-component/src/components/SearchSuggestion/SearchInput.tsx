import React, { useRef, useEffect } from 'react'

interface SearchInputProps {
  value: string
  onChange: (value: string) => void
  onFocus: () => void
  onBlur: () => void
  placeholder?: string
  selectedIndex: number
  resultCount: number
  onSelectChange: (index: number) => void
  onConfirm: () => void
  onClose: () => void
}

/**
 * 搜索输入框组件
 * 支持键盘导航
 */
export const SearchInput: React.FC<SearchInputProps> = ({
  value,
  onChange,
  onFocus,
  onBlur,
  placeholder = '请输入搜索关键词...',
  selectedIndex,
  resultCount,
  onSelectChange,
  onConfirm,
  onClose,
}) => {
  const inputRef = useRef<HTMLInputElement>(null)
  
  // 自动聚焦
  useEffect(() => {
    inputRef.current?.focus()
  }, [])
  
  // ============================================================
  // 🎯 【练手任务4】实现键盘导航功能
  // ============================================================
  
  /**
   * 处理键盘事件
   * 支持的按键：
   * - ⬆️ ArrowUp：选择上一项
   * - ⬇️ ArrowDown：选择下一项
   * - ↩️ Enter：确认选中
   * - ⎋ Escape：关闭建议列表
   */
  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    // ============================================================
    // 📝 实现思路：
    // ============================================================
    // 1. 使用switch语句判断按键类型（e.key）
    // 2. 不同按键执行不同操作
    // 3. 调用父组件传递的回调函数
    // 4. 注意边界情况（不能超出范围）
    // 5. 使用e.preventDefault()阻止默认行为
    //
    // ============================================================
    // 📌 知识点：selectedIndex的含义
    // ============================================================
    // selectedIndex是当前选中项的索引：
    //   -1：没有选中任何项
    //    0：选中第1项
    //    1：选中第2项
    //    ...
    //    resultCount-1：选中最后一项
    //
    // 例如：有5个结果
    //   selectedIndex可以是：-1, 0, 1, 2, 3, 4
    //
    // ============================================================
    // 📌 第1步：开始switch语句
    // ============================================================
    
    // 👉 你的代码：创建switch语句
    //
    // 语法：
    //   switch (e.key) {
    //     case 'ArrowUp':
    //       // 处理上箭头
    //       break
    //     case 'ArrowDown':
    //       // 处理下箭头
    //       break
    //     case 'Enter':
    //       // 处理回车
    //       break
    //     case 'Escape':
    //       // 处理ESC
    //       break
    //   }
    //
    // 写在这里：
    switch (e.key) {
      // ============================================================
      // 📌 第2步：处理ArrowUp（上箭头）
      // ============================================================
      case 'ArrowUp':
        // 👉 你的代码：实现上箭头逻辑
        //
        // 逻辑说明：
        //   - 按上箭头，selectedIndex应该减1
        //   - 但不能小于-1（-1表示不选中任何项）
        //
        // 步骤1：阻止默认行为（防止光标移动）
        //   e.preventDefault()
        //
        // 步骤2：计算新的index
        //   如果selectedIndex > -1，减1
        //   否则保持-1
        //
        // 步骤3：调用回调更新
        //   onSelectChange(newIndex)
        //
        // 代码示例：
        //   e.preventDefault()
        //   if (selectedIndex > -1) {
        //     onSelectChange(selectedIndex - 1)
        //   }
        //
        // 写在这里：
        e.preventDefault()
        if (selectedIndex > -1) {
          onSelectChange(selectedIndex - 1)
        }
        break
      
      // ============================================================
      // 📌 第3步：处理ArrowDown（下箭头）
      // ============================================================
      case 'ArrowDown':
        // 👉 你的代码：实现下箭头逻辑
        //
        // 逻辑说明：
        //   - 按下箭头，selectedIndex应该加1
        //   - 但不能超过resultCount-1（最后一项）
        //
        // 步骤1：阻止默认行为
        //   e.preventDefault()
        //
        // 步骤2：计算新的index
        //   如果selectedIndex < resultCount - 1，加1
        //   否则保持不变
        //
        // 步骤3：调用回调更新
        //   onSelectChange(newIndex)
        //
        // 代码示例：
        //   e.preventDefault()
        //   if (selectedIndex < resultCount - 1) {
        //     onSelectChange(selectedIndex + 1)
        //   }
        //
        // 写在这里：
        e.preventDefault()
        if (selectedIndex < resultCount - 1) {
          onSelectChange(selectedIndex + 1)
        }
        break
      
      // ============================================================
      // 📌 第4步：处理Enter（回车键）
      // ============================================================
      case 'Enter':
        // 👉 你的代码：实现回车逻辑
        //
        // 逻辑说明：
        //   - 如果有选中项（selectedIndex >= 0），确认选择
        //   - 如果没有选中项，可能是搜索
        //
        // 步骤1：阻止默认行为（防止表单提交）
        //   e.preventDefault()
        //
        // 步骤2：判断是否有选中项
        //   if (selectedIndex >= 0)
        //
        // 步骤3：如果有选中项，调用确认回调
        //   onConfirm()
        //
        // 代码示例：
        //   e.preventDefault()
        //   if (selectedIndex >= 0) {
        //     onConfirm()
        //   }
        //
        // 写在这里：
        e.preventDefault()
        if (selectedIndex >= 0) {
          onConfirm()
        }
        break
      
      // ============================================================
      // 📌 第5步：处理Escape（ESC键）
      // ============================================================
      case 'Escape':
        // 👉 你的代码：实现ESC逻辑
        //
        // 逻辑说明：
        //   - 按ESC关闭建议列表
        //   - 调用onClose回调
        //
        // 步骤1：阻止默认行为
        //   e.preventDefault()
        //
        // 步骤2：调用关闭回调
        //   onClose()
        //
        // 代码示例：
        //   e.preventDefault()
        //   onClose()
        //
        // 写在这里：
        e.preventDefault()
        onClose()
        break
    }
    
    // ============================================================
    // ✅ 完成！测试你的实现：
    // ============================================================
    //
    // 测试步骤：
    // 1. 启动应用：npm run dev
    // 2. 输入搜索关键词，显示建议列表
    // 3. 按⬇️键，第一项应该高亮
    // 4. 继续按⬇️键，高亮移动到下一项
    // 5. 按⬆️键，高亮移动到上一项
    // 6. 选中某一项，按Enter，应该触发选择
    // 7. 按ESC键，建议列表应该关闭
    // 8. 在第一项时按⬆️，应该取消选中（index变为-1）
    // 9. 在最后一项时按⬇️，应该保持在最后一项
    //
    // ============================================================
    // 💡 进阶优化（可选）：
    // ============================================================
    //
    // 1. 循环导航：
    //    在最后一项按⬇️，跳到第一项
    //    在第一项按⬆️，跳到最后一项
    //    提示：使用模运算 (index + 1) % resultCount
    //
    // 2. 支持更多快捷键：
    //    - Home：跳到第一项
    //    - End：跳到最后一项
    //    - PageUp/PageDown：快速翻页
    //
    // 3. 支持鼠标和键盘混合操作：
    //    鼠标hover时更新selectedIndex
    //
    // 4. 添加平滑滚动：
    //    当选中项不在可视区域时，自动滚动到可视区域
    //    提示：使用 element.scrollIntoView({ behavior: 'smooth' })
    //
    // 5. 支持Ctrl+K快捷键：
    //    像VSCode一样，用Ctrl+K唤起搜索
    //    提示：监听全局keydown事件
    //
    // ============================================================
  }
  
  /**
   * 处理输入变化
   */
  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    onChange(e.target.value)
  }
  
  /**
   * 处理清空按钮
   */
  const handleClear = () => {
    onChange('')
    inputRef.current?.focus()
  }
  
  return (
    <div className="search-input-wrapper">
      <div className="search-input-container">
        <span className="search-icon">🔍</span>
        
        <input
          ref={inputRef}
          type="text"
          className="search-input"
          value={value}
          onChange={handleChange}
          onKeyDown={handleKeyDown}
          onFocus={onFocus}
          onBlur={onBlur}
          placeholder={placeholder}
          autoComplete="off"
          spellCheck="false"
        />
        
        {value && (
          <button
            className="clear-button"
            onClick={handleClear}
            title="清空"
          >
            ✕
          </button>
        )}
      </div>
      
      {resultCount > 0 && (
        <div className="search-hint">
          <span className="hint-text">
            {selectedIndex >= 0 
              ? `${selectedIndex + 1} / ${resultCount}` 
              : `共${resultCount}个结果`}
          </span>
          <span className="hint-keys">
            ⬆️⬇️ 导航 · ↩️ 选择 · ⎋ 关闭
          </span>
        </div>
      )}
    </div>
  )
}

