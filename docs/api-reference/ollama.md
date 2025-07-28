# 🦙 API 参考: Ollama

本文档旨在说明如何通过 **AI Load** 与 **Ollama** 进行交互，将您本地运行的开源大模型无缝集成到统一的 API 网关中。

**核心原则：协议转换与本地集成**

对于 Ollama，**AI Load** 同样实现了**协议转换**。这意味着您可以使用任何与 OpenAI 兼容的 SDK 和工具链，向本地运行的 Ollama 模型（如 Llama 3, Qwen 等）发送请求。**AI Load** 会在后台自动将您的 OpenAI 格式请求转换为 Ollama API 所需的格式，并将 Ollama 的响应转换回 OpenAI 格式。

**这使得您可以将本地运行的、注重数据隐私的开源模型，无缝集成到 AI Load 的统一管理和路由体系中**，与云端模型（如 GPT, Gemini）一同被调用和管理。

**1. 端点 URL 配置**

- 您将继续使用与 OpenAI 兼容的标准端点。
- **AI Load** 会根据您请求的`model`名称（例如 `"llama3"` 或 `"qwen:7b"`）和 `routing` 配置，智能地将请求转发到在 `upstreams` 中定义的本地 Ollama 服务。
- **请求端点:** `http://localhost:14004/v1/chat/completions`

**2. 认证方式**

- 与 OpenAI 的使用方式完全相同。使用您在 `config.yaml` 的 `auth` 部分配置的**全局密钥**或**分组密钥**作为 `Bearer Token`。

**3. `curl` 请求示例**

- 您无需关心 Ollama 的原生 API 格式。只需像调用 OpenAI 模型一样调用本地模型即可。

  ```bash
  # 使用在 config.yaml 中配置的密钥
  export AILOAD_API_KEY="your_configured_key"

  curl http://localhost:14004/v1/chat/completions \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $AILOAD_API_KEY" \
    -d '{
      "model": "llama3", # 指定本地 Ollama 模型
      "messages": [{"role": "user", "content": "给我讲一个关于罗马帝国的趣闻。"}]
    }'
  ```

**4. Python SDK 使用示例**

- 同样，您可以使用 `openai` Python SDK 与本地 Ollama 模型交互。

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
              "content": "Python 和 Go 的主要区别是什么？",
          }
      ],
      model="llama3", # AI Load 会根据此模型路由到本地 Ollama
  )

  print(chat_completion.choices[0].message.content)
  ```

**5. 关键优势**

- **纳管本地模型**: 将您自己的、在本地运行的开源模型纳入统一的 API 管理。
- **统一的 API 接口**: 无论后端是云端闭源模型还是本地开源模型，您都只需要对接一套 API。
- **数据隐私与安全**: 请求和数据保留在您的本地网络中，不经过公共互联网，确保了最高的安全性和隐私性。
- **灵活的模型路由**: 可以通过简单的配置更改，实现云端与本地模型的混合路由、故障转移和负载均衡。
