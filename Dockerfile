# 使用官方 Golang 镜像作为构建环境
FROM golang:1.23 as builder

# 设置工作目录
WORKDIR /app

# 复制项目文件
COPY . .

# 安装依赖
RUN go mod download

# 构建应用
RUN CGO_ENABLED=0 go build -o /openai-forward main.go

# 使用轻量级 Alpine 镜像作为运行环境
FROM alpine:latest

# 安装 deb-init 安装必要的 CA 证书
RUN apk add --no-cache dumb-init ca-certificates

# 设置工作目录
WORKDIR /app

# 复制构建好的二进制文件和配置文件
COPY --from=builder /openai-forward .
COPY ./webroot /app/webroot

ENV OPENAI_API_KEY= \
    OPENAI_ORG_ID= \
    OPENAI_PROJECT_ID= \
    OPENAI_TARGET_BASE_URL=https://api.openai.com \
    PROXY_LISTEN_ADDR=:8080 \
    PROXY_WEB_ROOT=/app/webroot \
    PROXY_LOG_LEVEL=info

# 暴露代理服务端口
EXPOSE 8080


ENTRYPOINT ["/usr/bin/dumb-init", "--"]

# 启动应用
CMD ["./openai-forward"]