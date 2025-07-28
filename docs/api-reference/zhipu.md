### 📄 API 参考: Zhipu AI (智谱 AI)

本文档旨在说明如何通过 **AI Load** 与 **Zhipu AI** 的 API 进行交互。**Zhipu AI** 是一家领先的人工智能公司，提供包括 `glm-4`、`glm-3-turbo` 在内的多种先进语言模型。### 1. ⚙️ 前置条件在使用之前，您需要拥有一个 Zhipu AI 账户，并从其官方平台获取 **API Key**。

- **官方地址:** [https://open.bigmodel.cn/](https://open.bigmodel.cn/)

### 2. 🎨 Admin UI 配置

#### 2.1. 创建 Platform

1.  导航至 **Platforms** 管理页面。
2.  点击 **"New Platform"**。
3.  **Type** 选择 `Zhipu`。
4.  为该平台指定一个唯一的 **Name** (例如 `my-zhipu-platform`)。

#### 2.2. 添加 Token

1.  导航至 **Tokens** 管理页面。
2.  点击 **"New Token"**。
3.  **Platform** 选择您刚刚创建的 Zhipu 平台。4. 在 **Token** 字段中，填入您从 Zhipu 官方获取的 API Key。#### 2.3. 配置 Model1. 导航至 **Models** 管理页面。2. 点击 **"New Model"**。3. 为模型指定一个 **Name** (例如 `glm-4`)，这个名称将用于 API 请求中的 `model` 字段。
4.  **Platform** 选择您创建的 Zhipu 平台。
5.  **Model ID** 填入 Zhipu 官方的模型标识 (例如 `glm-4`)。### 3. 🚀 使用示例**核心原则：透明代理\*\***AI Load** 的核心设计哲学是**透明代理**。这意味着您无需更改任何现有的 SDK 或 `curl` 请求代码。只需将请求的目标 URL 指向 **AI Load** 服务实例，并使用在 **Tokens\*\* 中配置的密钥即可。

#### 3.1. `curl` 请求示例

- **原始请求:**

  ```bash
  curl https://open.bigmodel.cn/api/paas/v4/chat/completions \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ZHIPU_API_KEY" \
    -d '{
      "model": "glm-4",
      "messages": [{"role": "user", "content": "你好！"}]
    }'
  ```

- **通过 AI Load 的请求:**

  ```bash
  # 使用在 AI Load Admin UI 中配置的 Token
  export AILOAD_API_KEY="your_aiload_token"

  curl http://localhost:14004/v1/chat/completions \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $AILOAD_API_KEY" \
    -d '{
      "model": "glm-4", # AI Load 会根据此模型名称路由到正确的 Zhipu 平台
      "messages": [{"role": "user", "content": "你好！"}]
    }'
  ```

#### 3.2. Python SDK 使用示例

- **原始代码:**

  ```python
  from zhipuai import ZhipuAI

  client = ZhipuAI(
      api_key="YOUR_ZHIPU_API_KEY"
  )

  response = client.chat.completions.create(
      model="glm-4",
      messages=[
          {"role": "user", "content": "你好"},
      ],
  )
  print(response)
  ```

- **通过 AI Load 的代码:**

  由于 **AI Load** 兼容 OpenAI 的 API 格式，我们可以直接使用 `openai` SDK 进行调用，只需修改 `base_url` 和 `api_key`。

  ```python
  from openai import OpenAI

  client = OpenAI(
      # 指向 AI Load 服务地址, 并附带 /v1 后缀
      base_url="http://localhost:14004/v1",
      # 使用在 AI Load Admin UI 中配置的 Token
      api_key="your_aiload_token"
  )

  chat_completion = client.chat.completions.create(
      messages=[
          {
              "role": "user",
              "content": "你好",
          }
      ],
      model="glm-4", # AI Load 会根据此模型名称进行路由
  )
  ```
