# ♊️ API 参考: Google Gemini

本文档旨在说明如何通过 **AI Load** 与 **Google Gemini** API 进行交互。

**核心原则：协议转换与透明代理**

对于 Gemini，**AI Load** 不仅提供透明代理，还执行了**协议转换**。这意味着您仍然可以使用与 OpenAI 兼容的 SDK 和工具链向 Gemini 模型发送请求。**AI Load** 会在后台自动将您的 OpenAI 格式请求转换为 Gemini API 所需的格式，并将 Gemini 的响应转换回 OpenAI 格式。

这使得您可以在同一个应用中无缝切换和使用来自不同提供商的模型，而无需为每个提供商编写特定的代码。

**1. 端点 URL 配置**

- 您将继续使用与 OpenAI 兼容的端点。
- **AI Load** 会根据您请求的`model`名称（例如 `"gemini-pro"`）和 `routing` 配置，智能地将请求转发到在 `upstreams` 中定义的 Gemini 服务。
- **请求端点:** `http://localhost:14004/v1/chat/completions`

**2. 认证方式**

- 与 OpenAI 的使用方式完全相同。使用您在 `config.yaml` 的 `auth` 部分配置的**全局密钥**或**分组密钥**作为 `Bearer Token`。

**3. `curl` 请求示例**

- 您无需学习 Gemini 的特定 API 格式。只需像调用 OpenAI 模型一样调用 Gemini 模型即可。

  ```bash
  # 使用在 config.yaml 中配置的密钥
  export AILOAD_API_KEY="your_configured_key"

  curl http://localhost:14004/v1/chat/completions \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $AILOAD_API_KEY" \
    -d '{
      "model": "gemini-pro", # 指定 Gemini 模型
      "messages": [{"role": "user", "content": "Tell me a fun fact about the Roman Empire."}]
    }'
  ```

**4. Python SDK 使用示例**

- 同样，您可以使用 `openai` Python SDK 与 Gemini 模型交互。

  ```python
  from openai import OpenAI

  client = OpenAI(
      # 指向 AI Load 服务地址
      base_url="http://localhost:14004/v1",
      # 使用在 config.yaml 中配置的密钥
      api_key="your_configured_key"
  )

  chat_completion = client.chat.completions.create(
      messages=[
          {
              "role": "user",
              "content": "What are the main differences between Python and Go?",
          }
      ],
      model="gemini-pro", # AI Load 会根据此模型路由到 Gemini
  )

  print(chat_completion.choices[0].message.content)
  ```

**5. 关键优势**

- **统一的 API 接口**: 无论后端是 OpenAI、Gemini 还是其他服务，您都只需要对接一套 API。
- **简化的代码库**: 无需为不同厂商维护不同的 SDK 或请求逻辑。
- **灵活的模型路由**: 可以通过简单的配置更改，将流量从一个模型（如 `gpt-3.5-turbo`）切换到另一个（如 `gemini-pro`），而无需修改任何客户端代码。
