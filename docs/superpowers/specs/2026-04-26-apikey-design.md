# APIKey 设计调整方案

## 变更概述

| 项目 | 当前方案 | 新方案 |
|------|---------|--------|
| 格式 | RSA-2048 密钥对 | 随机字符串 `ark_xxxxxxxx` |
| 存储 | 公钥明文 + 私钥用 MasterKey RSA 加密 | API Key 明文用 MasterKey AES 加密存储 |
| 查看 | 创建时返回加密私钥 | 列表/详情可返回解密后明文 |

## 背景

当前 APIKey 使用 RSA-2048 密钥对方案：
- 公钥明文存储
- 私钥用 MasterKey RSA 加密后存储
- 用户拿到的是加密私钥，无法直接用于 API 认证

需求：用户可查看明文 API Key，且密文存储。

## 数据模型

### 数据库变更

**表 `iam_api_key` 调整：**

| 字段 | 变更 | 说明 |
|------|------|------|
| `key_prefix` | 保留 | 保留用于标识 |
| `public_key` | 删除 | 不再使用 RSA 公钥 |
| `encrypted_private_key` | 删除 | 不再使用 RSA 加密私钥 |
| `api_key` | 新增 | AES 加密后的完整 API Key 密文 |

**最终表结构：**

```go
type ApiKeyEntity struct {
    ID                  uint           `gorm:"column:id;type:bigint;autoIncrement;primaryKey"`
    TenantID            uint           `gorm:"column:tenant_id;type:bigint;not null;default 0;comment: 租户ID"`
    UserID              uint           `gorm:"column:user_id;type:bigint;not null;default 0;comment: 关联用户ID"`
    AppID               uint           `gorm:"column:app_id;type:bigint;not null;default 0;comment: 应用ID"`
    KeyName             string         `gorm:"column:key_name;type:varchar(64);not null;comment: 密钥名称"`
    KeyPrefix           string         `gorm:"column:key_prefix;type:varchar(16);not null;comment: 密钥前缀(ark_开头)"`
    ApiKey              string         `gorm:"column:api_key;type:text;not null;comment: AES加密的API Key"`
    AccessPolicy        ApiKeyAccessPolicy `gorm:"column:access_policy;type:varchar(16);default all;comment: 访问策略"`
    AllowedIPs          string         `gorm:"column:allowed_ips;type:text;;comment: 允许的IP列表(JSON)"`
    Scopes              string         `gorm:"column:scopes;type:varchar(255);;comment: 权限范围"`
    Status              ApiKeyStatus   `gorm:"column:status;type:varchar(16);;default enabled;comment: 状态"`
    LastUsedAt          string         `gorm:"column:last_used_at;type:datetime(3);;comment: 最后使用时间"`
    ExpiresAt           string         `gorm:"column:expires_at;type:datetime(3);;comment: 过期时间"`
    CreatedAt           string         `gorm:"column:created_at;type:datetime(3);;comment: 创建时间"`
    UpdatedAt           string         `gorm:"column:updated_at;type:datetime(3);;comment: 更新时间"`
    DeletedAt           gorm.DeletedAt `gorm:"column:deleted_at;type:datetime(3);;comment: 删除时间"`
    CreatedBy           uint           `gorm:"column:created_by;type:bigint;not null;default 0;comment: 创建人ID"`
    UpdatedBy           uint           `gorm:"column:updated_by;type:bigint;not null;default 0;comment: 更新人ID"`
    DeletedBy           uint           `gorm:"column:deleted_by;type:bigint;not null;default 0;comment: 删除人ID"`
}
```

## 接口变更

### 1. 创建 API Key

**请求：** 不变

**响应：** 新增 `apiKey` 明文字段

```go
type ApiKeyCreateResp struct {
    ID                 uint   `json:"id"`
    KeyPrefix          string `json:"keyPrefix"`           // "ark_xxxxxxxx"
    ApiKey             string `json:"apiKey"`              // 新增：完整 API Key 明文
    ExpiresAt          string `json:"expiresAt"`
}
```

### 2. 列表查询

**响应：** 新增 `apiKey` 明文字段

```go
type ApiKeyListItem struct {
    ID                 uint   `json:"id"`
    AppID              uint   `json:"appID"`
    KeyName            string `json:"keyName"`
    KeyPrefix          string `json:"keyPrefix"`
    ApiKey             string `json:"apiKey"`              // 新增：解密后的完整 API Key
    Scopes             string `json:"scopes"`
    AccessPolicy       string `json:"accessPolicy"`
    AllowedIPs         string `json:"allowedIPs"`
    Status             string `json:"status"`
    LastUsedAt         string `json:"lastUsedAt"`
    ExpiresAt          string `json:"expiresAt"`
    CreatedAt          int64  `json:"createdAt"`
}
```

### 3. 其他接口

- 删除、启用、禁用 — 不变

## 加密实现

### API Key 生成

```
前缀: "ark_"
随机字符: 16 字节，字符集 [a-z0-9]
最终格式: "ark_" + 16 随机字符，如 "ark_a1b2c3d4e5f6g7h8"
```

### 加密存储

```
密文 = AES-256-GCM(MasterKey, API Key 明文)
```

使用 golib 中已有的 AES 加密工具。

## 实现步骤

### Step 1: Model 层

1. 修改 `model/api_key.go`：
   - 删除 `PublicKey` 和 `EncryptedPrivateKey` 字段
   - 新增 `ApiKey` 字段
   - 保留 `KeyPrefix` 字段

### Step 2: DAO 层

1. 无需变更（DAO 使用通用结构）

### Step 3: Service 层

1. 修改 `service/svcapikey/api_key.go`：
   - 删除 RSA 密钥生成逻辑
   - 实现随机字符串生成逻辑
   - 实现 AES 加密存储逻辑
   - 实现 AES 解密返回明文逻辑

### Step 4: DTO 层

1. 修改 `dto/dtoapikey/request.go`：
   - `ApiKeyCreateResp` 新增 `ApiKey` 字段
   - `ApiKeyListItem` 新增 `ApiKey` 字段

### Step 5: 数据库迁移

1. 生成数据库迁移 SQL：
   - 删除 `public_key` 列
   - 删除 `encrypted_private_key` 列
   - 新增 `api_key` 列

## 涉及文件

| 层级 | 文件 | 变更 |
|------|------|------|
| Model | `apps/iam/model/api_key.go` | 删除/新增字段 |
| Service | `apps/iam/internal/service/svcapikey/api_key.go` | 重写加密逻辑 |
| DTO | `apps/iam/internal/dto/dtoapikey/request.go` | 新增字段 |
| - | 数据库迁移 | 新增迁移文件 |

## 安全考虑

1. **密文存储**：API Key 明文通过 AES-256-GCM 加密后存储，防止数据库泄露导致 Key 泄露
2. **传输安全**：列表/详情返回明文时走 HTTPS 传输
3. **密钥管理**：MasterKey 应通过环境变量或密钥管理服务配置，不硬编码