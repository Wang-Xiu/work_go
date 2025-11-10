import React, { useState } from 'react'
import { SearchSuggestion } from '../src/components/SearchSuggestion'
import { suggestions } from '../src/data/suggestions'
import { SuggestionItem } from '../src/types'
import './App.css'

/**
 * 示例应用
 * 展示搜索建议组件的使用方法
 */
function App() {
  const [selectedItem, setSelectedItem] = useState<SuggestionItem | null>(null)

  const handleSelect = (item: SuggestionItem) => {
    console.log('选中项:', item)
    setSelectedItem(item)
  }

  const handleChange = (keyword: string) => {
    console.log('搜索关键词:', keyword)
  }

  return (
    <div className="app">
      <header className="app-header">
        <h1>🔍 搜索建议组件</h1>
        <p className="subtitle">支持中文、拼音、首字母、模糊匹配</p>
      </header>

      <main className="app-main">
        <div className="search-container">
          <SearchSuggestion
            items={suggestions}
            onSelect={handleSelect}
            onChange={handleChange}
            placeholder="试试输入：iPhone、美食、pg、hd..."
            showCategoryFilter={true}
            config={{
              topN: 10,
              matchWeight: 0.6,
              hotWeight: 0.4,
              enablePinyin: true,
              enableFuzzy: true,
              minMatchScore: 0,
              debounceDelay: 300,
            }}
          />
        </div>

        {/* 选中结果显示 */}
        {selectedItem && (
          <div className="result-display">
            <h3>已选中：</h3>
            <div className="result-card">
              <div className="result-icon">{selectedItem.icon}</div>
              <div className="result-content">
                <div className="result-text">{selectedItem.text}</div>
                {selectedItem.description && (
                  <div className="result-description">{selectedItem.description}</div>
                )}
                <div className="result-meta">
                  <span className="result-category">{selectedItem.category}</span>
                  <span className="result-hot">🔥 {selectedItem.hotScore}</span>
                </div>
              </div>
            </div>
          </div>
        )}

        {/* 使用说明 */}
        <div className="instructions">
          <h3>📚 使用说明</h3>
          <div className="instruction-grid">
            <div className="instruction-item">
              <div className="instruction-title">🔤 文本搜索</div>
              <div className="instruction-content">直接输入中文或英文</div>
              <div className="instruction-example">例如：iPhone、华为、小米</div>
            </div>

            <div className="instruction-item">
              <div className="instruction-title">🅰️ 拼音搜索</div>
              <div className="instruction-content">输入完整拼音</div>
              <div className="instruction-example">例如：pingguo（苹果）</div>
            </div>

            <div className="instruction-item">
              <div className="instruction-title">🔡 首字母搜索</div>
              <div className="instruction-content">输入拼音首字母</div>
              <div className="instruction-example">例如：pg（苹果）、hd（火锅）</div>
            </div>

            <div className="instruction-item">
              <div className="instruction-title">🎯 模糊匹配</div>
              <div className="instruction-content">字符可以不连续</div>
              <div className="instruction-example">例如：iph（iPhone）</div>
            </div>

            <div className="instruction-item">
              <div className="instruction-title">⌨️ 键盘导航</div>
              <div className="instruction-content">⬆️⬇️ 选择 · ↩️ 确认 · ⎋ 关闭</div>
              <div className="instruction-example">支持纯键盘操作</div>
            </div>

            <div className="instruction-item">
              <div className="instruction-title">🏷️ 分类筛选</div>
              <div className="instruction-content">点击分类按钮过滤</div>
              <div className="instruction-example">电子产品、美食、旅游等</div>
            </div>
          </div>
        </div>

        {/* 技术特性 */}
        <div className="features">
          <h3>✨ 技术特性</h3>
          <ul className="feature-list">
            <li>🎯 <strong>智能排序</strong>：综合匹配度和热度评分</li>
            <li>⚡ <strong>高性能</strong>：LRU缓存 + 防抖优化</li>
            <li>🔤 <strong>中文友好</strong>：支持拼音和首字母搜索</li>
            <li>🎨 <strong>现代UI</strong>：流畅动画 + 响应式设计</li>
            <li>♿ <strong>可访问性</strong>：完整键盘导航支持</li>
            <li>🌙 <strong>深色模式</strong>：自动适配系统主题</li>
            <li>📱 <strong>移动端优化</strong>：触摸友好 + 自适应布局</li>
            <li>🔧 <strong>高度可配置</strong>：权重、TOP N、延迟等</li>
          </ul>
        </div>

        {/* 数据统计 */}
        <div className="statistics">
          <h3>📊 数据统计</h3>
          <div className="stat-grid">
            <div className="stat-item">
              <div className="stat-value">{suggestions.length}</div>
              <div className="stat-label">总数据量</div>
            </div>
            <div className="stat-item">
              <div className="stat-value">
                {new Set(suggestions.map(s => s.category)).size}
              </div>
              <div className="stat-label">分类数量</div>
            </div>
            <div className="stat-item">
              <div className="stat-value">5</div>
              <div className="stat-label">匹配算法</div>
            </div>
            <div className="stat-item">
              <div className="stat-value">300ms</div>
              <div className="stat-label">防抖延迟</div>
            </div>
          </div>
        </div>
      </main>

      <footer className="app-footer">
        <p>🎓 这是一个学习项目，包含5个练手任务</p>
        <p className="footer-hint">查看 <code>docs/PRACTICE_GUIDE.md</code> 开始学习</p>
      </footer>
    </div>
  )
}

export default App

