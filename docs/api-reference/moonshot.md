### 📄 API 参考: Moonshot (Kimi)

本文档旨在说明如何通过 **AI Load** 与 **Moonshot (Kimi)** API 进行交互。

#### 🌕 简介

**Moonshot (Kimi)** 是由月之暗面（Moonshot AI）公司开发的先进大型语言模型。它以其卓越的长文本处理能力而闻名，支持高达 20 万汉字的上下文窗口。- **主要模型**: `moonshot-v1-8k`, `moonshot-v1-32k`, `moonshot-v1-128k`- **核心优势**: 能够满足从快速问答到深度文档分析的各种需求。

**AI Load** 遵循与 OpenAI 兼容的 API 格式，让您可以无缝地将 Moonshot 集成到现有工作流中。#### 🔑 前置条件在开始之前，您需要访问 [Moonshot AI 开放平台](https://platform.moonshot.cn/) 并获取您的 API Key。

#### ⚙️ Admin UI 配置

通过 **AI Load** 的 Admin UI，您可以轻松地将 Moonshot API 集成进来。

**1. 创建 Platform**

- 在 Admin UI 的 `Platforms` 部分，点击 `Create`。
- **Name**: `moonshot` (或您喜欢的任何名称)
- **Provider**: `Moonshot`
- **Base URL**: `https://api.moonshot.cn`
- 保存更改。

**2. 添加 Token**

- 在 `Tokens` 部分，点击 `Create`。
- **Platform**: 选择您刚刚创建的 `moonshot` 平台。
- **API Key**: 粘贴您从 Moonshot 官方获取的 API Key。
- 保存更改。

**3. 配置 Model**

- 在 `Models` 部分，点击 `Create`。
- **Platform**: 选择 `moonshot` 平台。
- **Model ID**: `moonshot-v1-8k` (或其他模型，如 `moonshot-v1-32k`)- **Name**: `moonshot-v1-8k` (或您希望在请求中使用的自定义名称)
- **Status**: `Active`
- 保存更改。现在，您可以在 API 请求中使用 `moonshot-v1-8k` 这个模型了。

#### 🚀 使用示例

与 OpenAI 的集成类似，**AI Load** 作为透明代理，您只需修改请求的 `base_url` 和 `api_key`。

**`curl` 请求示例**

- **原始请求:**

  ```bash
  curl https://api.moonshot.cn/v1/chat/completions \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer YOUR_MOONSHOT_API_KEY" \
    -d '{
      "model": "moonshot-v1-8k",
      "messages": [
        {"role": "system", "content": "你是 Kimi，由 Moonshot AI 提供的人工智能助手。"},
        {"role": "user", "content": "你好，请给我介绍一下自己。"}
      ],
      "temperature": 0.3
    }'
  ```

- **通过 AI Load 的请求:**

  ```bash
  # 使用在 config.yaml 或 Admin UI 中配置的全局/分组密钥
  export AILOAD_API_KEY="your_configured_key"

  curl http://localhost:14004/v1/chat/completions \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $AILOAD_API_KEY" \
    -d '{
      "model": "moonshot-v1-8k", # AI Load 会根据此模型名称路由到正确的平台
      "messages": [
        {"role": "system", "content": "你是 Kimi，由 Moonshot AI 提供的人工智能助手。"},
        {"role": "user", "content": "你好，请给我介绍一下自己。"}
      ],
      "temperature": 0.3
    }'
  ```

**Python SDK 使用示例**

Moonshot 官方推荐使用与 OpenAI 兼容的 SDK。您只需修改 `base_url` 和 `api_key` 即可。

- **原始代码:**

  ```python
  from openai import OpenAI

  client = OpenAI(
      api_key = "YOUR_MOONSHOT_API_KEY",
      base_url = "https://api.moonshot.cn/v1",
  )

  completion = client.chat.completions.create(
    model="moonshot-v1-8k",
    messages=[
      {"role": "system", "content": "你是 Kimi，由 Moonshot AI 提供的人工智能助手。"},
      {"role": "user", "content": "你好，请给我介绍一下自己。"}
    ],
    temperature=0.3,
  )

  print(completion.choices[0].message.content)
  ```

- **通过 AI Load 的代码:**

  ```python
  from openai import OpenAI

  client = OpenAI(
      # 指向 AI Load 服务地址
      base_url="http://localhost:14004/v1",
      # 使用在 AI Load 中配置的全局/分组密钥
      api_key="your_configured_key"
  )
  chat_completion = client.chat.completions.create(
      messages=[
          {
              "role": "user",
              "content": "你好，请给我介绍一下自己。",
          }
      ],
      model="moonshot-v1-8k", # AI Load 会根据此模型名称进行路由
  )

  print(chat_completion.choices[0].message.content)
  ```
