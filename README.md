# Cloudflare DDNS Go

`cloudflare-ddns-go` 是一个使用 Go 语言编写的动态域名解析（DDNS）客户端，它可以自动更新 Cloudflare 上的 DNS 记录以匹配你的公网 IP 地址。

## 特性

- **自动更新**：当你的公网 IP 地址发生变化时，自动更新 Cloudflare 上的 DNS 记录。
- **简单配置**：通过环境变量进行配置，无需复杂的设置。
- **容器化部署**：支持 Docker，可以轻松部署在任何支持 Docker 的环境中。

## 快速开始

要使用 `cloudflare-ddns-go`，你需要有一个 Cloudflare 账户，并且在 Cloudflare 上有一个域名。

### 获取 API 令牌

登录你的 Cloudflare 账户，创建一个 API 令牌，授予它修改 DNS 记录的权限。

### 部署

使用以下命令启动 `cloudflare-ddns-go` 容器：

```bash
docker run --name cloudflare-ddns-go \
           -e CF_API_TOKEN=YOUR_CF_API_TOKEN \
           -e CF_DOMAIN_NAME=YOUR_CF_DOMAIN_NAME \
           -e CF_SUBDOMAIN_NAME=YOUR_CF_SUBDOMAIN_NAME \
           -e INTERVAL_MINUTES=5 \
           strangevy/cloudflare-ddns-go:latest
```
请将 YOUR_CF_API_TOKEN、YOUR_CF_DOMAIN_NAME 和 YOUR_CF_SUBDOMAIN_NAME 替换为你的实际 Cloudflare API 令牌、域名和子域名。

### 环境变量

- CF_API_TOKEN：你的 Cloudflare API 令牌。
- CF_DOMAIN_NAME：你想要更新的域名。如：example.com
- CF_SUBDOMAIN_NAME：你想要更新的子域名。如：sub，那你最终的解析的域名就是sub.example.com
- INTERVAL_MINUTES：定时执行频率，默认5（分钟）