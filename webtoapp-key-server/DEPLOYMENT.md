# WebToApp Key Server - 部署指南

## 快速开始

### 一键启动

```bash
bash docker/start.sh
```

该脚本会自动：
1. ✅ 检查 Docker 和 Docker Compose 依赖
2. ✅ 生成自签名 SSL 证书
3. ✅ 创建必要的目录结构
4. ✅ 启动 Nginx 反向代理和后端服务

### 访问服务

- **Web 管理界面**: https://localhost
- **HTTP 自动跳转**: http://localhost → https://localhost
- **API 端点**: https://localhost/api/admin
- **后端直连**: http://localhost:8080 (仅容器内部)

---

## 部署架构

```
┌─────────────────────────────────────────────┐
│           Internet / Local Machine          │
│         HTTP :80 / HTTPS :443               │
└──────────────────┬──────────────────────────┘
                   │
        ┌──────────▼──────────┐
        │  Nginx 反向代理    │
        │  (容器 nginx)      │
        │  SSL/TLS 处理       │
        │  CORS 支持          │
        └──────────┬──────────┘
                   │
        ┌──────────▼──────────┐
        │  Go 后端服务        │
        │  (容器 keyserver)   │
        │  :8080 内部端口     │
        │  SQLite 数据库      │
        └─────────────────────┘
```

### 网络配置

- **webtoapp-network**: Docker 桥接网络，用于容器间通信
- **Nginx**: 监听 0.0.0.0:80 和 0.0.0.0:443
- **Keyserver**: 只监听 0.0.0.0:8080，通过 Nginx 代理访问

---

## SSL/TLS 证书

### 自签名证书（开发/测试）

证书由 `docker/generate-certs.sh` 自动生成：

- **位置**: `docker/certs/`
- **有效期**: 365 天
- **加密方式**: RSA 4096-bit
- **主体名称**: localhost

### 生产环境

对于生产环境，使用真实 SSL 证书：

1. **Let's Encrypt** (推荐)
   ```bash
   # 使用 certbot 生成证书
   certbot certonly --webroot -w docker/html -d yourdomain.com
   
   # 将证书映射到 docker-compose.yml
   # volumes:
   #   - /etc/letsencrypt/live/yourdomain.com/fullchain.pem:/etc/nginx/certs/cert.pem
   #   - /etc/letsencrypt/live/yourdomain.com/privkey.pem:/etc/nginx/certs/key.pem
   ```

2. **商业 SSL 证书**
   - 从证书提供商获取 .pem 或 .crt 文件
   - 放入 `docker/certs/` 目录
   - 更新 `docker/nginx/conf.d/default.conf` 中的证书路径

---

## 常用命令

### 启动服务

```bash
# 自动启动（推荐）
bash docker/start.sh

# 或手动启动
docker-compose up -d
```

### 查看状态

```bash
# 列出所有容器
docker-compose ps

# 查看日志
docker-compose logs -f

# 查看特定服务日志
docker-compose logs -f nginx      # Nginx 日志
docker-compose logs -f keyserver  # 后端日志
```

### 停止/重启服务

```bash
# 停止所有服务
docker-compose stop

# 删除容器（保留数据）
docker-compose down

# 完全清理（删除容器和网络）
docker-compose down -v

# 重启服务
docker-compose restart
```

### 进入容器调试

```bash
# 进入后端容器
docker-compose exec keyserver sh

# 进入 Nginx 容器
docker-compose exec nginx sh

# 查看后端进程
docker-compose exec keyserver ps aux

# 查看 Nginx 配置
docker-compose exec nginx cat /etc/nginx/nginx.conf
```

---

## 数据持久化

### 数据库

- **位置**: `data/keyserver.db` (本地目录)
- **映射**: `data:/opt/app/data` (容器内)
- **备份**: 定期备份 `data/` 目录

```bash
# 备份数据库
cp -r data data.backup.$(date +%Y%m%d)

# 恢复数据库
cp data.backup.20240101/keyserver.db data/
docker-compose restart keyserver
```

### 日志

- **Nginx 日志**: 输出到 Docker 日志系统
- **应用日志**: stdout/stderr → Docker logs

```bash
# 查看完整日志历史
docker-compose logs --tail=100
```

---

## 健康检查

Docker Compose 配置了自动健康检查：

```bash
# 检查服务健康状态
docker-compose exec keyserver curl http://localhost:8080/health

# 通过 Nginx 检查
curl -k https://localhost/health
```

### 健康指标

- **Keyserver**: `/health` 端点，返回 200 OK 和 JSON 状态
- **Nginx**: TCP 连接检查到 `:80` 和 `:443`

---

## 环境配置

### 环境变量

在 `docker-compose.yml` 中修改或创建 `.env` 文件：

```bash
# .env
DB_PATH=/opt/app/data/keyserver.db
SERVER_PORT=8080
SERVER_HOST=0.0.0.0
LOG_LEVEL=info
```

