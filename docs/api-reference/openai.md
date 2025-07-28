### 📄 API 参考: OpenAI

本文档旨在说明如何通过 **AI Load** 与兼容 **OpenAI** 的 API 进行交互。

**核心原则：透明代理**

**AI Load** 的核心设计哲学是**透明代理**。这意味着您无需更改任何现有的 OpenAI SDK 或 `curl` 请求代码。只需将请求的目标 URL 指向 **AI Load** 服务实例，并使用在 `config.yaml` 中配置的密钥即可。

**1. 端点 URL 配置**

- **原始请求:**
  ```
  https://api.openai.com/v1/chat/completions
  ```
- **通过 AI Load 的请求:**
  - 将 `https://api.openai.com` 替换为您的 **AI Load** 服务地址，例如 `http://localhost:8080`。
  - 最终端点为: `http://localhost:8080/v1/chat/completions`

**2. 认证方式**

- 使用您在 `config.yaml` 的 `auth` 部分配置的**全局密钥**或**分组密钥**作为 `Bearer Token`。

**3. `curl` 请求示例**

- **原始请求:**

  ```bash
  curl https://api.openai.com/v1/chat/completions \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $OPENAI_API_KEY" \
    -d '{
      "model": "gpt-4",
      "messages": [{"role": "user", "content": "Hello!"}]
    }'
  ```

- **通过 AI Load 的请求:**

  ```bash
  # 使用在 config.yaml 中配置的密钥
  export AILOAD_API_KEY="your_configured_key"

  curl http://localhost:8080/v1/chat/completions \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $AILOAD_API_KEY" \
    -d '{
      "model": "gpt-4",
      "messages": [{"role": "user", "content": "Hello!"}]
    }'
  ```

**4. Python SDK 使用示例**

- **原始代码:**

  ```python
  from openai import OpenAI

  client = OpenAI(
      api_key="YOUR_OPENAI_API_KEY"
  )

  chat_completion = client.chat.completions.create(
      messages=[
          {
              "role": "user",
              "content": "Say this is a test",
          }
      ],
      model="gpt-4",
  )
  ```

- **通过 AI Load 的代码:**

  - 只需修改 `base_url` 和 `api_key` 即可。

  ```python
  from openai import OpenAI

  client = OpenAI(
      # 指向 AI Load 服务地址
      base_url="http://localhost:8080/v1",
      # 使用在 config.yaml 中配置的密钥
      api_key="your_configured_key"
  )

  chat_completion = client.chat.completions.create(
      messages=[
          {
              "role": "user",
              "content": "Say this is a test",
          }
      ],
      model="gpt-4", # AI Load 会根据此模型名称进行路由
  )
  ```

**5. 支持的端点**

**AI Load** 透明支持所有标准的 OpenAI API v1 端点，包括但不限于：

- `/v1/chat/completions`
- `/v1/completions` (Legacy)
- `/v1/embeddings`
- `/v1/images/generations`
