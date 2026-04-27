# IAM 业务流程文档

## 一、项目概述

### 1.1 系统定位

IAM（Identity and Access Management）是 GoArk 平台的身份认证与访问管理服务，为多租户 SaaS 系统提供用户身份认证、权限管理、单点登录（SSO）能力。

### 1.2 核心概念

| 概念 | 说明 | 业务特点 |
|------|------|----------|
| **Organization** | 组织 | 顶层实体，代表企业/机构，通过域名（Domain）区分 |
| **Tenant** | 租户 | 属于组织，可支持多租户场景，含租户路径（tenant_path）支持层级 |
| **Person** | 自然人 | 平台级用户实体，存储身份信息（手机、邮箱、密码哈希），一个自然人可关联多个租户的用户账号 |
| **User** | 用户账号 | 租户级实体，代表某自然人在某租户下的登录账号 |
| **Department** | 部门 | 属于租户，支持树形层级（parent_id + dept_path） |
| **Role** | 角色 | 属于租户，支持数据权限范围（all/dept_and_sub/dept/self/custom） |
| **Menu** | 菜单 | 属于租户，支持访问策略控制（public/authorized/org_admin/tenant_admin） |
| **Application** | 应用 | OIDC 客户端应用配置，用于第三方应用 SSO |

### 1.3 模块架构图

```mermaid
graph TB
    subgraph "展示层"
        WebApp["Web前端"]
        MobileApp["移动App"]
        ThirdPartyApp["第三方应用"]
    end

    subgraph "IAM服务"
        Auth["认证模块<br/>/auth/*"]
        OIDC["OIDC/SSO模块<br/>/oidc/*"]
        User["用户模块<br/>/user/*"]
        Org["组织模块<br/>/organization/*<br/>/tenant/*<br/>/department/*"]
        Permission["权限模块<br/>/menu/*<br/>/role/*"]
        App["应用模块<br/>/application/*"]
    end

    subgraph "数据层"
        Person["Person表"]
        UserTable["User表"]
        Tenant["Tenant表"]
        OrgTable["Organization表"]
        Dept["Department表"]
        Role["Role表"]
        Menu["Menu表"]
        AppTable["Application表"]
    end

    WebApp --> Auth
    MobileApp --> Auth
    ThirdPartyApp --> OIDC

    Auth --> UserTable
    Auth --> Person
    Auth --> Tenant
    Auth --> OrgTable

    User --> UserTable
    User --> Person
    User --> Dept
    User --> Role

    Org --> OrgTable
    Org --> Tenant
    Org --> Dept

    Permission --> Role
    Permission --> Menu
    Permission --> UserTable

    App --> AppTable
    OIDC --> AppTable
    OIDC --> Person
```

---

## 二、数据模型关系

### 2.1 Entity 关系图

```mermaid
erDiagram
    ORGANIZATION ||--o{ TENANT : "has"
    ORGANIZATION {
        uint id PK
        string domain UK
        string org_name
        OrgStatus status
    }

    TENANT ||--o{ USER : "contains"
    TENANT ||--o{ DEPARTMENT : "contains"
    TENANT ||--o{ ROLE : "contains"
    TENANT ||--o{ MENU : "contains"
    TENANT ||--o{ APPLICATION : "registers"
    TENANT {
        uint id PK
        uint org_id FK
        string tenant_name
        string tenant_path
        TenantStatus status
    }

    PERSON ||--o{ USER : "identifies"
    PERSON {
        uint id PK
        string mobile
        string email
        string password_hash
        string real_name
    }

    USER ||--o{ USER_DEPARTMENT : "belongs_to"
    USER ||--o{ USER_ROLE : "has"
    USER {
        uint id PK
        uint tenant_id FK
        uint person_id FK
        uint dept_id
        string username
        UserStatus status
        UserType user_type
    }

    DEPARTMENT ||--o{ USER_DEPARTMENT : "contains"
    DEPARTMENT {
        uint id PK
        uint tenant_id FK
        uint parent_id
        string dept_path
    }

    USER_DEPARTMENT {
        uint id PK
        uint user_id FK
        uint dept_id FK
        UserDeptType dept_type
    }

    ROLE ||--o{ USER_ROLE : "assigned_to"
    ROLE ||--o{ ROLE_MENU : "contains"
    ROLE {
        uint id PK
        uint tenant_id FK
        string role_code
        RoleDataScope data_scope
    }

    USER_ROLE {
        uint id PK
        uint user_id FK
        uint role_id FK
    }

    MENU ||--o{ ROLE_MENU : "granted_to"
    MENU {
        uint id PK
        uint tenant_id FK
        string menu_code
        MenuType menu_type
        AccessPolicy access_policy
    }

    ROLE_MENU {
        uint id PK
        uint role_id FK
        uint menu_id FK
    }

    APPLICATION {
        uint id PK
        uint tenant_id FK
        string client_id
        string client_secret
        string allowed_callbacks
    }
```

