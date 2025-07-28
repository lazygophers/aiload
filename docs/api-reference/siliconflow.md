# ✨ API 参考: Siliconflow

本文档旨在说明如何通过 **AI Load** 与 **Siliconflow** API 进行交互。

**核心原则：兼容性与统一接口**

对于 Siliconflow，**AI Load** 提供了一个强大的**协议转换层**。这意味着您无需更改代码，可以继续使用您熟悉的 OpenAI 格式请求和 SDK 与 Siliconflow 提供的模型进行通信。**AI Load** 会在后台无缝处理请求和响应的转换。

这使得将 Siliconflow 模型集成到现有应用中变得异常简单，并保持了代码库的整洁和一致性。

**1. 端点 URL 配置**

- 您将继续使用标准的 OpenAI 兼容端点。
- **AI Load** 会根据您请求中的`model`名称（例如 `Qwen/Qwen2-7B-Instruct`）和 `routing` 配置，智能地将请求转发到在 `upstreams` 中定义的 Siliconflow 服务。
- **请求端点:** `/v1/chat/completions`

**2. 认证方式**

- 与 AI Load 的标准认证方式完全相同。使用您在 `config.yaml` 的 `auth` 部分配置的**全局密钥**或**分组密钥**作为 `Bearer Token`。

**3. `curl` 请求示例**

- 无需学习 Siliconflow 的特定 API 格式，像调用任何 OpenAI 模型一样即可。

  ```bash
  # 使用在 config.yaml 中配置的密钥
  export AILOAD_API_KEY="your_configured_key"

  curl http://localhost:14004/v1/chat/completions \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $AILOAD_API_KEY" \
    -d '{
      "model": "Qwen/Qwen2-7B-Instruct",
      "messages": [{"role": "user", "content": "写一首关于星际旅行的五言绝句。"}]
    }'
  ```

**4. Python SDK 使用示例**

- 同样，您可以使用 `openai` Python SDK，只需修改 `base_url` 即可与 Siliconflow 模型交互。

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
              "content": "What is the capital of California?",
          }
      ],
      model="Qwen/Qwen2-7B-Instruct", # AI Load 会根据此模型路由到 Siliconflow
  )

  print(chat_completion.choices[0].message.content)
  ```

**5. 关键优势**

- **统一的 API 接口**: 无论后端是 OpenAI、Siliconflow 还是其他厂商，您都只需要对接一套 API。
- **简化的代码库**: 无需为不同厂商维护不同的 SDK 或请求逻辑，降低了维护成本。
- **灵活的模型路由**: 可以通过简单的配置更改，在不同模型间无缝切换，而无需修改任何客户端代码。
