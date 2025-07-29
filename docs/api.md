# 🚀 AI Load 统一 API 参考文档

欢迎使用 **AI Load**！本文档是您与 AI Load 服务进行交互的终极指南。

## 核心理念：智能路由与协议转换

**AI Load** 的核心价值在于其强大的**智能路由**与**协议转换**能力。

1.  **透明代理 (Transparent Proxy)**
    对于本身就兼容 OpenAI API 格式的服务（如 Groq, Moonshot, Together AI 等），AI Load 扮演**透明代理**的角色。您无需更改任何代码，只需将请求指向 AI Load，即可享受统一的认证、日志、限流和重试等中间件服务。

2.  **协议转换 (Protocol Conversion)**
    对于 API 格式与 OpenAI **不兼容**的服务（如 Google Gemini, Anthropic Claude, Zhipu AI），AI Load 会执行**自动的协议转换**。这意味着：
    - **您发送的请求**：始终使用标准的 OpenAI 格式。
    - **AI Load 的工作**：在后台将您的请求实时转换为目标服务（如 Gemini）所需的原生格式。
    - **返回给您的响应**：AI Load 会再将目标服务的原生响应格式转换回标准的 OpenAI 格式。

**最终实现的效果是**：无论后端模型千变万化，您的客户端代码**永远保持统一和简洁**。您可以使用任何熟悉的 OpenAI-Compatible SDK 或工具，无缝调用所有已接入的大模型，而无需为每个厂商编写和维护特定的适配代码。

---

## API 详情

### 1. 基础 URL

所有 API 请求的基础路径 (Base URL) 如下：

```
http://localhost:14004/v1
```

### 2. 认证方式

所有请求都需要在 `Header` 中携带认证信息。

- **Header**: `Authorization`
- **格式**: `Bearer <YOUR_AILOAD_API_KEY>`

这里的 `<YOUR_AILOAD_API_KEY>` 是您在 **AI Load** 的 `config.yaml` 或 Admin UI 中配置的**全局密钥**或**分组密钥**，而非任何特定模型供应商的原始 API Key。

### 3. 主要端点

#### `/chat/completions`

这是最核心、最常用的端点，用于与所有支持的聊天模型进行对话。

- **Method**: `POST`
- **Path**: `/v1/chat/completions`
- **Headers**:

  - `Content-Type: application/json`
  - `Authorization: Bearer <YOUR_AILOAD_API_KEY>`

- **Body (请求体)**:

  请求体遵循 OpenAI `chat.completions` 的标准格式。

| 字段                 | 类型      | 是否必须 | 描述                                                                            |
| :------------------- | :-------- | :------- | :------------------------------------------------------------------------------ |
| `model`              | `string`  | **是**   | 您希望调用的模型名称。AI Load 会根据此名称自动路由到后台配置的正确模型供应商。  |
| `messages`           | `array`   | **是**   | 一个消息对象数组，用于描述对话上下文。                                          |
| `messages[].role`    | `string`  | **是**   | 消息发送者的角色，可以是 `user`, `assistant`, 或 `system`。                     |
| `messages[].content` | `string`  | **是**   | 消息的具体内容。                                                                |
| `stream`             | `boolean` | 否       | 是否使用流式传输。如果为 `true`，响应将以 Server-Sent Events (SSE) 的形式返回。 |
| `temperature`        | `number`  | 否       | 控制生成文本的随机性，介于 0 和 2 之间。                                        |
| `max_tokens`         | `integer` | 否       | 生成响应的最大 token 数量。                                                     |

- **`curl` 示例**:

  ```bash
  # 使用在 AI Load 中配置的密钥
  export AILOAD_API_KEY="your_configured_key"

  curl http://localhost:14004/v1/chat/completions \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $AILOAD_API_KEY" \
    -d '{
      "model": "gpt-4", # 可以是任何已在 AI Load 中配置的模型
      "messages": [{"role": "user", "content": "你好，AI Load！"}]
    }'
  ```

- **Python SDK 示例**:

  ```python
  from openai import OpenAI

  client = OpenAI(
      # 指向 AI Load 服务地址
      base_url="http://localhost:14004/v1",
      # 使用在 AI Load 中配置的密钥
      api_key="your_configured_key"
  )

  chat_completion = client.chat.completions.create(
      messages=[
          {
              "role": "user",
              "content": "请用 Python 写一个 Hello World。",
          }
      ],
      model="gpt-4", # 同样，可以是任何已配置的模型
  )

  print(chat_completion.choices[0].message.content)
  ```

