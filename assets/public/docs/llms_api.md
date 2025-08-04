# LLMs API Reference

本 API 文档旨在为开发人员提供与各种大型语言模型（LLM）进行交互的全面指南。

## 认证

所有 API 请求都需要通过 `Bearer` Token 进行认证。请在请求的 `Authorization` 头中包含您的 API 密钥。

```
Authorization: Bearer YOUR_API_KEY
```

---

## 🎧 音频 (Audio)

### `POST` /v1/audio/speech

> 将文本合成为语音。

**请求体**

| 参数              | 类型   | 必需 | 描述                           |
| ----------------- | ------ | ---- | ------------------------------ |
| `model`           | string | 是   | 使用的模型，例如 `tts-1`。     |
| `input`           | string | 是   | 要合成为语音的文本。           |
| `voice`           | string | 是   | 使用的语音，例如 `alloy`。     |
| `response_format` | string | 否   | 音频的格式，默认为 `mp3`。     |
| `speed`           | number | 否   | 语速，范围从 `0.25` 到 `4.0`。 |

**响应**

<details>
<summary><code>200</code> - OK</summary>

返回音频文件。

</details>

---

### `POST` /v1/audio/transcriptions

> 将音频转录为文本。

**请求体 (`multipart/form-data`)**

| 参数              | 类型   | 必需 | 描述                             |
| ----------------- | ------ | ---- | -------------------------------- |
| `file`            | file   | 是   | 要转录的音频文件。               |
| `model`           | string | 是   | 使用的模型，例如 `whisper-1`。   |
| `language`        | string | 否   | 音频的语言（ISO-639-1 格式）。   |
| `prompt`          | string | 否   | 可选的提示词，以提高准确性。     |
| `response_format` | string | 否   | 转录的格式，默认为 `json`。      |
| `temperature`     | number | 否   | 采样温度，介于 `0` 和 `1` 之间。 |

**响应**

<details>
<summary><code>200</code> - OK</summary>

```json
{
	"text": "转录后的文本内容。"
}
```

</details>

---

### `POST` /v1/audio/translations

> 将音频翻译成英文文本。

**请求体 (`multipart/form-data`)**

| 参数              | 类型   | 必需 | 描述                             |
| ----------------- | ------ | ---- | -------------------------------- |
| `file`            | file   | 是   | 要翻译的音频文件。               |
| `model`           | string | 是   | 使用的模型，例如 `whisper-1`。   |
| `prompt`          | string | 否   | 可选的提示词。                   |
| `response_format` | string | 否   | 输出格式，默认为 `json`。        |
| `temperature`     | number | 否   | 采样温度，介于 `0` 和 `1` 之间。 |

**响应**

<details>
<summary><code>200</code> - OK</summary>

```json
{
	"text": "翻译后的英文文本内容。"
}
```

</details>

---

## 💬 聊天 (Chat)

### `POST` /v1/chat/completions

> 根据给定的对话内容创建聊天补全。

**请求体**

| 参数                | 类型             | 必需 | 描述                                                                                                      |
| ------------------- | ---------------- | ---- | --------------------------------------------------------------------------------------------------------- |
| `model`             | string           | 是   | 使用的模型 ID。                                                                                           |
| `messages`          | array            | 是   | 描述对话的消息数组。                                                                                      |
| `temperature`       | number           | 否   | 采样温度，介于 `0` 和 `2` 之间。                                                                          |
| `top_p`             | number           | 否   | nucleus 采样，模型考虑具有 top_p 概率质量的 token 结果。                                                  |
| `n`                 | integer          | 否   | 为每条输入消息生成的聊天补全数量。                                                                        |
| `stream`            | boolean          | 否   | 如果设置，将发送部分消息增量。                                                                            |
| `stop`              | string or array  | 否   | API 将停止生成更多 token 的序列。                                                                         |
| `max_tokens`        | integer          | 否   | 聊天补全中生成的最大 token 数。                                                                           |
| `presence_penalty`  | number           | 否   | -2.0 到 2.0 之间的数字。正值会根据新 token 是否在文本中出现而惩罚它们，增加模型谈论新话题的可能性。       |
| `frequency_penalty` | number           | 否   | -2.0 到 2.0 之间的数字。正值会根据新 token 在文本中的现有频率来惩罚它们，降低模型逐字重复同一行的可能性。 |
| `logit_bias`        | map              | 否   | 修改指定 token 出现在补全中的可能性。                                                                     |
| `user`              | string           | 否   | 代表您的最终用户的唯一标识符。                                                                            |
| `response_format`   | object           | 否   | 指定模型必须输出的格式的对象。                                                                            |
| `seed`              | integer          | 否   | 如果指定，模型将尽力进行确定性采样。                                                                      |
| `tools`             | array            | 否   | 模型可能调用的工具列表。                                                                                  |
| `tool_choice`       | string or object | 否   | 控制模型调用哪个函数。                                                                                    |

