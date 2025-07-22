# OpenAI API 代理服务
[![GitHub license](https://img.shields.io/github/license/mmhk/openai-forward?style=flat-square)](https://github.com/mmhk/openai-forward/blob/main/LICENSE)
[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/mmhk/openai-forward/go.yml?branch=main&style=flat-square)](https://github.com/mmhk/openai-forward/actions/workflows/go.yml)
[![Docker Image Version (latest by date)](https://img.shields.io/docker/v/mmhk/openai-forward?style=flat-square)](https://hub.docker.com/r/mmhk/openai-forward)
[![Go Report Card](https://goreportcard.com/badge/github.com/mmhk/openai-forward?style=flat-square)](https://goreportcard.com/report/github.com/mmhk/openai-forward)
[![Go Version](https://img.shields.io/badge/go-1.21-blue?style=flat-square&logo=go)](https://go.dev/dl/)

一个轻量级的代理服务，用于转发所有 OpenAI API 请求。

## 项目特点
- 轻量级服务
- 高效转发所有 OpenAI API 接口
- 支持通过 `.env` 文件或环境变量配置 API 密钥、组织 ID 和项目 ID
- 最终编译成 Docker 镜像，支持多服务器部署
- 开发过程中引入日志处理，通过日志级别控制控制台日志输出

## 安装与使用

### 本地开发环境

1. 安装 Go 1.21
2. 克隆项目仓库
3. 安装依赖:
   ```bash
   go mod tidy
   ```
4. 启动服务:
   ```bash
   go run main.go
   ```

### Docker 部署

1. 构建 Docker 镜像:
   ```bash
   docker build -t openai-forward .
   ```
2. 运行 Docker 容器:
   ```bash
   docker run -d -p 8080:8080 openai-forward
   ```

## 配置说明

- `OPENAI_API_KEY`: OpenAI 的 API 密钥
- `OPENAI_ORG_ID`: OpenAI 的组织 ID (可选)
- `OPENAI_PROJECT_ID`: OpenAI 的项目 ID (可选)
- `PROXY_LISTEN_ADDR`: 代理服务监听地址 (默认: `:8080`)
- `PROXY_LOG_LEVEL`: 日志级别 (默认: `info`, 可选: `debug`)

## 目录结构
```
openai-forward/
├── main.go
├── go.mod
├── .env
├── config/
│   └── config.go
├── proxy/
│   └── proxy.go
├── logging/
│   └── logger.go
├── Dockerfile
└── README.md
```

## 开发计划
请参考 [development_plan.md](development_plan.md) 文件。