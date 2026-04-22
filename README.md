# GoLDAP - OpenLDAP 管理系统

提供基于 Web 的用户、组和 sudo 规则管理功能的一个完整的 OpenLDAP 管理系统。

## 功能特性

- **用户管理**：创建、更新、删除 LDAP 用户，支持 Unix 属性（posixAccount）
- **组管理**：管理 LDAP 组和 posixGroups
- **Sudo 管理**：支持 sudoRole 的集中式 sudo 规则管理
- **SSH 密钥管理**：在 LDAP 中存储和管理 SSH 公钥
- **TLS 支持**：支持 TLS/LDAPS 安全连接
- **Web 界面**：现代化的 Vue.js 前端，便于管理
- **REST API**：完整的 REST API，支持自动化操作

## 架构

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│   Vue.js Web    │────▶│  Go Backend     │────▶│   OpenLDAP      │
│   前端界面       │     │  (端口 8150)    │     │   (端口 389)    │
└─────────────────┘     └────────┬────────┘     └─────────────────┘
                                 │
                                 ▼
                        ┌─────────────────┐
                        │     MySQL       │
                        │   (端口 3306)   │
                        └─────────────────┘
```

## 快速开始

### 前置要求

- Docker 和 Docker Compose
- OpenSSL（用于生成 TLS 证书）

### 部署步骤

1. **克隆并部署**：
   ```bash
   cd <项目目录>
   chmod +x scripts/*.sh openldap/*.sh
   ./scripts/deploy.sh
   ```

2. **访问 Web 界面**：
   - 访问地址：http://localhost:8150
   - 默认管理员：`admin` / `admin123`（⚠️ 首次登录后请立即修改）

3. **配置客户端机器**：
   ```bash
   # 复制到客户端机器并运行
   sudo LDAP_SERVER=ldap://your-server-address ./scripts/client-setup.sh
   ```

## 手动部署

### 1. 生成 TLS 证书

```bash
./scripts/generate-certs.sh
```

### 2. 启动服务

```bash
cd docker-compose
docker compose up -d
```

### 3. 加载 Sudo Schema

```bash
./openldap/load_sudo_schema.sh
```

### 4. 初始化 LDAP 结构

```bash
LDAP_HOST=localhost ./openldap/init_ldap_structure.sh
```

## 配置说明

### 服务器配置（`goldap/server/config.yml`）

| 配置项 | 设置项 | 说明 |
|---------|---------|-------------|
| system | mode | debug/release/test（运行模式） |
| system | port | 服务器端口（默认：8150） |
| ldap | url | LDAP 服务器地址 |
| ldap | base-dn | LDAP 基础 DN |
| ldap | admin-dn | 管理员 DN，用于修改操作 |
| mysql | host/port | MySQL 连接设置 |

### 环境变量

| 变量名 | 默认值 | 说明 |
|----------|---------|-------------|
| LDAP_DOMAIN | example.com | LDAP 域（请根据您的环境自定义） |
| LDAP_ADMIN_PASSWORD | - | LDAP 管理员密码（通过密钥设置） |
| MYSQL_ROOT_PASSWORD | - | MySQL root 密码（通过密钥设置） |
| MYSQL_USER | - | MySQL 用户（通过密钥设置） |
| MYSQL_PASSWORD | - | MySQL 用户密码（通过密钥设置） |

> **安全提示**：切勿将实际密码提交到版本控制系统。请使用环境变量或密钥管理系统（如 Docker secrets、Vault、Kubernetes secrets）。

## 客户端配置

### Ubuntu/Debian

使用您的服务器信息运行客户端配置脚本：

```bash
sudo LDAP_SERVER=ldap://YOUR_SERVER_IP \
     LDAP_BASE_DN=dc=example,dc=com \
     ./scripts/client-setup.sh
```

### 手动配置 SSSD

1. 安装软件包：
   ```bash
   apt-get install sssd sssd-ldap libnss-sss libpam-sss ldap-utils
   ```

2. 配置 `/etc/sssd/sssd.conf`（将占位符替换为您的实际值）：
   ```ini
   [sssd]
   services = nss, pam, ssh, sudo
   config_file_version = 2
   domains = example.com
   # debug_level = 9   # 调试时可开启

   [domain/example.com]
   id_provider = ldap
   auth_provider = ldap
   chpass_provider = ldap
   sudo_provider = ldap

   ldap_uri = ldap://YOUR_SERVER_IP
   ldap_search_base = dc=example,dc=com

   # SUDO 配置
   ldap_sudo_search_base = ou=sudoers,dc=example,dc=com
   ldap_sudo_full_refresh_interval = 86400
   ldap_sudo_smart_refresh_interval = 3600

   # 缓存配置
   cache_credentials = true
   enumerate = false
   entry_cache_timeout = 5400
   entry_cache_user_timeout = 5400
   entry_cache_group_timeout = 5400
   entry_cache_netgroup_timeout = 5400
   entry_cache_service_timeout = 5400
   entry_cache_sudo_timeout = 5400

   # 排除本地用户
   filter_users = root
   filter_groups = root

   # TLS/证书配置
   ldap_tls_reqcert = allow

   ldap_id_use_start_tls = true
   ldap_tls_cacertdir = /etc/openldap/cacerts

   # SSH 公钥支持
   ldap_user_ssh_public_key = sshPublicKey
   ```

3. 配置 `/etc/nsswitch.conf`：
   ```
   passwd:         files systemd sss
   group:          files systemd sss
   shadow:         files systemd sss
   gshadow:        files systemd
   sudoers:        files sss

   hosts:          files mdns4_minimal [NOTFOUND=return] dns
   networks:       files

   protocols:      db files
   services:       db files sss
   ethers:         db files
   rpc:             db files

   netgroup:       nis sss
   automount:      sss
   ```

4. 重启 SSSD：
   ```bash
   systemctl restart sssd
   ```

## TLS 配置

### 服务端 TLS

TLS 在 docker-compose 中自动配置。证书存储在 `certs/` 目录下：
- `ca.crt` - CA 证书
- `ldap.crt` - 服务器证书
- `ldap.key` - 服务器私钥

### 客户端 TLS

1. 将 `ca.crt` 复制到客户端：`/etc/openldap/cacerts/ca.crt`

2. 在 `/etc/ldap/ldap.conf` 中配置 SSSD 的 TLS：
   ```ini
   #TLS_CACERT      /etc/ssl/certs/ca-certificates.crt
   TLS_CACERT       /etc/openldap/cacerts/ca.crt
   TLS_REQCERT      allow
   ```

## Sudo 配置

### 通过 Web 界面创建 Sudo 规则

1. 登录 Web 界面
2. 导航到 "Sudo Rules"
3. 创建新规则，设置：
   - **sudoUser**：用户名或 "ALL"
   - **sudoHost**：主机名或 "ALL"
   - **sudoCommand**：命令或 "ALL"
   - **sudoOption**：输入 "!authenticate" 表示免密（NOPASSWD）

### 示例：授予所有用户免密 sudo 权限

```ldif
dn: cn=all-nopasswd,ou=sudoers,dc=example,dc=com
objectClass: top
objectClass: sudoRole
cn: all-nopasswd
sudoUser: ALL
sudoHost: ALL
sudoCommand: ALL
sudoOption: !authenticate
```

## API 参考

### 认证

```bash
curl -X POST http://localhost:8150/api/base/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"your_password"}'
```

### 用户管理

```bash
# 列出用户
curl -H "Authorization: Bearer TOKEN" http://localhost:8150/api/user/list

# 创建用户
curl -X POST http://localhost:8150/api/user/add \
  -H "Authorization: Bearer TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"username":"newuser","nickname":"New User",...}'
```

## 目录结构

```
.
├── docker-compose/      # Docker compose 配置
├── goldap/
│   ├── client/         # Vue.js 前端
│   └── server/         # Go 后端
├── initdb/             # MySQL 初始化脚本
├── openldap/           # OpenLDAP 配置和脚本
├── scripts/            # 部署脚本
│   ├── deploy.sh       # 主部署脚本
│   ├── generate-certs.sh  # TLS 证书生成
│   └── client-setup.sh # 客户端配置脚本
└── certs/              # TLS 证书（自动生成）
```

## 故障排查

### LDAP 连接问题

```bash
# 测试 LDAP 连接
ldapsearch -x -H ldap://YOUR_SERVER_IP:389 -b "dc=example,dc=com"

# 使用 TLS 测试
ldapsearch -x -H ldap://YOUR_SERVER_IP:389 -ZZ -b "dc=example,dc=com"
```

### SSSD 问题

```bash
# 检查 SSSD 状态
systemctl status sssd

# 清除 SSSD 缓存
sss_cache -E
systemctl restart sssd

# 测试用户查询
getent passwd username
id username
```

### Sudo 问题

```bash
# 列出用户的 sudo 规则
sudo -l -U username

# 检查 sudo LDAP 配置
ldapsearch -x -H ldap://YOUR_SERVER_IP -b "ou=sudoers,dc=example,dc=com" "(objectClass=sudoRole)"
```

## 安全最佳实践

1. **修改默认凭据**：部署后立即修改默认管理员密码
2. **使用强 TLS**：确保正确的证书管理和定期更新
3. **网络隔离**：将 LDAP 服务器部署在安全的网络段
4. **定期备份**：实施定期 LDAP 数据库备份
5. **访问控制**：通过 VPN 或防火墙规则限制 Web 界面访问
6. **审计日志**：定期监控和审查 LDAP 访问日志

## 许可证

MIT License
