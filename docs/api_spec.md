# 📖 API 文档 (aiload.proto)

本文档根据 `aiload.proto` 文件自动生成，详细描述了其中定义的所有服务接口、数据模型和枚举类型。

---

## 🔮 接口 (RPC) 文档

本部分详细说明了 `aiload` 服务提供的所有 RPC 方法。

### `POST /AddUserAdmin`

-   **描述**: 添加用户 (管理员)
-   **权限**: `admin`
-   **请求参数**: [`AddUserAdminReq`](#AddUserAdminReq)
-   **响应参数**: [`AddUserAdminRsp`](#AddUserAdminRsp)

### `POST /GetUser`

-   **描述**: 获取当前用户信息
-   **权限**: `user` (默认)
-   **请求参数**: [`GetUserReq`](#GetUserReq)
-   **响应参数**: [`GetUserRsp`](#GetUserRsp)

### `POST /GetUserAdmin`

-   **描述**: 获取指定用户信息 (管理员)
-   **权限**: `admin`
-   **请求参数**: [`GetUserAdminReq`](#GetUserAdminReq)
-   **响应参数**: [`GetUserAdminRsp`](#GetUserAdminRsp)

### `POST /ListUserAdmin`

-   **描述**: 列出所有用户 (管理员)
-   **权限**: `admin`
-   **请求参数**: [`ListUserAdminReq`](#ListUserAdminReq)
-   **响应参数**: [`ListUserAdminRsp`](#ListUserAdminRsp)

### `POST /SetUser`

-   **描述**: 更新当前用户信息
-   **权限**: `user` (默认)
-   **请求参数**: [`SetUserReq`](#SetUserReq)
-   **响应参数**: [`SetUserRsp`](#SetUserRsp)

### `POST /SetUserAdmin`

-   **描述**: 更新指定用户信息 (管理员)
-   **权限**: `admin`
-   **请求参数**: [`SetUserAdminReq`](#SetUserAdminReq)
-   **响应参数**: [`SetUserAdminRsp`](#SetUserAdminRsp)

### `POST /DelUserAdmin`

-   **描述**: 删除指定用户 (管理员)
-   **权限**: `admin`
-   **请求参数**: [`DelUserAdminReq`](#DelUserAdminReq)
-   **响应参数**: [`DelUserAdminRsp`](#DelUserAdminRsp)

### `POST /Login`

-   **描述**: 用户登录
-   **权限**: `public`
-   **请求参数**: [`LoginReq`](#LoginReq)
-   **响应参数**: [`LoginRsp`](#LoginRsp)

---

## 🧱 数据模型 (Message) 文档

本部分清晰地展示了每个 `message` 的数据结构。

### `ModelChannel`

| 字段         | 类型                    | 描述                                           |
| :----------- | :---------------------- | :--------------------------------------------- |
| `id`         | `uint64`                | 唯一标识符                                     |
| `created_at` | `int64`                 | 创建时间                                       |
| `updated_at` | `int64`                 | 更新时间                                       |
| `deleted_at` | `int64`                 | 删除时间，`@gorm: index: idx_channel,unique`   |
| `name`       | `string`                | 渠道名称                                       |
| `token`      | `string`                | 渠道 Token, `@gorm: index: idx_channel,unique` |
| `platform`   | [`Platform`](#Platform) | 平台类型, `@gorm: index: idx_channel,unique`   |

### `ModelChannelAccess`

| 字段         | 类型     | 描述                                                |
| :----------- | :------- | :-------------------------------------------------- |
| `id`         | `uint64` | 唯一标识符                                          |
| `created_at` | `int64`  | 创建时间                                            |
| `updated_at` | `int64`  | 更新时间                                            |
| `deleted_at` | `int64`  | 删除时间, `@gorm: index: idx_channel_access,unique` |
| `channel_id` | `uint64` | 渠道 ID, `@gorm: index: idx_channel_access,unique`  |
| `model`      | `string` | 模型名称, `@gorm: index: idx_channel_access,unique` |

### `ModelModelAlias`

| 字段          | 类型     | 描述                                      |
| :------------ | :------- | :---------------------------------------- |
| `id`          | `uint64` | 唯一标识符                                |
| `created_at`  | `int64`  | 创建时间                                  |
| `updated_at`  | `int64`  | 更新时间                                  |
| `deleted_at`  | `int64`  | 删除时间, `@gorm: index:idx_alias,unique` |
| `model`       | `string` | 模型名称, `@gorm: index:idx_alias,unique` |
| `model_alias` | `string` | 模型别名, `@gorm: index:idx_alias,unique` |

### `ModelUser`

| 字段         | 类型                    | 描述                                     |
| :----------- | :---------------------- | :--------------------------------------- |
| `id`         | `uint64`                | 唯一标识符                               |
| `created_at` | `int64`                 | 创建时间                                 |
| `updated_at` | `int64`                 | 更新时间                                 |
| `deleted_at` | `int64`                 | 删除时间, `@gorm: index:idx_user,unique` |
| `username`   | `string`                | 用户名, `@gorm: index:idx_user,unique`   |
| `password`   | `string`                | 密码                                     |
| `name`       | `string`                | 用户昵称                                 |
| `role`       | [`UserRole`](#UserRole) | 用户角色                                 |

### `ModelUserToken`

| 字段         | 类型     | 描述                                        |
| :----------- | :------- | :------------------------------------------ |
| `id`         | `uint64` | 唯一标识符                                  |
| `created_at` | `int64`  | 创建时间                                    |
| `updated_at` | `int64`  | 更新时间                                    |
| `deleted_at` | `int64`  | 删除时间, `@gorm: index:idx_token,unique`   |
| `user_id`    | `uint64` | 用户 ID, `@gorm: index:idx_token,unique`    |
| `token`      | `string` | 用户 Token, `@gorm: index:idx_token,unique` |
| `limit`      | `int64`  | 限制                                        |

### `ModelUserAccess`

| 字段         | 类型     | 描述                                            |
| :----------- | :------- | :---------------------------------------------- |
| `id`         | `uint64` | 唯一标识符                                      |
| `created_at` | `int64`  | 创建时间                                        |
| `updated_at` | `int64`  | 更新时间                                        |
| `deleted_at` | `int64`  | 删除时间, `@gorm: index:idx_user_access,unique` |
| `token`      | `uint64` | Token ID, `@gorm: index:idx_user_access,unique` |
| `model`      | `string` | 模型名称, `@gorm: index:idx_user_access,unique` |

### `AddUserAdminReq`

| 字段   | 类型                      | 描述                            |
| :----- | :------------------------ | :------------------------------ |
| `user` | [`ModelUser`](#ModelUser) | 用户信息, `@validate: required` |

### `AddUserAdminRsp`

| 字段   | 类型                      | 描述               |
| :----- | :------------------------ | :----------------- |
| `user` | [`ModelUser`](#ModelUser) | 创建成功的用户信息 |

### `GetUserReq`

空请求。

### `GetUserRsp`

| 字段   | 类型                      | 描述         |
| :----- | :------------------------ | :----------- |
| `user` | [`ModelUser`](#ModelUser) | 当前用户信息 |

### `GetUserAdminReq`

| 字段 | 类型     | 描述                           |
| :--- | :------- | :----------------------------- |
| `id` | `uint64` | 用户 ID, `@validate: required` |

### `GetUserAdminRsp`

| 字段   | 类型                      | 描述             |
| :----- | :------------------------ | :--------------- |
| `user` | [`ModelUser`](#ModelUser) | 获取到的用户信息 |

### `ListUserAdminReq`

| 字段          | 类型                               | 描述                            |
| :------------ | :--------------------------------- | :------------------------------ |
| `list_option` | `lazygophers.lrpc.core.ListOption` | 列表选项, `@validate: required` |

`ListUserAdminReq.ListOption` 枚举:
| 名称 | 值 | 描述 |
| :--- | :-: | :--- |
| `ListOptionNil` | 0 | 默认值 |
| `ListOptionName` | 1 | 按名称模糊搜索 |
| `ListOptionUsername` | 2 | 按用户名模糊搜索 |

### `ListUserAdminRsp`

| 字段       | 类型                             | 描述     |
| :--------- | :------------------------------- | :------- |
| `paginate` | `lazygophers.lrpc.core.Paginate` | 分页信息 |
| `list`     | `repeated ModelUser`             | 用户列表 |

### `SetUserReq`

| 字段   | 类型                      | 描述                                      |
| :----- | :------------------------ | :---------------------------------------- |
| `user` | [`ModelUser`](#ModelUser) | 需要更新的用户信息, `@validate: required` |

### `SetUserRsp`

| 字段   | 类型                      | 描述             |
| :----- | :------------------------ | :--------------- |
| `user` | [`ModelUser`](#ModelUser) | 更新后的用户信息 |

### `SetUserAdminReq`

| 字段   | 类型                      | 描述                                      |
| :----- | :------------------------ | :---------------------------------------- |
| `user` | [`ModelUser`](#ModelUser) | 需要更新的用户信息, `@validate: required` |

### `SetUserAdminRsp`

| 字段   | 类型                      | 描述             |
| :----- | :------------------------ | :--------------- |
| `user` | [`ModelUser`](#ModelUser) | 更新后的用户信息 |

### `DelUserAdminReq`

| 字段 | 类型     | 描述                                   |
| :--- | :------- | :------------------------------------- |
| `id` | `uint64` | 要删除的用户 ID, `@validate: required` |

### `DelUserAdminRsp`

空响应。

### `LoginReq`

| 字段       | 类型     | 描述                          |
| :--------- | :------- | :---------------------------- |
| `username` | `string` | 用户名, `@validate: required` |
| `password` | `string` | 密码, `@validate: required`   |

### `LoginRsp`

| 字段         | 类型                      | 描述           |
| :----------- | :------------------------ | :------------- |
| `user`       | [`ModelUser`](#ModelUser) | 用户信息       |
| `token`      | `string`                  | 登录凭证 Token |
| `expires_at` | `int64`                   | Token 过期时间 |

---

## 🔢 枚举 (Enum) 文档

本部分将每个 `enum` 的值和描述以表格形式呈现。

### `ErrCode`

| 名称                        | 值    | 描述              |
| :-------------------------- | :---- | :---------------- |
| `Success`                   | 0     | 成功              |
| `ModelAliasNotFound`        | 10000 | 模型别名未找到    |
| `ModelAliasDuplicateKey`    | 10001 | 模型别名重复      |
| `ChannelNotFound`           | 10002 | 渠道未找到        |
| `ChannelDuplicateKey`       | 10003 | 渠道重复          |
| `UserTokenNotFound`         | 10004 | 用户 token 未找到 |
| `UserTokenDuplicateKey`     | 10005 | 用户 token 重复   |
| `UserDuplicateKey`          | 10006 | 用户重复          |
| `UserNotFound`              | 10007 | 用户未找到        |
| `UserAccessNotFound`        | 10008 | 用户访问未找到    |
| `UserAccessDuplicateKey`    | 10009 | 用户访问重复      |
| `ChannelAccessNotFound`     | 10010 | 渠道访问未找到    |
| `ChannelAccessDuplicateKey` | 10011 | 渠道访问重复      |

### `UserRole`

| 名称     | 值  | 描述     |
| :------- | :-: | :------- |
| `Public` |  0  | 公共用户 |
| `User`   |  1  | 普通用户 |
| `Admin`  |  2  | 管理员   |
| `System` |  3  | 系统     |

### `Platform`

| 名称               | 值  | 描述        |
| :----------------- | :-: | :---------- |
| `Nil`              |  0  | 未指定      |
| `OpenAiCompatible` |  1  | OpenAI 兼容 |
| `Siliconflow`      |  2  | Siliconflow |
| `Gemini`           |  3  | Gemini      |
| `Ollama`           |  4  | Ollama      |
| `OpenAi`           |  5  | OpenAI      |