---

### Embeddings API

通过 `embeddings` 接口，您可以将文本转换为固定长度的密集向量（embeddings），用于文本检索、聚类、分类等多种下游任务。AI Load 遵循 OpenAI 的 Embeddings API 格式。

**端点**

```
POST /v1/embeddings
```

**请求体 (Request Body)**

| 参数              | 类型            | 必需 | 描述                                                         |
| ----------------- | --------------- | ---- | ------------------------------------------------------------ |
| `model`           | string          | 是   | 用于生成嵌入向量的模型 ID。例如 `text-embedding-ada-002`。   |
| `input`           | string 或 array | 是   | 需要进行嵌入的输入文本。可以是单个字符串或字符串数 ​​ 组。   |
| `encoding_format` | string          | 否   | 嵌入向量的格式。可以是 `float` 或 `base64`。默认为 `float`。 |
| `user`            | string          | 否   | 代表最终用户的唯一标识符，用于帮助监控和检测滥用行为。       |

**`curl` 示例**

```bash
# 使用在 config.yaml 中配置的密钥
export AILOAD_API_KEY="your_configured_key"

curl http://localhost:14004/v1/embeddings \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $AILOAD_API_KEY" \
  -d '{
    "model": "text-embedding-ada-002",
    "input": "The food was delicious and the waiter..."
  }'
```

**Python SDK 示例**

```python
from openai import OpenAI

client = OpenAI(
    base_url="http://localhost:14004/v1",
    api_key="your_configured_key"
)

response = client.embeddings.create(
    model="text-embedding-ada-002",
    input="The food was delicious and the waiter...",
    encoding_format="float"
)

print(response.data[0].embedding)
```

---

### Rerank API

`rerank` 接口用于对一组文档，根据给定的查询进行相关性重排。它能显著提升检索结果的质量。

**端点**

```
POST /v1/rerank
```

**请求体 (Request Body)**

| 参数               | 类型    | 必需 | 描述                                                 |
| ------------------ | ------- | ---- | ---------------------------------------------------- |
| `model`            | string  | 是   | 用于执行重排任务的模型 ID。                          |
| `query`            | string  | 是   | 用于对文档进行排序的查询字符串。                     |
| `documents`        | array   | 是   | 一个包含多个待排序文档的字符串数 ​​ 组。             |
| `top_n`            | integer | 否   | 返回的重排后文档的数量。如果未指定，将返回所有文档。 |
| `return_documents` | boolean | 否   | 是否在响应中返回文档内容。默认为 `false`。           |

**`curl` 示例**

```bash
# 使用在 config.yaml 中配置的密钥
export AILOAD_API_KEY="your_configured_key"

curl http://localhost:14004/v1/rerank \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $AILOAD_API_KEY" \
  -d '{
    "model": "rerank-english-v2.0",
    "query": "What is the capital of France?",
    "documents": [
      "Paris is the capital of France.",
      "The Eiffel Tower is a famous landmark in Paris.",
      "France is a country in Western Europe."
    ],
    "top_n": 2
  }'
```

---

## 支持的模型 (示例)

您可以在 `model` 字段中指定任何已在 **AI Load** 中配置和路由的模型。以下是一些示例：

- **OpenAI**: `gpt-4`, `gpt-3.5-turbo`
- **Google Gemini**: `gemini-pro`
- **Anthropic Claude**: `claude-3-opus-20240229`, `claude-3-haiku-20240307`
- **Groq**: `llama3-8b-8192`, `mixtral-8x7b-32768`
- **Moonshot (Kimi)**: `moonshot-v1-8k`, `moonshot-v1-128k`
- **DeepSeek**: `deepseek-chat`, `deepseek-coder`
- **Siliconflow**: `Qwen/Qwen2-7B-Instruct`
- **Together AI**: `meta-llama/Llama-3-70b-chat-hf`
- **Zhipu AI (智谱)**: `glm-4`, `glm-3-turbo`
- **Ollama (本地模型)**: `llama3`, `qwen:7b`

**请注意**: 可用模型列表取决于您在 **AI Load** Admin UI 或配置文件中的具体设置。

---

## 🖼️ Image API

### 1. Create image

创建一个图片。

