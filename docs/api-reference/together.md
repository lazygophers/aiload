### 📄 API 参考: Together AI

本文档旨在说明如何通过 **AI Load** 与 **Together AI** 的 API 进行高效交互。

**核心理念：无缝集成**

**AI Load** 致力于提供与原生 API 一致的调用体验。您无需修改现有的 `curl` 或 SDK 代码逻辑，只需将请求指向 **AI Load** 服务，并使用其管理的认证信息，即可无缝切换至 Together AI 的强大模型。

---

### 1. 📖 简介

[Together AI](https://www.together.ai/) 是一个领先的云平台，专注于为开发者提供快速、高效的开源大语言模型推理服务。它提供了对众多顶级开源模型（如 Llama、Mistral、Mixtral 等）的访问能力，并以其卓越的性能和成本效益而著称。### 2. 🔑 前置条件在开始之前，您需要拥有一个 **Together AI** 账户，并从其官方网站获取您的 **API Key**。

- **获取地址**: [Together AI API Keys](https://api.together.ai/settings/api-keys)

---

### 3. ⚙️ Admin UI 配置

为了让 **AI Load** 能够代理您的请求，您需要在 Admin UI 中完成以下三个核心配置：**Platform**、**Token** 和 **Model**。

#### 3.1. 创建 Platform

**Platform** 用于定义一个上游 API 服务的基本信息。

1.  访问 **AI Load** 的 Admin UI。
2.  在 `Platforms` 管理页面，点击 `+ New Platform`。3. **Name**: 填入一个易于识别的名称，例如 `Together AI`。4. **Upstream URL**: 填入 Together AI 的官方 API 端点：`https://api.together.ai/v1`
3.  **Type**: 选择 `openai` 类型，因为 Together AI 的 API 格式与 OpenAI 兼容。
4.  点击 `Save` 保存。

#### 3.2. 添加 Token

**Token** 用于存储您从 Together AI 获取的 API Key。

1.  在 `Tokens` 管理页面，点击 `+ New Token`。2. **Name**: 填入一个描述性名称，例如 `Together AI Key 1`。
2.  **Platform**: 从下拉列表中选择您刚刚创建的 `Together AI` 平台。
3.  **Token**: 填入您从 Together AI 官网获取的 API Key。
4.  点击 `Save` 保存。

#### 3.3. 配置 Model

**Model** 用于将一个模型名称（例如 `llama3-70b-8192`) 映射到您配置的平台和 Token。

1.  在 `Models` 管理页面，点击 `+ New Model`。
2.  **Model ID**: 填入您希望在请求中使用的模型名称，例如 `meta-llama/Llama-3-70b-chat-hf`。**这个 ID 必须与 Together AI 支持的模型标识符完全一致**。
3.  **Platform**: 从下拉列表中选择 `Together AI`。
4.  **Tokens**: 选择您刚刚添加的 `Together AI Key 1`。
5.  **Description**: (可选) 添加模型描述。6. 点击 `Save` 保存。您可以重复此步骤，添加多个不同的 Together AI 模型，例如：

- `mistralai/Mixtral-8x7B-Instruct-v0.1`_ `Qwen/Qwen1.5-72B-Chat`_ `google/gemma-7b-it`---### 4. 🚀 使用示例

配置完成后，您可以像调用标准 OpenAI API 一样，通过 **AI Load** 使用 Together AI 的模型。

#### 4.1. `curl` 请求示例

- **原始请求:**

  ```bash
  curl https://api.together.ai/v1/chat/completions \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOGETHER_API_KEY" \
    -d '{
      "model": "meta-llama/Llama-3-70b-chat-hf",
      "messages": [{"role": "user", "content": "Hello, how are you?"}]
    }'
  ```

- **通过 AI Load 的请求:**
  只需将 `base_url` 替换为 **AI Load** 地址，并使用 **AI Load** 的认证密钥。

  ```bash
  # 使用在 config.yaml 中配置的 AI Load 全局密钥
  export AILOAD_API_KEY="your_aiload_auth_key"

  curl http://localhost:14004/v1/chat/completions \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $AILOAD_API_KEY" \
    -d '{
      "model": "meta-llama/Llama-3-70b-chat-hf",
      "messages": [{"role": "user", "content": "Hello, how are you?"}]
    }'
  ```

#### 4.2. Python SDK 使用示例

- **原始代码:**

  ```python
  from openai import OpenAI

  client = OpenAI(
      api_key="YOUR_TOGETHER_API_KEY",
      base_url="https://api.together.ai/v1",
  )

  chat_completion = client.chat.completions.create(
      messages=[{"role": "user", "content": "Tell me a joke."}],
      model="mistralai/Mixtral-8x7B-Instruct-v0.1",
  )
  ```

- **通过 AI Load 的代码:**
  修改 `base_url` 指向 **AI Load** 服务，并使用其 `api_key`。

  ```python
  from openai import OpenAI

  client = OpenAI(
      # 指向 AI Load 服务地址
      base_url="http://localhost:14004/v1",
      # 使用在 config.yaml 中配置的 AI Load 密钥
      api_key="your_aiload_auth_key"
  )

  chat_completion = client.chat.completions.create(
      messages=[{"role": "user", "content": "Tell me a joke."}],
      # AI Load 会根据此模型 ID 自动路由到 Together AI 平台
  ```
