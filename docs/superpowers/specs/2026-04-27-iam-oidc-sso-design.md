# IAM OAuth2 + OIDC SSO 架构设计

## 一、整体架构

```
┌─────────────────────────────────────────────────────────────────┐
│                        用户浏览器                                │
└─────────────────┬───────────────────────────────────────────────┘
                  │
    ┌─────────────┴─────────────┐
    │                           │
┌───▼────────┐           ┌──────▼──────┐
│  应用 A      │           │  应用 B      │
│ app-a.com   │           │ app-b.com   │
└─────┬───────┘           └──────┬──────┘
      │                          │
      │  3. 重定向                 │  3. 重定向
      │  /authorize               │  /authorize
      ▼                          ▼
┌─────────────────────────────────────────┐
│               IAM 系统                    │
│            iam.platform.com              │
│                                          │
│  ┌─────────────────────────────────────┐ │
│  │ OAuth2 / OIDC 端点                  │ │
│  │  /v1/iam/oidc/authorize            │ │
│  │  /v1/iam/oidc/token                │ │
│  │  /v1/iam/oidc/userinfo             │ │
│  │  /v1/iam/oidc/logout               │ │
│  │  /v1/iam/oidc/revoke               │ │
│  └─────────────────────────────────────┘ │
│                                          │
│  ┌─────────────────────────────────────┐ │
│  │ SSO Session 管理 (Redis + DB)       │ │
│  └─────────────────────────────────────┘ │
│                                          │
│  ┌─────────────────────────────────────┐ │
│  │ 权限管理                            │ │
│  │  /v1/iam/user/permissions          │ │
│  └─────────────────────────────────────┘ │
└─────────────────────────────────────────┘
```

## 二、OAuth2 授权流程

### 2.1 首次登录（无 SSO 会话）

```
1. 用户访问 app-a.com → 未登录 → 重定向到 IAM
   GET /v1/iam/oidc/authorize?
       response_type=code&
       client_id=app_a_client_id&
       redirect_uri=https://app-a.com/callback&
       scope=openid%20profile%20offline_access&
       state=random_state&
       code_challenge=pkce_challenge&
       code_challenge_method=S256

2. IAM 检查 SSO Session (Redis + DB)
   - 无会话 → 显示登录页

3. 用户登录成功后（SSO Session 创建），重定向回应用
   GET https://app-a.com/callback?code=xxx&state=xxx

4. 应用后端用 Code 换 Token
   POST /v1/iam/oidc/token
   Headers: Authorization: Basic base64(client_id:client_secret)
   Body: grant_type=authorization_code&code=xxx&redirect_uri=xxx&code_verifier=xxx

5. IAM 返回 Token
   {
     "access_token": "xxx",
     "refresh_token": "xxx",
     "id_token": "xxx",
     "token_type": "Bearer",
     "expires_in": 3600
   }

6. 应用用 access_token 调用权限接口
   GET /v1/iam/user/permissions
   Headers: Authorization: Bearer access_token

7. IAM 返回该用户在当前应用下的权限
   {
     "tenant_id": 1,
     "app_code": "app_a",
     "roles": [...],
     "menus": [...],
     "permissions": [...]
   }
```

### 2.2 静默登录（有 SSO 会话）

```
用户 → 应用B → IAM(检测SSO有效) → 直接重定向Code → 应用B换Token → 获取权限
```

### 2.3 单点登出

```
用户 → 应用A登出 → 调用IAM logout → 清除SSO Session + Token → 重定向IAM登出页
用户 → 应用B → 检测未登录 → 重新走授权流程
```

## 三、数据模型

### 3.1 `iam_application` 表更新 - 新增字段

在现有 `iam_application` 表中增加 OAuth Client 相关字段：

| 字段 | 类型 | 说明 |
|------|------|------|
| client_id | varchar(64) | OAuth2 Client ID（唯一） |
| client_secret | varchar(255) | OAuth2 Client Secret（加密存储） |
| client_type | varchar(16) | 客户端类型: web/app/spa/mini |
| pkce_required | tinyint(1) | 是否强制 PKCE：0-否 1-是 |
| allowed_scopes | varchar(255) | 允许的 scopes，逗号分隔 |
| allowed_callbacks | text | 允许的重定向 URI（JSON 数组） |