### 2.2 核心表字段说明

#### Person（自然人表）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | bigint | 主键 |
| mobile | varchar(16) | 手机号（登录账号） |
| email | varchar(64) | 邮箱（登录账号） |
| password_hash | varchar(128) | 密码哈希 |
| real_name | varchar(32) | 真实姓名 |
| avatar_url | varchar(255) | 头像URL |

#### User（用户账号表）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | bigint | 主键 |
| tenant_id | bigint | 所属租户ID |
| person_id | bigint | 关联自然人ID |
| dept_id | bigint | 主部门ID（冗余） |
| username | varchar(32) | 用户名（租户内唯一） |
| status | varchar(16) | enabled/locked/disabled |
| user_type | varchar(16) | normal/tenant_admin/platform_admin |
| login_fail_count | int | 连续登录失败次数 |
| locked_until | datetime | 账户锁定截止时间 |

#### Tenant（租户表）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | bigint | 主键 |
| org_id | bigint | 所属组织ID |
| tenant_name | varchar(128) | 租户名称 |
| tenant_code | varchar(32) | 租户编码 |
| tenant_path | varchar(512) | 租户路径（如 /1/2/3/） |
| status | varchar(16) | enabled/trial/expired/disabled |

#### Department（部门表）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | bigint | 主键 |
| tenant_id | bigint | 所属租户ID |
| parent_id | bigint | 父部门ID（0表示根部门） |
| dept_path | varchar(512) | 部门路径（如 /1/2/3/） |
| dept_level | int | 部门层级 |

---

## 三、认证授权模块

### 3.1 业务概念

#### 3.1.1 登录方式
| 方式 | 说明 |
|------|------|
| **密码登录** | 用户名（手机/邮箱）+ 密码登录，需先选择租户 |
| **OIDC/SSO** | 第三方应用通过 OAuth2/OIDC 授权码模式接入 |

#### 3.1.2 Token 类型
| 类型 | 过期时间 | 说明 |
|------|----------|------|
| **TempToken** | 10分钟 | 临时令牌，登录后选择租户前使用 |
| **AuthToken** | 30分钟 | 访问令牌（AuthToken） |
| **RefreshToken** | 7天 | 刷新令牌 |

#### 3.1.3 用户状态
| 状态 | 说明 |
|------|------|
| enabled | 正常 |
| locked | 锁定（连续登录失败后自动锁定） |
| disabled | 禁用（管理员禁用） |

### 3.2 密码登录 + 选择租户流程

#### 3.2.1 时序图

