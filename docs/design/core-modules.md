# 核心模块设计

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

- **请求监听与接收**：监听指定端口（如 `:14004`），接收来自客户端的 HTTP/HTTPS 请求。
- **请求校验**：对请求头、请求体进行基础的格式校验，确保其符合 OpenAPI 规范。
- **流量控制**：实现限流 (Rate Limiting) 和熔断 (Circuit Breaking) 机制，保护后端服务。
- **协议转换**：可将外部请求（如 RESTful）转换为内部统一的 RPC 调用格式。

### 关键实现

- **Web 框架**：基于高性能的 Go Web 框架（如 `Gin` 或 `Echo`）构建。
- **中间件**：通过一系列中间件（Middleware）实现校验、日志、认证等功能，实现逻辑解耦。

### 交互接口

- **下游调用**：将校验通过的请求传递给 **路由与分发 (Router & Dispatcher)** 模块进行处理。
- **协同模块**：与 **认证与授权 (Auth & AuthZ)** 模块交互，完成身份验证。

---

## 路由与分发 (Router & Dispatcher)

### 主要职责

- **动态路由**：根据请求中的模型名称（如 `model: "gpt-4o"`）或元数据，匹配到 `config.yaml` 中定义的相应 `provider`。
- **负载均衡**：当一个模型配置了多个 `provider` 实例时，根据预设策略（如轮询、权重）进行负载均衡。
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

- **接口驱动设计**：定义一个统一的 `Provider` 接口，所有适配器都实现该接口。
  ```go
  // Provider 定义了所有 AI 服务提供商适配器必须实现的接口
  type Provider interface {
      // ChatCompletions a model with the given request.
      ChatCompletions(ctx context.Context, request *types.ChatRequest) (*types.ChatResponse, error)
      // Embeddings creates embeddings for the given request.
      Embeddings(ctx context.Context, request *types.EmbeddingRequest) (*types.EmbeddingResponse, error)
  }
  ```
- **策略模式**：每个适配器都是一个独立的策略实现，可以被 **路由与分发** 模块动态选择和替换。

### 交互接口

- **上游承接**：接收来自 **路由与分发 (Router & Dispatcher)** 模块的调用。
- **外部通信**：通过 HTTP/HTTPS 与外部 AI 服务提供商的 API 端点进行通信。

---

## 认证与授权 (Auth & AuthZ)

### 主要职责

- **密钥认证**：验证请求头中 `Authorization: Bearer <token>` 的有效性。
- **权限校验**：根据认证通过的身份，检查其是否有权限访问请求的目标模型或功能。
- **安全存储**：确保 API 密钥等敏感信息被安全地存储和管理。

### 关键实现

- **密钥比较**：使用恒定时间比较算法（`subtle.ConstantTimeCompare`）来防止时序攻击。
- **缓存机制**：对验证通过的密钥进行短时缓存（如使用 `LRU` 缓存），减少重复验证的开销。
- **配置集成**：认证规则和密钥列表通过 **配置管理器 (Config Manager)** 加载。

### 交互接口

- **服务模块**：作为中间件被 **API 网关 (API Gateway)** 调用。
- **配置依赖**：从 **配置管理器 (Config Manager)** 获取密钥和权限配置。

---

## 日志 (Logging)

### 主要职责

- **结构化日志**：记录详细的请求/响应日志，包括请求 ID、延迟、状态码、模型名称等，格式为 JSON 以便于机器解析。
- **健康检查**：提供 `/healthz` 端点，用于系统健康状态检查。

### 关键实现

- **日志库**：使用 `github.com/lazygophers/log` 进行结构化日志记录。

### 交互接口

- **集成方式**：通过中间件的形式无侵入地集成到 **API 网关 (API Gateway)** 的请求处理流程中。

---
