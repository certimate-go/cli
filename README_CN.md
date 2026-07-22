# certimate-cli

[English](./README.md) | 简体中文

[Certimate](https://github.com/certimate-go/certimate) SSL 证书管理的命令行工具。

[![Go Report Card](https://goreportcard.com/badge/github.com/certimate-go/cli)](https://goreportcard.com/report/github.com/certimate-go/cli)
[![License](https://img.shields.io/github/license/certimate-go/cli)](LICENSE)

## 特性

- **Agent 优先设计**：内置 SKILL.md 集成，支持 Claude Code、OpenClaw 等各种 AI Agent
- 工作流管理：列出、查看和执行证书工作流
- 证书管理：查看和下载多种格式的 SSL 证书
- 访问凭证管理：创建、编辑和删除服务商凭证
- 多种输出格式：JSON（默认）和表格输出
- 多配置支持：管理多个服务器配置

## 安装

### 从源码安装

```bash
go install github.com/certimate-go/cli@latest
```

### 下载二进制文件

从 [GitHub Releases](https://github.com/certimate-go/cli/releases) 下载最新版本。

### 使用 Homebrew

```bash
brew install certimate-go/tap/certimate
```

## 快速开始

### 1. 配置连接

```bash
certimate config set --server http://127.0.0.1:8090 --token YOUR_API_TOKEN
```

从 PocketBase 管理后台获取 API Token：
1. 打开 PocketBase 管理后台（如 `http://127.0.0.1:8090/_/`）
2. 进入 Collections > `_superusers` > 选择超级用户
3. 点击 "Impersonate" 下拉菜单生成不可续期的认证 Token

> **注意：** PocketBase 使用超级用户模拟 Token 而非传统的 API Key。
> 详见 [PocketBase 认证文档](https://pocketbase.io/docs/authentication/#api-keys)。

### 2. 列出工作流

```bash
certimate workflow list
```

### 3. 执行工作流

```bash
# 异步执行
certimate workflow run WORKFLOW_ID

# 等待完成
certimate workflow run WORKFLOW_ID --wait
```

### 4. 列出证书

```bash
certimate certificate list

# 表格输出
certimate certificate list --output table
```

## 命令

### 配置

```bash
# 设置配置
certimate config set --server URL --token TOKEN [--profile NAME]

# 查看当前配置
certimate config get

# 显示当前配置文件
certimate config current

# 列出所有配置文件
certimate config list
```

### 工作流

```bash
# 列出工作流
certimate workflow list [--filter EXPR] [--limit N]

# 获取工作流详情
certimate workflow get WORKFLOW_ID

# 执行工作流
certimate workflow run WORKFLOW_ID [--wait] [--timeout SECONDS]

# 取消运行中的工作流
certimate workflow cancel WORKFLOW_ID RUN_ID

# 列出执行历史
certimate workflow runs WORKFLOW_ID [--limit N]
```

### 证书

```bash
# 列出证书
certimate certificate list [--filter EXPR] [--limit N]

# 获取证书详情
certimate certificate get CERTIFICATE_ID

# 下载证书
certimate certificate download CERTIFICATE_ID --format PEM|PFX|JKS [--dest FILE]
```

### 访问凭证

```bash
# 列出访问凭证
certimate access list [--reveal]

# 获取访问凭证详情
certimate access get ACCESS_ID [--reveal]

# 创建新访问凭证
certimate access create --name NAME --provider PROVIDER --config JSON

# 编辑访问凭证
certimate access edit ACCESS_ID --name NAME --config JSON

# 删除访问凭证
certimate access delete ACCESS_ID
```

### 其他命令

```bash
# 显示版本
certimate version

# 生成 Shell 补全
certimate completion bash|zsh|fish
```

## 全局参数

| 参数 | 说明 |
|------|------|
| `--config` | 配置文件路径（默认 `~/.config/certimate-cli/config.yaml`） |
| `--debug` | 启用调试输出 |
| `-o, --output` | 输出格式：`json` 或 `table`（默认 `json`） |
| `--profile` | 配置文件名（默认 `default`） |

## 输出格式

所有命令支持 JSON（默认）和表格输出：

```bash
# JSON 输出（默认）
certimate workflow list

# 表格输出
certimate workflow list --output table
```

## 环境变量

| 变量 | 说明 |
|------|------|
| `CERTIMATE_CLI_TOKEN` | API Token（优先级高于配置文件） |
| `CERTIMATE_CLI_SERVER` | 服务器 URL（优先级高于配置文件） |
| `CERTIMATE_CLI_CONFIG_DIR` | 自定义配置目录 |

## 退出码

| 退出码 | 含义 |
|--------|------|
| 0 | 成功 |
| 1 | 一般错误 |
| 2 | 参数无效 |
| 3 | 认证错误 |
| 4 | 网络错误 |
| 5 | 未找到 |

## AI Agent Skills

Certimate CLI **专为 AI Agent 设计**。内置 Agent Skills（SKILL.md 文件），可与各种 AI Agent 框架无缝集成：

- **Claude Code** - Anthropic 的 CLI Agent
- **OpenClaw** - 开源 Agent 框架
- **任何 Claw 兼容的 Agent** - Skills 使用标准元数据格式

### 可用 Skills

| Skill | 说明 |
|-------|------|
| `ctm-shared` | 通用模式、认证和 CLI 设置 |
| `ctm-workflow` | 证书工作流操作（列出、运行、取消） |
| `ctm-workflow-create` | 创建和配置新工作流 |
| `ctm-certificate` | 证书查看和下载 |
| `ctm-access` | 服务商凭证管理 |

### 安装

#### Claude Code

```bash
# 一次性安装所有 skills
npx skills add https://github.com/certimate-go/cli

# 或按需安装
npx skills add https://github.com/certimate-go/cli/tree/main/skills/ctm-workflow
```

#### OpenClaw / 其他 Claw 兼容 Agent

Skills 包含 `openclaw` 元数据，支持自动发现：

```yaml
metadata:
  openclaw:
    category: "devops"
    requires:
      bins: ["certimate"]
```

将 skills 目录克隆或软链接到你的 Agent skills 路径即可。

#### 手动安装

```bash
cp -r skills/* ~/.claude/skills/
```

### 为什么使用 Agent Skills？

- **零学习成本**：AI Agent 立即理解 CLI 能力
- **类型安全操作**：Skills 记录精确的命令语法和输出
- **错误处理**：内置常见问题的故障排除模式
- **多格式输出**：默认 JSON 便于机器解析，表格便于人类阅读

## 开发

```bash
# 构建
make build

# 本地安装
make install

# 运行测试
make test

# 构建所有平台
make build-all

# 生成 Shell 补全
make completion
```

## 许可证

[MIT License](LICENSE)

## 相关项目

- [Certimate](https://github.com/certimate-go/certimate) - SSL 证书管理系统
