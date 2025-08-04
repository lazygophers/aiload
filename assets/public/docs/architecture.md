# 🏛️ 系统架构设计文档

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

```mermaid
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
```

### 2.2. 组件详解

- **1. 接收与解析请求 (Request Handler):** 服务入口，监听并解析 HTTP/HTTPS 请求。
- **2. 缓存查询 (Cache Lookup):** 在请求处理的早期阶段，优先查询 **bbolt** 缓存。
- **3. 认证中间件 (Auth Middleware):** 从 `Authorization` 头中提取密钥，在数据库中进行哈希校验。
- **4. 路由与模型匹配 (Router):** 根据请求的 `model`，从数据库查询并匹配到具体的上游服务平台和密钥组。
- **5. 负载均衡器 (Load Balancer):** 根据**加权轮询**策略，从密钥组中选择一个健康的实例。
- **6. 上游请求分发 (Upstream Dispatcher):** 通过**协议转换层 (Protocol Adapter)** 兼容多厂商 API，并向上游服务发起实际调用。
- **7. 响应处理与日志 (Response & Logging):** 将上游响应返回给客户端，并记录日志。成功处理的响应会**异步更新缓存**。

---

## 3. 数据持久化层 (SQLite)

### 3.1. 数据库文件

- **路径:** `data/aiload.db`
- **描述:** 系统所有持久化数据的“单一事实来源”。

### 3.2. 数据表模式 (Schema)

#### 🔑 a. `ApiKeys` (API 密钥表)

| 字段名               | 数据类型   | 约束/索引                   | 描述与设计 rationale                                                                              |
| -------------------- | ---------- | --------------------------- | ------------------------------------------------------------------------------------------------- |
| `id`                 | `INTEGER`  | `PRIMARY KEY AUTOINCREMENT` | 唯一密钥 ID                                                                                       |
| `user_id`            | `INTEGER`  | `FOREIGN KEY(Users.id)`     | 关联的用户 ID                                                                                     |
| **`key_prefix`**     | `TEXT`     | `NOT NULL, INDEX`           | **设计 rationale:** 密钥前缀 (例如 `sk-abc...`)，用于快速识别和日志追溯，**不用于安全校验**。     |
| **`key_hash`**       | `TEXT`     | `UNIQUE NOT NULL`           | **设计 rationale:** 使用 **SHA-256** 哈希后的完整密钥，用于安全校验，确保数据库中不存储明文密钥。 |
| `status`             | `TEXT`     | `NOT NULL`                  | 密钥状态 (`active`, `inactive`, `revoked`)                                                        |
| `rate_limit_per_min` | `INTEGER`  | `NOT NULL, DEFAULT 60`      | 每分钟允许的请求次数                                                                              |
| `expires_at`         | `DATETIME` | `NULL`                      | 密钥过期时间 (UTC)，`NULL` 表示永不过期                                                           |
| `created_at`         | `DATETIME` | `NOT NULL`                  | 密钥创建时间 (UTC)                                                                                |

#### 👤 b. `Users` (用户表)

| 字段名            | 数据类型   | 约束/索引                   | 描述                       |
| ----------------- | ---------- | --------------------------- | -------------------------- |
| `id`              | `INTEGER`  | `PRIMARY KEY AUTOINCREMENT` | 唯一用户 ID                |
| `username`        | `TEXT`     | `UNIQUE NOT NULL`           | 用户名，唯一               |
| `email`           | `TEXT`     | `UNIQUE NOT NULL`           | 电子邮箱，唯一             |
| `hashed_password` | `TEXT`     | `NOT NULL`                  | 加密后的用户密码           |
| `created_at`      | `DATETIME` | `NOT NULL`                  | 用户创建时间 (UTC)         |
| `updated_at`      | `DATETIME` | `NOT NULL`                  | 用户信息最后更新时间 (UTC) |

#### 🤖 c. `Models` (模型信息表)

| 字段名        | 数据类型  | 约束/索引                | 描述                       |
| ------------- | --------- | ------------------------ | -------------------------- |
| `id`          | `TEXT`    | `PRIMARY KEY`            | 模型 ID (例如 `gpt-4`)     |
| `provider`    | `TEXT`    | `NOT NULL`               | 模型提供商 (例如 `OpenAI`) |
| `description` | `TEXT`    | `NULL`                   | 模型的简短描述             |
| `is_active`   | `BOOLEAN` | `NOT NULL, DEFAULT TRUE` | 模型当前是否可用           |

#### 📊 d. `UsageLogs` (使用日志表)

| 字段名              | 数据类型   | 约束/索引                   | 描述                                          |
| ------------------- | ---------- | --------------------------- | --------------------------------------------- |
| `id`                | `INTEGER`  | `PRIMARY KEY AUTOINCREMENT` | 唯一日志 ID                                   |
| `api_key_id`        | `INTEGER`  | `FOREIGN KEY(ApiKeys.id)`   | 关联的 API 密钥 ID                            |
| `endpoint`          | `TEXT`     | `NOT NULL, INDEX`           | 请求的 API 端点 (例如 `/v1/chat/completions`) |
| `http_status`       | `INTEGER`  | `NOT NULL, INDEX`           | HTTP 响应状态码 (例如 200, 429)               |
| `response_time_ms`  | `INTEGER`  | `NOT NULL`                  | 从收到请求到发出响应的总耗时（毫秒）          |
| `request_timestamp` | `DATETIME` | `NOT NULL`                  | 请求到达服务器的时间 (UTC)                    |

---

## 4. 高性能缓存层 (bbolt)

为了最大限度地减少对数据库的直接访问并提升响应速度，我们引入了基于 **bbolt** 的嵌入式缓存层。

### 4.1. 缓存文件

- **路径:** `data/cache.db`
- **描述:** 用于存储高频访问数据的 Key-Value 数据库。

### 4.2. 缓存策略

#### a. 缓存内容与粒度

我们主要缓存两类数据：

1.  **API 密钥校验结果:**

    - **Key:** `token:<SHA-256-hash-of-the-full-key>`
    - **Value:** 序列化后的 `ApiKey` 对象（包含状态、速率限制等信息）。
    - **目的:** 避免每次请求都进行数据库哈希比对，这是最核心的性能优化点。

2.  **/v1/models 接口响应:**
    - **Key:** `models:list`
    - **Value:** 序列化后的模型列表 JSON 响应。
    - **目的:** `/v1/models` 是一个高频访问但内容不常变化的接口，对其进行缓存收益极高。

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
