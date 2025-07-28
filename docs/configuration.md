# ⚙️ 详细配置指南 (Configuration Guide)

本文档将深入介绍 `config.yaml` 文件中的所有配置选项。

---

## 🔢 1. 全局配置 (`server`)

这是服务级别的核心配置。

| 配置项         | 类型     | 是否必需  | 默认值    | 描述                                       |
|-------------|--------|-------|--------|------------------------------------------|
| `port`      | number | **是** | -      | 服务监听的主端口。                                |
| `log_level` | string | 否     | `info` | 日志级别 (`debug`, `info`, `warn`, `error`)。 |

```yaml
# 示例:
server:
    port: 8080
    log_level: "debug"
```

---


---

## 🗺️ 4. 负载均衡与路由 (`routing`)

路由规则定义了如何根据请求的模型名称将其分发到指定的上游服务组。规则按顺序匹配。

- `model`: 匹配模型名称，支持使用 `*` 作为通配符。
- `upstreams`: 一个或多个上游服务的 `name` 列表，请求将根据权重在这些服务之间进行负载均衡。

```yaml
# 示例:
routing:
    -   model: "gpt-4-*"
        upstreams:
            - "openai_us_east"
            - "openai_eu_west"
    -   model: "gemini-pro"
        upstreams:
            - "gemini_pro_main"
    -   model: "claude-3-haiku-*"
        upstreams:
            - "claude_haiku"
    -   model: "Qwen/Qwen2-7B-Instruct"
        upstreams:
            - "siliconflow_qwen"
    -   model: "llama3"
        upstreams:
            - "local_llama3"
```

---

## 📊 5. 监控与日志 (`monitoring`)

配置可观测性相关的选项。

| 配置项                   | 类型      | 是否必需 | 默认值     | 描述                                 |
|-----------------------|---------|------|---------|------------------------------------|
| `prometheus_enabled`  | boolean | 否    | `false` | 是否启用 Prometheus 指标端点 (`/metrics`)。 |
| `prometheus_port`     | number  | 否    | `9090`  | Prometheus 端点的独立端口。                |
| `request_log_enabled` | boolean | 否    | `true`  | 是否记录详细的请求/响应日志。                    |

```yaml
# 示例:
monitoring:
    prometheus_enabled: true
    prometheus_port: 9091
    request_log_enabled: true
```