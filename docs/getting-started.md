# 🚀 快速上手指南 (Getting Started)

**1. 预备环境 (Prerequisites):**

- 列出运行 AI Load 所需的最低环境要求。
- **Go**: `1.21` 或更高版本。
- **Node.js**: `18.x` 或更高版本（用于前端开发）。
- **Git**: 用于克隆项目。

**2. 下载与安装 (Download and Installation):**

- **步骤 1: 克隆仓库**
  - 提供使用 `git clone` 命令从 GitHub 克隆项目的完整指令。
  ```bash
  git clone https://github.com/lazygophers/aiload.git
  cd aiload
  ```
- **步骤 2: 构建后端**
  - 提供编译和构建 Go 后端服务的命令。
  ```bash
  go build -o aiload_server ./cmd/server
  ```
- **步骤 3: 构建前端**
  - 提供进入 `assets` 目录、安装依赖并打包前端应用的完整命令。
  ```bash
  cd assets
  npm install
  npm run build
  cd ..
  ```

**3. 基础配置 (Basic Configuration):**

- 解释需要创建一个基础的配置文件，例如 `config.yaml`。
- 提供一个最简化的 `config.yaml` 示例，包含服务端口和至少一个 AI 模型的上游配置。

  ```yaml
  server:
    port: 8080

  auth:
    # 全局密钥，适用于所有未指定分组密钥的请求
    keys:
      - "sk-global-xxxxxxxxxxxxxxxxxxxxxxxx"

  upstreams:
    - name: "openai_default"
      # AI 服务类型
      provider: "openai"
      # 负载均衡权重
      weight: 100
      # 你的 OpenAI API 密钥
      api_key: "sk-your-openai-api-key"
  ```

- 提醒用户将示例中的 `api_key` 替换为自己的真实密钥。

**4. 启动服务 (Launch the Service):**

- 提供启动后端服务的命令。
  ```bash
  ./aiload_server --config config.yaml
  ```
- 描述服务成功启动后，在终端会看到的提示信息（例如：`Server started at :8080`）。

**5. 验证服务 (Verify the Service):**

- 提供一个使用 `curl` 命令向代理服务发起请求的示例，以验证服务是否正常工作。
  ```bash
  curl http://localhost:8080/v1/chat/completions \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer sk-global-xxxxxxxxxxxxxxxxxxxxxxxx" \
    -d '{
      "model": "gpt-4",
      "messages": [{"role": "user", "content": "你好，世界！"}]
    }'
  ```
- 解释预期的成功响应是什么样的。
