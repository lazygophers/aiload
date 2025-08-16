# 工具使用指南

本文档为 roo 系统提供简洁实用的工具选择和应用指南。

---

## 工具快速参考

### 核心工具分类

#### 🔍 探索与分析

- **[`list_files`](tool-usage.md:12)** - 浏览项目结构，支持递归
- **[`read_file`](tool-usage.md:13)** - 读取文件内容（最多 5 个）
- **[`list_code_definition_names`](tool-usage.md:14)** - 获取代码定义概览
- **[`search_files`](tool-usage.md:15)** - 正则搜索文件内容

#### ✏️ 文件编辑

- **[`apply_diff`](tool-usage.md:18)** - 精确修改代码段（优先使用）
- **[`insert_content`](tool-usage.md:19)** - 在指定行插入内容
- **[`search_and_replace`](tool-usage.md:20)** - 批量查找替换
- **[`write_to_file`](tool-usage.md:21)** - 创建新文件或完全重写

#### 🤝 交互决策

- **[`ask_followup_question`](tool-usage.md:24)** - 请求用户决策（提供选项）
- **[`update_todo_list`](tool-usage.md:25)** - 管理任务清单
- **[`new_task`](tool-usage.md:26)** - 委派给专业模式
- **[`sequentialthinking`](tool-usage.md:27)** - 复杂问题分析（MCP）

#### 🛠️ 系统操作

- **[`execute_command`](tool-usage.md:30)** - 执行命令行操作
- **[`use_mcp_tool`](tool-usage.md:31)** - 调用 MCP 服务工具
- **[`attempt_completion`](tool-usage.md:32)** - 完成任务交付

### 工具选择决策树

```
任务开始
├─ 需要了解项目？→ list_files + read_file + list_code_definition_names
├─ 需要用户决策？→ ask_followup_question（L1级决策）
├─ 任务复杂度高？→ sequentialthinking → 任务分解
├─ 修改文件？
│  ├─ 新建 → write_to_file
│  ├─ 精确修改 → apply_diff
│  ├─ 插入内容 → insert_content
│  └─ 批量替换 → search_and_replace
├─ 需要专业能力？→ new_task（委派专业模式）
└─ 完成任务 → update_todo_list + attempt_completion
```

---

## 场景化工具组合

### 常用工作流模式

#### 🚀 项目初始化

```yaml
流程: list_files → read_file(配置文件) → list_code_definition_names → update_todo_list
用途: 快速了解项目结构和依赖
```

#### 🐛 Bug 修复

```yaml
流程: read_file → search_files(影响范围) → apply_diff → execute_command(测试)
用途: 定位问题、修复代码、验证结果
```

#### 🏗️ 功能开发

```yaml
流程: sequentialthinking(设计) → write_to_file(新文件) → insert_content(添加功能) → execute_command(构建)
用途: 从设计到实现的完整开发流程
```

#### ♻️ 代码重构

```yaml
流程: list_code_definition_names → search_files → apply_diff(批量) → execute_command(测试)
用途: 系统化改进代码结构
```

#### 📝 文档更新

```yaml
流程: list_files → new_task(doc-writer) → write_to_file/apply_diff
用途: 委派给专业文档模式处理
```

### 工具组合速查

| 目标         | 推荐组合                                       | 关键点         |
| ------------ | ---------------------------------------------- | -------------- |
| **探索理解** | `list_files` → `read_file` → `search_files`    | 从整体到细节   |
| **精确修改** | `read_file` → `apply_diff` → `execute_command` | 先读后改再验证 |
| **批量操作** | `search_files` → `search_and_replace`          | 先搜索确认范围 |
| **创建功能** | `write_to_file` → `insert_content`             | 先框架后细节   |
| **复杂决策** | `sequentialthinking` → `ask_followup_question` | 分析后确认     |

---

## MCP 服务使用

### Sequential Thinking

**用途**：复杂问题分析、方案设计、任务分解  
**触发条件**：任务复杂度高、需要多步骤规划、问题原因不明

```xml
<use_mcp_tool>
  <server_name>sequentialthinking</server_name>
  <tool_name>sequentialthinking</tool_name>
  <arguments>{"thought": "分析步骤...", "nextThoughtNeeded": true, "thoughtNumber": 1, "totalThoughts": 5}</arguments>
</use_mcp_tool>
```

### Context7 & GitHub

**用途**：查询技术文档、代码示例、GitHub 操作  
**使用时机**：需要最新文档、外部代码参考、仓库管理

---

## 最佳实践

### ✅ 推荐做法

1. **分层决策**：L1 级决策必须请求用户确认
2. **批量优化**：合并多个修改到单个 `apply_diff`
3. **验证先行**：修改前用 `list_files` 确认文件存在
4. **及时更新**：每个关键节点更新 `todo_list`
5. **工具优先级**：
   - 文件修改：`apply_diff` > `search_and_replace` > `write_to_file`
   - 文件探索：内建工具 > `execute_command`

### ❌ 避免反模式

| 反模式                               | 问题               | 正确做法                       |
| ------------------------------------ | ------------------ | ------------------------------ |
| 用 `write_to_file` 做小修改          | 效率低、风险高     | 使用 `apply_diff`              |
| 用 `execute_command` 执行 `ls`/`cat` | 平台依赖、解析困难 | 使用 `list_files`/`read_file`  |
| 无规划连续调用工具                   | 逻辑混乱、效率低   | 先用 `sequentialthinking` 规划 |
| 开放式 `ask_followup_question`       | 负担转移给用户     | 提供具体选项供选择             |

### 异常处理

| 异常       | 处理工具                               | 策略             |
| ---------- | -------------------------------------- | ---------------- |
| 文件不存在 | `list_files` → `ask_followup_question` | 确认路径或创建   |
| 权限不足   | `ask_followup_question`                | 请求权限或换方案 |
| 信息不足   | `ask_followup_question`                | 收集必要参数     |
| 复杂度过高 | `sequentialthinking` → `new_task`      | 分解或委派       |