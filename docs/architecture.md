### 🏗️ 技术架构 (Technical Architecture)

本文档提供了 **AI Load** 服务的高级技术概览，旨在帮助开发者和贡献者理解其核心组件和内部工作流程。

**1. 核心设计理念**

- **高性能:** 采用 **Golang** 作为后端开发语言，充分利用其并发优势和高效的网络处理能力。
- **可扩展性:** 模块化的设计使得添加新的 AI 服务提供商或负载均衡策略变得简单。
- **易用性:** 通过统一的 API 端点和动态配置加载，最大程度地简化用户的接入和维护成本。
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

    subgraph "配置与监控"
        J[config.yaml] -->(动态加载) B;
        B -->(指标) K[Prometheus];
        B -->(日志) L[Request Logs];
    end

    subgraph "上游 AI 服务"
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
  - 在 `config.yaml` 定义的**全局密钥池**和**分组密钥池**中进行快速查找和匹配。
  - 如果验证失败，立即返回 `401 Unauthorized` 错误。

- **路由与模型匹配 (Router):**

  - 根据请求中指定的 `model` 名称，查询 `config.yaml` 中的 `routing` 规则。
  - 找到匹配的模型规则（支持通配符，如 `gpt-4-*`），确定该请求应该由哪个上游服务组 (`upstreams`) 来处理。
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
  - 记录详细的请求/响应日志（如果 `request_log_enabled` 为 `true`），用于调试和审计。
  - 更新 Prometheus 指标（如果 `prometheus_enabled` 为 `true`），如请求计数、延迟、错误率等。

**4. 动态配置加载**

- **AI Load** 启动时会加载 `config.yaml` 文件。
- 同时，它会使用文件系统监视（如 `fsnotify`）来监听 `config.yaml` 的变化。
- 当文件被修改并保存后，服务会**在不中断现有连接的情况下，平滑地重新加载所有配置**，包括密钥、上游服务和路由规则。这实现了真正的热重载，无需重启服务。

请确保图表和文字描述准确反映了服务的设计，并突出其关键优势。
