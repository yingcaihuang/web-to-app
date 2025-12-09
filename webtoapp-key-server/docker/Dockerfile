# 使用官方 Go 基础镜像
FROM golang:1.24 AS builder

# 安装依赖
RUN apt-get update && apt-get install -y --no-install-recommends \
    git gcc sqlite3 libsqlite3-dev \
    && rm -rf /var/lib/apt/lists/*

# 设置工作目录
WORKDIR /build

# 复制 go.mod 和 go.sum
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 编译应用
RUN CGO_ENABLED=1 GOOS=linux go build -o keyserver ./cmd/main.go

# 最终镜像
FROM debian:bookworm-slim

# 安装 CA 证书和 SQLite 库
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates libsqlite3-0 \
    && rm -rf /var/lib/apt/lists/*

# 设置工作目录
WORKDIR /app

# 创建数据目录
RUN mkdir -p /app/data

# 从构建镜像复制二进制文件
COPY --from=builder /build/keyserver .

# 复制配置文件
COPY configs ./configs

# 复制前端文件
COPY web ./web

# 复制初始化脚本
COPY docker/entrypoint.sh ./entrypoint.sh
RUN chmod +x ./entrypoint.sh

# 暴露端口
EXPOSE 8080

# 启动应用（通过初始化脚本）
ENTRYPOINT ["./entrypoint.sh"]
