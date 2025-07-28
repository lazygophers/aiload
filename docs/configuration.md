# ⚙️ 详细配置指南 (Configuration Guide)

本文档将深入介绍 `config.yaml` 文件中的所有配置选项。

---

## 🔢 1. 全局配置 (`server`)

这是服务级别的核心配置。

| 配置项                      | 类型   | 是否必需 | 默认值 | 描述                                          |
| --------------------------- | ------ | -------- | ------ | --------------------------------------------- |
| `port`                      | number | **是**   | -      | 服务监听的主端口。                            |
| `admin_port`                | number | 否       | -      | 管理后台 API 的独立端口，建议设置以隔离流量。 |
| `log_level`                 | string | 否       | `info` | 日志级别 (`debug`, `info`, `warn`, `error`)。 |
| `graceful_shutdown_timeout` | number | 否       | `30`   | 优雅关闭的等待超时时间（秒）。                |

```yaml
# 示例:
server:
  port: 8080
  admin_port: 8081
  log_level: "debug"
  graceful_shutdown_timeout: 60
```

---

## 🔑 2. 认证配置 (`auth`)

管理所有传入请求的 API 密钥认证。

### **全局密钥 (`keys`)**

这是一组适用于所有请求的默认密钥池。如果请求没有匹配到任何特定的 `groups`，将使用此处的密钥进行验证。

```yaml
# 示例:
auth:
  keys:
    - "sk-default-key-1"
    - "sk-default-key-2"
```

### **分组密钥 (`groups`)**

您可以为特定的用户、团队或应用场景创建独立的认证分组，以实现更精细的访问控制。

- `name`: 分组的唯一名称。
- `keys`: 该分组专属的密钥列表。

```yaml
# 示例:
auth:
  groups:
    - name: "team_a"
      keys:
        - "team_a_key1"
        - "team_a_key2"
    - name: "high_priority_users"
      keys:
        - "hp_key1"
```

---

## ☁️ 3. 上游服务配置 (`upstreams`)

这里定义了所有后端 AI 服务的连接池。您可以配置一个或多个提供商，并为每个提供商设置多个实例。

### **通用字段**

| 配置项        | 类型   | 是否必需 | 默认值 | 描述                                                                       |
| ------------- | ------ | -------- | ------ | -------------------------------------------------------------------------- |
| `name`        | string | **是**   | -      | 上游服务的唯一标识，用于在路由中引用。                                     |
| `provider`    | string | **是**   | -      | AI 服务提供商 (`openai`, `gemini`, `claude`, `siliconflow`, `ollama` 等)。 |
| `api_key`     | string | **是**   | -      | 对应服务商的 API 密钥。                                                    |
| `weight`      | number | 否       | `100`  | 用于加权负载均衡的权重，值越大，优先级越高。                               |
| `timeout`     | number | 否       | `120`  | 请求超时时间（秒）。                                                       |
| `max_retries` | number | 否       | `2`    | 请求失败后的最大重试次数。                                                 |

### **完整示例**

这是一个包含多种服务和多个实例的复杂配置，展示了如何结合使用不同的提供商和权重。

```yaml
upstreams:
  - name: "openai_us_east"
    provider: "openai"
    api_key: "sk-openai-east-xxxx"
    weight: 200
    timeout: 180
  - name: "openai_eu_west"
    provider: "openai"
    api_key: "sk-openai-west-xxxx"
    weight: 100
  - name: "gemini_pro_main"
    provider: "gemini"
    api_key: "ai-gemini-xxxx"
    weight: 150
    max_retries: 3
  - name: "claude_haiku"
    provider: "claude"
    api_key: "sk-claude-xxxx"
  - name: "siliconflow_qwen"
    provider: "siliconflow"
    api_key: "sk-siliconflow-xxxx"
    weight: 120
  - name: "local_llama3"
    provider: "ollama"
    # Ollama 通常在本地运行，因此 api_key 不是必需的
    # 但需要配置 base_url 来指定 Ollama 服务的地址
    base_url: "http://localhost:11434"
    api_key: "ollama" # 占位，可任意填写
```

---

## 🗺️ 4. 负载均衡与路由 (`routing`)

路由规则定义了如何根据请求的模型名称将其分发到指定的上游服务组。规则按顺序匹配。

- `model`: 匹配模型名称，支持使用 `*` 作为通配符。
- `upstreams`: 一个或多个上游服务的 `name` 列表，请求将根据权重在这些服务之间进行负载均衡。

```yaml
# 示例:
routing:
  - model: "gpt-4-*"
    upstreams:
      - "openai_us_east"
      - "openai_eu_west"
  - model: "gemini-pro"
    upstreams:
      - "gemini_pro_main"
  - model: "claude-3-haiku-*"
    upstreams:
      - "claude_haiku"
  - model: "Qwen/Qwen2-7B-Instruct"
    upstreams:
      - "siliconflow_qwen"
  - model: "llama3"
    upstreams:
      - "local_llama3"
```

---

## 📊 5. 监控与日志 (`monitoring`)

配置可观测性相关的选项。

| 配置项                | 类型    | 是否必需 | 默认值  | 描述                                        |
| --------------------- | ------- | -------- | ------- | ------------------------------------------- |
| `prometheus_enabled`  | boolean | 否       | `false` | 是否启用 Prometheus 指标端点 (`/metrics`)。 |
| `prometheus_port`     | number  | 否       | `9090`  | Prometheus 端点的独立端口。                 |
| `request_log_enabled` | boolean | 否       | `true`  | 是否记录详细的请求/响应日志。               |

```yaml
# 示例:
monitoring:
  prometheus_enabled: true
  prometheus_port: 9091
  request_log_enabled: true
```