```mermaid
sequenceDiagram
    participant U as 用户
    participant C as 前端
    participant IAM as IAM服务
    participant DB as 数据库

    U->>C: 输入账号密码
    C->>IAM: POST /auth/loginByPassword<br/>{account, password}
    IAM->>DB: 查询Organization（按域名）
    DB-->>IAM: 返回Organization
    IAM->>DB: 查询Person（按手机/邮箱）
    DB-->>IAM: 返回Person
    IAM->>DB: 查询User列表（按PersonID + OrgID过滤）
    DB-->>IAM: 返回User列表
    IAM->>IAM: 检查用户锁定状态
    IAM->>DB: 验证密码哈希
    alt 密码错误
        IAM->>DB: login_fail_count++
        IAM-->>C: 密码错误
    else 密码正确
        IAM->>DB: 清零login_fail_count
        IAM->>IAM: 生成TempToken
        IAM-->>C: 返回TempToken + 租户列表<br/>{tempToken, tenantList[]}
    end

    U->>C: 选择租户
    C->>IAM: POST /auth/selectTenant<br/>{tempToken, tenantID}
    IAM->>DB: 验证TempToken
    IAM->>DB: 查询Tenant + User
    IAM->>DB: 查询用户部门、角色
    IAM->>IAM: 生成TokenPair<br/>(AuthToken + RefreshToken)
    IAM->>DB: 更新last_login_at/IP/count
    IAM-->>C: 返回TokenPair<br/>{token, refreshToken}
```

#### 3.2.2 登录流程（不含租户选择）

```mermaid
flowchart TD
    A[开始] --> B{检查账号锁定状态}
    B -->|已锁定且未过期| C[返回账户已锁定]
    B -->|未锁定或已过期| D{验证密码}
    D -->|密码错误| E[记录登录失败]
    E --> F{失败次数≥5?}
    F -->|是| G[锁定账户<br/>locked_until=now+5min]
    F -->|否| H[返回密码错误]
    G --> H
    D -->|密码正确| I[清零失败计数]
    I --> J[生成TempToken]
    J --> K[构建租户列表]
    K --> L[返回TempToken + 租户列表]
```

#### 3.2.3 API 接口

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | `/v1/iam/auth/loginByPassword` | 密码登录 | 否 |
| POST | `/v1/iam/auth/selectTenant` | 选择租户获取Token | TempToken |
| POST | `/v1/iam/auth/refreshToken` | 刷新Token | RefreshToken |
| POST | `/v1/iam/auth/logout` | 登出 | AuthToken |
| POST | `/v1/iam/auth/unlockAccount` | 自助解锁 | 否 |

### 3.3 OIDC/SSO 授权码流程

#### 3.3.1 业务概念

OIDC（OpenID Connect）基于 OAuth2 的身份认证协议，支持第三方应用通过授权码模式实现单点登录。

#### 3.3.2 时序图

```mermaid
sequenceDiagram
    participant U as 用户
    participant App as 第三方应用
    participant IAM as IAM服务
    participant DB as 数据库

    U->>App: 访问受保护资源
    App->>App: 生成state + codeVerifier
    App->>U: 重定向到 IAM授权页<br/>GET /oidc/authorize?client_id=xxx&redirect_uri=xxx&state=xxx
    U->>IAM: 打开授权页
    IAM->>U: 显示登录页（若无SSO Session）
    U->>IAM: 输入账号密码登录
    IAM->>DB: 验证登录
    IAM->>DB: 创建SSO Session
    IAM-->>App: 携带code重定向<br/>Callback?code=xxx&state=xxx
    App->>IAM: POST /oidc/token<br/>{code, client_secret, redirect_uri}
    IAM->>DB: 验证AuthCode
    IAM->>DB: 标记AuthCode已使用
    IAM->>App: 返回TokenPair<br/>{accessToken, refreshToken}
    App->>IAM: GET /oidc/userinfo<br/>Authorization: Bearer xxx
    IAM-->>App: 返回用户信息<br/>{sub, name, email, phone}
```

#### 3.3.3 授权页面流程

```mermaid
flowchart TD
    A[开始] --> B{检查X-Sso-Session-Id Header}
    B -->|有Session| C{Session有效?}
    C -->|有效| D[验证ClientID + RedirectURI]
    D -->|验证通过| E[生成AuthCode]
    E --> F[重定向到redirect_uri?code=xxx&state=xxx]
    C -->|无效| G[继续到登录]
    B -->|无Session| G
    G[要求登录] --> H[登录成功后创建SSO Session]
    H --> F
```

