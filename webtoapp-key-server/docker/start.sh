#!/bin/bash

###############################################################################
# WebToApp Key Server - Docker Compose 启动脚本
# 自动生成 SSL 证书并启动所有服务
###############################################################################

set -e

PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$PROJECT_DIR"

echo ""
echo "╔════════════════════════════════════════════════════════════════╗"
echo "║     WebToApp Key Server - Docker Compose 启动                 ║"
echo "╚════════════════════════════════════════════════════════════════╝"
echo ""

# 1. 检查依赖
echo "1️⃣  检查依赖..."
if ! command -v docker &> /dev/null; then
    echo "✗ Docker 未安装"
    exit 1
fi
echo "✓ Docker 已安装"

if ! command -v docker-compose &> /dev/null; then
    echo "⚠️  Docker Compose 未安装，尝试使用 docker compose..."
    COMPOSE_CMD="docker compose"
else
    COMPOSE_CMD="docker-compose"
fi
echo "✓ 使用命令: $COMPOSE_CMD"
echo ""

# 2. 生成 SSL 证书
echo "2️⃣  生成 SSL 证书..."
if [ -f "docker/generate-certs.sh" ]; then
    bash docker/generate-certs.sh
else
    echo "⚠️  generate-certs.sh 未找到，跳过证书生成"
fi
echo ""

# 3. 创建必要的目录
echo "3️⃣  创建目录结构..."
mkdir -p data
mkdir -p docker/certs
mkdir -p docker/nginx/conf.d
mkdir -p docker/html
echo "✓ 目录结构已创建"
echo ""

# 4. 显示配置信息
echo "4️⃣  配置信息:"
echo "  • HTTP 端口: 80"
echo "  • HTTPS 端口: 443"
echo "  • 后端服务: http://localhost:8080"
echo "  • Web 管理界面: https://localhost"
echo ""

# 5. 启动服务
echo "5️⃣  启动 Docker 容器..."
$COMPOSE_CMD up -d

if [ $? -eq 0 ]; then
    echo ""
    echo "✓ 容器启动成功!"
    echo ""
    echo "╔════════════════════════════════════════════════════════════════╗"
    echo "║                       服务已启动                              ║"
    echo "╠════════════════════════════════════════════════════════════════╣"
    echo "║ HTTP:  http://localhost                                       ║"
    echo "║ HTTPS: https://localhost                                      ║"
    echo "║                                                                ║"
    echo "║ API:   https://localhost/api/admin                            ║"
    echo "║ 后端:  http://localhost:8080                                  ║"
    echo "╚════════════════════════════════════════════════════════════════╝"
    echo ""
    echo "📝 查看日志: $COMPOSE_CMD logs -f"
    echo "🛑 停止服务: $COMPOSE_CMD down"
    echo "🔄 重启服务: $COMPOSE_CMD restart"
    echo ""
else
    echo "✗ 容器启动失败"
    exit 1
fi