- **Endpoint:** `POST /v1/images/generations`
- **Description:** Creates an image given a prompt.

**Request Body:**

| Field    | Type    | Required | Description                                                                            |
| :------- | :------ | :------- | :------------------------------------------------------------------------------------- |
| `prompt` | string  | Yes      | A text description of the desired image(s). The maximum length is 1000 characters.     |
| `model`  | string  | No       | The model to use for image generation. Defaults to `dall-e-2`.                         |
| `n`      | integer | No       | The number of images to generate. Must be between 1 and 10.                            |
| `size`   | string  | No       | The size of the generated images. Must be one of `256x256`, `512x512`, or `1024x1024`. |

**Example Request:**

```shell
curl --request POST \
  --url {{baseURL}}/v1/images/generations \
  --header 'Authorization: Bearer {{token}}' \
  --header 'Content-Type: application/json' \
  --data '{
  "prompt": "A cute baby sea otter",
  "n": 2,
  "size": "1024x1024"
}'
```

```python
import openai

openai.api_key = "YOUR_API_KEY"

response = openai.Image.create(
  prompt="a white siamese cat",
  n=1,
  size="1024x1024"
)
image_url = response['data']['url']
print(image_url)
```

### 2. Create image edit

根据文本指令，对给定图像进行编辑。

- **Endpoint:** `POST /v1/images/edits`
- **Description:** Creates an edited or extended image given an original image and a prompt.

**Request Body (multipart/form-data):**

| Field    | Type    | Required | Description                                                                                |
| :------- | :------ | :------- | :----------------------------------------------------------------------------------------- |
| `image`  | file    | Yes      | The image to edit. Must be a valid PNG file, less than 4MB, and square.                    |
| `prompt` | string  | Yes      | A text description of the desired edits. The maximum length is 1000 characters.            |
| `mask`   | file    | No       | An additional image whose fully transparent areas indicate where `image` should be edited. |
| `n`      | integer | No       | The number of images to generate. Must be between 1 and 10.                                |
| `size`   | string  | No       | The size of the generated images. Must be one of `256x256`, `512x512`, or `1024x1024`.     |

**Example Request:**

```shell
curl --request POST \
  --url {{baseURL}}/v1/images/edits \
  --header 'Authorization: Bearer {{token}}' \
  --header 'Content-Type: multipart/form-data' \
  -F image=@otter.png \
  -F prompt="A cute baby sea otter wearing a beret"
```

### 3. Create image variation

创建给定图像的变体。

- **Endpoint:** `POST /v1/images/variations`
- **Description:** Creates a variation of a given image.

**Request Body (multipart/form-data):**

| Field   | Type    | Required | Description                                                                             |
| :------ | :------ | :------- | :-------------------------------------------------------------------------------------- |
| `image` | file    | Yes      | The image to create variations of. Must be a valid PNG file, less than 4MB, and square. |
| `n`     | integer | No       | The number of images to generate. Must be between 1 and 10.                             |
| `size`  | string  | No       | The size of the generated images. Must be one of `256x256`, `512x512`, or `1024x1024`.  |

**Example Request:**

```shell
curl --request POST \
  --url {{baseURL}}/v1/images/variations \
  --header 'Authorization: Bearer {{token}}' \
  --header 'Content-Type: multipart/form-data' \
  -F image=@otter.png
```

---

## 🎧 Audio API

### 1. Create speech

将文本生成为音频。

- **Endpoint:** `POST /v1/audio/speech`
- **Description:** Generates audio from the input text.

**Request Body:**

| Field   | Type   | Required | Description                                                                                                               |
| :------ | :----- | :------- | :------------------------------------------------------------------------------------------------------------------------ |
| `model` | string | Yes      | One of the available TTS models: `tts-1` or `tts-1-hd`.                                                                   |
| `input` | string | Yes      | The text to generate audio for. The maximum length is 4096 characters.                                                    |
| `voice` | string | Yes      | The voice to use when generating the audio. Supported voices are `alloy`, `echo`, `fable`, `onyx`, `nova`, and `shimmer`. |
| `speed` | number | No       | The speed of the generated audio. Select a value from `0.25` to `4.0`. `1.0` is the default.                              |

**Example Request:**

```shell
curl --request POST \
  --url {{baseURL}}/v1/audio/speech \
  --header 'Authorization: Bearer {{token}}' \
  --header 'Content-Type: application/json' \
  --data '{
  "model": "tts-1",
  "input": "The quick brown fox jumped over the lazy dog.",
  "voice": "alloy"
}'
```

