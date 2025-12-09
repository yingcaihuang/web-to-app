#!/bin/bash

# WebToApp Key Server 初始化脚本
# 在容器启动时运行，生成默认 Admin API Key

set -e

echo "🔧 WebToApp Key Server 初始化..."
echo ""

# 检查数据库是否存在
if [ ! -f "/app/data/keyserver.db" ]; then
    echo "📊 首次启动，初始化数据库..."
    # 数据库会在应用启动时自动创建
else
    echo "📊 数据库已存在，跳过初始化"
fi

echo ""
echo "✅ 初始化完成，启动应用..."
echo ""

# 启动应用
exec /app/keyserver
