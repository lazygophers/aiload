### 📄 API 参考: DeepSeek

本文档旨在说明如何通过 **AI Load** 与 **DeepSeek** API 进行交互。

#### **简介**

**DeepSeek** 是由深度求索公司开发的一系列强大的大型语言模型。它以其卓越的代码生成和对话能力而闻名，主要模型包括 `deepseek-chat` (通用对话模型) 和 `deepseek-coder` (代码生成模型)。

**AI Load** 同样通过**透明代理**的方式支持 DeepSeek，允许您在不修改现有代码逻辑的情况下，无缝切换和使用 DeepSeek 的模型。

#### **1. 前置条件**

在使用之前，您必须拥有一个有效的 **DeepSeek API Key**。请前往 [DeepSeek 官方平台](https://platform.deepseek.com/) 注册并获取您的 API Key。

#### **2. Admin UI 配置**

为了让 **AI Load** 能够代理您的请求到 DeepSeek，您需要进行如下配置：

- **创建平台 (Platform)**:

  1.  进入 **AI Load** 管理后台，在“平台”或“Providers”部分，创建一个新平台。 2. 选择 `DeepSeek` 作为平台类型。 3. 为该平台指定一个唯一的别名，例如 `my-deepseek-platform`。

- **添加令牌 (Token)**: 1. 在您刚刚创建的 `my-deepseek-platform` 平台下，进入“令牌”或“Tokens”管理页面。 2. 添加一个新的令牌，并将您从 DeepSeek 官方获取的 API Key 粘贴于此。\* **配置模型 (Model)**: 1. 导航至“模型”或“Routing”配置部分。 2. 创建新的模型路由规则，将您希望使用的模型名称（例如 `deepseek-chat`）与 `my-deepseek-platform` 平台关联起来。

> **核心思想**: 经过以上配置，当 **AI Load** 收到一个请求，其 `model` 字段为 `deepseek-chat` 时，它会自动将此请求路由到 `my-deepseek-platform` 平台，并使用您配置的 DeepSeek API Key 进行认证。

#### **3. `curl` 请求示例**

- **原始请求:**

  ```bash
  curl https://api.deepseek.com/v1/chat/completions \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $DEEPSEEK_API_KEY" \
    -d '{
      "model": "deepseek-chat",
      "messages": [{"role": "user", "content": "Hello!"}]
    }'
  ```

- **通过 AI Load 的请求:**

  ```bash
  # 使用您在 AI Load 中配置的全局或分组密钥
  export AILOAD_API_KEY="your_configured_aiload_key"

  curl http://localhost:14004/v1/chat/completions \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $AILOAD_API_KEY" \
    -d '{
      "model": "deepseek-chat", # AI Load 会根据此模型名称路由到 DeepSeek
      "messages": [{"role": "user", "content": "Hello!"}]
    }'
  ```

#### **4. Python SDK 使用示例**

DeepSeek API 与 OpenAI 的官方 Python SDK 完全兼容，您无需引入新的库。

- **原始代码 (直接调用 DeepSeek):**

  ```python
  from openai import OpenAI

  client = OpenAI(
      api_key="YOUR_DEEPSEEK_API_KEY",
      base_url="https://api.deepseek.com/v1"
  )

  chat_completion = client.chat.completions.create(
      messages=[
          {
              "role": "user",
              "content": "你好，请介绍一下自己。",
          }
      ],
      model="deepseek-chat",
  )
  ```

- **通过 AI Load 的代码:**

  您只需将 `base_url` 指向 **AI Load** 服务地址，并使用其认证密钥即可。

  ```python
  from openai import OpenAI

  client = OpenAI(
      # 指向您的 AI Load 服务地址
      base_url="http://localhost:14004/v1",
      # 使用您在 AI Load 中配置的密钥 (而非原始 DeepSeek Key)
      api_key="your_configured_aiload_key"
  )

  chat_completion = client.chat.completions.create(
      messages=[
          {
              "role": "user",
              "content": "你好，请介绍一下自己。",
          }
      ],
      model="deepseek-chat", # AI Load 会根据此模型名称进行智能路由
  )
  ```
