### 📄 API 参考: Groq

本文档旨在说明如何通过 **AI Load** 与 **Groq** API 进行交互。

#### 🚀 Groq 平台简介

Groq 是一个专注于高性能推理的云平台，其核心是自研的 **LPU（Language Processing Unit）** 推理引擎。LPU 旨在为大型语言模型（LLM）提供超低延迟和极高的吞吐量，使其在实时应用中表现出色。#### 1. 前置条件在使用前，您需要拥有一个 GroqCloud 账户，并从中获取您的 API 密钥。

- **获取 API Key**：请访问 [GroqCloud 控制台](https://console.groq.com/keys) 创建并获取您的 API Key。

#### 2. Admin UI 配置

请按照以下步骤在 **AI Load** 的 Admin UI 中配置 Groq 平台、Token 和模型。

##### a. 创建 Groq 平台 (Platform)

1.  在 Admin UI 的 "Platforms" 页面，点击 "Create"。2. **Platform** 选择 `groq`。3. **Name** 字段可以自定义，例如 `my-groq-platform`。
2.  保存平台。

##### b. 添加 API Token

1.  在 "Tokens" 页面，点击 "Create"。
2.  **Platform** 选择您刚刚创建的 Groq 平台。3. **Token** 字段中填入您从 GroqCloud 获取的 API Key。
3.  保存 Token。

##### c. 配置模型 (Model)

1.  在 "Models" 页面，点击 "Create"。
2.  **Platform** 选择您的 Groq 平台。
3.  **Name** 字段填写您希望在请求中使用的模型名称，例如 `llama3-8b-8192`。4. **Model** 字段填写 Groq 官方对应的模型 ID，例如 `llama3-8b-8192`。5. **Type** 选择 `chat`。6. 保存模型。**可用模型示例：\*** `llama3-8b-8192`_ `llama3-70b-8192`_ `mixtral-8x7b-32768`\* `gemma-7b-it`#### 3. 使用示例

由于 Groq API 兼容 OpenAI 的格式，您可以像使用 OpenAI API 一样使用它。

##### a. `curl` 请求示例

- **原始请求 (直接调用 Groq):**

  ```bash
  curl -X POST "https://api.groq.com/openai/v1/chat/completions" \
    -H "Authorization: Bearer $GROQ_API_KEY" \
    -H "Content-Type: application/json" \
    -d '{
      "messages": [{"role": "user", "content": "Explain the importance of low latency LLMs"}],
      "model": "llama3-8b-8192"
    }'
  ```

- **通过 AI Load 的请求:**

  ```bash
  # 使用在 AI Load 中配置的全局或分组密钥
  export AILOAD_API_KEY="your_aiload_configured_key"

  curl http://localhost:14004/v1/chat/completions \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $AILOAD_API_KEY" \
    -d '{
      "model": "llama3-8b-8192", # AI Load 会根据此模型名称路由到 Groq 平台      "messages": [{"role": "user", "content": "Explain the importance of low latency LLMs"}]
    }'
  ```

##### b. Python SDK 使用示例

- **通过 AI Load 的代码:**

  ```python
  from openai import OpenAI

  client = OpenAI(
      # 指向 AI Load 服务地址
      base_url="http://localhost:14004/v1",
      # 使用在 AI Load 中配置的全局或分组密钥
      api_key="your_aiload_configured_key"
  )

  chat_completion = client.chat.completions.create(
      messages=[
          {
              "role": "user",
              "content": "Explain the importance of low latency LLMs",
          }
      ],
      model="llama3-8b-8192", # AI Load 会根据此模型名称进行路由
  )

  print(chat_completion.choices[0].message.content)
  ```
