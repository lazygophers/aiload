### 🏗️ 技术架构 (Technical Architecture)

本文档提供了 **AI Load** 服务的高级技术概览，旨在帮助开发者和贡献者理解其核心组件和内部工作流程。

**1. 核心设计理念**

- **高性能:** 采用 **Golang** 作为后端开发语言，充分利用其并发优势和高效的网络处理能力。
- **可扩展性:** 模块化的设计使得添加新的 AI 服务提供商或负载均衡策略变得简单。
- **易用性:** 通过统一的 API 端点和数据库驱动的动态配置，最大程度地简化用户的接入和维护成本。
- **状态无感知:** 服务本身是无状态的，易于水平扩展和部署。

**2. 架构图**

```mermaid
graph TD
    subgraph "用户请求"
        A[Client SDK / curl] --> B{AI Load Service};
    end

    subgraph "AI Load 核心服务 (Golang)"
        B --> C[1. 接收与解析请求];
        C --> D[2. 认证中间件];
        D --> E[3. 路由与模型匹配];
        E --> F[4. 负载均衡器];
        F --> G[5. 协议转换层];
        G --> H[6. 上游请求分发];
        H --> I[7. 响应处理与日志];
    end

    subgraph "数据与配置"
        J[config.yaml] -->(服务启动) B;
        DB[(Database)] -->(数据源) D;
        DB -->(数据源) E;
        DB -->(数据源) F;
    end

    B -->(日志) L[Console Output];

    subgraph "上游 AI 服务 (平台)"
        H --> M[OpenAI API];
        H --> N[Google Gemini API];
        H --> O[Anthropic Claude API];
        H --> P[Siliconflow API];
        H --> Q[Ollama (Local)];
    end

    I --> A;
```

**3. 组件详解**

- **接收与解析请求 (Request Handler):**

  - 作为服务的入口，监听指定端口上的 HTTP/HTTPS 请求。
  - 解析传入的请求头、路径和 Body，提取模型名称、消息内容等关键信息。

- **认证中间件 (Auth Middleware):**

  - 在请求到达核心逻辑之前，首先进行身份验证。
  - 从 `Authorization` 头中提取 `Bearer Token`。
  - 在数据库中查询匹配的 `Token` 信息进行验证。
  - 如果验证失败，立即返回 `401 Unauthorized` 错误。

- **路由与模型匹配 (Router):**

  - 根据请求的 `model`，查询数据库中的 `平台 (Platform)` -> `密钥 (Token)` -> `模型 (Model)` 关联关系，找到对应的上游服务进行处理。
  - 如果没有找到匹配的路由规则，返回错误。

- **负载均衡器 (Load Balancer):**

  - 获取到目标上游服务组后，根据预设的策略选择一个具体的上游实例。
  - 当前支持**加权轮询 (Weighted Round-Robin)** 策略，根据 `weight` 字段分配请求概率。
  - 未来可扩展支持其他策略，如最少连接数、最快响应时间等。

- **协议转换层 (Protocol Adapter):**

  - **这是实现多供应商兼容的核心**。

  * 如果目标上游是 OpenAI，则该层是透明的，直接转发请求。
  * 如果目标上游是 Gemini 或 Claude，该层会将 OpenAI 格式的请求体 (`messages` 结构等) 转换为目标供应商所需的格式。
  * 同样，在收到上游响应后，会将其转换回标准的 OpenAI API 响应格式。

- **上游请求分发 (Upstream Dispatcher):**

  - 使用 Golang 的 `net/http` 客户端，向上游服务实例发起实际的 API 调用。
  - 处理连接、超时 (`timeout`) 和重试 (`max_retries`) 逻辑。

- **响应处理与日志 (Response & Logging):**
  - 将从上游服务收到的（可能已经过协议转换的）响应返回给原始客户端。
  - 记录详细的请求/响应日志到**控制台输出 (Console Output)**，用于调试和审计。

**4. 核心数据模型与动态性**

> **AI Load** 的核心是围绕一个三层数据模型构建的，所有动态配置均存储在数据库中，并可通过管理后台实时更新，无需重启服务。
>
> - **平台 (Platform)**: 代表一个上游的 AI 服务提供商，如 OpenAI、Google 等。
> - **密钥 (Token)**: 关联到某一个**平台**，是实际用于请求认证的凭证。一个平台可以拥有多个密钥，便于轮换和管理。每个密钥都有独立的可用状态。
> - **模型 (Model)**: 关联到某一个或多个**密钥**，定义了具体的模型名称（如 `gpt-4-turbo`）。
>
> 这种设计使得路由和负载均衡完全由数据库中的数据驱动，实现了真正的动态化和高可用性。