#### 3.3.4 Token 交换流程

```mermaid
flowchart TD
    A[开始] --> B{验证grant_type=authorization_code}
    B -->|否| C[返回不支持的grant_type]
    B -->|是| D{验证AuthCode}
    D -->|Code不存在/已使用/过期| E[返回无效Code]
    D -->|有效| F{验证RedirectURI]
    F -->|不匹配| G[返回URI不匹配]
    F -->|匹配| H{验证ClientSecret}
    H -->|不匹配| I[返回无效Client]
    H -->|匹配| J{检查PKCE}
    J -->|需要但未提供| K[返回PKCE Required]
    J -->|通过| L[标记AuthCode已使用]
    L --> M[生成AccessToken + RefreshToken]
    M --> N[返回TokenPair]
```

#### 3.3.5 API 接口

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `/v1/iam/oidc/authorize` | 授权入口 | Session/登录 |
| POST | `/v1/iam/oidc/token` | 获取Token | ClientSecret |
| POST | `/v1/iam/oidc/refreshToken` | 刷新Token | RefreshToken |
| GET | `/v1/iam/oidc/userinfo` | 获取用户信息 | AccessToken |
| POST | `/v1/iam/oidc/logout` | OIDC登出 | AuthToken |

### 3.4 Token 管理

#### 3.4.1 Token 黑名单机制

```mermaid
flowchart TD
    A[请求] --> B{检查Token类型}
    B -->|AccessToken| C[添加到AccessToken黑名单]
    B -->|RefreshToken| D[添加到RefreshToken黑名单]
    C --> E[设置过期时间]
    D --> E
    E --> F[返回成功]
```

#### 3.4.2 Token 刷新流程

```mermaid
sequenceDiagram
    participant C as 前端
    participant IAM as IAM服务
    participant Redis as Redis

    C->>IAM: POST /auth/refreshToken<br/>{refreshToken}
    IAM->>IAM: 验证RefreshToken签名
    IAM->>Redis: 检查是否在黑名单
    alt 在黑名单
        IAM-->>C: Token无效
    else 不在黑名单
        IAM->>Redis: 添加到黑名单
        IAM->>IAM: 生成新TokenPair
        IAM-->>C: 返回新TokenPair
    end
```

### 3.5 用户注册流程

#### 3.5.1 时序图

```mermaid
sequenceDiagram
    participant U as 用户
    participant C as 前端
    participant IAM as IAM服务
    participant DB as 数据库

    U->>C: 填写注册信息
    C->>IAM: POST /auth/register<br/>{username, password, email, mobile, tenantName, tenantCode}
    IAM->>DB: 查询Organization（按域名）
    IAM->>DB: 检查注册是否开放
    IAM->>DB: 验证身份类型（邮箱/手机/两者）
    alt 需要审核
        UserStatus = disabled
    else 直接启用
        UserStatus = enabled
    end
    IAM->>DB: 生成密码哈希
    IAM->>DB: 事务创建:
    Note over DB: 1. Tenant<br/>2. Department<br/>3. Person<br/>4. User<br/>5. UserDepartment
    IAM-->>C: 返回结果<br/>{tenantID, userID, status, message}
```

#### 3.5.2 API 接口

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | `/v1/iam/auth/register` | 用户注册 | 否 |

---

## 四、用户管理模块

### 4.1 业务概念

#### 4.1.1 用户类型
| 类型 | 说明 |
|------|------|
| normal | 普通用户 |
| tenant_admin | 租户管理员 |
| platform_admin | 平台管理员（可管理所有租户） |

#### 4.1.2 用户与部门关系
- 一个用户只能有一个主部门（primary）
- 一个用户可以有多个副部门（secondary）
- 关系存储在 `iam_user_department` 表

#### 4.1.3 用户与角色关系
- 用户角色分配为全量替换模式
- 关系存储在 `iam_user_role` 表

### 4.2 创建用户流程

