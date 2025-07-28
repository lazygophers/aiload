# 🚀 AI Load - 下一代 AI 接口透明代理服务

[![Go Report Card](https://goreportcard.com/badge/github.com/lazygophers/aiload)](https://goreportcard.com/report/github.com/lazygophers/aiload)
[![Build Status](https://img.shields.io/github/actions/workflow/status/lazygophers/aiload/go.yml?branch=main)](https://github.com/lazygophers/aiload/actions)
[![License](https://img.shields.io/github/license/lazygophers/aiload)](https://github.com/lazygophers/aiload/blob/main/LICENSE)
[![Release](https://img.shields.io/github/v/release/lazygophers/aiload)](https://github.com/lazygophers/aiload/releases)

**AI Load** 是一款专为需要集成多种 AI 服务的企业和开发者设计的高性能、高可用的 AI 接口透明代理服务。它以**零侵入、高效率**为核心，帮助您轻松驾驭复杂的 AI 服务矩阵。

---

## ✨ 功能特性

- **🔄 透明代理**: **完全保留**原生 API 格式，无缝对接 OpenAI、Google Gemini、Anthropic Claude、**Siliconflow** 以及本地运行的 **Ollama** 等多种服务，无需修改现有代码。
- **🔑 智能密钥管理**: 高性能**密钥池**，支持分组管理、自动轮换和故障恢复，告别手动维护的烦恼。
- **⚖️ 负载均衡**: 支持多上游端点的**加权负载均衡**，智能分配请求，显著提升服务可用性和稳定性。
- **🛡️ 智能故障处理**: 自动化的**密钥黑名单**和恢复机制，主动规避故障节点，确保业务连续性。
- **⚙️ 动态配置**: 系统设置和分组配置**支持热重载**，修改配置无需重启服务，即刻生效。
- **🏢 企业级架构**: 专为生产环境设计，支持**分布式主从部署**，轻松实现水平扩展和高可用。
- **🖥️ 现代化管理**: 基于 **React + Ant Design** 的现代化 Web 管理界面，所有操作直观易用。
- **📊 全面监控**: 提供实时**统计面板**、健康检查和详细的请求日志，让服务状态尽在掌握。
- **⚡ 高性能设计**: 采用**零拷贝流式传输**、连接池复用和原子操作等技术，最大化处理性能。
- **🔐 双重认证体系**: **管理端与代理端认证分离**，代理认证支持全局和分组级别密钥，安全可控。

## 🤖 支持的 AI 服务

AI Load 作为透明代理服务，完整保留了各大 AI 服务商的原生 API 格式，包括：

- **OpenAI 格式**: 官方 OpenAI API、Azure OpenAI 及其他兼容服务。
- **Google Gemini 格式**: Gemini Pro、Gemini Pro Vision 等原生 API。
- **Anthropic Claude 格式**: Claude 系列模型的高质量对话与文本生成 API。
- **Siliconflow 格式**: 兼容 OpenAI 格式的 Siliconflow 云端服务。
- **Ollama (本地)**: 支持在本地环境中运行的 Ollama 模型，保障数据私密性。

## 🛠️ 技术栈

- **后端**: Golang
- **前端**: React + Ant Design

## 🚀 快速开始

请参阅我们的 **[快速上手指南](./docs/getting-started.md)** 来快速部署和使用 AI Load。

## 📚 文档

更详细的文档请访问 [docs](./docs) 目录。

- [项目介绍](./docs/introduction.md)
- [部署指南](./docs/deployment.md)
- [配置说明](./docs/configuration.md)
- [API 参考](./docs/api-reference)

## 🤝 贡献

我们欢迎任何形式的贡献！请阅读 **[贡献指南](./docs/contributing.md)** 来了解如何参与项目。

## 📄 许可证

本项目基于 [MIT](LICENSE) 许可证。