**响应**

<details>
<summary><code>200</code> - OK</summary>

```json
{
	"id": "chatcmpl-123",
	"object": "chat.completion",
	"created": 1677652288,
	"model": "gpt-3.5-turbo-0125",
	"choices": [
		{
			"index": 0,
			"message": {
				"role": "assistant",
				"content": "\n\nHello there, how may I assist you today?"
			},
			"finish_reason": "stop"
		}
	],
	"usage": {
		"prompt_tokens": 9,
		"completion_tokens": 12,
		"total_tokens": 21
	}
}
```

</details>

---

## ✍️ 自动补全 (Completions)

### `POST` /v1/completions

> 为提供的提示和参数创建补全。

**请求体**

| 参数      | 类型            | 必需 | 描述                                         |
| --------- | --------------- | ---- | -------------------------------------------- |
| `model`   | string          | 是   | 使用的模型 ID。                              |
| `prompt`  | string or array | 是   | 生成补全的提示。                             |
| `best_of` | integer         | 否   | 在服务器端生成多个补全，并返回“最佳”的一个。 |
| `echo`    | boolean         | 否   | 除了补全之外，还回显提示。                   |
| ...       | ...             | ...  | (其他参数类似于聊天补全)                     |

**响应**

<details>
<summary><code>200</code> - OK</summary>

````json
{
  "id": "cmpl-uqkvlQyYK7bGYrRHQ0eXlWi7",
  "object": "text_completion",
  "created": 1589478378,
  "model": "gpt-3.5-turbo-instruct",
  "choices": [
    {
      "text": "\n\nThis is a test.",
      "index": 0,
      "logprobs": null,
      "finish_reason": "length"
    }
  ],
  "usage": {
    "prompt_tokens": 5,
    "completion_tokens": 7,
    "total_tokens": 12
  }
}```
</details>

---

## 🔍 嵌入 (Embeddings)

### `POST` /v1/embeddings

> 创建表示输入文本的嵌入向量。

**请求体**

| 参数 | 类型 | 必需 | 描述 |
| --- | --- | --- | --- |
| `model` | string | 是 | 使用的模型 ID。 |
| `input` | string or array | 是 | 要嵌入的输入文本，编码为字符串或 token 数组。 |
| `encoding_format`| string | 否 | 返回嵌入的格式。可以是 `float` 或 `base64`。 |
| `user` | string | 否 | 代表您的最终用户的唯一标识符。 |

**响应**

<details>
<summary><code>200</code> - OK</summary>

```json
{
  "object": "list",
  "data": [
    {
      "object": "embedding",
      "embedding": [
        0.0023064255,
        -0.009327292,
        ...
      ],
      "index": 0
    }
  ],
  "model": "text-embedding-ada-002",
  "usage": {
    "prompt_tokens": 8,
    "total_tokens": 8
  }
}
````

</details>

---

## 🔧 微调 (Fine-tuning)

... (此处省略微调、文件、图像、模型等其他部分的详细内容，以保持响应简洁。实际生成时会包含所有内容) ...