```mermaid
sequenceDiagram
    participant A as 管理员
    participant C as 前端
    participant IAM as IAM服务
    participant DB as 数据库

    A->>C: 填写用户信息
    C->>IAM: POST /user/create<br/>{username, mobile, email, deptID, roleIDs[], ...}
    IAM->>DB: 检查用户名唯一性
    IAM->>DB: 获取或创建主部门
    alt 管理员指定密码
        IAM->>IAM: 验证密码复杂度
    else 系统生成密码
        IAM->>IAM: 生成12位随机密码
    end
    IAM->>IAM: 生成密码哈希
    IAM->>DB: 事务创建:
    Note over DB: 1. Person<br/>2. User<br/>3. UserDepartment<br/>4. UserRole
    IAM-->>C: 返回结果<br/>{userID, personID, password}
```

### 4.3 分配部门流程

```mermaid
flowchart TD
    A[开始] --> B{验证用户存在}
    B -->|不存在| C[返回用户不存在]
    B -->|存在| D{验证主部门存在且属于当前租户]
    D -->|不通过| E[返回部门不存在/范围错误]
    D -->|通过| F{验证副部门存在且属于当前租户]
    F -->|存在无效| G[返回部门不存在/范围错误]
    F -->|全部通过| H[开启事务]
    H --> I[删除现有部门关联]
    I --> J[插入主部门关联]
    J --> K[插入副部门关联]
    K --> L[提交事务]
    L --> M[更新用户主部门字段]
```

### 4.4 分配角色流程

```mermaid
flowchart TD
    A[开始] --> B[验证用户存在]
    B -->|不存在| C[返回错误]
    B -->|存在| D[开启事务]
    D --> E[删除该用户所有角色关联]
    E --> F{遍历请求的RoleIDs}
    F --> G[插入新的UserRole记录]
    G --> H{还有更多角色?}
    H -->|是| F
    H -->|否| I[提交事务]
    I --> J[返回成功]
```

### 4.5 个人中心流程

#### 5.5.1 获取当前用户信息

```mermaid
sequenceDiagram
    participant U as 用户
    participant C as 前端
    participant IAM as IAM服务
    participant DB as 数据库

    U->>C: 点击个人中心
    C->>IAM: GET /user/getCurrentUserInfo<br/>Authorization: Bearer xxx
    IAM->>DB: 获取当前用户信息
    IAM->>DB: 获取当前用户角色列表
    IAM->>DB: 获取当前用户部门列表
    IAM-->>C: 返回用户信息<br/>{userInfo, roles[], depts[]}
```

#### 5.5.2 修改密码

```mermaid
flowchart TD
    A[开始] --> B{验证旧密码}
    B -->|错误| C[返回密码错误]
    B -->|正确| D{验证新密码复杂度}
    D -->|不通过| E[返回密码复杂度错误]
    D -->|通过| F[生成新密码哈希]
    F --> G[更新Person表password_hash]
    G --> H[返回成功]
```

### 4.6 API 接口

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | `/v1/iam/user/create` | 创建用户 | AuthToken |
| POST | `/v1/iam/user/delete` | 删除用户 | AuthToken |
| POST | `/v1/iam/user/update` | 更新用户 | AuthToken |
| GET | `/v1/iam/user/detail` | 获取用户详情 | AuthToken |
| POST | `/v1/iam/user/pageList` | 分页查询用户 | AuthToken |
| POST | `/v1/iam/user/assignDepartments` | 分配部门 | AuthToken |
| GET | `/v1/iam/user/listDepartments` | 获取用户部门列表 | AuthToken |
| POST | `/v1/iam/user/assignRoles` | 分配角色 | AuthToken |
| GET | `/v1/iam/user/listRoles` | 获取用户角色列表 | AuthToken |
| GET | `/v1/iam/user/getCurrentUserInfo` | 获取当前用户信息 | AuthToken |
| POST | `/v1/iam/user/updateProfile` | 更新个人资料 | AuthToken |
| POST | `/v1/iam/user/changePassword` | 修改密码 | AuthToken |
| GET | `/v1/iam/user/loginHistory` | 登录历史 | AuthToken |
| POST | `/v1/iam/user/logout` | 登出 | AuthToken |

