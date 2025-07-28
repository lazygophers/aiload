# ⚙️ 详细配置指南 (Configuration Guide)

本文档将深入介绍 `config.yaml` 文件中的所有配置选项。

---

> **重要提示**: 从 `v2.0` 版本开始，系统的核心动态配置，包括**用户认证、上游平台管理、模型路由与负载均衡**等，已全部迁移至**数据库**，并通过可视化的**管理后台**进行动态配置。`config.yaml` 文件仅保留了服务启动所必需的基础静态配置。

## 🔢 1. 全局配置 (`server`)

这是服务级别的核心配置。

| 配置项                      | 类型   | 是否必需 | 默认值  | 描述                                          |
| --------------------------- | ------ | -------- | ------- | --------------------------------------------- |
| `port`                      | number | **是**   | `14004` | 服务监听的主端口。                            |
| `admin_port`                | number | 否       | -       | 管理后台 API 的独立端口，建议设置以隔离流量。 |
| `log_level`                 | string | 否       | `info`  | 日志级别 (`debug`, `info`, `warn`, `error`)。 |
| `graceful_shutdown_timeout` | number | 否       | `30`    | 优雅关闭的等待超时时间（秒）。                |

```yaml
# 示例:
server:
  port: 14004
  admin_port: 8081
  log_level: "debug"
  graceful_shutdown_timeout: 60
```