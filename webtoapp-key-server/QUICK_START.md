# 🎉 WebToApp Key Server - 成功部署！

## ✅ 部署完成清单

- ✅ Go 后端应用编译 (Go 1.24)
- ✅ Docker 镜像构建成功
- ✅ Nginx 反向代理启动
- ✅ SSL/TLS 证书生成 (www.yingcai.com)
- ✅ 前端文件加载
- ✅ 数据库初始化
- ✅ 审计日志系统
- ✅ 统计分析功能
- ✅ 所有容器健康检查就绪

---

## 🚀 立即访问

### 使用 localhost（本地测试）

```bash
# HTTPS (忽略证书警告)
https://localhost/index.html

# 或使用 curl
curl -k https://localhost/index.html
```

### 使用自定义域名（推荐）

编辑 `/etc/hosts`：

```bash
# 打开编辑器
sudo nano /etc/hosts

# 添加这行
127.0.0.1  www.yingcai.com

# 保存并清空 DNS 缓存
sudo dscacheutil -flushcache

# 然后访问
https://www.yingcai.com/index.html
```

---

## 📊 系统架构

```
┌─────────────────────────────────────┐
│     Client (Browser/API)            │
│     https://www.yingcai.com         │
└────────────────┬────────────────────┘
                 │ HTTPS:443 / HTTP:80
                 │
        ┌────────▼────────┐
        │  Nginx Proxy    │
        │  (Port 80/443)  │
        │  SSL/TLS        │
        └────────┬────────┘
                 │ Internal
                 │ http://localhost:8080
        ┌────────▼────────┐
        │ Go Backend      │
        │ (Port 8080)     │
        │ • 18+ APIs      │
        │ • Auth          │
        │ • Audit Logs    │
        │ • Statistics    │
        └────────┬────────┘
                 │
        ┌────────▼────────┐
        │ SQLite Database │
        │ data/keyserver  │
        └─────────────────┘

Network: webtoapp-network (Docker Bridge)
```

---

## 🔧 常用命令

### 启动/停止

```bash
# 启动
docker-compose up -d

# 停止
docker-compose stop

# 重启
docker-compose restart

# 查看日志
docker-compose logs -f
docker-compose logs -f keyserver
docker-compose logs -f nginx
```

### 调试

```bash
# 进入容器
docker-compose exec keyserver bash
docker-compose exec nginx bash

# 查看容器状态
docker-compose ps

# 检查证书
openssl x509 -in docker/certs/cert.pem -text -noout
```

### 清理

```bash
# 停止并删除容器（保留数据）
docker-compose down

# 完全清理（删除所有）
docker-compose down -v

# 清理 Docker 镜像
docker image prune -a

# 清理 Docker 缓存
docker builder prune -af
```

---

## 📋 关键配置

### 域名与证书

```
主域名:       www.yingcai.com
备选域名:     yingcai.com, *.yingcai.com, localhost
证书类型:     自签名 (RSA 4096-bit)
有效期:       365 天 (2025-12-09 ~ 2026-12-09)
证书位置:     docker/certs/cert.pem
密钥位置:     docker/certs/key.pem
```

### 端口配置

```
80/TCP   → Nginx HTTP (自动重定向 → HTTPS)
443/TCP  → Nginx HTTPS
8080/TCP → Go 后端 (仅容器内部)
```

### 数据持久化

```
./data/                 → 数据库文件
./docker/certs/         → SSL 证书
./docker/nginx/         → Nginx 配置
./configs/              → 应用配置
./web/                  → 前端文件
```

---

## 🎯 API 端点

### 公开端点

```
GET     /                      # 重定向到 index.html
GET     /index.html            # 管理后台
GET     /login.html            # 登录页面
GET     /css/*                 # CSS 文件
GET     /js/*                  # JavaScript 文件
GET     /static/*              # 静态资源
```

### 受保护的管理端点 (需要 API Key)

```
POST    /api/admin/api-keys                 # 生成 API Key
GET     /api/admin/api-keys                 # 列出 API Keys
GET     /api/admin/api-keys/:id             # 获取单个 API Key
PUT     /api/admin/api-keys/:id             # 更新 API Key
DELETE  /api/admin/api-keys/:id             # 撤销 API Key

GET     /api/admin/statistics               # 统计数据
GET     /api/admin/statistics/dashboard     # 仪表板
GET     /api/admin/statistics/apps/:app_id  # 应用统计
GET     /api/admin/statistics/apps/:app_id/trends  # 趋势数据

GET     /api/admin/logs                     # 审计日志
GET     /api/admin/health                   # 健康检查
```

### 激活 API (Public)

```
POST    /api/activation/generate            # 生成激活码
POST    /api/activation/verify              # 验证激活码
GET     /api/activation/:app_id             # 获取统计
DELETE  /api/activation/:app_id/:code       # 撤销激活码
```

---

## 🔐 认证方式

### Bearer Token

```bash
curl -k https://www.yingcai.com/api/admin/api-keys \
  -H "Authorization: Bearer YOUR_API_KEY"
```

### API Key 格式

```
YOUR_API_KEY = UUID.HexString

示例:
11e4c383-ba20-412b-8726-f439c39b0214.ae188e6c-102d-4109-bb24-7c02049b1b7f
```