---

## 五、组织架构模块

### 5.1 业务概念

#### 5.1.1 层级关系

```mermaid
graph TD
    Org1["Organization (组织)"]
    Org1 --> Tenant1["Tenant (租户)"]
    Org1 --> Tenant2["Tenant (租户)"]
    Tenant1 --> Dept1["Department (部门)"]
    Tenant1 --> Dept2["Department (部门)"]
    Dept1 --> Dept3["Department (子部门)"]
```

#### 5.1.2 租户路径（tenant_path）
- 格式：`/{parent_id}/{id}/`
- 示例：`/1/5/` 表示 ID 为 5 的租户，上级为 ID=1 的租户
- 用于支持多级租户和快速查询子租户

#### 5.1.3 部门路径（dept_path）
- 格式：`/{parent_id}/{id}/`
- 示例：`/1/3/7/` 表示 ID 为 7 的部门，路径为 根部门->3->7
- 用于支持多级部门和快速查询子部门

### 5.2 API 接口

#### Organization（组织管理）
| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | `/v1/iam/organization/create` | 创建组织 | AuthToken |
| POST | `/v1/iam/organization/delete` | 删除组织 | AuthToken |
| POST | `/v1/iam/organization/update` | 更新组织 | AuthToken |
| GET | `/v1/iam/organization/detail` | 获取组织详情 | AuthToken |
| POST | `/v1/iam/organization/pageList` | 分页查询组织 | AuthToken |
| GET | `/v1/iam/organization/getOrgConfig` | 获取组织配置 | AuthToken |
| GET | `/v1/iam/organization/listConfigDefinitions` | 获取配置定义 | AuthToken |

#### Tenant（租户管理）
| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | `/v1/iam/tenant/create` | 创建租户 | AuthToken |
| POST | `/v1/iam/tenant/delete` | 删除租户 | AuthToken |
| POST | `/v1/iam/tenant/update` | 更新租户 | AuthToken |
| GET | `/v1/iam/tenant/detail` | 获取租户详情 | AuthToken |
| POST | `/v1/iam/tenant/pageList` | 分页查询租户 | AuthToken |

#### Department（部门管理）
| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | `/v1/iam/department/create` | 创建部门 | AuthToken |
| POST | `/v1/iam/department/delete` | 删除部门 | AuthToken |
| POST | `/v1/iam/department/update` | 更新部门 | AuthToken |
| GET | `/v1/iam/department/detail` | 获取部门详情 | AuthToken |
| POST | `/v1/iam/department/pageList` | 分页查询部门 | AuthToken |
| GET | `/v1/iam/department/tree` | 获取部门树 | AuthToken |

---

## 六、权限管理模块

### 6.1 业务概念

#### 6.1.1 角色类型
| 类型 | 说明 |
|------|------|
| custom | 自定义角色 |
| system | 系统内置角色 |

#### 6.1.2 数据权限范围
| 范围 | 说明 |
|------|------|
| all | 全部数据 |
| dept_and_sub | 本部门及所有子部门 |
| dept | 仅本部门 |
| self | 仅本人数据 |
| custom | 自定义数据权限 |

#### 6.1.3 菜单类型
| 类型 | 说明 |
|------|------|
| directory | 目录（不可点击） |
| menu | 菜单（可点击） |
| button | 按钮（权限点） |

#### 6.1.4 访问策略
| 策略 | 说明 |
|------|------|
| public | 所有人均可访问 |
| authorized | 需要登录 |
| org_admin | 组织管理员 |
| tenant_admin | 租户管理员 |

### 6.2 角色分配菜单流程

```mermaid
flowchart TD
    A[开始] --> B[验证角色存在]
    B -->|不存在| C[返回角色不存在]
    B -->|存在| D[开启事务]
    D --> E[删除该角色所有菜单关联]
    E --> F{遍历请求的MenuIDs}
    F --> G[插入新的RoleMenu记录]
    G --> H{还有更多菜单?}
    H -->|是| F
    H -->|否| I[提交事务]
    I --> J[返回成功]
```