### 配置文件

#### Nginx 主配置
- 文件: `docker/nginx/nginx.conf`
- 包含: Worker 进程数、日志格式、Gzip 压缩、SSL 设置

#### Nginx 站点配置
- 文件: `docker/nginx/conf.d/default.conf`
- 包含: 上游服务器、路由规则、安全头、CORS 设置

---

## 安全建议

### 1. SSL/TLS 配置

✅ **已实施**:
- HTTPS 强制 (HTTP 自动重定向)
- SSL/TLS 1.2+ 最小版本
- 强密码套件
- HSTS 头部 (strict-transport-security)

### 2. 安全头部

✅ **已实施**:
- `X-Frame-Options: DENY` - 防止点击劫持
- `X-Content-Type-Options: nosniff` - 防止 MIME 嗅探
- `X-XSS-Protection: 1; mode=block` - XSS 防护
- `Content-Security-Policy` - 内容安全策略
- `Referrer-Policy: strict-origin-when-cross-origin`

### 3. CORS 配置

✅ **已配置**:
- 允许来源: 任何来源 (可根据需要限制)
- 允许方法: GET, POST, PUT, DELETE, OPTIONS
- 允许头部: Content-Type, Authorization

### 4. API 认证

- 使用 Bearer Token 验证
- API Key SHA256 哈希存储
- 所有管理 API 操作记录审计日志

### 5. 防火墙规则

```bash
# 仅允许必要的端口
# 允许 80 (HTTP)
# 允许 443 (HTTPS)
# 限制 8080 (后端) 仅本地访问

# 示例 (iptables)
iptables -A INPUT -p tcp --dport 80 -j ACCEPT
iptables -A INPUT -p tcp --dport 443 -j ACCEPT
iptables -A INPUT -p tcp --dport 8080 -i lo -j ACCEPT
iptables -A INPUT -p tcp --dport 8080 -j DROP
```

---

## 故障排除

### 问题: 容器无法启动

```bash
# 查看错误日志
docker-compose logs keyserver
docker-compose logs nginx

# 检查端口占用
lsof -i :80
lsof -i :443
lsof -i :8080

# 解决: 更改 docker-compose.yml 中的端口
# ports:
#   - "8000:80"    # 改为 8000
#   - "8443:443"   # 改为 8443
```

### 问题: SSL 证书过期

```bash
# 重新生成证书
bash docker/generate-certs.sh

# 重启 Nginx
docker-compose restart nginx
```

### 问题: 连接超时

```bash
# 检查网络连通性
docker-compose exec nginx ping keyserver

# 检查防火墙
docker-compose exec keyserver netstat -tulpn

# 查看 Nginx 错误日志
docker-compose logs nginx | grep error
```

### 问题: 数据库锁定

```bash
# 检查数据库连接
docker-compose exec keyserver lsof data/keyserver.db

# 重新启动服务
docker-compose restart keyserver
```

---

## 监控和维护

### 定期检查清单

- [ ] 检查容器运行状态: `docker-compose ps`
- [ ] 查看系统资源使用: `docker stats`
- [ ] 验证日志无错误: `docker-compose logs`
- [ ] 备份数据库: `cp -r data data.backup`
- [ ] 检查 SSL 证书有效期: `docker-compose exec nginx openssl x509 -in /etc/nginx/certs/cert.pem -noout -dates`
- [ ] 测试 API 端点: `curl -k https://localhost/health`

### 日志分析

```bash
# 查看最后 50 行日志
docker-compose logs --tail=50

# 按时间过滤
docker-compose logs --since 1h

# 保存日志到文件
docker-compose logs > app.log 2>&1
```

---

## 升级和更新

### 更新后端代码

```bash
# 1. 重新构建镜像
docker-compose build keyserver

# 2. 重启服务
docker-compose up -d keyserver

# 3. 验证
docker-compose logs keyserver | tail -20
```

### 更新 Nginx 配置

```bash
# 1. 编辑 docker/nginx/conf.d/default.conf
nano docker/nginx/conf.d/default.conf

# 2. 重载 Nginx 配置 (无停机)
docker-compose exec nginx nginx -s reload

# 或重启
docker-compose restart nginx
```

---

## 生产部署清单

- [ ] 使用真实 SSL 证书 (Let's Encrypt 或商业)
- [ ] 配置 DNS 记录
- [ ] 设置备份策略
- [ ] 配置日志收集 (ELK, Splunk 等)
- [ ] 设置监控告警
- [ ] 配置 CI/CD 管道
- [ ] 执行安全审计
- [ ] 压力测试和性能优化
- [ ] 制定灾难恢复计划
- [ ] 实施变更管理流程

---

## 联系和支持

- 📧 邮件: support@webtoapp.com
- 📱 问题: GitHub Issues
- 📖 文档: https://webtoapp.com/docs
- 🔧 技术支持: tech@webtoapp.com

---

**最后更新**: 2024-01-01
**版本**: 1.0.0
