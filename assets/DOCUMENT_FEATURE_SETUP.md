# 文档详情页面功能设置说明

## 已完成的功能

1. ✅ **MarkdownRenderer.tsx 组件**
   - 使用 react-markdown 渲染 Markdown 内容
   - 支持代码高亮（通过动态导入 react-syntax-highlighter）
   - 自动生成目录（TOC）
   - 标题锚点功能
   - 表格、链接等元素的样式优化

2. ✅ **Document.tsx 页面组件**
   - 根据路由参数显示对应文档内容
   - 响应式布局（左侧内容，右侧目录）
   - 返回文档列表功能
   - 滚动时自动更新活动锚点
   - 支持URL哈希直接跳转到指定标题

3. ✅ **集成功能**
   - 与现有设计风格保持一致
   - 支持国际化
   - 完整的错误处理

## 需要手动完成的步骤

### 1. 安装依赖

由于系统中没有包管理器，需要手动安装以下依赖：

```bash
# 在 assets 目录下执行
npm install react-syntax-highlighter @types/react-syntax-highlighter
# 或使用
pnpm add react-syntax-highlighter @types/react-syntax-highlighter
# 或使用
yarn add react-syntax-highlighter @types/react-syntax-highlighter
```

### 2. 更新 docs.tsx 中的导航逻辑

将 `handleDocumentClick` 函数中的注释代码取消注释：

```typescript
// 将这行
// navigate(`/docs/${docId}`);

// 改为
navigate(`/docs/${docId}`);
```

### 3. 更新 App.tsx 添加路由

在 App.tsx 中添加 Document 页面的路由：

```typescript
// 添加导入
import Document from './pages/Document';

// 在 Routes 组件中添加路由
<Route path="/docs/:docId" element={<Document />} />
```

## 功能特性

### 代码高亮
- 支持多种编程语言的语法高亮
- 自动检测代码块语言
- 深色/浅色主题自适应

### 目录导航
- 自动从 Markdown 内容中提取标题生成目录
- 点击目录项平滑滚动到对应位置
- 滚动时自动高亮当前阅读的章节

### 响应式设计
- 移动端自动隐藏目录侧边栏
- 内容区域自适应宽度
- 优化的小屏幕阅读体验

### 性能优化
- 动态导入语法高亮库，减少首包大小
- 错误降级处理，即使依赖加载失败也能正常显示
- 防抖滚动事件处理

## 使用说明

1. 访问文档列表页：`/docs`
2. 点击任意文档卡片进入详情页：`/docs/:docId`
3. 使用右侧目录快速导航
4. 点击标题前的 # 符号可以复制链接

## 注意事项

- 确保安装了所有必需的依赖
- 如果 react-syntax-highlighter 加载失败，代码块会以普通 `<pre>` 标签显示
- 文档内容存储在 `assets/src/pages/docs.ts` 文件中