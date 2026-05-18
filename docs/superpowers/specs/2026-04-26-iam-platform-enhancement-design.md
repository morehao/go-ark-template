# IAM 平台能力补全设计方案

**日期**: 2026-04-26
**版本**: v1.0.0

---

## 一、背景与目标

当前 IAM 平台实现了一部分基础能力，但仍存在：
1. 数据表设计规范不统一
2. 已实现功能存在 Bug
3. 缺少日志服务
4. 缺少 API 密钥管理

本设计方案旨在补全这些能力，建立规范的 IAM 基础设施。

---

## 二、实施方案

### 2.1 统一表设计规范

**原则**: 以 `iam_tenant` 表为准，统一所有 Model 定义。

#### 2.1.1 修正 SQL 初始化脚本

**文件**: `apps/iam/scripts/sql/v1.0.0_2__insert_data.sql`

| 问题 | 修正 |
|------|------|
| `iam_org` | → `iam_organization` |
| `'active'` | → `'enabled'` |
| `tenant_code` | 确认字段名正确 |

#### 2.1.2 统一 Model 定义

| 文件 | 字段 | 当前 | 修正 |
|------|------|------|------|
| `model/person.go` | Remark | VARCHAR(255) | VARCHAR(500) |
| `model/user.go` | LastLoginAt | `time.Time` | `*time.Time` |
| 所有 model | gorm.Model | 隐式 | 显式字段定义 |

#### 2.1.3 状态值统一

```
所有 status 字段统一:
  - enabled: 启用
  - disabled: 停用

租户额外状态:
  - trial: 试用
  - expired: 已过期

用户额外状态:
  - locked: 锁定
```

---

### 2.2 修复现有 Bug

| # | 文件 | 问题 | 修复方案 |
|---|------|------|----------|
| 1 | `svcuser/user.go:185` | 调试代码 `RedisCli.Set` | 删除该行 |
| 2 | `core/user/user.go` | 用户名唯一性未检查 | 创建用户前校验租户内用户名唯一 |
| 3 | `svcauth/auth.go:Register` | 租户 `tenant_path` 未更新 | 添加路径更新逻辑 |
| 4 | `svcpermission/role.go` | 创建角色缺 `tenant_id` | 从 context 获取 tenant_id |

---

### 2.3 操作日志 + 登录日志服务

#### 2.3.1 表结构（已有）

```sql
-- 操作日志表
CREATE TABLE iam_operation_log (...);

-- 登录日志表
CREATE TABLE iam_login_log (...);
```

#### 2.3.2 服务接口

**OperationLogSvc**
```go
type OperationLogSvc interface {
    Create(ctx *gin.Context, req *OperationLogCreateReq) error
    PageList(ctx *gin.Context, req *OperationLogPageListReq) (*OperationLogPageListResp, error)
}
```

**LoginLogSvc**
```go
type LoginLogSvc interface {
    Create(ctx *gin.Context, req *LoginLogCreateReq) error
    PageList(ctx *gin.Context, req *LoginLogPageListReq) (*LoginLogPageListResp, error)
}
```

#### 2.3.3 记录点

| 日志类型 | 记录位置 |
|----------|----------|
| 登录日志 | AuthSvc.LoginByPassword, Logout, RefreshToken |
| 操作日志 | 各 Service 的 Create/Update/Delete 方法 |

---

### 2.4 API 密钥管理

#### 2.4.1 表结构

```sql
CREATE TABLE iam_api_key (
    id BIGINT AUTO_INCREMENT COMMENT 'ID',
    tenant_id BIGINT NOT NULL COMMENT '租户ID',
    user_id BIGINT NOT NULL COMMENT '关联用户ID',
    app_id BIGINT NOT NULL DEFAULT 0 COMMENT '应用ID',
    key_name VARCHAR(64) NOT NULL COMMENT '密钥名称',
    key_prefix VARCHAR(16) NOT NULL COMMENT '密钥前缀',
    public_key TEXT NOT NULL COMMENT '公钥',
    encrypted_private_key TEXT NOT NULL COMMENT '加密私钥',
    access_policy VARCHAR(16) DEFAULT 'all' COMMENT '访问策略',
    allowed_ips TEXT COMMENT '允许的IP列表',
    scopes VARCHAR(255) COMMENT '权限范围',
    status VARCHAR(16) DEFAULT 'enabled' COMMENT '状态',
    last_used_at DATETIME(3) COMMENT '最后使用时间',
    expires_at DATETIME(3) COMMENT '过期时间',
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at DATETIME(3) NULL,
    created_by BIGINT NOT NULL DEFAULT 0,
    updated_by BIGINT NOT NULL DEFAULT 0,
    deleted_by BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (id),
    UNIQUE KEY uk_tenant_key_prefix (tenant_id, key_prefix),
    INDEX idx_tenant_user (tenant_id, user_id),
    INDEX idx_status (status)
) COMMENT='API密钥表';
```

#### 2.4.2 密钥生成流程

1. 生成 RSA 2048-bit 密钥对
2. 公钥明文存储
3. 私钥使用主密钥加密后存储
4. 返回 `key_prefix + encrypted_private_key` 给用户（一次性显示）

#### 2.4.3 验证流程

1. 客户端请求携带 `X-Api-Key: {key_prefix}{random_part}`
2. 服务端根据 prefix 找到 ApiKey 记录
3. 解密私钥
4. 使用私钥验签
5. 验证 IP、Scopes、过期时间

#### 2.4.4 服务接口

```go
type ApiKeySvc interface {
    Create(ctx *gin.Context, req *ApiKeyCreateReq) (*ApiKeyCreateResp, error)
    Delete(ctx *gin.Context, req *ApiKeyDeleteReq) error
    List(ctx *gin.Context, req *ApiKeyListReq) (*ApiKeyListResp, error)
    Disable(ctx *gin.Context, req *ApiKeyDisableReq) error
    Enable(ctx *gin.Context, req *ApiKeyEnableReq) error
}
```

---

## 三、实现顺序

| 顺序 | 任务 | 说明 |
|------|------|------|
| 1 | 统一表设计规范 | 基础设施 |
| 2 | 修复现有 Bug | 在规范基础上修复 |
| 3 | 操作日志 + 登录日志 | 基于修复后的代码 |
| 4 | API 密钥管理 | 开放平台扩展 |

---

## 四、验收标准

- [ ] SQL 初始化脚本可正常执行
- [ ] 所有 Model 与数据库表结构一致
- [ ] Bug 修复后功能正常
- [ ] 日志服务可正常记录
- [ ] API 密钥可正常创建和验证
- [ ] `make test APP=iam` 通过
- [ ] `make lint` 通过
