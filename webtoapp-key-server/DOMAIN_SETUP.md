# WebToApp Key Server - 域名配置指南

## 📝 域名配置

所有配置已更新以支持以下域名：

- **主域名**: `www.yingcai.com`
- **根域名**: `yingcai.com`
- **通配符**: `*.yingcai.com`
- **本地**: `localhost`

---

## 🔐 SSL 证书信息

### 证书生成配置

证书文件已配置以支持所有域名变体：

```
主体名称 (CN):     www.yingcai.com
证书位置:          docker/certs/cert.pem
私钥位置:          docker/certs/key.pem
有效期:            365 天
加密方式:          RSA 4096-bit

主体备选名称 (SAN):
  - www.yingcai.com
  - yingcai.com
  - *.yingcai.com
  - localhost
```

---

## 🖥️ 本地测试设置

### 步骤 1: 配置 hosts 文件

编辑 `/etc/hosts` 文件（需要 sudo 权限）：

```bash
sudo nano /etc/hosts
```

在文件末尾添加以下行：

```
127.0.0.1  www.yingcai.com
127.0.0.1  yingcai.com
127.0.0.1  localhost
```

**macOS 用户**:
```bash
# 编辑
sudo nano /etc/hosts

# 或使用 vi
sudo vi /etc/hosts
```

**验证配置**:
```bash
# 检查 hosts 文件
cat /etc/hosts | grep yingcai

# 验证 DNS 解析
ping www.yingcai.com
# 应该返回: PING www.yingcai.com (127.0.0.1)
```

### 步骤 2: 清空 DNS 缓存

```bash
# macOS
sudo dscacheutil -flushcache

# Linux
sudo systemctl restart systemd-resolved

# 或
sudo resolvectl flush-caches
```

---

## 🚀 启动服务

### 前置条件

确保 Docker 和 Docker Compose 已安装：

```bash
docker --version
docker-compose --version
```

### 生成 SSL 证书

```bash
bash docker/generate-certs.sh
```

**输出示例**:
```
🔐 生成自签名 SSL 证书...
✓ SSL 证书生成成功!
  证书文件: ./docker/certs/cert.pem
  密钥文件: ./docker/certs/key.pem
  主域名: www.yingcai.com
  备选域名: yingcai.com, *.yingcai.com, localhost

⚠️  注意: 这是自签名证书，用于开发/测试环境
⚠️  生产环境建议使用正式的 SSL 证书（如 Let's Encrypt）
```

### 启动 Docker 容器

```bash
# 方法 1: 使用启动脚本（推荐）
bash docker/start.sh

# 方法 2: 直接使用 docker-compose
docker-compose up -d
```

---

## 🌐 访问服务

### 通过域名访问

使用任何配置的域名都可以访问服务：

```
https://www.yingcai.com      # 推荐使用
https://yingcai.com           # 根域名
https://localhost             # 本地测试
```

### API 端点

```
https://www.yingcai.com/api/admin/keys        # API Key 管理
https://www.yingcai.com/api/admin/statistics  # 统计数据
https://www.yingcai.com/api/admin/logs        # 审计日志
https://www.yingcai.com/health                # 健康检查
```

### Web 管理界面

```
https://www.yingcai.com/index.html
```

---

## 🔍 验证和调试

### 验证 SSL 证书

```bash
# 查看证书详情
openssl x509 -in docker/certs/cert.pem -text -noout

# 查看证书有效期
openssl x509 -in docker/certs/cert.pem -noout -dates

# 查看证书主体和 SAN
openssl x509 -in docker/certs/cert.pem -noout -subject -ext subjectAltName
```

**示例输出**:
```
Subject: C = CN, ST = Guangdong, L = Guangzhou, O = Yingcai, CN = www.yingcai.com
X509v3 Subject Alternative Name:
    DNS:www.yingcai.com, DNS:yingcai.com, DNS:*.yingcai.com, DNS:localhost
```

### 测试 HTTPS 连接

```bash
# 使用 curl（忽略自签名证书警告）
curl -k https://www.yingcai.com/health

# 使用 openssl
openssl s_client -connect www.yingcai.com:443

# 查看证书信息
openssl s_client -connect www.yingcai.com:443 -showcerts
```

### 检查容器状态

