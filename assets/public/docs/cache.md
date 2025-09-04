# 缓存架构设计

本文档详细阐述了 `AI Load` 项目的缓存架构设计，旨在为系统提供一套高效、可扩展且易于维护的缓存解决方案。其设计原则与 `docs/database.md` 中定义的核心规范保持严格一致。

## 1. 核心思想

通过定义一个统一的 `Cache` 接口，解耦业务逻辑与底层缓存实现。这使得我们能够根据不同的部署环境和性能要求，灵活地切换缓存后端（例如，单机环境下的 `bbolt` 和分布式环境下的 `Redis`），而无需修改任何业务代码。

## 2. `Cache` 接口定义

核心的 `Cache` 接口定义了所有缓存操作的基础契约。

| 方法 (Method) | 签名 (Signature)                                                  | 描述 (Description)                                          |
| :------------ | :---------------------------------------------------------------- | :---------------------------------------------------------- |
| `Get`         | `(key string) (interface{}, error)`                               | 从缓存中检索一个项目。如果未找到，则返回 "not found" 错误。 |
| `Set`         | `(key string, value interface{}, expiration time.Duration) error` | 向缓存中添加一个项目，并可设置可选的过期时间。              |
| `Delete`      | `(key string) error`                                              | 从缓存中删除一个项目。                                      |
| `Exists`      | `(key string) bool`                                               | 检查缓存中是否存在指定的键。                                |

## 3. 键名规范 (Key Naming Convention)

为了确保缓存键的**可读性**、**唯一性**并避免冲突，我们制定了以下与数据库实体严格对应的命名规范。

**核心格式：**

```
<service>:<entity>:<identifier>[:<field>]
```

-   **`service`**: 服务/模块名称，固定为 `aiload`。
-   **`entity`**: 实体的名称（**snake_case**），与数据库表名一致，例如 `user`, `user_token`, `model`。
-   **`identifier`**: 实体的唯一标识符，使用 **snake_case** 占位符，例如 `{user_id}`, `{token}`, `{alias_name}`。
-   **`field`**: (可选) 具体缓存的字段或属性，用于缓存单个字段或进行关系映射，例如 `details`, `mapping`。

## 4. 缓存实体与键定义

下表详细定义了每个核心实体的缓存结构，严格遵循 `docs/database.md` 的命名与类型规范。

| 实体 (Entity)        | 缓存键 (Cache Key)                           | 数据结构 (Type) | 描述 (Description)                          |
| :------------------- | :------------------------------------------- | :-------------- | :------------------------------------------ |
| `ModelUser`          | `aiload:user:{user_id}`                      | `HASH`          | 缓存用户的核心信息。                        |
| `ModelUserToken`     | `aiload:user_token:{token}`                  | `STRING`        | 将 `token` 字符串映射到其对应的 `user_id`。 |
| `ModelUserAccess`    | `aiload:user_access:{user_id}:{model}`       | `SET`           | 缓存用户有权访问的模型列表。                |
| `ModelChannel`       | `aiload:channel:{channel_id}`                | `HASH`          | 缓存渠道的详细信息。                        |
| `ModelChannelModel` | `aiload:channel_access:{channel_id}:{model}` | `SET`           | 缓存渠道有权访问的模型列表。                |
| `ModelModelAlias`    | `aiload:model_alias:{alias_name}`            | `STRING`        | 将模型别名映射到实际的模型名称。            |

**字段类型约定:**

-   所有 ID 类型，如 `user_id`, `platform_id`, `model_id` 等，均为 **`uint64`**。
-   所有时间戳，如 `created_at`, `updated_at` 等，均为 **`int64`** (Unix Timestamp)。

## 5. 技术选型

系统将内置对以下两种缓存实现的适配，以满足不同场景的需求：

-   **`bbolt`**:
    -   **类型**: 嵌入式键/值数据库。
    -   **场景**: 适用于单机部署或开发环境。它以文件的形式存在，无需单独的服务进程，配置简单，性能出色。
-   **`Redis`**:
    -   **类型**: 内存数据结构存储，用作数据库、缓存和消息代理。
    -   **场景**: 适用于分布式系统或需要共享缓存的场景。它提供高性能的读写能力和丰富的数据结构，是构建可扩展服务的理想选择。

通过 `Cache` 接口，上层应用可以无缝地在这两种实现之间切换。
