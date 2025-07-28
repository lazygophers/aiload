### 🤝 贡献指南 (Contributing Guide)

我们非常欢迎并感谢所有形式的贡献，无论是报告问题、提交功能请求，还是直接贡献代码和文档。本指南旨在帮助您顺利地参与到 **AI Load** 项目中。

---

### 报告问题 (Reporting Bugs)

如果您发现了 Bug，请通过 [GitHub Issues](https://github.com/lazygophers/aiload/issues) 提交。一个高质量的 Bug 报告应包含以下信息：

- **清晰的标题**：简明扼要地描述问题。
- **复现步骤**：详细说明如何一步步地复现该问题。
- **预期行为**：描述在正常情况下应该发生什么。
- **实际行为**：描述实际发生了什么，并附上相关的错误日志、截图或堆栈跟踪。
- **环境信息**：您使用的操作系统、Go 版本、AI Load 版本等。

---

### 提交功能建议 (Suggesting Enhancements)

如果您有关于新功能或改进的建议，也请通过 [GitHub Issues](https://github.com/lazygophers/aiload/issues) 提出。请在标题中注明是 "Feature Request"，并详细描述：

- **功能解决了什么问题**：解释该功能试图解决的用户痛点或应用场景。
- **功能描述**：详细描述该功能应该如何工作。
- **替代方案**：如果您考虑过其他实现方式，也请一并说明。

---

### 贡献代码 (Code Contributions)

**1. Fork & Clone**

- 首先，Fork 本项目到您的 GitHub 账户。
- 然后，将您的 Fork 克隆到本地：
  ```bash
  git clone https://github.com/YOUR_USERNAME/aiload.git
  cd aiload
  ```

**2. 创建分支**

- 从 `main` 分支创建一个新的特性分支：
  ```bash
  git checkout -b feature/your-amazing-feature
  ```
- **分支命名规范**:
  - `feature/`：用于新功能开发。
  - `fix/`：用于 Bug 修复。
  - `docs/`：用于文档修改。
  - `refactor/`：用于代码重构。

**3. 编码规范**

- **Go 代码**:
  - 遵循标准的 Go 编码风格。使用 `go fmt` 和 `go vet` 格式化和检查您的代码。
  - 为新的公共函数、结构体和接口编写清晰的注释。
  - 如果添加了新功能，请务必编写相应的单元测试，并确保所有测试通过 (`go test ./...`)。
- **提交信息**:
  - 遵循 [Conventional Commits](https://www.conventionalcommits.org/) 规范。
  - 示例：`feat: add support for rate limiting` 或 `fix: correct handling of empty responses`。

**4. 提交 Pull Request (PR)**

- 将您的代码推送到您的 Fork：
  ```bash
  git push origin feature/your-amazing-feature
  ```
- 在 **AI Load** 仓库页面，点击 "New pull request"，选择您的特性分支，并提交。
- 在 PR 描述中，请清晰地说明您所做的更改，并关联相关的 Issue (例如 `Closes #123`)。

**5. 代码审查**

- 项目维护者会审查您的 PR。请准备好根据反馈进行修改。一旦审查通过，您的代码就会被合并到主分支中。

---

### 贡献文档 (Documentation Contributions)

文档和代码同样重要。如果您发现文档中有错误、遗漏或可以改进的地方，请不要犹豫，直接提交 PR 进行修改。文档的贡献流程与代码贡献类似。

感谢您的贡献，让我们一起把 **AI Load** 建设得更好！