### 2. Create transcription

将音频转录为文本。

- **Endpoint:** `POST /v1/audio/transcriptions`
- **Description:** Transcribes audio into the input language.

**Request Body (multipart/form-data):**

| Field      | Type   | Required | Description                                                                                                                |
| :--------- | :----- | :------- | :------------------------------------------------------------------------------------------------------------------------- |
| `file`     | file   | Yes      | The audio file object to transcribe. The file format must be one of `mp3`, `mp4`, `mpeg`, `mpga`, `m4a`, `wav`, or `webm`. |
| `model`    | string | Yes      | ID of the model to use. Only `whisper-1` is currently available.                                                           |
| `prompt`   | string | No       | An optional text to guide the model's style or continue a previous audio segment.                                          |
| `language` | string | No       | The language of the input audio in ISO-639-1 format.                                                                       |

**Example Request:**

```shell
curl --request POST \
  --url {{baseURL}}/v1/audio/transcriptions \
  --header 'Authorization: Bearer {{token}}' \
  --header 'Content-Type: multipart/form-data' \
  -F file=@speech.mp3 \
  -F model=whisper-1
```

### 3. Create translation

将音频翻译为英语文本。

- **Endpoint:** `POST /v1/audio/translations`
- **Description:** Translates audio into English.

**Request Body (multipart/form-data):**

| Field    | Type   | Required | Description                                                                                                               |
| :------- | :----- | :------- | :------------------------------------------------------------------------------------------------------------------------ |
| `file`   | file   | Yes      | The audio file object to translate. The file format must be one of `mp3`, `mp4`, `mpeg`, `mpga`, `m4a`, `wav`, or `webm`. |
| `model`  | string | Yes      | ID of the model to use. Only `whisper-1` is currently available.                                                          |
| `prompt` | string | No       | An optional text to guide the model's style.                                                                              |

**Example Request:**

```shell
curl --request POST \
  --url {{baseURL}}/v1/audio/translations \
  --header 'Authorization: Bearer {{token}}' \
  --header 'Content-Type: multipart/form-data' \
  -F file=@german_speech.mp3 \
  -F model=whisper-1
```

---

## 🎬 Video API

### 1. Create video generation

根据文本提示生成视频。

- **Endpoint:** `POST /v1/video/generations`
- **Description:** Creates a video given a prompt.

**Request Body:**

| Field      | Type    | Required | Description                                                                     |
| :--------- | :------ | :------- | :------------------------------------------------------------------------------ |
| `prompt`   | string  | Yes      | A text description of the desired video. The maximum length is 4000 characters. |
| `model`    | string  | No       | The model to use for video generation. e.g., `sora-1`.                          |
| `n`        | integer | No       | The number of videos to generate. Defaults to 1.                                |
| `size`     | string  | No       | The dimensions of the video. e.g., `1920x1080`.                                 |
| `duration` | integer | No       | The duration of the video in seconds.                                           |

**Example Request:**

```shell
curl --request POST \
  --url {{baseURL}}/v1/video/generations \
  --header 'Authorization: Bearer {{token}}' \
  --header 'Content-Type: application/json' \
  --data '{
  "prompt": "A cinematic shot of a puppy playing in the autumn leaves.",
  "model": "sora-1",
  "n": 1,
  "size": "1920x1080",
  "duration": 15
}'
```

```python
import openai

openai.api_key = "YOUR_API_KEY"

response = openai.Video.create(
  prompt="A majestic eagle soaring over the mountains.",
  model="sora-1"
)
video_url = response['data']['url']
print(video_url)
```

---

## Models

### List models

Retrieves a list of available models.

**GET** `/v1/models`

#### Headers

- `Authorization: Bearer {{token}}`

#### Response

```json
{
  "object": "list",
  "data": [
    {
      "id": "dall-e-3",
      "object": "model",
      "created": 1698785685,
      "owned_by": "system"
    },
    {
      "id": "tts-1-hd",
      "object": "model",
      "created": 1699053533,
      "owned_by": "system"
    },
    {
      "id": "whisper-1",
      "object": "model",
      "created": 1677532384,
      "owned_by": "system"
    },
    {
      "id": "sora-1",
      "object": "model",
      "created": 1708623111,
      "owned_by": "system"
    }
  ]
}
```
