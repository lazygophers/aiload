### Anthropic API 参考: Anthropic Claude

本文档旨在说明如何通过 **AI Load** 与 **Anthropic Claude** 系列模型进行交互。

**核心原则：协议转换与统一接口**

与 Gemini 类似，**AI Load** 通过强大的**协议转换**层，实现了对 Claude API 的无缝兼容。您无需关心 Anthropic 的原生 SDK 或其独特的 API 格式。只需通过标准的 OpenAI 格式发送请求，**AI Load** 将负责所有后台的翻译工作。

这使您能够将 Claude 模型（如 `claude-3-haiku` 或 `claude-3-opus`）集成到现有工作流中，而无需进行任何代码层面的修改。

**1. 端点 URL 配置**

- 您将继续使用与 OpenAI 兼容的端点。
- **AI Load** 会根据您请求的`model`名称（例如 `"claude-3-opus-20240229"`）和 `routing` 配置，智能地将请求转发到在 `upstreams` 中定义的 Claude 服务。
- **请求端点:** `http://localhost:14004/v1/chat/completions`

**2. 认证方式**

- 与 OpenAI 和 Gemini 的使用方式完全相同。使用您在 `config.yaml` 的 `auth` 部分配置的**全局密钥**或**分组密钥**作为 `Bearer Token`。

**3. `curl` 请求示例**

- 您无需理会 Claude API 的具体实现。只需像调用任何 OpenAI 模型一样调用 Claude 模型即可。

  ```bash
  # 使用在 config.yaml 中配置的密钥
  export AILOAD_API_KEY="your_configured_key"

  curl http://localhost:14004/v1/chat/completions \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $AILOAD_API_KEY" \
    -d '{
      "model": "claude-3-haiku-20240307", # 指定 Claude 模型
      "messages": [{"role": "user", "content": "Explain the concept of federalism in simple terms."}]
    }'
  ```

**4. Python SDK 使用示例**

- 同样，您可以使用 `openai` Python SDK 与 Claude 模型交互，无需引入 `anthropic` SDK。

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
              "content": "Write a short poem about the beauty of cherry blossoms.",
          }
      ],
      model="claude-3-opus-20240229", # AI Load 会根据此模型路由到 Claude
  )

  print(chat_completion.choices[0].message.content)
  ```

**5. 关键优势**

- **供应商无关性**: 您的应用程序与底层 AI 提供商完全解耦。
- **极致的灵活性**: 可以在不同厂商的旗舰模型之间进行 A/B 测试或动态切换，而无需重写任何代码。
- **统一的管理与监控**: 所有请求，无论最终由哪个模型处理，都通过 **AI Load** 进行统一的认证、日志记录和监控。