```bash
# 列出运行中的容器
docker-compose ps

# 查看日志
docker-compose logs -f

# 查看特定服务日志
docker-compose logs -f nginx
docker-compose logs -f keyserver

# 进入容器调试
docker-compose exec nginx sh
docker-compose exec keyserver sh
```

### 验证 Nginx 配置

```bash
# 检查 Nginx 配置语法
docker-compose exec nginx nginx -t

# 查看 Nginx 配置
docker-compose exec nginx cat /etc/nginx/conf.d/default.conf

# 查看 Nginx 日志
docker-compose exec nginx tail -f /var/log/nginx/access.log
```

---

## 🌍 DNS 和网络配置

### 本地开发环境

当前配置适用于本地开发，使用 hosts 文件映射域名到 localhost。

### 生产环境部署

若要在互联网上使用 `www.yingcai.com`，需要：

1. **注册域名**
   - 购买 `yingcai.com` 域名
   - 从域名注册商获得 DNS 管理权限

2. **配置 DNS 记录**
   ```dns
   Type    Host              Value
   A       www.yingcai.com   your.server.ip.address
   A       yingcai.com       your.server.ip.address
   CNAME   *.yingcai.com     www.yingcai.com
   ```

3. **获取正式 SSL 证书** (Let's Encrypt)
   ```bash
   # 使用 certbot
   sudo apt-get install certbot python3-certbot-nginx
   sudo certbot certonly --nginx -d www.yingcai.com -d yingcai.com
   
   # 将证书路径更新到 docker-compose.yml
   ```

4. **更新 docker-compose.yml**
   ```yaml
   volumes:
     - /etc/letsencrypt/live/www.yingcai.com/fullchain.pem:/etc/nginx/certs/cert.pem:ro
     - /etc/letsencrypt/live/www.yingcai.com/privkey.pem:/etc/nginx/certs/key.pem:ro
   ```

5. **重启 Nginx**
   ```bash
   docker-compose restart nginx
   ```

---

## 📋 常见问题

### Q: 为什么访问 HTTPS 时浏览器显示安全警告？

**A**: 因为使用的是自签名证书。这对开发/测试是正常的。
- 点击"继续" → "高级" → "继续访问"
- 或使用 `curl -k` 忽略证书警告

### Q: 如何在生产环境中使用正式证书？

**A**: 使用 Let's Encrypt 获取免费正式证书：
```bash
# 停止 Docker
docker-compose down

# 生成证书
sudo certbot certonly --standalone -d www.yingcai.com -d yingcai.com

# 更新 docker-compose.yml 中的证书路径
# 重启服务
docker-compose up -d
```

### Q: hosts 文件配置后为什么仍然无法解析？

**A**: 需要清空 DNS 缓存：
```bash
# macOS
sudo dscacheutil -flushcache

# Linux
sudo systemctl restart systemd-resolved
```

### Q: 如何添加更多域名？

**A**: 修改以下文件并重新生成证书：

1. **docker/generate-certs.sh** - 更新 `-addext "subjectAltName=..."` 行
2. **docker/nginx/conf.d/default.conf** - 更新 `server_name` 指令
3. **重新生成证书**: `bash docker/generate-certs.sh`
4. **重启 Nginx**: `docker-compose restart nginx`

---

## 📊 配置总结

| 配置项 | 值 |
|------|-----|
| 主域名 | www.yingcai.com |
| 根域名 | yingcai.com |
| 通配符 | *.yingcai.com |
| 本地访问 | localhost |
| HTTP 端口 | 80 |
| HTTPS 端口 | 443 |
| 证书类型 | 自签名 (开发/测试) |
| 证书有效期 | 365 天 |
| SSL/TLS 版本 | TLSv1.2, TLSv1.3 |
| 加密算法 | RSA 4096-bit |

---

## 🔗 相关文件

- 证书生成脚本: `docker/generate-certs.sh`
- Nginx 主配置: `docker/nginx/nginx.conf`
- 站点配置: `docker/nginx/conf.d/default.conf`
- Docker Compose: `docker-compose.yml`
- 启动脚本: `docker/start.sh`
- 部署指南: `DEPLOYMENT.md`

---

**最后更新**: 2024-01-01  
**版本**: 1.0.0  
**域名**: www.yingcai.com
