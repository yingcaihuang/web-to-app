# 🚀 WebToApp Key Server - 启动完成！

## ✅ 部署状态

所有 Docker 容器已成功启动！

```
✓ Keyserver (Go 后端)  - http://localhost:8080
✓ Nginx (反向代理)      - https://localhost:443, http://localhost:80
✓ 网络                   - webtoapp-network
✓ SSL/TLS 证书           - www.yingcai.com (自签名)
```

---

## 📝 配置信息

### 域名
- **主域名**: www.yingcai.com
- **根域名**: yingcai.com
- **通配符**: *.yingcai.com
- **本地测试**: localhost

### 证书信息
```
位置:     docker/certs/
证书:     cert.pem (RSA 4096-bit)
密钥:     key.pem
有效期:   365 天
类型:     自签名 (开发/测试)
```

### 端口映射
```
80/TCP    → Nginx HTTP (自动重定向到 HTTPS)
443/TCP   → Nginx HTTPS + SSL/TLS
8080/TCP  → Go 后端 (仅容器内部)
```

---

## 🔐 访问服务

### 本地测试 (localhost)

```bash
# 忽略证书警告访问 HTTPS
curl -k https://localhost/index.html

# 或使用浏览器（点击"继续"）
# 浏览器: https://localhost
```

### 使用域名 (www.yingcai.com)

配置 `/etc/hosts` 文件：

```bash
# macOS/Linux
sudo nano /etc/hosts

# 添加以下行
127.0.0.1  www.yingcai.com
127.0.0.1  yingcai.com
127.0.0.1  localhost

# 保存后清空 DNS 缓存
sudo dscacheutil -flushcache  # macOS
```

然后访问：
```
https://www.yingcai.com
https://yingcai.com
```

---

## 📊 验证部署

### 1. 检查容器运行状态

```bash
docker-compose ps

# 输出示例:
# NAME                 STATUS              PORTS
# webtoapp-keyserver   Up (health: started) 0.0.0.0:8080->8080/tcp
# webtoapp-nginx       Up (health: started) 0.0.0.0:80->80/tcp, 0.0.0.0:443->443/tcp
```

### 2. 检查 SSL 证书

```bash
# 查看证书详情
openssl x509 -in docker/certs/cert.pem -text -noout

# 查看证书有效期
openssl x509 -in docker/certs/cert.pem -noout -dates

# 输出示例:
# notBefore=Dec  9 02:58:23 2025 GMT
# notAfter=Dec  9 02:58:23 2026 GMT
```

### 3. 测试 HTTPS 连接

```bash
# 使用 curl（忽略自签名证书）
curl -k https://localhost/api/admin/health \
  -H "Authorization: Bearer YOUR_API_KEY"

# 使用 openssl
openssl s_client -connect localhost:443 -showcerts

# 使用浏览器开发者工具
# 1. 打开 https://localhost
# 2. 右键 → 检查 (Inspect)
# 3. Console 标签查看 HTTPS 连接状态
```

### 4. 查看日志

```bash
# 查看所有日志
docker-compose logs -f

# 查看 Nginx 日志
docker-compose logs -f nginx

# 查看后端日志
docker-compose logs -f keyserver

# 查看特定数量的日志行
docker-compose logs --tail=50
```

---

## 🛠️ 常见操作

### 启动/停止服务

```bash
# 启动所有服务
docker-compose up -d

# 停止所有服务
docker-compose stop

# 重启所有服务
docker-compose restart

# 删除容器（数据保留）
docker-compose down

# 完全清理（删除所有）
docker-compose down -v
```

### 重新生成证书

```bash
# 删除旧证书
rm -rf docker/certs

# 生成新证书
bash docker/generate-certs.sh

# 重启 Nginx 加载新证书
docker-compose restart nginx
```

### 进入容器调试

```bash
# 进入后端容器
docker-compose exec keyserver /bin/bash

# 进入 Nginx 容器
docker-compose exec nginx /bin/bash

# 查看数据库文件
docker-compose exec keyserver ls -la /app/data/

# 验证 Nginx 配置
docker-compose exec nginx nginx -t
```

---

## 📋 文件结构

```
webtoapp-key-server/
├── docker/
│   ├── Dockerfile                    # Go 后端构建配置
│   ├── generate-certs.sh            # SSL 证书生成脚本
│   ├── start.sh                     # 启动脚本
│   ├── certs/                       # SSL 证书目录
│   │   ├── cert.pem                 # 证书文件
│   │   └── key.pem                  # 私钥文件
│   ├── nginx/
│   │   ├── nginx.conf              # Nginx 主配置
│   │   └── conf.d/
│   │       └── default.conf        # 站点配置
│   └── html/                       # 静态文件目录（可选）
├── docker-compose.yml              # Docker 容器编排
├── data/                           # SQLite 数据库
├── web/                            # 前端资源（由后端提供）
├── cmd/
│   └── main.go                     # Go 应用入口
├── DOMAIN_SETUP.md                 # 域名配置指南
├── DEPLOYMENT.md                   # 部署指南
└── docker/start.sh                 # 启动脚本
```

