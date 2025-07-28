### 🚀 部署指南 (Deployment Guide)

本文档提供了将 **AI Load** 服务部署到生产环境的多种方法。

**前置条件:**

- 一台拥有网络访问权限的服务器或本地机器。
- 准备好的 `config.yaml` 配置文件。

---

### 方式一：直接运行二进制文件 (推荐)

这是最简单直接的部署方式，适合快速启动和大多数常见场景。

**1. 下载预编译的二进制文件**

- 从项目的 [GitHub Releases](https://github.com/lazygophers/aiload/releases) 页面下载与您服务器操作系统和架构相匹配的最新版本。
- 例如，对于 `Linux x86_64`，您可以下载 `aiload-linux-amd64`。

**2. 上传文件**

- 将下载的二进制文件和您准备好的 `config.yaml` 上传到服务器的同一个目录下，例如 `/opt/aiload/`。

**3. 赋予执行权限**

```bash
chmod +x aiload-linux-amd64
```

**4. 启动服务**

- 直接在前台启动（用于测试）：
  ```bash
  ./aiload-linux-amd64 --config config.yaml
  ```
- **(推荐)** 使用 `nohup` 在后台持久化运行：
  ```bash
  nohup ./aiload-linux-amd64 --config config.yaml &
  ```
  使用 `nohup` 在后台持久化运行，日志将直接输出到当前控制台。

**5. (可选) 使用 Systemd 进行管理**

- 为了实现开机自启和更专业的服务管理，建议创建一个 `systemd` 服务文件。
- 创建 `/etc/systemd/system/aiload.service`:

  ```ini
  [Unit]
  Description=AI Load Proxy Service
  After=network.target

  [Service]
  Type=simple
  ExecStart=/opt/aiload/aiload-linux-amd64 --config /opt/aiload/config.yaml
  Restart=on-failure
  RestartSec=5s

  [Install]
  WantedBy=multi-user.target
  ```

- 重新加载 `systemd` 并启动服务：
  ```bash
  sudo systemctl daemon-reload
  sudo systemctl start aiload
  sudo systemctl enable aiload  # 设置开机自启
  sudo systemctl status aiload # 查看服务状态
  ```

---

### 方式二：从源码构建

如果您想使用最新的未发布功能或进行自定义修改，可以从源码构建。

**1. 安装 Go 环境**

- 请确保您已安装 Go 1.21 或更高版本。

**2. 克隆仓库**

```bash
git clone https://github.com/lazygophers/aiload.git
cd aiload
```

**3. 构建后端**

```bash
go build -o aiload_server ./cmd/server
```

这将在项目根目录下生成一个名为 `aiload_server` 的可执行文件。

**4. 构建前端 (如果需要)**

- 如果您修改了前端代码 (`assets` 目录)，需要重新构建并嵌入。
  ```bash
  # (进入 assets 目录安装依赖并构建)
  cd assets
  npm install
  npm run build
  cd ..
  # (使用 go:embed 重新打包)
  # ... 具体指令取决于项目实现
  ```

**5. 运行**

- 构建完成后，即可按照**方式一**中的步骤运行 `aiload_server` 二进制文件。

---

请确保所有路径、命令和文件名都清晰准确，并为不同技术水平的用户提供易于理解的指导。