---

## 🌐 生产部署注意事项

### 1. 使用真实域名
```bash
# 购买域名并配置 DNS 记录
A       www.yingcai.com    YOUR_SERVER_IP
A       yingcai.com        YOUR_SERVER_IP
CNAME   *.yingcai.com      www.yingcai.com
```

### 2. 获取正式 SSL 证书
```bash
# 使用 Let's Encrypt (免费)
certbot certonly --nginx -d www.yingcai.com

# 更新 docker-compose.yml 的证书路径
```

### 3. 配置环境变量
```bash
# 创建 .env 文件
JWT_SECRET=your-secure-random-secret
DB_PATH=/app/data/keyserver.db
ENV=production
```

### 4. 备份数据
```bash
# 定期备份数据库
cp -r data data.backup.$(date +%Y%m%d)
```

### 5. 监控和日志
```bash
# 持续监控容器
docker stats

# 收集和分析日志
docker-compose logs --tail=1000 > app.log
```

---

## 📊 性能指标

### 系统要求

| 指标 | 最小值 | 推荐值 |
|------|--------|--------|
| CPU | 1核 | 2核+ |
| 内存 | 512MB | 2GB+ |
| 硬盘 | 2GB | 10GB+ |
| 网络 | 10Mbps | 100Mbps+ |

### 容器资源

| 容器 | CPU | 内存 |
|------|-----|------|
| keyserver | ~50m | ~150MB |
| nginx | ~10m | ~50MB |

---

## 🐛 故障排查

### 问题：浏览器显示"不是私密连接"

**原因**: 自签名证书  
**解决**: 点击"继续"或使用 `curl -k`

### 问题：502 Bad Gateway

**原因**: 后端服务未启动  
**解决**:
```bash
docker-compose logs keyserver
docker-compose restart keyserver
```

### 问题：SSL 证书验证失败

**原因**: 域名不匹配或证书过期  
**解决**:
```bash
# 检查证书有效期
openssl x509 -in docker/certs/cert.pem -noout -dates

# 重新生成证书
rm -rf docker/certs
bash docker/generate-certs.sh
docker-compose restart nginx
```

### 问题：数据库锁定

**原因**: 多个进程访问数据库  
**解决**:
```bash
docker-compose restart keyserver
```

### 问题：超时错误

**原因**: 代理超时设置过短  
**解决**: 增加 `docker/nginx/conf.d/default.conf` 中的 `proxy_read_timeout`

---

## 📚 相关文档

- [DOMAIN_SETUP.md](./DOMAIN_SETUP.md) - 域名和 SSL 配置指南
- [DEPLOYMENT.md](./DEPLOYMENT.md) - 详细部署指南
- [STARTUP_COMPLETE.md](./STARTUP_COMPLETE.md) - 启动完成说明
- [README.md](./README.md) - 项目概述

---

## 🎓 学习资源

### Docker & Docker Compose

- [Docker 官方文档](https://docs.docker.com/)
- [Docker Compose 官方文档](https://docs.docker.com/compose/)

### Nginx

- [Nginx 官方文档](http://nginx.org/en/docs/)
- [Nginx 反向代理配置](https://nginx.org/en/docs/http/ngx_http_proxy_module.html)

### SSL/TLS

- [Let's Encrypt 文档](https://letsencrypt.org/zh-cn/docs/)
- [OpenSSL 官方文档](https://www.openssl.org/docs/)

### Go

- [Go 官方网站](https://golang.org/)
- [Gin Web 框架](https://github.com/gin-gonic/gin)
- [GORM 文档](https://gorm.io/zh_CN/)

---

## 💬 技术支持

### 常见问题反馈

遇到问题？请检查：

1. ✅ Docker 和 Docker Compose 版本
2. ✅ 容器日志 (`docker-compose logs`)
3. ✅ 端口占用 (`lsof -i :80,8080,443`)
4. ✅ 磁盘空间 (`df -h`)
5. ✅ DNS 配置 (`cat /etc/hosts`)

### 获取帮助

```bash
# 检查系统状态
docker-compose ps
docker-compose logs --tail=100

# 进行诊断
docker-compose exec keyserver curl http://localhost:8080/api/admin/health \
  -H "Authorization: Bearer YOUR_API_KEY"

# 导出日志用于分析
docker-compose logs > diagnose.log 2>&1
```

---

## 🎊 庆祝成功！

```
  _____ _                _           _      
 / ____| |              | |         | |     
| (___ | |__  _   _ ___| |_   _  __| | ____ 
 \___ \| '_ \| | | / __| __| (_)/ _` |/ _  |
 ____) | | | | |_| \__ \ |_   _/ (_| | (_) |
|_____/|_| |_|\__,_|___/\__| |_|\__,_|\__  |
                                       __/ |
                                      |___/ 
   WebToApp Key Server 成功部署！
   
   访问: https://www.yingcai.com
   时间: 2025-12-09 11:02 CST
   状态: ✅ 就绪
```

---

**恭祝部署成功！** 🎉

最后更新: 2025-12-09  
版本: 1.0.0  
状态: 生产就绪 ✅