### 3.2 `iam_sso_session` 表

| 字段 | 类型 | 说明 |
|------|------|------|
| session_id | varchar(64) | SSO 会话 ID（唯一） |
| person_id | bigint | 自然人 ID |
| org_id | bigint | 组织 ID |
| login_time | datetime(3) | 登录时间 |
| last_active_time | datetime(3) | 最后活跃时间 |
| expires_at | datetime(3) | 过期时间 |

### 3.3 `iam_auth_code` 表

| 字段 | 类型 | 说明 |
|------|------|------|
| code | varchar(64) | 授权码（唯一） |
| client_id | varchar(64) | Client ID |
| person_id | bigint | 自然人 ID |
| tenant_id | bigint | 租户 ID |
| org_id | bigint | 组织 ID |
| redirect_uri | varchar(255) | 重定向 URI |
| scope | varchar(255) | 请求的 scope |
| state | varchar(128) | state 参数，防 CSRF |
| code_challenge | varchar(64) | PKCE code_challenge |
| code_challenge_method | varchar(8) | PKCE challenge 方法 |
| expires_at | datetime(3) | 过期时间（10 分钟） |
| used | tinyint(1) | 是否已使用 |

### 3.4 `iam_token` 表

| 字段 | 类型 | 说明 |
|------|------|------|
| token_id | varchar(64) | Token 唯一标识 |
| person_id | bigint | 自然人 ID |
| user_id | bigint | 用户 ID |
| client_id | varchar(64) | Client ID |
| tenant_id | bigint | 租户 ID |
| org_id | bigint | 组织 ID |
| token_type | varchar(16) | Token 类型: access/refresh/id |
| access_token_hash | varchar(128) | Access Token 哈希 |
| refresh_token_hash | varchar(128) | Refresh Token 哈希 |
| scopes | varchar(255) | 授权的 scopes |
| expires_at | datetime(3) | 过期时间 |
| revoked | tinyint(1) | 是否撤销 |

### 3.5 `iam_login_log` 表更新 - 增加字段

| 字段 | 类型 | 说明 |
|------|------|------|
| logout_time | datetime(3) | 退出时间 |
| session_id | varchar(64) | 关联的 SSO 会话 ID |

## 四、API 端点设计

### 4.1 OAuth2 / OIDC 端点

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/v1/iam/oidc/authorize` | 授权入口 |
| POST | `/v1/iam/oidc/token` | 获取 Token |
| POST | `/v1/iam/oidc/token/refresh` | 刷新 Token |
| GET | `/v1/iam/oidc/userinfo` | 获取用户信息 |
| POST | `/v1/iam/oidc/logout` | 单点登出 |
| POST | `/v1/iam/oidc/revoke` | 撤销 Token |

### 4.2 用户权限端点

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/v1/iam/user/permissions` | 获取当前用户权限 |
| GET | `/v1/iam/user/info` | 获取当前用户信息 |

## 五、Scope 设计

| Scope | 说明 |
|-------|------|
| openid | OIDC 必选，表示这是 OIDC 请求 |
| profile | 获取用户基本信息（name, email） |
| offline_access | 获取 refresh_token |

## 六、Token 设计

### Access Token Claims

```json
{
  "iss": "iam.platform.com",
  "sub": "person_123",
  "aud": "app_a_client_id",
  "exp": 1699999999,
  "iat": 1699996399,
  "jti": "token_xxx",
  "client_id": "app_a_client_id",
  "person_id": 123,
  "org_id": 1,
  "tenant_id": 1,
  "user_id": 456,
  "scope": "openid profile"
}
```

### ID Token Claims

```json
{
  "iss": "iam.platform.com",
  "sub": "person_123",
  "aud": "app_a_client_id",
  "exp": 1699999999,
  "iat": 1699996399,
  "name": "张三",
  "email": "zhangsan@example.com"
}
```