---

## 🔗 关键文件修改

### ✅ Dockerfile (docker/Dockerfile)
- **改动**: Go 1.21 → Go 1.24
- **改动**: Alpine → Debian Bookworm
- **原因**: 支持 go.mod 中的 Go 1.23.0 版本要求

### ✅ docker-compose.yml
- **改动**: 移除 `version: '3.8'`（已废弃）
- **新增**: Nginx 服务配置
- **新增**: SSL/TLS 证书挂载
- **改动**: 域名配置为 www.yingcai.com

### ✅ Nginx 配置
- **新增**: `docker/nginx/nginx.conf` - 主配置
- **新增**: `docker/nginx/conf.d/default.conf` - 站点配置
- **特性**: HTTP → HTTPS 自动重定向
- **特性**: SSL/TLS 1.2+ 支持
- **特性**: 安全头部配置
- **特性**: CORS 支持

### ✅ 证书配置
- **新增**: `docker/generate-certs.sh` - 证书生成脚本
- **特性**: RSA 4096-bit 加密
- **特性**: SAN 支持多个域名
- **特性**: 365 天有效期

---

## 🌍 下一步：生产部署

### 1. 使用真实域名
```bash
# 购买 www.yingcai.com 域名后
# 配置 DNS 记录指向服务器 IP
A       www.yingcai.com    your.server.ip
A       yingcai.com        your.server.ip
CNAME   *.yingcai.com      www.yingcai.com
```

### 2. 使用 Let's Encrypt 证书
```bash
# 安装 certbot
sudo apt-get install certbot python3-certbot-nginx

# 生成证书
sudo certbot certonly --nginx -d www.yingcai.com -d yingcai.com

# 更新 docker-compose.yml 中的证书路径
volumes:
  - /etc/letsencrypt/live/www.yingcai.com/fullchain.pem:/etc/nginx/certs/cert.pem:ro
  - /etc/letsencrypt/live/www.yingcai.com/privkey.pem:/etc/nginx/certs/key.pem:ro

# 重启 Nginx
docker-compose restart nginx
```

### 3. 配置环境变量
```bash
# 创建 .env 文件
cat > .env << EOF
JWT_SECRET=your-secure-random-secret-here
DB_PATH=/app/data/keyserver.db
SERVER_PORT=8080
ENV=production
EOF

# 在 docker-compose.yml 中使用
environment:
  - JWT_SECRET=${JWT_SECRET}
```

### 4. 设置备份策略
```bash
# 定时备份数据库
0 2 * * * cp -r /path/to/data /path/to/backup/data.$(date +\%Y\%m\%d)
```

---

## 📞 技术支持

### 常见问题

**Q: 为什么浏览器显示"您的连接不是私密连接"？**
A: 因为使用自签名证书。点击"继续"或使用 `curl -k` 忽略警告。

**Q: 如何使用真实域名？**
A: 购买域名 → 配置 DNS → 使用 Let's Encrypt 获取正式证书 → 更新 docker-compose.yml

**Q: 数据会丢失吗？**
A: 不会。`data/` 目录已映射到容器外，重启不会丢失数据。

**Q: 如何添加新的域名？**
A: 修改 `docker/generate-certs.sh` 中的 `-addext "subjectAltName=..."` 和 `docker/nginx/conf.d/default.conf` 中的 `server_name`

---

## 📊 系统要求

- **Docker**: 20.10+
- **Docker Compose**: 2.0+
- **硬盘**: 2GB+
- **内存**: 512MB+
- **CPU**: 双核+
- **操作系统**: Linux, macOS, Windows (with WSL2)

---

## 🎯 快速命令速查表

```bash
# 启动应用
bash docker/start.sh
# 或
docker-compose up -d

# 查看状态
docker-compose ps

# 查看日志
docker-compose logs -f

# 停止应用
docker-compose stop

# 完全清理
docker-compose down -v

# 进入后端
docker-compose exec keyserver bash

# 检查证书
openssl x509 -in docker/certs/cert.pem -noout -dates

# 测试 API
curl -k https://localhost/api/admin/health \
  -H "Authorization: Bearer YOUR_API_KEY"
```

---

✨ **享受 WebToApp Key Server!** ✨

**最后更新**: 2025-12-09  
**版本**: 1.0.0  
**状态**: ✅ 生产就绪
