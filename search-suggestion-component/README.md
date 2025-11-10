# 🔍 智能搜索建议组件

> 一个支持中文、拼音搜索的React搜索建议组件，包含1000条测试数据，80%已实现，20%留给你练手！

## ✨ 特性

- 🚀 **多种匹配模式**：前缀匹配、包含匹配、模糊匹配、拼音匹配
- 🎯 **智能排序**：结合热门度和匹配度的综合排序算法
- 🏷️ **分类筛选**：支持按分类过滤搜索结果
- ⌨️ **键盘导航**：支持上下键选择、Enter确认、ESC关闭
- 💡 **关键词高亮**：匹配的关键词自动高亮显示
- 🎨 **防抖优化**：300ms防抖，避免频繁搜索
- 📦 **LRU缓存**：缓存搜索结果，提升性能
- 🇨🇳 **中文友好**：完整的中文和拼音搜索支持

## 📊 数据规模

- **1000条测试数据**
- **10个分类**：电子产品、美食餐饮、旅游景点、在线学习、影视娱乐、网购商品、健康医疗、运动健身、汽车服务、金融理财、生活服务

## 🎯 练手任务

| 任务 | 难度 | 预计时间 | 文件位置 |
|------|------|---------|---------|
| 模糊匹配算法 | ⭐⭐⭐ | 1-2小时 | `src/core/algorithms/matcher.ts` |
| 分类筛选功能 | ⭐⭐ | 30分钟 | `src/core/SearchEngine.ts` |
| LRU缓存机制 | ⭐⭐⭐ | 1小时 | `src/utils/cache.ts` |
| 键盘导航 | ⭐⭐ | 30分钟 | `src/components/SearchSuggestion/SearchInput.tsx` |
| 关键词高亮 | ⭐⭐ | 30分钟 | `src/components/SearchSuggestion/SuggestionItem.tsx` |

**总计：3-4小时完成所有练手任务**

## 🚀 快速开始

### 1. 安装依赖

```bash
cd search-suggestion-component
npm install
```

### 2. 启动开发服务器

```bash
npm run dev
```

浏览器自动打开 `http://localhost:3000`

### 3. 开始练手

查看 `docs/PRACTICE_GUIDE.md` 获取详细的练手指导。

## 📁 项目结构

```
search-suggestion-component/
├── src/
│   ├── types/              # TypeScript类型定义
│   ├── data/               # 1000条测试数据
│   ├── utils/              # 工具函数（拼音、缓存）
│   ├── core/               # 核心搜索引擎
│   │   ├── algorithms/     # 匹配和评分算法
│   │   └── SearchEngine.ts # 搜索引擎类
│   ├── components/         # React组件
│   │   └── SearchSuggestion/
│   ├── hooks/              # 自定义Hooks
│   └── examples/           # 示例页面
├── docs/
│   ├── PRACTICE_GUIDE.md   # 📖 练手指导文档
│   └── API.md              # API文档
└── package.json
```

## 🎨 使用示例

```tsx
import { SearchSuggestion } from './components/SearchSuggestion'
import { mockSuggestions } from './data/suggestions'

function App() {
  const handleSelect = (item) => {
    console.log('选中：', item.text)
  }

  return (
    <SearchSuggestion
      data={mockSuggestions}
      onSelect={handleSelect}
      placeholder="搜索商品、地点、课程..."
    />
  )
}
```

## 📖 文档

- [练手指导文档](docs/PRACTICE_GUIDE.md) - 详细的任务说明和实现提示
- [API文档](docs/API.md) - 组件API和配置说明

## 🔧 技术栈

- **React 18** - UI框架
- **TypeScript** - 类型安全
- **Vite** - 构建工具
- **pinyin-pro** - 拼音转换
- **Vitest** - 单元测试

## 📝 脚本命令

```bash
# 开发
npm run dev

# 构建
npm run build

# 测试
npm run test

# 代码检查
npm run lint
```

## 🎓 学习目标

通过完成这个项目，你将学会：

1. **算法实现**：模糊匹配、编辑距离、LRU缓存
2. **React开发**：Hooks、事件处理、组件设计
3. **TypeScript**：泛型、接口、类型推导
4. **性能优化**：防抖、缓存、虚拟滚动
5. **中文处理**：拼音转换、中文分词

## 🤝 贡献指南

欢迎提交Issue和Pull Request！

## 📄 License

MIT License

---

**开始你的练手之旅吧！** 🚀

查看 [练手指导文档](docs/PRACTICE_GUIDE.md) 获取详细说明。

