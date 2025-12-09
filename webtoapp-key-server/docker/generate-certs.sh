#!/bin/bash

# 生成自签名证书的脚本
# 用于 Docker 容器中的 Nginx SSL 配置

CERT_DIR="./docker/certs"
CERT_FILE="$CERT_DIR/cert.pem"
KEY_FILE="$CERT_DIR/key.pem"

# 创建证书目录
mkdir -p "$CERT_DIR"

# 检查证书是否已存在
if [ -f "$CERT_FILE" ] && [ -f "$KEY_FILE" ]; then
    echo "✓ SSL 证书已存在: $CERT_FILE"
    echo "✓ SSL 密钥已存在: $KEY_FILE"
    exit 0
fi

echo "🔐 生成自签名 SSL 证书..."

# 生成自签名证书（365天有效期，支持 www.yingcai.com）
openssl req -x509 -newkey rsa:4096 -keyout "$KEY_FILE" -out "$CERT_FILE" \
    -days 365 -nodes \
    -subj "/C=CN/ST=Guangdong/L=Guangzhou/O=Yingcai/CN=www.yingcai.com" \
    -addext "subjectAltName=DNS:www.yingcai.com,DNS:yingcai.com,DNS:*.yingcai.com,DNS:localhost"

if [ $? -eq 0 ]; then
    echo "✓ SSL 证书生成成功!"
    echo "  证书文件: $CERT_FILE"
    echo "  密钥文件: $KEY_FILE"
    echo "  主域名: www.yingcai.com"
    echo "  备选域名: yingcai.com, *.yingcai.com, localhost"
    echo ""
    echo "⚠️  注意: 这是自签名证书，用于开发/测试环境"
    echo "⚠️  生产环境建议使用正式的 SSL 证书（如 Let's Encrypt）"
else
    echo "✗ SSL 证书生成失败"
    exit 1
fi