### Token 过期时间

- Access Token：1 小时
- Refresh Token：7 天
- SSO Session：24 小时

## 七、核心流程时序

### 7.1 首次登录（无 SSO 会话）

```
用户 → 应用A → IAM登录页 → 登录成功 → IAM创建SSO Session → 重定向Code → 应用A换Token → 获取权限
```

### 7.2 静默登录（有 SSO 会话）

```
用户 → 应用B → IAM(检测SSO有效) → 直接重定向Code → 应用B换Token → 获取权限
```

### 7.3 单点登出

```
用户 → 应用A登出 → 调用IAM logout → 清除SSO Session + Token + 清除登录失败计数 → 重定向IAM登出页
用户 → 应用B → 检测未登录 → 重新走授权流程
```

## 八、安全考虑

1. **PKCE**：前端场景必须使用，防止 Code 拦截攻击
2. **redirect_uri 校验**：必须与注册时完全匹配
3. **state 参数**：防止 CSRF 攻击
4. **client_secret**：必须安全存储，服务端不暴露
5. **Token 撤销**：支持主动撤销
6. **SSO Session**：存储于 Redis + DB，支持分布式

## 九、数据库变更汇总

```sql
-- 1. iam_application 表新增字段
ALTER TABLE iam_application
ADD COLUMN client_id VARCHAR(64) UNIQUE COMMENT 'OAuth2 Client ID',
ADD COLUMN client_secret VARCHAR(255) COMMENT 'OAuth2 Client Secret',
ADD COLUMN client_type VARCHAR(16) DEFAULT 'web' COMMENT '客户端类型: web/app/spa/mini',
ADD COLUMN pkce_required TINYINT(1) DEFAULT 0 COMMENT '是否强制PKCE',
ADD COLUMN allowed_scopes VARCHAR(255) DEFAULT 'openid,profile' COMMENT '允许的scopes',
ADD COLUMN allowed_callbacks TEXT COMMENT '允许的重定向URI，JSON数组';

-- 2. 新建 iam_sso_session 表
CREATE TABLE IF NOT EXISTS iam_sso_session (
    id BIGINT AUTO_INCREMENT COMMENT '主键ID',
    session_id VARCHAR(64) UNIQUE NOT NULL COMMENT 'SSO会话ID',
    person_id BIGINT NOT NULL DEFAULT 0 COMMENT '自然人ID',
    org_id BIGINT NOT NULL DEFAULT 0 COMMENT '组织ID',
    login_time DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '登录时间',
    last_active_time DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '最后活跃时间',
    expires_at DATETIME(3) NOT NULL COMMENT '过期时间',
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at DATETIME(3) NULL COMMENT '删除时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_session_id (session_id),
    INDEX idx_person_id (person_id),
    INDEX idx_org_id (org_id),
    INDEX idx_expires_at (expires_at),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='SSO会话表';

-- 3. 新建 iam_auth_code 表
CREATE TABLE IF NOT EXISTS iam_auth_code (
    id BIGINT AUTO_INCREMENT COMMENT '主键ID',
    code VARCHAR(64) UNIQUE NOT NULL COMMENT '授权码',
    client_id VARCHAR(64) NOT NULL COMMENT 'Client ID',
    person_id BIGINT NOT NULL DEFAULT 0 COMMENT '自然人ID',
    tenant_id BIGINT NOT NULL DEFAULT 0 COMMENT '租户ID',
    org_id BIGINT NOT NULL DEFAULT 0 COMMENT '组织ID',
    redirect_uri VARCHAR(255) NOT NULL COMMENT '重定向URI',
    scope VARCHAR(255) DEFAULT 'openid,profile' COMMENT '请求的scope',
    state VARCHAR(128) COMMENT 'state参数，防CSRF',
    code_challenge VARCHAR(64) COMMENT 'PKCE code_challenge',
    code_challenge_method VARCHAR(8) DEFAULT 'S256' COMMENT 'PKCE challenge方法',
    expires_at DATETIME(3) NOT NULL COMMENT '过期时间',
    used TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否已使用：0-否 1-是',
    used_at DATETIME(3) NULL COMMENT '使用时间',
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at DATETIME(3) NULL COMMENT '删除时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_code (code),
    INDEX idx_client_id (client_id),
    INDEX idx_person_id (person_id),
    INDEX idx_expires_at (expires_at),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='授权码表';

-- 4. 新建 iam_token 表
CREATE TABLE IF NOT EXISTS iam_token (
    id BIGINT AUTO_INCREMENT COMMENT '主键ID',
    token_id VARCHAR(64) UNIQUE NOT NULL COMMENT 'Token唯一标识',
    person_id BIGINT NOT NULL DEFAULT 0 COMMENT '自然人ID',
    user_id BIGINT NOT NULL DEFAULT 0 COMMENT '用户ID',
    client_id VARCHAR(64) NOT NULL COMMENT 'Client ID',
    tenant_id BIGINT NOT NULL DEFAULT 0 COMMENT '租户ID',
    org_id BIGINT NOT NULL DEFAULT 0 COMMENT '组织ID',
    token_type VARCHAR(16) NOT NULL COMMENT 'Token类型: access/refresh/id',
    access_token_hash VARCHAR(128) COMMENT 'Access Token哈希',
    refresh_token_hash VARCHAR(128) COMMENT 'Refresh Token哈希',
    scopes VARCHAR(255) DEFAULT 'openid,profile' COMMENT '授权的scopes',
    expires_at DATETIME(3) NOT NULL COMMENT '过期时间',
    revoked TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否撤销：0-否 1-是',
    revoked_at DATETIME(3) NULL COMMENT '撤销时间',
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at DATETIME(3) NULL COMMENT '删除时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_token_id (token_id),
    INDEX idx_person_id (person_id),
    INDEX idx_client_id (client_id),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_expires_at (expires_at),
    INDEX idx_revoked (revoked),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Token表';

-- 5. iam_login_log 表新增字段
ALTER TABLE iam_login_log
ADD COLUMN logout_time DATETIME(3) NULL COMMENT '退出时间',
ADD COLUMN session_id VARCHAR(64) COMMENT '关联的SSO会话ID';
```

