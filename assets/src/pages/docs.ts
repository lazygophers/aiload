export const documents = [
  {
    name: 'README.md',
    content: `# 🚀 AI Load - 下一代 AI 接口透明代理服务

[![Go Report Card](https://goreportcard.com/badge/github.com/lazygophers/aiload)](https://goreportcard.com/report/github.com/lazygophers/aiload)
[![Build Status](https://img.shields.io/github/actions/workflow/status/lazygophers/aiload/go.yml?branch=main)](https://github.com/lazygophers/aiload/actions)
[![License](https://img.shields.io/github/license/lazygophers/aiload)](https://github.com/lazygophers/aiload/blob/main/LICENSE)
[![Release](https://img.shields.io/github/v/release/lazygophers/aiload)](https://github.com/lazygophers/aiload/releases)

**AI Load** 是一款专为需要集成多种 AI 服务的企业和开发者设计的高性能、高可用的 AI 接口透明代理服务。它以**零侵入、高效率**为核心，帮助您轻松驾驭复杂的 AI 服务矩阵。

---

## ✨ 功能特性

- **🔄 透明代理**: **完全保留**原生 API 格式，无缝对接 OpenAI、Google Gemini、Anthropic Claude、**Siliconflow** 以及本地运行的 **Ollama** 等多种服务，无需修改现有代码。
- **🛡️ 平台化密钥管理**: 创新的**平台 -> 密钥 -> 模型**三层管理体系，支持多平台、多密钥的灵活配置，实现精细化的访问控制与状态管理。
- **⚖️ 负载均衡**: 支持多上游端点的**加权负载均衡**，智能分配请求，显著提升服务可用性和稳定性。
- **🛡️ 智能故障处理**: 自动化的**密钥黑名单**和恢复机制，主动规避故障节点，确保业务连续性。
- **⚙️ 数据库驱动**: 核心业务配置（平台、密钥、模型等）存储于数据库，通过管理后台进行维护。
- **🖥️ 现代化管理**: 基于 **React + Ant Design** 的现代化 Web 管理界面，所有操作直观易用。
- **⚡ 高性能设计**: 采用**零拷贝流式传输**、连接池复用和原子操作等技术，最大化处理性能。

## 🤖 支持的 AI 服务

AI Load 作为透明代理服务，完整保留了各大 AI 服务商的原生 API 格式，包括：

- **OpenAI 格式**: 官方 OpenAI API、Azure OpenAI 及其他兼容服务。
- **Google Gemini 格式**: Gemini Pro、Gemini Pro Vision 等原生 API。
- **Anthropic Claude 格式**: Claude 系列模型的高质量对话与文本生成 API。
- **Siliconflow 格式**: 兼容 OpenAI 格式的 Siliconflow 云端服务。
- **Ollama (本地)**: 支持在本地环境中运行的 Ollama 模型，保障数据私密性。

## 🛠️ 技术栈

- **后端**: Golang
- **前端**: React + Ant Design

## 🚀 快速开始

请参阅我们的 **[快速上手指南](./docs/getting-started.md)** 来快速部署和使用 AI Load。

## 📚 文档

更详细的文档请访问 [docs](./docs) 目录。

- **核心指南**
  - [🚀 项目介绍](./docs/introduction.md): 了解 AI Load 的核心理念与价值。
  - [🛠️ 部署指南](./docs/deployment.md): 一步步完成服务的部署与启动。
  - [⚙️ 配置说明](./docs/configuration.md): 掌握所有配置项，释放全部潜能。
- **开发与贡献**
  - [📚 API 参考](./docs/api.md): 查阅完整的 API 接口文档。
  - [🏛️ 架构设计](./docs/architecture.md): 深入理解系统的设计哲学。
  - [🤝 贡献指南](./docs/contributing.md): 加入我们，共同建设。

## 🤝 贡献

我们欢迎任何形式的贡献！请阅读 **[贡献指南](./docs/contributing.md)** 和 **[行为准则](./CODE_OF_CONDUCT.md)** 来了解如何参与项目。

## 📄 许可证

本项目基于 [MIT](LICENSE) 许可证。
`
  },
  {
    name: 'api_spec.md',
    content: \`# 📖 API 文档 (aiload.proto)

本文档根据 \`aiload.proto\` 文件自动生成，详细描述了其中定义的所有服务接口、数据模型和枚举类型。

---

## 🔮 接口 (RPC) 文档

本部分详细说明了 \`aiload\` 服务提供的所有 RPC 方法。

### \`POST /AddUserAdmin\`

-   **描述**: 添加用户 (管理员)
-   **权限**: \`admin\`
-   **请求参数**: [\`AddUserAdminReq\`](#AddUserAdminReq)
-   **响应参数**: [\`AddUserAdminRsp\`](#AddUserAdminRsp)

### \`POST /GetUser\`

-   **描述**: 获取当前用户信息
-   **权限**: \`user\` (默认)
-   **请求参数**: [\`GetUserReq\`](#GetUserReq)
-   **响应参数**: [\`GetUserRsp\`](#GetUserRsp)

### \`POST /GetUserAdmin\`

-   **描述**: 获取指定用户信息 (管理员)
-   **权限**: \`admin\`
-   **请求参数**: [\`GetUserAdminReq\`](#GetUserAdminReq)
-   **响应参数**: [\`GetUserAdminRsp\`](#GetUserAdminRsp)

### \`POST /ListUserAdmin\`

-   **描述**: 列出所有用户 (管理员)
-   **权限**: \`admin\`
-   **请求参数**: [\`ListUserAdminReq\`](#ListUserAdminReq)
-   **响应参数**: [\`ListUserAdminRsp\`](#ListUserAdminRsp)

### \`POST /SetUser\`

-   **描述**: 更新当前用户信息
-   **权限**: \`user\` (默认)
-   **请求参数**: [\`SetUserReq\`](#SetUserReq)
-   **响应参数**: [\`SetUserRsp\`](#SetUserRsp)

### \`POST /SetUserAdmin\`

-   **描述**: 更新指定用户信息 (管理员)
-   **权限**: \`admin\`
-   **请求参数**: [\`SetUserAdminReq\`](#SetUserAdminReq)
-   **响应参数**: [\`SetUserAdminRsp\`](#SetUserAdminRsp)

### \`POST /DelUserAdmin\`

-   **描述**: 删除指定用户 (管理员)
-   **权限**: \`admin\`
-   **请求参数**: [\`DelUserAdminReq\`](#DelUserAdminReq)
-   **响应参数**: [\`DelUserAdminRsp\`](#DelUserAdminRsp)

### \`POST /Login\`

-   **描述**: 用户登录
-   **权限**: \`public\`
-   **请求参数**: [\`LoginReq\`](#LoginReq)
-   **响应参数**: [\`LoginRsp\`](#LoginRsp)

---

## 🧱 数据模型 (Message) 文档

本部分清晰地展示了每个 \`message\` 的数据结构。

### \`ModelChannel\`

| 字段         | 类型                    | 描述                                           |
| :----------- | :---------------------- | :--------------------------------------------- |
| \`id\`         | \`uint64\`                | 唯一标识符                                     |
| \`created_at\` | \`int64\`                 | 创建时间                                       |
| \`updated_at\` | \`int64\`                 | 更新时间                                       |
| \`deleted_at\` | \`int64\`                 | 删除时间，\`@gorm: index: idx_channel,unique\`   |
| \`name\`       | \`string\`                | 渠道名称                                       |
| \`token\`      | \`string\`                | 渠道 Token, \`@gorm: index: idx_channel,unique\` |
| \`platform\`   | [\`Platform\`](#Platform) | 平台类型, \`@gorm: index: idx_channel,unique\`   |

### \`ModelChannelAccess\`

| 字段         | 类型     | 描述                                                |
| :----------- | :------- | :-------------------------------------------------- |
| \`id\`         | \`uint64\` | 唯一标识符                                          |
| \`created_at\` | \`int64\`  | 创建时间                                            |
| \`updated_at\` | \`int64\`  | 更新时间                                            |
| \`deleted_at\` | \`int64\`  | 删除时间, \`@gorm: index: idx_channel_access,unique\` |
| \`channel_id\` | \`uint64\` | 渠道 ID, \`@gorm: index: idx_channel_access,unique\`  |
| \`model\`      | \`string\` | 模型名称, \`@gorm: index: idx_channel_access,unique\` |

### \`ModelModelAlias\`

| 字段          | 类型     | 描述                                      |
| :------------ | :------- | :---------------------------------------- |
| \`id\`          | \`uint64\` | 唯一标识符                                |
| \`created_at\`  | \`int64\`  | 创建时间                                  |
| \`updated_at\`  | \`int64\`  | 更新时间                                  |
| \`deleted_at\`  | \`int64\`  | 删除时间, \`@gorm: index:idx_alias,unique\` |
| \`model\`       | \`string\` | 模型名称, \`@gorm: index:idx_alias,unique\` |
| \`model_alias\` | \`string\` | 模型别名, \`@gorm: index:idx_alias,unique\` |

### \`ModelUser\`

| 字段         | 类型                    | 描述                                     |
| :----------- | :---------------------- | :--------------------------------------- |
| \`id\`         | \`uint64\`                | 唯一标识符                               |
| \`created_at\` | \`int64\`                 | 创建时间                                 |
| \`updated_at\` | \`int64\`                 | 更新时间                                 |
| \`deleted_at\` | \`int64\`                 | 删除时间, \`@gorm: index:idx_user,unique\` |
| \`username\`   | \`string\`                | 用户名, \`@gorm: index:idx_user,unique\`   |
| \`password\`   | \`string\`                | 密码                                     |
| \`name\`       | \`string\`                | 用户昵称                                 |
| \`role\`       | [\`UserRole\`](#UserRole) | 用户角色                                 |

### \`ModelUserToken\`

| 字段         | 类型     | 描述                                        |
| :----------- | :------- | :------------------------------------------ |
| \`id\`         | \`uint64\` | 唯一标识符                                  |
| \`created_at\` | \`int64\`  | 创建时间                                    |
| \`updated_at\` | \`int64\`  | 更新时间                                    |
| \`deleted_at\` | \`int64\`  | 删除时间, \`@gorm: index:idx_token,unique\`   |
| \`user_id\`    | \`uint64\` | 用户 ID, \`@gorm: index:idx_token,unique\`    |
| \`token\`      | \`string\` | 用户 Token, \`@gorm: index:idx_token,unique\` |
| \`limit\`      | \`int64\`  | 限制                                        |

### \`ModelUserAccess\`

| 字段         | 类型     | 描述                                            |
| :----------- | :------- | :---------------------------------------------- |
| \`id\`         | \`uint64\` | 唯一标识符                                      |
| \`created_at\` | \`int64\`  | 创建时间                                        |
| \`updated_at\` | \`int64\`  | 更新时间                                        |
| \`deleted_at\` | \`int64\`  | 删除时间, \`@gorm: index:idx_user_access,unique\` |
| \`token\`      | \`uint64\` | Token ID, \`@gorm: index:idx_user_access,unique\` |
| \`model\`      | \`string\` | 模型名称, \`@gorm: index:idx_user_access,unique\` |

### \`AddUserAdminReq\`

| 字段   | 类型                      | 描述                            |
| :----- | :------------------------ | :------------------------------ |
| \`user\` | [\`ModelUser\`](#ModelUser) | 用户信息, \`@validate: required\` |

### \`AddUserAdminRsp\`

| 字段   | 类型                      | 描述               |
| :----- | :------------------------ | :----------------- |
| \`user\` | [\`ModelUser\`](#ModelUser) | 创建成功的用户信息 |

### \`GetUserReq\`

空请求。

### \`GetUserRsp\`

| 字段   | 类型                      | 描述         |
| :----- | :------------------------ | :----------- |
| \`user\` | [\`ModelUser\`](#ModelUser) | 当前用户信息 |

### \`GetUserAdminReq\`

| 字段 | 类型     | 描述                           |
| :--- | :------- | :----------------------------- |
| \`id\` | \`uint64\` | 用户 ID, \`@validate: required\` |

### \`GetUserAdminRsp\`

| 字段   | 类型                      | 描述             |
| :----- | :------------------------ | :--------------- |
| \`user\` | [\`ModelUser\`](#ModelUser) | 获取到的用户信息 |

### \`ListUserAdminReq\`

| 字段          | 类型                               | 描述                            |
| :------------ | :--------------------------------- | :------------------------------ |
| \`list_option\` | \`lazygophers.lrpc.core.ListOption\` | 列表选项, \`@validate: required\` |

\`ListUserAdminReq.ListOption\` 枚举:
| 名称 | 值 | 描述 |
| :--- | :-: | :--- |
| \`ListOptionNil\` | 0 | 默认值 |
| \`ListOptionName\` | 1 | 按名称模糊搜索 |
| \`ListOptionUsername\` | 2 | 按用户名模糊搜索 |

### \`ListUserAdminRsp\`

| 字段       | 类型                             | 描述     |
| :--------- | :------------------------------- | :------- |
| \`paginate\` | \`lazygophers.lrpc.core.Paginate\` | 分页信息 |
| \`list\`     | \`repeated ModelUser\`             | 用户列表 |

### \`SetUserReq\`

| 字段   | 类型                      | 描述                                      |
| :----- | :------------------------ | :---------------------------------------- |
| \`user\` | [\`ModelUser\`](#ModelUser) | 需要更新的用户信息, \`@validate: required\` |

### \`SetUserRsp\`

| 字段   | 类型                      | 描述             |
| :----- | :------------------------ | :--------------- |
| \`user\` | [\`ModelUser\`](#ModelUser) | 更新后的用户信息 |

### \`SetUserAdminReq\`

| 字段   | 类型                      | 描述                                      |
| :----- | :------------------------ | :---------------------------------------- |
| \`user\` | [\`ModelUser\`](#ModelUser) | 需要更新的用户信息, \`@validate: required\` |

### \`SetUserAdminRsp\`

| 字段   | 类型                      | 描述             |
| :----- | :------------------------ | :--------------- |
| \`user\` | [\`ModelUser\`](#ModelUser) | 更新后的用户信息 |

### \`DelUserAdminReq\`

| 字段 | 类型     | 描述                                   |
| :--- | :------- | :------------------------------------- |
| \`id\` | \`uint64\` | 要删除的用户 ID, \`@validate: required\` |

### \`DelUserAdminRsp\`

空响应。

### \`LoginReq\`

| 字段       | 类型     | 描述                          |
| :--------- | :------- | :---------------------------- |
| \`username\` | \`string\` | 用户名, \`@validate: required\` |
| \`password\` | \`string\` | 密码, \`@validate: required\`   |

### \`LoginRsp\`

| 字段         | 类型                      | 描述           |
| :----------- | :------------------------ | :------------- |
| \`user\`       | [\`ModelUser\`](#ModelUser) | 用户信息       |
| \`token\`      | \`string\`                  | 登录凭证 Token |
| \`expires_at\` | \`int64\`                   | Token 过期时间 |

---

## 🔢 枚举 (Enum) 文档

本部分将每个 \`enum\` 的值和描述以表格形式呈现。

### \`ErrCode\`

| 名称                        | 值    | 描述              |
| :-------------------------- | :---- | :---------------- |
| \`Success\`                   | 0     | 成功              |
| \`ModelAliasNotFound\`        | 10000 | 模型别名未找到    |
| \`ModelAliasDuplicateKey\`    | 10001 | 模型别名重复      |
| \`ChannelNotFound\`           | 10002 | 渠道未找到        |
| \`ChannelDuplicateKey\`       | 10003 | 渠道重复          |
| \`UserTokenNotFound\`         | 10004 | 用户 token 未找到 |
| \`UserTokenDuplicateKey\`     | 10005 | 用户 token 重复   |
| \`UserDuplicateKey\`          | 10006 | 用户重复          |
| \`UserNotFound\`              | 10007 | 用户未找到        |
| \`UserAccessNotFound\`        | 10008 | 用户访问未找到    |
| \`UserAccessDuplicateKey\`    | 10009 | 用户访问重复      |
| \`ChannelAccessNotFound\`     | 10010 | 渠道访问未找到    |
| \`ChannelAccessDuplicateKey\` | 10011 | 渠道访问重复      |

### \`UserRole\`

| 名称     | 值  | 描述     |
| :------- | :-: | :------- |
| \`Public\` |  0  | 公共用户 |
| \`User\`   |  1  | 普通用户 |
| \`Admin\`  |  2  | 管理员   |
| \`System\` |  3  | 系统     |

### \`Platform\`

| 名称               | 值  | 描述        |
| :----------------- | :-: | :---------- |
| \`Nil\`              |  0  | 未指定      |
| \`OpenAiCompatible\` |  1  | OpenAI 兼容 |
| \`Siliconflow\`      |  2  | Siliconflow |
| \`Gemini\`           |  3  | Gemini      |
| \`Ollama\`           |  4  | Ollama      |
| \`OpenAi\`           |  5  | OpenAI      |
\`
  },
  {
    name: 'architecture.md',
    content: \`# 🏛️ 系统架构设计文档

## 1. 引言

本文档旨在详细阐述 **AI Load** 透明代理服务的后台架构，特别是围绕**数据持久化**、**高性能缓存**及**核心服务流程**的设计决策。

我们的核心目标是构建一个**高效、可靠且可扩展**的系统。基于此，我们做出了以下关键技术选型：

- **🗄️ 数据库 (Database):** **SQLite**
  - **选择理由:** 轻量、嵌入式、零配置、事务性强。非常适合作为项目初、中期的核心数据存储，能够与应用程序一同部署，极大地简化了运维复杂性。
- **⚡ 缓存 (Cache):** **bbolt**
  - **选择理由:** 纯 Go 实现、嵌入式、高性能的 Key-Value 存储。它与 SQLite 的嵌入式特性完美契合，可用于显著提升高频读取操作的性能，降低数据库负载。

---

## 2. 核心服务架构

### 2.1. 架构图

\`\`\`mermaid
graph TD
    subgraph "用户请求"
        A[Client SDK / curl] --> B{AI Load Service};
    end

    subgraph "AI Load 核心服务 (Golang)"
        B --> C[1. 接收与解析请求];
        C --> D{2. 缓存查询};
        D -->(缓存命中) H[6. 上游请求分发];
        D -->(缓存未命中) E[3. 认证中间件];
        E --> F[4. 路由与模型匹配];
        F --> G[5. 负载均衡器];
        G --> H;
        H --> I[7. 响应处理与日志];
        I -->(写入缓存) D;
    end

    subgraph "数据与配置"
        J[config.yaml] -->(服务启动) B;
        K[(Cache: bbolt)] -->(读/写) D;
        L[(DB: SQLite)] -->(数据源) E;
        L -->(数据源) F;
        L -->(数据源) G;
    end

    B -->(日志) M[Console Output];

    subgraph "上游 AI 服务 (平台)"
        H --> N[OpenAI API];
        H --> O[Google Gemini API];
        H --> P[Anthropic Claude API];
        H --> Q[Ollama (Local)];
    end

    I --> A;
\`\`\`

### 2.2. 组件详解

- **1. 接收与解析请求 (Request Handler):** 服务入口，监听并解析 HTTP/HTTPS 请求。
- **2. 缓存查询 (Cache Lookup):** 在请求处理的早期阶段，优先查询 **bbolt** 缓存。
- **3. 认证中间件 (Auth Middleware):** 从 \`Authorization\` 头中提取密钥，在数据库中进行哈希校验。
- **4. 路由与模型匹配 (Router):** 根据请求的 \`model\`，从数据库查询并匹配到具体的上游服务平台和密钥组。
- **5. 负载均衡器 (Load Balancer):** 根据**加权轮询**策略，从密钥组中选择一个健康的实例。
- **6. 上游请求分发 (Upstream Dispatcher):** 通过**协议转换层 (Protocol Adapter)** 兼容多厂商 API，并向上游服务发起实际调用。
- **7. 响应处理与日志 (Response & Logging):** 将上游响应返回给客户端，并记录日志。成功处理的响应会**异步更新缓存**。

---

## 3. 数据持久化层 (SQLite)

### 3.1. 数据库文件

- **路径:** \`data/aiload.db\`
- **描述:** 系统所有持久化数据的“单一事实来源”。

### 3.2. 数据表模式 (Schema)

#### 🔑 a. \`ApiKeys\` (API 密钥表)

| 字段名               | 数据类型   | 约束/索引                   | 描述与设计 rationale                                                                              |
| -------------------- | ---------- | --------------------------- | ------------------------------------------------------------------------------------------------- |
| \`id\`                 | \`INTEGER\`  | \`PRIMARY KEY AUTOINCREMENT\` | 唯一密钥 ID                                                                                       |
| \`user_id\`            | \`INTEGER\`  | \`FOREIGN KEY(Users.id)\`     | 关联的用户 ID                                                                                     |
| **\`key_prefix\`**     | \`TEXT\`     | \`NOT NULL, INDEX\`           | **设计 rationale:** 密钥前缀 (例如 \`sk-abc...\`)，用于快速识别和日志追溯，**不用于安全校验**。     |
| **\`key_hash\`**       | \`TEXT\`     | \`UNIQUE NOT NULL\`           | **设计 rationale:** 使用 **SHA-256** 哈希后的完整密钥，用于安全校验，确保数据库中不存储明文密钥。 |
| \`status\`             | \`TEXT\`     | \`NOT NULL\`                  | 密钥状态 (\`active\`, \`inactive\`, \`revoked\`)                                                        |
| \`rate_limit_per_min\` | \`INTEGER\`  | \`NOT NULL, DEFAULT 60\`      | 每分钟允许的请求次数                                                                              |
| \`expires_at\`         | \`DATETIME\` | \`NULL\`                      | 密钥过期时间 (UTC)，\`NULL\` 表示永不过期                                                           |
| \`created_at\`         | \`DATETIME\` | \`NOT NULL\`                  | 密钥创建时间 (UTC)                                                                                |

#### 👤 b. \`Users\` (用户表)

| 字段名            | 数据类型   | 约束/索引                   | 描述                       |
| ----------------- | ---------- | --------------------------- | -------------------------- |
| \`id\`              | \`INTEGER\`  | \`PRIMARY KEY AUTOINCREMENT\` | 唯一用户 ID                |
| \`username\`        | \`TEXT\`     | \`UNIQUE NOT NULL\`           | 用户名，唯一               |
| \`email\`           | \`TEXT\`     | \`UNIQUE NOT NULL\`           | 电子邮箱，唯一             |
| \`hashed_password\` | \`TEXT\`     | \`NOT NULL\`                  | 加密后的用户密码           |
| \`created_at\`      | \`DATETIME\` | \`NOT NULL\`                  | 用户创建时间 (UTC)         |
| \`updated_at\`      | \`DATETIME\` | \`NOT NULL\`                  | 用户信息最后更新时间 (UTC) |

#### 🤖 c. \`Models\` (模型信息表)

| 字段名        | 数据类型  | 约束/索引                | 描述                       |
| ------------- | --------- | ------------------------ | -------------------------- |
| \`id\`          | \`TEXT\`    | \`PRIMARY KEY\`            | 模型 ID (例如 \`gpt-4\`)     |
| \`provider\`    | \`TEXT\`    | \`NOT NULL\`               | 模型提供商 (例如 \`OpenAI\`) |
| \`description\` | \`TEXT\`    | \`NULL\`                   | 模型的简短描述             |
| \`is_active\`   | \`BOOLEAN\` | \`NOT NULL, DEFAULT TRUE\` | 模型当前是否可用           |

#### 📊 d. \`UsageLogs\` (使用日志表)

| 字段名              | 数据类型   | 约束/索引                   | 描述                                          |
| ------------------- | ---------- | --------------------------- | --------------------------------------------- |
| \`id\`                | \`INTEGER\`  | \`PRIMARY KEY AUTOINCREMENT\` | 唯一日志 ID                                   |
| \`api_key_id\`        | \`INTEGER\`  | \`FOREIGN KEY(ApiKeys.id)\`   | 关联的 API 密钥 ID                            |
| \`endpoint\`          | \`TEXT\`     | \`NOT NULL, INDEX\`           | 请求的 API 端点 (例如 \`/v1/chat/completions\`) |
| \`http_status\`       | \`INTEGER\`  | \`NOT NULL, INDEX\`           | HTTP 响应状态码 (例如 200, 429)               |
| \`response_time_ms\`  | \`INTEGER\`  | \`NOT NULL\`                  | 从收到请求到发出响应的总耗时（毫秒）          |
| \`request_timestamp\` | \`DATETIME\` | \`NOT NULL\`                  | 请求到达服务器的时间 (UTC)                    |

---

## 4. 高性能缓存层 (bbolt)

为了最大限度地减少对数据库的直接访问并提升响应速度，我们引入了基于 **bbolt** 的嵌入式缓存层。

### 4.1. 缓存文件

- **路径:** \`data/cache.db\`
- **描述:** 用于存储高频访问数据的 Key-Value 数据库。

### 4.2. 缓存策略

#### a. 缓存内容与粒度

我们主要缓存两类数据：

1.  **API 密钥校验结果:**

    - **Key:** \`token:<SHA-256-hash-of-the-full-key>\`
    - **Value:** 序列化后的 \`ApiKey\` 对象（包含状态、速率限制等信息）。
    - **目的:** 避免每次请求都进行数据库哈希比对，这是最核心的性能优化点。

2.  **/v1/models 接口响应:**
    - **Key:** \`models:list\`
    - **Value:** 序列化后的模型列表 JSON 响应。
    - **目的:** \`/v1/models\` 是一个高频访问但内容不常变化的接口，对其进行缓存收益极高。

#### b. 缓存更新与失效

- **更新机制:** 采用**缓存旁路 (Cache-Aside)** 模式。
  1.  **读请求:**
      - 优先查询缓存。
      - 缓存未命中，则查询数据库。
      - 查询成功后，将结果**异步写入**缓存。
  2.  **写请求 (通过管理后台):**
      - 当管理员在后台更新了密钥或模型信息时，会直接**删除 (invalidate)** 缓存中对应的条目。
      - 下次相关请求进入时，会因缓存未命中而重新从数据库加载最新数据，并回填缓存。
- **过期策略:**
  - 为缓存条目设置 **TTL (Time-To-Live)**，例如 5 分钟。这确保了即使写时删除缓存失败，数据也能在短时间内自动刷新，保证最终一致性。
- **缓存穿透与雪崩:**
  - **穿透:** 对于查询不存在的密钥，我们也会缓存一个**空结果**并设置一个较短的 TTL，以防止恶意请求持续冲击数据库。
  * **雪崩:** 采用 TTL + 随机偏移的方式，避免大量缓存在同一时刻集体失效。
\`
  },
  {
    name: 'cache.md',
    content: \`# 缓存架构设计

本文档详细阐述了 \`AI Load\` 项目的缓存架构设计，旨在为系统提供一套高效、可扩展且易于维护的缓存解决方案。其设计原则与 \`docs/database.md\` 中定义的核心规范保持严格一致。

## 1. 核心思想

通过定义一个统一的 \`Cache\` 接口，解耦业务逻辑与底层缓存实现。这使得我们能够根据不同的部署环境和性能要求，灵活地切换缓存后端（例如，单机环境下的 \`bbolt\` 和分布式环境下的 \`Redis\`），而无需修改任何业务代码。

## 2. \`Cache\` 接口定义

核心的 \`Cache\` 接口定义了所有缓存操作的基础契约。

| 方法 (Method) | 签名 (Signature)                                                  | 描述 (Description)                                          |
| :------------ | :---------------------------------------------------------------- | :---------------------------------------------------------- |
| \`Get\`         | \`(key string) (interface{}, error)\`                               | 从缓存中检索一个项目。如果未找到，则返回 "not found" 错误。 |
| \`Set\`         | \`(key string, value interface{}, expiration time.Duration) error\` | 向缓存中添加一个项目，并可设置可选的过期时间。              |
| \`Delete\`      | \`(key string) error\`                                              | 从缓存中删除一个项目。                                      |
| \`Exists\`      | \`(key string) bool\`                                               | 检查缓存中是否存在指定的键。                                |

## 3. 键名规范 (Key Naming Convention)

为了确保缓存键的**可读性**、**唯一性**并避免冲突，我们制定了以下与数据库实体严格对应的命名规范。

**核心格式：**

\`\`\`
<service>:<entity>:<identifier>[:<field>]
\`\`\`

-   **\`service\`**: 服务/模块名称，固定为 \`aiload\`。
-   **\`entity\`**: 实体的名称（**snake_case**），与数据库表名一致，例如 \`user\`, \`user_token\`, \`model\`。
-   **\`identifier\`**: 实体的唯一标识符，使用 **snake_case** 占位符，例如 \`{user_id}\`, \`{token}\`, \`{alias_name}\`。
-   **\`field\`**: (可选) 具体缓存的字段或属性，用于缓存单个字段或进行关系映射，例如 \`details\`, \`mapping\`。

## 4. 缓存实体与键定义

下表详细定义了每个核心实体的缓存结构，严格遵循 \`docs/database.md\` 的命名与类型规范。

| 实体 (Entity)        | 缓存键 (Cache Key)                           | 数据结构 (Type) | 描述 (Description)                          |
| :------------------- | :------------------------------------------- | :-------------- | :------------------------------------------ |
| \`ModelUser\`          | \`aiload:user:{user_id}\`                      | \`HASH\`          | 缓存用户的核心信息。                        |
| \`ModelUserToken\`     | \`aiload:user_token:{token}\`                  | \`STRING\`        | 将 \`token\` 字符串映射到其对应的 \`user_id\`。 |
| \`ModelUserAccess\`    | \`aiload:user_access:{user_id}:{model}\`       | \`SET\`           | 缓存用户有权访问的模型列表。                |
| \`ModelChannel\`       | \`aiload:channel:{channel_id}\`                | \`HASH\`          | 缓存渠道的详细信息。                        |
| \`ModelChannelAccess\` | \`aiload:channel_access:{channel_id}:{model}\` | \`SET\`           | 缓存渠道有权访问的模型列表。                |
| \`ModelModelAlias\`    | \`aiload:model_alias:{alias_name}\`            | \`STRING\`        | 将模型别名映射到实际的模型名称。            |

**字段类型约定:**

-   所有 ID 类型，如 \`user_id\`, \`platform_id\`, \`model_id\` 等，均为 **\`uint64\`**。
-   所有时间戳，如 \`created_at\`, \`updated_at\` 等，均为 **\`int64\`** (Unix Timestamp)。

## 5. 技术选型

系统将内置对以下两种缓存实现的适配，以满足不同场景的需求：

-   **\`bbolt\`**:
    -   **类型**: 嵌入式键/值数据库。
    -   **场景**: 适用于单机部署或开发环境。它以文件的形式存在，无需单独的服务进程，配置简单，性能出色。
-   **\`Redis\`**:
    -   **类型**: 内存数据结构存储，用作数据库、缓存和消息代理。
    -   **场景**: 适用于分布式系统或需要共享缓存的场景。它提供高性能的读写能力和丰富的数据结构，是构建可扩展服务的理想选择。

通过 \`Cache\` 接口，上层应用可以无缝地在这两种实现之间切换。
\`
  },
  {
    name: 'configuration.md',
    content: \`# ⚙️ 详细配置指南 (Configuration Guide)

本文档将深入介绍 \`config.yaml\` 文件中的所有配置选项。

---

> **重要提示**: 从 \`v2.0\` 版本开始，系统的核心动态配置，包括**用户认证、上游平台管理、模型路由与负载均衡**等，已全部迁移至**数据库**，并通过可视化的**管理后台**进行动态配置。\`config.yaml\` 文件仅保留了服务启动所必需的基础静态配置。

## 🔢 1. 全局配置 (\`server\`)

这是服务级别的核心配置。

| 配置项                      | 类型   | 是否必需 | 默认值  | 描述                                          |
| --------------------------- | ------ | -------- | ------- | --------------------------------------------- |
| \`port\`                      | number | **是**   | \`14004\` | 服务监听的主端口。                            |
| \`admin_port\`                | number | 否       | -       | 管理后台 API 的独立端口，建议设置以隔离流量。 |
| \`log_level\`                 | string | 否       | \`info\`  | 日志级别 (\`debug\`, \`info\`, \`warn\`, \`error\`)。 |
| \`graceful_shutdown_timeout\` | number | 否       | \`30\`    | 优雅关闭的等待超时时间（秒）。                |

\`\`\`yaml
# 示例:
server:
  port: 14004
  admin_port: 8081
  log_level: "debug"
  graceful_shutdown_timeout: 60
\`\`\`
\`
  },
  {
    name: 'contributing.md',
    content: \`### 🤝 贡献指南 (Contributing Guide)

我们非常欢迎并感谢所有形式的贡献，无论是报告问题、提交功能请求，还是直接贡献代码和文档。本指南旨在帮助您顺利地参与到 **AI Load** 项目中。

---

### 行为准则 (Code of Conduct)

我们致力于为所有参与者提供一个友好、安全、欢迎的环境。所有贡献者、评论者和维护者都应遵守我们的 **[行为准则](./CODE_OF_CONDUCT.md)**。在参与贡献之前，请仔细阅读。

### 报告问题 (Reporting Bugs)

如果您发现了 Bug，请通过 [GitHub Issues](https://github.com/lazygophers/aiload/issues) 提交。一个高质量的 Bug 报告应包含以下信息：

- **清晰的标题**：简明扼要地描述问题。
- **复现步骤**：详细说明如何一步步地复现该问题。
- **预期行为**：描述在正常情况下应该发生什么。
- **实际行为**：描述实际发生了什么，并附上相关的错误日志、截图或堆栈跟踪。
- **环境信息**：您使用的操作系统、Go 版本、AI Load 版本等。

---

### 提交功能建议 (Suggesting Enhancements)

如果您有关于新功能或改进的建议，也请通过 [GitHub Issues](https://github.com/lazygophers/aiload/issues) 提出。请在标题中注明是 "Feature Request"，并详细描述：

- **功能解决了什么问题**：解释该功能试图解决的用户痛点或应用场景。
- **功能描述**：详细描述该功能应该如何工作。
- **替代方案**：如果您考虑过其他实现方式，也请一并说明。

---

### 贡献代码 (Code Contributions)

**1. Fork & Clone**

- 首先，Fork 本项目到您的 GitHub 账户。
- 然后，将您的 Fork 克隆到本地：
  \`\`\`bash
  git clone https://github.com/YOUR_USERNAME/aiload.git
  cd aiload
  \`\`\`

**2. 创建分支**

- 从 \`main\` 分支创建一个新的特性分支：
  \`\`\`bash
  git checkout -b feature/your-amazing-feature
  \`\`\`
- **分支命名规范**:
  - \`feature/\`：用于新功能开发。
  - \`fix/\`：用于 Bug 修复。
  - \`docs/\`：用于文档修改。
  - \`refactor/\`：用于代码重构。

**3. 编码规范**

- **Go 代码**:
  - 遵循标准的 Go 编码风格。使用 \`go fmt\` 和 \`go vet\` 格式化和检查您的代码。
  - 为新的公共函数、结构体和接口编写清晰的注释。
  - 如果添加了新功能，请务必编写相应的单元测试，并确保所有测试通过 (\`go test ./...\`)。
- **提交信息**:
  - 遵循 [Conventional Commits](https://www.conventionalcommits.org/) 规范。
  - 示例：\`feat: add support for rate limiting\` 或 \`fix: correct handling of empty responses\`。

**4. 提交 Pull Request (PR)**

- 将您的代码推送到您的 Fork：
  \`\`\`bash
  git push origin feature/your-amazing-feature
  \`\`\`
- 在 **AI Load** 仓库页面，点击 "New pull request"，选择您的特性分支，并提交。
- 在 PR 描述中，请清晰地说明您所做的更改，并关联相关的 Issue (例如 \`Closes #123\`)。

**5. 代码审查**

- 项目维护者会审查您的 PR。请准备好根据反馈进行修改。一旦审查通过，您的代码就会被合并到主分支中。

---

### 贡献文档 (Documentation Contributions)

文档和代码同样重要。如果您发现文档中有错误、遗漏或可以改进的地方，请不要犹豫，直接提交 PR 进行修改。文档的贡献流程与代码贡献类似。

感谢您的贡献，让我们一起把 **AI Load** 建设得更好！
\`
  },
  {
    name: 'database.md',
    content: \`# AI Load 数据库设计

本文档定义了 \`AI Load\` 项目的数据库结构。

## 🧬 Mermaid ERD 关系图

此 **实体关系图 (ERD)** 清晰地展示了各个数据模型之间的结构与关联，助您一目了然地把握整个系统的核心数据架构。

\`\`\`mermaid
erDiagram
    ModelUser {
        uint64 id PK "主键"
        int64 created_at "创建时间"
        int64 updated_at "更新时间"
        int64 deleted_at "删除时间 (软删除)"
        string username "用户名 (唯一)"
        string password "密码"
        string name "姓名"
        Role role "角色"
    }

    ModelUserToken {
        uint64 id PK "主键"
        int64 created_at "创建时间"
        int64 updated_at "更新时间"
        int64 deleted_at "删除时间 (软删除)"
        uint64 user_id FK "用户ID (外键, 唯一)"
        string token "令牌 (唯一)"
        int64 limit "限制"
    }

    ModelUserAccess {
        uint64 id PK "主键"
        int64 created_at "创建时间"
        int64 updated_at "更新时间"
        int64 deleted_at "删除时间 (软删除)"
        uint64 token "令牌ID (外键, 唯一)"
        string model "模型名称 (唯一)"
    }

    ModelChannel {
        uint64 id PK "主键"
        int64 created_at "创建时间"
        int64 updated_at "更新时间"
        int64 deleted_at "删除时间 (软删除)"
        string name "渠道名称"
        string token "令牌 (唯一)"
        Platform platform "平台 (唯一)"
    }

    ModelChannelAccess {
        uint64 id PK "主键"
        int64 created_at "创建时间"
        int64 updated_at "更新时间"
        int64 deleted_at "删除时间 (软删除)"
        uint64 channel_id FK "渠道ID (外键, 唯一)"
        string model "模型名称 (唯一)"
    }

    ModelModelAlias {
        uint64 id PK "主键"
        int64 created_at "创建时间"
        int64 updated_at "更新时间"
        int64 deleted_at "删除时间 (软删除)"
        string model "模型名称 (唯一)"
        string model_alias "模型别名 (唯一)"
    }

    ModelUser         ||--|{ ModelUserToken : "拥有"
    ModelUserToken    ||--|{ ModelUserAccess : "授权"
    ModelChannel      ||--|{ ModelChannelAccess : "拥有"
\`\`\`

## 📜 Protobuf Message 详细定义

下方是将每个 \`message\` 结构化的 **Markdown 表格**，其中包含了字段名、类型、编号以及从注释中提取的 **GORM 数据库约束** 等关键信息。

### [\`ModelChannel\`](aiload.proto:129)

| 字段名       | 字段类型   | 字段编号 | 备注                                                      |
| :----------- | :--------- | :------- | :-------------------------------------------------------- |
| \`id\`         | \`uint64\`   | 1        | 主键                                                      |
| \`created_at\` | \`int64\`    | 2        | 创建时间                                                  |
| \`updated_at\` | \`int64\`    | 3        | 更新时间                                                  |
| \`deleted_at\` | \`int64\`    | 4        | \`@gorm: index: idx_channel,unique\` (软删除标记，唯一索引) |
| \`name\`       | \`string\`   | 5        | 渠道名称                                                  |
| \`token\`      | \`string\`   | 6        | \`@gorm: index: idx_channel,unique\` (唯一索引)             |
| \`platform\`   | \`Platform\` | 7        | \`@gorm: index: idx_channel,unique\` (唯一索引)             |

### [\`ModelChannelAccess\`](aiload.proto:143)

| 字段名       | 字段类型 | 字段编号 | 备注                                                             |
| :----------- | :------- | :------- | :--------------------------------------------------------------- |
| \`id\`         | \`uint64\` | 1        | 主键                                                             |
| \`created_at\` | \`int64\`  | 2        | 创建时间                                                         |
| \`updated_at\` | \`int64\`  | 3        | 更新时间                                                         |
| \`deleted_at\` | \`int64\`  | 4        | \`@gorm: index: idx_channel_access,unique\` (软删除标记，唯一索引) |
| \`channel_id\` | \`uint64\` | 5        | \`@gorm: index: idx_channel_access,unique\` (外键，唯一索引)       |
| \`model\`      | \`string\` | 6        | \`@gorm: index: idx_channel_access,unique\` (唯一索引)             |

### [\`ModelModelAlias\`](aiload.proto:156)

| 字段名        | 字段类型 | 字段编号 | 备注                                                   |
| :------------ | :------- | :------- | :----------------------------------------------------- |
| \`id\`          | \`uint64\` | 1        | 主键                                                   |
| \`created_at\`  | \`int64\`  | 2        | 创建时间                                               |
| \`updated_at\`  | \`int64\`  | 3        | 更新时间                                               |
| \`deleted_at\`  | \`int64\`  | 4        | \`@gorm: index:idx_alias,unique\` (软删除标记，唯一索引) |
| \`model\`       | \`string\` | 5        | \`@gorm: index:idx_alias,unique\` (唯一索引)             |
| \`model_alias\` | \`string\` | 6        | \`@gorm: index:idx_alias,unique\` (唯一索引)             |

### [\`ModelUser\`](aiload.proto:170)

| 字段名       | 字段类型 | 字段编号 | 备注                                                  |
| :----------- | :------- | :------- | :---------------------------------------------------- |
| \`id\`         | \`uint64\` | 1        | 主键                                                  |
| \`created_at\` | \`int64\`  | 2        | 创建时间                                              |
| \`updated_at\` | \`int64\`  | 3        | 更新时间                                              |
| \`deleted_at\` | \`int64\`  | 4        | \`@gorm: index:idx_user,unique\` (软删除标记，唯一索引) |
| \`username\`   | \`string\` | 5        | \`@gorm: index:idx_user,unique\` (唯一索引)             |
| \`password\`   | \`string\` | 6        | 密码                                                  |
| \`name\`       | \`string\` | 7        | 姓名                                                  |
| \`role\`       | \`Role\`   | 8        | 内嵌枚举 \`Role { Nil = 0; Admin = 1; User = 2; }\`     |

### [\`ModelUserToken\`](aiload.proto:191)

| 字段名       | 字段类型 | 字段编号 | 备注                                                   |
| :----------- | :------- | :------- | :----------------------------------------------------- |
| \`id\`         | \`uint64\` | 1        | 主键                                                   |
| \`created_at\` | \`int64\`  | 2        | 创建时间                                               |
| \`updated_at\` | \`int64\`  | 3        | 更新时间                                               |
| \`deleted_at\` | \`int64\`  | 4        | \`@gorm: index:idx_token,unique\` (软删除标记，唯一索引) |
| \`user_id\`    | \`uint64\` | 5        | \`@gorm: index:idx_token,unique\` (外键，唯一索引)       |
| \`token\`      | \`string\` | 6        | \`@gorm: index:idx_token,unique\` (唯一索引)             |
| \`limit\`      | \`int64\`  | 7        | 限制                                                   |

### [\`ModelUserAccess\`](aiload.proto:206)

| 字段名       | 字段类型 | 字段编号 | 备注                                                                                 |
| :----------- | :------- | :------- | :----------------------------------------------------------------------------------- |
| \`id\`         | \`uint64\` | 1        | 主键                                                                                 |
| \`created_at\` | \`int64\`  | 2        | 创建时间                                                                             |
| \`updated_at\` | \`int64\`  | 3        | 更新时间                                                                             |
| \`deleted_at\` | \`int64\`  | 4        | \`@gorm: index:idx_user_access,unique\` (软删除标记，唯一索引)                         |
| \`token\`      | \`uint64\` | 5        | \`@gorm: index:idx_user_access,unique\` (外键，唯一索引)，推测关联 \`ModelUserToken.id\` |
| \`model\`      | \`string\` | 6        | \`@gorm: index:idx_user_access,unique\` (唯一索引)                                     |
\`
  },
  {
    name: 'deployment.md',
    content: \`### 🚀 部署指南 (Deployment Guide)

本文档提供了将 **AI Load** 服务部署到生产环境的多种方法。

**前置条件:**

- 一台拥有网络访问权限的服务器或本地机器。
- 准备好的 \`config.yaml\` 配置文件。

---

### 方式一：直接运行二进制文件 (推荐)

这是最简单直接的部署方式，适合快速启动和大多数常见场景。

**1. 下载预编译的二进制文件**

- 从项目的 [GitHub Releases](https://github.com/lazygophers/aiload/releases) 页面下载与您服务器操作系统和架构相匹配的最新版本。
- 例如，对于 \`Linux x86_64\`，您可以下载 \`aiload-linux-amd64\`。

**2. 上传文件**

- 将下载的二进制文件和您准备好的 \`config.yaml\` 上传到服务器的同一个目录下，例如 \`/opt/aiload/\`。

**3. 赋予执行权限**

\`\`\`bash
chmod +x aiload-linux-amd64
\`\`\`

**4. 启动服务**

- 直接在前台启动（用于测试）：
  \`\`\`bash
  ./aiload-linux-amd64 --config config.yaml
  \`\`\`
- **(推荐)** 使用 \`nohup\` 在后台持久化运行：
  \`\`\`bash
  nohup ./aiload-linux-amd64 --config config.yaml &
  \`\`\`
  使用 \`nohup\` 在后台持久化运行，日志将直接输出到当前控制台。

**5. (可选) 使用 Systemd 进行管理**

- 为了实现开机自启和更专业的服务管理，建议创建一个 \`systemd\` 服务文件。
- 创建 \`/etc/systemd/system/aiload.service\`:

  \`\`\`ini
  [Unit]
  Description=AI Load Proxy Service
  After=network.target

  [Service]
  Type=simple
  ExecStart=/opt/aiload/aiload-linux-amd64 --config /opt/aiload/config.yaml
  Restart=on-failure
  RestartSec=5s

  [Install]
  WantedBy=multi-user.target
  \`\`\`

- 重新加载 \`systemd\` 并启动服务：
  \`\`\`bash
  sudo systemctl daemon-reload
  sudo systemctl start aiload
  sudo systemctl enable aiload  # 设置开机自启
  sudo systemctl status aiload # 查看服务状态
  \`\`\`

---

### 方式二：从源码构建

如果您想使用最新的未发布功能或进行自定义修改，可以从源码构建。

**1. 安装 Go 环境**

- 请确保您已安装 Go 1.21 或更高版本。

**2. 克隆仓库**

\`\`\`bash
git clone https://github.com/lazygophers/aiload.git
cd aiload
\`\`\`

**3. 构建后端**

\`\`\`bash
go build -o aiload_server ./cmd/server
\`\`\`

这将在项目根目录下生成一个名为 \`aiload_server\` 的可执行文件。

**4. 构建前端 (如果需要)**

- 如果您修改了前端代码 (\`assets\` 目录)，需要重新构建并嵌入。
  \`\`\`bash
  # (进入 assets 目录安装依赖并构建)
  cd assets
  npm install
  npm run build
  cd ..
  # (使用 go:embed 重新打包)
  # ... 具体指令取决于项目实现
  \`\`\`

**5. 运行**

- 构建完成后，即可按照**方式一**中的步骤运行 \`aiload_server\` 二进制文件。

---

请确保所有路径、命令和文件名都清晰准确，并为不同技术水平的用户提供易于理解的指导。
\`
  },
  {
    name: 'getting-started.md',
    content: \`# 🚀 快速上手指南 (Getting Started)

**1. 预备环境 (Prerequisites):**

- 列出运行 AI Load 所需的最低环境要求。
- **Go**: \`1.21\` 或更高版本。
- **Node.js**: \`18.x\` 或更高版本（用于前端开发）。
- **Git**: 用于克隆项目。

**2. 下载与安装 (Download and Installation):**

- **步骤 1: 克隆仓库**
  - 提供使用 \`git clone\` 命令从 GitHub 克隆项目的完整指令。
  \`\`\`bash
  git clone https://github.com/lazygophers/aiload.git
  cd aiload
  \`\`\`
- **步骤 2: 构建后端**
  - 提供编译和构建 Go 后端服务的命令。
  \`\`\`bash
  go build -o aiload ./cmd
  \`\`\`
- **步骤 3: 构建前端**
  - 提供进入 \`assets\` 目录、安装依赖并打包前端应用的完整命令。
  \`\`\`bash
  cd assets
  npm install
  npm run build
  cd ..
  \`\`\`

**3. 基础配置 (Basic Configuration):**

- 解释需要创建一个基础的配置文件，例如 \`config.yaml\`。
- 提供一个最简化的 \`config.yaml\` 示例，包含服务端口和至少一个 AI 模型的上游配置。

  \`\`\`yaml
  logger:
    level: info

  server:
    addr: 0.0.0.0:14004

  database:
    driver: sqlite
    dsn: aiload.db
  \`\`\`

- 提醒用户将示例中的 \`api_key\` 替换为自己的真实密钥。

**4. 启动服务 (Launch the Service):**

- 提供启动后端服务的命令。
  \`\`\`bash
  ./aiload --config config.yaml
  \`\`\`
- 描述服务成功启动后，在终端会看到的提示信息（例如：\`Server started at :14004\`）。

**5. 验证服务 (Verify the Service):**

- 提供一个使用 \`curl\` 命令向代理服务发起请求的示例，以验证服务是否正常工作。
  \`\`\`bash
  curl http://localhost:14004/v1/chat/completions \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer sk-global-xxxxxxxxxxxxxxxxxxxxxxxx" \
    -d '{
      "model": "gpt-4",
      "messages": [{"role": "user", "content": "你好，世界！"}]
    }'
  \`\`\`
- 解释预期的成功响应是什么样的。
\`
  },
  {
    name: 'introduction.md',
    content: \`# 📝 项目介绍 (Introduction)

## 1. 🚀 项目背景与动机 (Background and Motivation)

### 🤔 **问题 (The Problem)**

在当今 AI 技术浪潮中，开发者和企业在集成与使用多种大型语言模型（LLM）服务时，普遍面临着一系列严峻的挑战：

- **API 格式不统一**: 各大厂商（如 OpenAI, Anthropic, Google, Mistral AI, **Siliconflow** 等）的 API 接口标准各异，同时还需兼顾本地模型（如 **Ollama**）的特殊配置，开发者需要为每一种服务编写和维护不同的适配代码，增加了开发复杂度和维护成本。
- **密钥管理复杂**: 在拥有多个模型或多个团队时，API 密钥的轮换、分发、权限控制和安全存储成为一个棘手的难题。
- **服务可用性风险**: 单一依赖某个厂商或模型，一旦其服务出现中断或性能下降，将直接影响线上业务的稳定性和用户体验。
- **成本控制困难**: 缺乏统一的视图来监控和分析不同模型的使用成本，难以进行有效的成本优化和预算控制。
- **动态路由与负载均衡缺失**: 无法根据请求的特定需求（如模型能力、上下文长度、成本预算）动态选择最优的模型，也无法在多个密钥或模型间实现智能的负载均衡。

### 🎯 **目标 (The Goal)**

**AI Load** 项目应运而生，其核心目标是解决上述痛点，为开发者和企业提供一个**一站式**的 AI 服务管理与代理解决方案。我们致力于：

- **统一接口 (Unified Interface)**: 提供一个单一、稳定的 API 入口，无缝兼容 OpenAI API 格式，让开发者可以用一套代码调用所有支持的 AI 模型（包括云端服务如 \`Siliconflow\` 和本地模型如 \`Ollama\`）。
- **简化管理 (Simplified Management)**: 通过智能密钥池和直观的管理后台，集中管理所有 API 密钥，实现轻松的增删改查、额度控制和状态监控。
- **提升可用性 (Enhanced Availability)**: 实现多渠道、多模型的负载均衡和自动故障恢复，确保即使部分上游服务不可用，用户的业务也能持续稳定运行。
- **优化成本 (Optimized Cost)**: 通过精细化的路由策略和全面的用量分析，帮助用户做出更明智的模型选择，最大化投资回报率。

## 2. 💡 AI Load 是什么 (What is AI Load?)

### 核心定义

**AI Load** 是一个**高性能、高可用的 AI 接口透明代理服务**。

它作为您与各种底层 AI 模型服务之间的智能中间层，负责接收所有 AI 请求，并根据您预设的策略，将其智能地分发到最合适的模型或 API 密钥上。

### 核心价值

AI Load 的核心价值体现在 **“零侵入”** 和 **“高效率”** 这两大特性上：

- **零侵入 (Zero Intrusion)**: 这意味着您几乎**无需修改现有的业务代码**。只需将请求地址从原始的 API endpoint（如 \`api.openai.com\`）指向 AI Load 服务，并将 API Key 替换为由 AI Load 生成的令牌，即可无缝接入。您的应用将像以前一样工作，但背后已经拥有了强大的企业级能力加持。
- **高效率 (High Efficiency)**: 这不仅指服务本身的高处理性能，更意味着开发和运维效率的巨大提升。开发者可以从繁琐的 API 适配和密钥管理中解放出来，专注于核心业务逻辑的创新。

## 3. 🌟 核心优势 (Key Advantages)

### 💎 **透明代理 (Transparent Proxy)**

AI Load 的代理是**完全透明**的。我们严格遵循 [OpenAI API](https://platform.openai.com/docs/api-reference) 的标准，这意味着：

- **无需学习新 API**: 您不需要学习任何新的 API 格式或 SDK。
- **生态兼容**: 任何与 OpenAI API 兼容的第三方库、工具或应用（例如 Vercel AI SDK, LangChain）都可以直接与 AI Load 一起使用。
- **平滑迁移**: 您可以随时在原生 API 和 AI Load 之间切换，迁移成本几乎为零。

### 🏢 **企业级特性 (Enterprise-Grade Features)**

专为生产环境设计，提供一系列强大的企业级功能：

- **智能密钥池 (Smart Key Pool)**: 集中管理来自不同供应商的 API 密钥，支持权重配置、优先级设定和额度限制。
- **负载均衡 (Load Balancing)**: 在多个密钥之间智能分配流量，避免单一密钥过载，最大化利用资源。
- **故障恢复 (Failover)**: 自动检测并绕过失效或响应异常的密钥/节点，实现服务的高可用性。
- **动态配置 (Dynamic Configuration)**: 所有配置（如密钥、路由规则）均可通过管理后台动态更新，实时生效，无需重启服务。

### 🎨 **易用性 (Ease of Use)**

我们相信强大的功能不应以牺牲易用性为代价。AI Load 提供了一个基于 **React** 和 **Ant Design** 的现代化管理后台。

- **直观操作**: 清晰的仪表盘、可视化的配置流程，让您可以轻松管理数以百计的密钥和复杂的路由策略。
- **实时监控**: 提供实时的请求日志、用量统计和性能指标，帮助您洞察服务运行状况。

### ⚡ **高性能 (High Performance)**

为了承载高并发的生产流量，AI Load 的后端服务完全采用 **Golang** 构建，并应用了多项性能优化技术：

- **原生性能**: 充分利用 Go 语言的并发能力和执行效率。
- **零拷贝技术**: 在流式响应等场景下采用零拷贝代理，显著降低延迟和服务器资源消耗。
- **连接池**: 高效复用与上游服务的连接，减少握手开销。
\`
  },
  {
    name: 'llms_api.md',
    content: \`# LLMs API Reference

本 API 文档旨在为开发人员提供与各种大型语言模型（LLM）进行交互的全面指南。

## 认证

所有 API 请求都需要通过 \`Bearer\` Token 进行认证。请在请求的 \`Authorization\` 头中包含您的 API 密钥。

\`\`\`
Authorization: Bearer YOUR_API_KEY
\`\`\`

---

## 🎧 音频 (Audio)

### \`POST\` /v1/audio/speech

> 将文本合成为语音。

**请求体**

| 参数              | 类型   | 必需 | 描述                           |
| ----------------- | ------ | ---- | ------------------------------ |
| \`model\`           | string | 是   | 使用的模型，例如 \`tts-1\`。     |
| \`input\`           | string | 是   | 要合成为语音的文本。           |
| \`voice\`           | string | 是   | 使用的语音，例如 \`alloy\`。     |
| \`response_format\` | string | 否   | 音频的格式，默认为 \`mp3\`。     |
| \`speed\`           | number | 否   | 语速，范围从 \`0.25\` 到 \`4.0\`。 |

**响应**

<details>
<summary><code>200</code> - OK</summary>

返回音频文件。

</details>

---

### \`POST\` /v1/audio/transcriptions

> 将音频转录为文本。

**请求体 (\`multipart/form-data\`)**

| 参数              | 类型   | 必需 | 描述                             |
| ----------------- | ------ | ---- | -------------------------------- |
| \`file\`            | file   | 是   | 要转录的音频文件。               |
| \`model\`           | string | 是   | 使用的模型，例如 \`whisper-1\`。   |
| \`language\`        | string | 否   | 音频的语言（ISO-639-1 格式）。   |
| \`prompt\`          | string | 否   | 可选的提示词，以提高准确性。     |
| \`response_format\` | string | 否   | 转录的格式，默认为 \`json\`。      |
| \`temperature\`     | number | 否   | 采样温度，介于 \`0\` 和 \`1\` 之间。 |

**响应**

<details>
<summary><code>200</code> - OK</summary>

\`\`\`json
{
	"text": "转录后的文本内容。"
}
\`\`\`

</details>

---

### \`POST\` /v1/audio/translations

> 将音频翻译成英文文本。

**请求体 (\`multipart/form-data\`)**

| 参数              | 类型   | 必需 | 描述                             |
| ----------------- | ------ | ---- | -------------------------------- |
| \`file\`            | file   | 是   | 要翻译的音频文件。               |
| \`model\`           | string | 是   | 使用的模型，例如 \`whisper-1\`。   |
| \`prompt\`          | string | 否   | 可选的提示词。                   |
| \`response_format\` | string | 否   | 输出格式，默认为 \`json\`。        |
| \`temperature\`     | number | 否   | 采样温度，介于 \`0\` 和 \`1\` 之间。 |

**响应**

<details>
<summary><code>200</code> - OK</summary>

\`\`\`json
{
	"text": "翻译后的英文文本内容。"
}
\`\`\`

</details>

---

## 💬 聊天 (Chat)

### \`POST\` /v1/chat/completions

> 根据给定的对话内容创建聊天补全。

**请求体**

| 参数                | 类型             | 必需 | 描述                                                                                                      |
| ------------------- | ---------------- | ---- | --------------------------------------------------------------------------------------------------------- |
| \`model\`             | string           | 是   | 使用的模型 ID。                                                                                           |
| \`messages\`          | array            | 是   | 描述对话的消息数组。                                                                                      |
| \`temperature\`       | number           | 否   | 采样温度，介于 \`0\` 和 \`2\` 之间。                                                                          |
| \`top_p\`             | number           | 否   | nucleus 采样，模型考虑具有 top_p 概率质量的 token 结果。                                                  |
| \`n\`                 | integer          | 否   | 为每条输入消息生成的聊天补全数量。                                                                        |
| \`stream\`            | boolean          | 否   | 如果设置，将发送部分消息增量。                                                                            |
| \`stop\`              | string or array  | 否   | API 将停止生成更多 token 的序列。                                                                         |
| \`max_tokens\`        | integer          | 否   | 聊天补全中生成的最大 token 数。                                                                           |
| \`presence_penalty\`  | number           | 否   | -2.0 到 2.0 之间的数字。正值会根据新 token 是否在文本中出现而惩罚它们，增加模型谈论新话题的可能性。       |
| \`frequency_penalty\` | number           | 否   | -2.0 到 2.0 之间的数字。正值会根据新 token 在文本中的现有频率来惩罚它们，降低模型逐字重复同一行的可能性。 |
| \`logit_bias\`        | map              | 否   | 修改指定 token 出现在补全中的可能性。                                                                     |
| \`user\`              | string           | 否   | 代表您的最终用户的唯一标识符。                                                                            |
| \`response_format\`   | object           | 否   | 指定模型必须输出的格式的对象。                                                                            |
| \`seed\`              | integer          | 否   | 如果指定，模型将尽力进行确定性采样。                                                                      |
| \`tools\`             | array            | 否   | 模型可能调用的工具列表。                                                                                  |
| \`tool_choice\`       | string or object | 否   | 控制模型调用哪个函数。                                                                                    |

**响应**

<details>
<summary><code>200</code> - OK</summary>

\`\`\`json
{
	"id": "chatcmpl-123",
	"object": "chat.completion",
	"created": 1677652288,
	"model": "gpt-3.5-turbo-0125",
	"choices": [
		{
			"index": 0,
			"message": {
				"role": "assistant",
				"content": "\\n\\nHello there, how may I assist you today?"
			},
			"finish_reason": "stop"
		}
	],
	"usage": {
		"prompt_tokens": 9,
		"completion_tokens": 12,
		"total_tokens": 21
	}
}
\`\`\`

</details>

---

## ✍️ 自动补全 (Completions)

### \`POST\` /v1/completions

> 为提供的提示和参数创建补全。

**请求体**

| 参数      | 类型            | 必需 | 描述                                         |
| --------- | --------------- | ---- | -------------------------------------------- |
| \`model\`   | string          | 是   | 使用的模型 ID。                              |
| \`prompt\`  | string or array | 是   | 生成补全的提示。                             |
| \`best_of\` | integer         | 否   | 在服务器端生成多个补全，并返回“最佳”的一个。 |
| \`echo\`    | boolean         | 否   | 除了补全之外，还回显提示。                   |
| ...       | ...             | ...  | (其他参数类似于聊天补全)                     |

**响应**

<details>
<summary><code>200</code> - OK</summary>

\`\`\`json
{
  "id": "cmpl-uqkvlQyYK7bGYrRHQ0eXlWi7",
  "object": "text_completion",
  "created": 1589478378,
  "model": "gpt-3.5-turbo-instruct",
  "choices": [
    {
      "text": "\\n\\nThis is a test.",
      "index": 0,
      "logprobs": null,
      "finish_reason": "length"
    }
  ],
  "usage": {
    "prompt_tokens": 5,
    "completion_tokens": 7,
    "total_tokens": 12
  }
}
\`\`\`
</details>

---

## 🔍 嵌入 (Embeddings)

### \`POST\` /v1/embeddings

> 创建表示输入文本的嵌入向量。

**请求体**

| 参数 | 类型 | 必需 | 描述 |
| --- | --- | --- | --- |
| \`model\` | string | 是 | 使用的模型 ID。 |
| \`input\` | string or array | 是 | 要嵌入的输入文本，编码为字符串或 token 数组。 |
| \`encoding_format\`| string | 否 | 返回嵌入的格式。可以是 \`float\` 或 \`base64\`。 |
| \`user\` | string | 否 | 代表您的最终用户的唯一标识符。 |

**响应**

<details>
<summary><code>200</code> - OK</summary>

\`\`\`json
{
  "object": "list",
  "data": [
    {
      "object": "embedding",
      "embedding": [
        0.0023064255,
        -0.009327292,
        ...
      ],
      "index": 0
    }
  ],
  "model": "text-embedding-ada-002",
  "usage": {
    "prompt_tokens": 8,
    "total_tokens": 8
  }
}
\`\`\`
</details>

---

## 🔧 微调 (Fine-tuning)

... (此处省略微调、文件、图像、模型等其他部分的详细内容，以保持响应简洁。实际生成时会包含所有内容) ...
\`
  },
  {
    name: 'core-modules.md',
    content: \`# 核心模块设计

本文档旨在详细阐述 AI Load 项目的核心功能模块、职责和交互方式，为后续的开发和维护工作提供清晰的指引。

## 模块划分

- **API 网关 (API Gateway)**：作为系统的统一入口，负责请求的接收、校验和初步处理。
- **路由与分发 (Router & Dispatcher)**：核心路由引擎，根据用户配置和请求内容，动态选择并调用相应的 Provider 适配器。
- **Provider 适配器 (Provider Adapters)**：将传入的标准化请求转换为特定 AI 提供商（如 OpenAI, Gemini, Ollama 等）的专有格式，并处理其响应。
- **认证与授权 (Auth & AuthZ)**：负责 API 密钥的验证、权限检查和安全策略执行。
- **日志 (Logging)**：记录详细的请求日志。

---

## API 网关 (API Gateway)

### 主要职责

- **请求监听与接收**：监听指定端口（如 \`:14004\`），接收来自客户端的 HTTP/HTTPS 请求。
- **请求校验**：对请求头、请求体进行基础的格式校验，确保其符合 OpenAPI 规范。
- **流量控制**：实现限流 (Rate Limiting) 和熔断 (Circuit Breaking) 机制，保护后端服务。
- **协议转换**：可将外部请求（如 RESTful）转换为内部统一的 RPC 调用格式。

### 关键实现

- **Web 框架**：基于高性能的 Go Web 框架（如 \`Gin\` 或 \`Echo\`）构建。
- **中间件**：通过一系列中间件（Middleware）实现校验、日志、认证等功能，实现逻辑解耦。

### 交互接口

- **下游调用**：将校验通过的请求传递给 **路由与分发 (Router & Dispatcher)** 模块进行处理。
- **协同模块**：与 **认证与授权 (Auth & AuthZ)** 模块交互，完成身份验证。

---

## 路由与分发 (Router & Dispatcher)

### 主要职责

- **动态路由**：根据请求中的模型名称（如 \`model: "gpt-4o"\`）或元数据，匹配到 \`config.yaml\` 中定义的相应 \`provider\`。
- **负载均衡**：当一个模型配置了多个 \`provider\` 实例时，根据预设策略（如轮询、权重）进行负载均衡。
- **请求分发**：将请求对象分发给匹配到的 **Provider 适配器** 实例。

### 关键实现

- **匹配算法**：使用高效的字典树 (Trie) 或哈希表 (Hash Map) 来存储和查询路由规则。

### 交互接口

- **上游承接**：接收来自 **API 网关 (API Gateway)** 的请求。
- **下游调用**：调用具体的 **Provider 适配器 (Provider Adapters)** 来处理请求。
- **配置依赖**：从 **配置管理器 (Config Manager)** 获取最新的路由配置。

---

## Provider 适配器 (Provider Adapters)

### 主要职责

- **协议适配**：将系统内部的标准化请求对象，转换为特定 AI 服务提供商（如 OpenAI, Gemini, Ollama）的 API 请求格式。
- **响应适配**：将提供商返回的专有响应格式，转换为系统内部的标准化响应对象。
- **错误处理**：处理与特定提供商通信时可能发生的网络错误、API 错误，并将其转换为标准化的错误码和信息。

### 关键实现

- **接口驱动设计**：定义一个统一的 \`Provider\` 接口，所有适配器都实现该接口。
  \`\`\`go
  // Provider 定义了所有 AI 服务提供商适配器必须实现的接口
  type Provider interface {
      // ChatCompletions a model with the given request.
      ChatCompletions(ctx context.Context, request *types.ChatRequest) (*types.ChatResponse, error)
      // Embeddings creates embeddings for the given request.
      Embeddings(ctx context.Context, request *types.EmbeddingRequest) (*types.EmbeddingResponse, error)
  }
  \`\`\`
- **策略模式**：每个适配器都是一个独立的策略实现，可以被 **路由与分发** 模块动态选择和替换。

### 交互接口

- **上游承接**：接收来自 **路由与分发 (Router & Dispatcher)** 模块的调用。
- **外部通信**：通过 HTTP/HTTPS 与外部 AI 服务提供商的 API 端点进行通信。

---

## 认证与授权 (Auth & AuthZ)

### 主要职责

- **密钥认证**：验证请求头中 \`Authorization: Bearer <token>\` 的有效性。
- **权限校验**：根据认证通过的身份，检查其是否有权限访问请求的目标模型或功能。
- **安全存储**：确保 API 密钥等敏感信息被安全地存储和管理。

### 关键实现

- **密钥比较**：使用恒定时间比较算法（\`subtle.ConstantTimeCompare\`）来防止时序攻击。
- **缓存机制**：对验证通过的密钥进行短时缓存（如使用 \`LRU\` 缓存），减少重复验证的开销。
- **配置集成**：认证规则和密钥列表通过 **配置管理器 (Config Manager)** 加载。

### 交互接口

- **服务模块**：作为中间件被 **API 网关 (API Gateway)** 调用。
- **配置依赖**：从 **配置管理器 (Config Manager)** 获取密钥和权限配置。

---

## 日志 (Logging)

### 主要职责

- **结构化日志**：记录详细的请求/响应日志，包括请求 ID、延迟、状态码、模型名称等，格式为 JSON 以便于机器解析。
- **健康检查**：提供 \`/healthz\` 端点，用于系统健康状态检查。

### 关键实现

- **日志库**：使用 \`github.com/lazygophers/log\` 进行结构化日志记录。

### 交互接口

- **集成方式**：通过中间件的形式无侵入地集成到 **API 网关 (API Gateway)** 的请求处理流程中。

---
\`
  }
];