### 6.3 API 接口

#### Menu（菜单管理）
| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | `/v1/iam/menu/create` | 创建菜单 | AuthToken |
| POST | `/v1/iam/menu/delete` | 删除菜单 | AuthToken |
| POST | `/v1/iam/menu/update` | 更新菜单 | AuthToken |
| GET | `/v1/iam/menu/detail` | 获取菜单详情 | AuthToken |
| POST | `/v1/iam/menu/pageList` | 分页查询菜单 | AuthToken |
| GET | `/v1/iam/menu/tree` | 获取菜单树 | AuthToken |

#### Role（角色管理）
| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | `/v1/iam/role/create` | 创建角色 | AuthToken |
| POST | `/v1/iam/role/delete` | 删除角色 | AuthToken |
| POST | `/v1/iam/role/update` | 更新角色 | AuthToken |
| GET | `/v1/iam/role/detail` | 获取角色详情 | AuthToken |
| POST | `/v1/iam/role/pageList` | 分页查询角色 | AuthToken |
| POST | `/v1/iam/role/assignMenus` | 分配菜单 | AuthToken |
| GET | `/v1/iam/role/listMenus` | 获取角色菜单列表 | AuthToken |

---

## 七、应用与 API Key 模块

### 7.1 业务概念

#### 7.1.1 Application（应用）
- 用于配置 OIDC 客户端信息
- 支持 PKCE 可选启用
- 支持配置允许的重定向 URI

#### 7.1.2 API Key
- 用于服务端 API 调用身份认证
- 支持设置过期时间

### 7.2 API 接口

#### Application（应用管理）
| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | `/v1/iam/application/create` | 创建应用 | AuthToken |
| POST | `/v1/iam/application/delete` | 删除应用 | AuthToken |
| POST | `/v1/iam/application/update` | 更新应用 | AuthToken |
| GET | `/v1/iam/application/detail` | 获取应用详情 | AuthToken |
| POST | `/v1/iam/application/pageList` | 分页查询应用 | AuthToken |

---

## 八、日志模块

### 8.1 API 接口

#### 登录日志
| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | `/v1/iam/auth/loginLog/create` | 创建登录日志 | 否 |
| GET | `/v1/iam/auth/loginLog/pageList` | 分页查询登录日志 | AuthToken |

#### 操作日志
| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | `/v1/iam/auth/operationLog/create` | 创建操作日志 | 否 |
| GET | `/v1/iam/auth/operationLog/pageList` | 分页查询操作日志 | AuthToken |

---

## 附录

### A. 错误码说明

| 错误码 | 说明 |
|--------|------|
| AuthLoginError | 登录失败 |
| AuthPasswordError | 密码错误 |
| AuthOrgNotFoundError | 组织不存在 |
| AuthPersonNotFoundError | 自然人不存在 |
| AuthNoTenantError | 没有可用的租户 |
| AuthTenantSelectError | 租户选择失败 |
| AuthRefreshTokenInvalidError | RefreshToken无效 |
| OIDCClientInvalidError | OIDC客户端无效 |
| OIDCInvalidCodeError | 授权码无效 |
| UserNotExistError | 用户不存在 |
| UsernameDuplicateError | 用户名重复 |
| DepartmentNotExistError | 部门不存在 |
| PasswordComplexityError | 密码复杂度不足 |

### B. 密码安全策略配置项

| 配置项 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| password_min_length | int | 8 | 最小长度 |
| password_require_uppercase | bool | true | 必须包含大写字母 |
| password_require_lowercase | bool | true | 必须包含小写字母 |
| password_require_number | bool | true | 必须包含数字 |
| password_require_special | bool | false | 必须包含特殊字符 |
| login_max_fail_count | int | 5 | 登录失败最大次数 |
| login_lock_duration | int | 300 | 锁定时长/秒 |