## 十、模块结构

```
apps/iam/
├── internal/
│   ├── service/svcuser/
│   │   └── user.go              # 用户相关服务
│   ├── service/svcauth/
│   │   └── auth.go              # 密码登录服务（保留现有）
│   ├── service/svcpermission/
│   │   └── role.go              # 权限相关服务
│   └── NEW: service/svcoidc/    # OIDC 服务
│       ├── authorize.go         # 授权端点服务
│       ├── token.go             # Token 服务
│       ├── userinfo.go          # UserInfo 服务
│       ├── logout.go            # 登出服务
│       └── session.go           # SSO Session 服务
├── model/
│   ├── application.go           # 更新：增加 OAuth 字段
│   ├── user.go
│   └── NEW: sso_session.go      # SSO Session 模型
│   └── NEW: auth_code.go        # Auth Code 模型
│   └── NEW: token.go            # Token 模型
├── dao/
│   ├── user.go
│   └── NEW: sso_session.go      # SSO Session DAO
│   └── NEW: auth_code.go        # Auth Code DAO
│   └── NEW: token.go            # Token DAO
├── dto/dtoauth/                 # 保留现有
├── dto/dtoauth/                 # 保留现有
│   └── NEW: dtooidc/            # OIDC DTO
│       ├── authorize.go         # 授权请求/响应
│       ├── token.go             # Token 请求/响应
│       ├── userinfo.go          # UserInfo 响应
│       └── logout.go            # 登出请求/响应
├── controller/ctroidc/         # OIDC 控制器
│   ├── authorize.go             # 授权端点
│   ├── token.go                 # Token 端点
│   ├── userinfo.go             # UserInfo 端点
│   ├── logout.go               # 登出端点
│   └── revoke.go               # 撤销端点
├── router/
│   └── oidc.go                  # OIDC 路由注册
```
