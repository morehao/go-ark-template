# APIKey 设计调整实现计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 将 APIKey 从 RSA 密钥对方案改为随机字符串方案，支持用户查看明文

**Architecture:** APIKey 使用随机字符串格式，通过 AES-256-GCM 加密存储，MasterKey 从配置获取

**Tech Stack:** Gin + GORM + gcrypto(AES) + MySQL

---

## 文件变更总览

| 层级 | 文件 | 变更 |
|------|------|------|
| Model | `apps/iam/model/api_key.go` | 删除字段、新增字段 |
| DTO | `apps/iam/internal/dto/dtoapikey/request.go` | 新增 `ApiKey` 字段 |
| Service | `apps/iam/internal/service/svcapikey/api_key.go` | 重写加密逻辑 |

---

## Task 1: 修改 Model 层

**Files:**
- Modify: `apps/iam/model/api_key.go`

- [ ] **Step 1: 修改 ApiKeyEntity 结构**

旧结构（删除 `PublicKey` 和 `EncryptedPrivateKey`，新增 `ApiKey`）：

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

- [ ] **Step 2: 验证 model 编译**

Run: `go build ./apps/iam/model/...`
Expected: 编译通过

- [ ] **Step 3: Commit**

```bash
git add apps/iam/model/api_key.go
git commit -m "refactor(iam): update ApiKeyEntity - replace RSA fields with ApiKey field"
```

---

## Task 2: 修改 DTO 层

**Files:**
- Modify: `apps/iam/internal/dto/dtoapikey/request.go`

- [ ] **Step 1: 修改 ApiKeyCreateResp**

将 `EncryptedPrivateKey` 替换为 `ApiKey`：

```go
type ApiKeyCreateResp struct {
    ID                 uint   `json:"id"`
    KeyPrefix          string `json:"keyPrefix"`
    ApiKey             string `json:"apiKey"`    // 新增：完整 API Key 明文
    ExpiresAt          string `json:"expiresAt"`
}
```

- [ ] **Step 2: 修改 ApiKeyListItem**

新增 `ApiKey` 字段：

```go
type ApiKeyListItem struct {
    ID                 uint   `json:"id"`
    AppID              uint   `json:"appID"`
    KeyName            string `json:"keyName"`
    KeyPrefix          string `json:"keyPrefix"`
    ApiKey             string `json:"apiKey"`    // 新增：解密后的完整 API Key
    Scopes             string `json:"scopes"`
    AccessPolicy       string `json:"accessPolicy"`
    AllowedIPs         string `json:"allowedIPs"`
    Status             string `json:"status"`
    LastUsedAt         string `json:"lastUsedAt"`
    ExpiresAt          string `json:"expiresAt"`
    CreatedAt          int64  `json:"createdAt"`
}
```

- [ ] **Step 3: 验证 dto 编译**

Run: `go build ./apps/iam/internal/dto/...`
Expected: 编译通过

- [ ] **Step 4: Commit**

```bash
git add apps/iam/internal/dto/dtoapikey/request.go
git commit -m "feat(iam): update DTO - add ApiKey field to resp and list item"
```

---

## Task 3: 重写 Service 层

**Files:**
- Modify: `apps/iam/internal/service/svcapikey/api_key.go`

- [ ] **Step 1: 修改 import**

删除不再使用的 `crypto/rsa`, `crypto/x509`, `encoding/pem`，保留 `crypto/rand`（用于随机字符串生成），使用 `gcrypto.NewAES`

- [ ] **Step 2: 修改 Create 方法**

替换 RSA 密钥生成逻辑为随机字符串生成和 AES 加密：

```go
func (svc *apiKeySvc) Create(ctx *gin.Context, req *dtoapikey.ApiKeyCreateReq) (*dtoapikey.ApiKeyCreateResp, error) {
    tenantID := gincontext.GetTenantID(ctx)
    userID := gincontext.GetUserID(ctx)
    operatorID := gincontext.GetUserID(ctx)

    // 生成随机 API Key
    apiKeyPlain := svc.generateApiKey()

    // AES 加密存储
    aesEncryptor, err := gcrypto.NewAES(config.Conf.MasterKey)
    if err != nil {
        glog.Errorf(ctx, "[svcapikey.Create] NewAES fail, err:%v", err)
        return nil, code.GetError(code.ApiKeyCreateError)
    }

    encryptedApiKey, err := aesEncryptor.EncryptString(apiKeyPlain)
    if err != nil {
        glog.Errorf(ctx, "[svcapikey.Create] Encrypt API Key fail, err:%v", err)
        return nil, code.GetError(code.ApiKeyCreateError)
    }

    keyPrefix := "ark_" + apiKeyPlain[4:12] // ark_ + 8字符

    insertEntity := &model.ApiKeyEntity{
        TenantID:     tenantID,
        UserID:       userID,
        AppID:        req.AppID,
        KeyName:      req.KeyName,
        KeyPrefix:    keyPrefix,
        ApiKey:       encryptedApiKey,
        AccessPolicy: model.ApiKeyAccessPolicy(req.AccessPolicy),
        AllowedIPs:   req.AllowedIPs,
        Scopes:       req.Scopes,
        ExpiresAt:    req.ExpiresAt,
        Status:       model.ApiKeyStatusEnabled,
        CreatedBy:    operatorID,
        UpdatedBy:    operatorID,
    }

    if err := dao.NewApiKeyDao().Insert(ctx, insertEntity); err != nil {
        glog.Errorf(ctx, "[svcapikey.Create] Insert fail, err:%v, req:%s", err, gutil.ToJsonString(req))
        return nil, code.GetError(code.ApiKeyCreateError)
    }

    return &dtoapikey.ApiKeyCreateResp{
        ID:        insertEntity.ID,
        KeyPrefix: keyPrefix,
        ApiKey:    apiKeyPlain,
        ExpiresAt: req.ExpiresAt,
    }, nil
}
```

- [ ] **Step 3: 修改 List 方法**

新增解密 API Key 并返回明文：

```go
func (svc *apiKeySvc) List(ctx *gin.Context, req *dtoapikey.ApiKeyListReq) (*dtoapikey.ApiKeyListResp, error) {
    tenantID := gincontext.GetTenantID(ctx)

    cond := &dao.ApiKeyCond{
        BaseCond: &genericdao.BaseCond{
            Page:     req.Page,
            PageSize: req.PageSize,
        },
        TenantID: tenantID,
        AppID:    req.AppID,
        KeyName:  req.KeyName,
        Status:   model.ApiKeyStatus(req.Status),
    }

    apiKeyList, total, err := dao.NewApiKeyDao().GetPageListByCond(ctx, cond)
    if err != nil {
        glog.Errorf(ctx, "[svcapikey.List] GetPageListByCond fail, err:%v, req:%s", err, gutil.ToJsonString(req))
        return nil, code.GetError(code.ApiKeyGetPageListError)
    }

    // AES 解密器
    aesDecryptor, err := gcrypto.NewAES(config.Conf.MasterKey)
    if err != nil {
        glog.Errorf(ctx, "[svcapikey.List] NewAES fail, err:%v", err)
        return nil, code.GetError(code.ApiKeyGetPageListError)
    }

    list := make([]dtoapikey.ApiKeyListItem, 0, len(apiKeyList))
    for _, v := range apiKeyList {
        // 解密 API Key
        apiKeyPlain, err := aesDecryptor.DecryptString(v.ApiKey)
        if err != nil {
            glog.Errorf(ctx, "[svcapikey.List] Decrypt API Key fail, err:%v, id:%d", err, v.ID)
            apiKeyPlain = "" // 解密失败时返回空
        }

        list = append(list, dtoapikey.ApiKeyListItem{
            ID:           v.ID,
            AppID:        v.AppID,
            KeyName:      v.KeyName,
            KeyPrefix:    v.KeyPrefix,
            ApiKey:       apiKeyPlain,
            Scopes:       v.Scopes,
            AccessPolicy: string(v.AccessPolicy),
            AllowedIPs:   v.AllowedIPs,
            Status:       string(v.Status),
            LastUsedAt:   v.LastUsedAt,
            ExpiresAt:    v.ExpiresAt,
        })
    }

    return &dtoapikey.ApiKeyListResp{
        List:  list,
        Total: total,
    }, nil
}
```

- [ ] **Step 4: 修改 generateKeyPrefix 方法为 generateApiKey**

替换原有方法，生成随机字符串格式 API Key：

```go
func (svc *apiKeySvc) generateApiKey() string {
    const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
    result := make([]byte, 16)
    for i := range result {
        b := make([]byte, 1)
        rand.Read(b)
        result[i] = charset[int(b[0])%len(charset)]
    }
    return "ark_" + string(result)
}
```

- [ ] **Step 5: 删除无用方法**

删除 `encodePrivateKeyToPEM` 和 `encodePublicKeyToPEM` 函数

- [ ] **Step 6: 验证 service 编译**

Run: `go build ./apps/iam/internal/service/svcapikey/...`
Expected: 编译通过

- [ ] **Step 7: Commit**

```bash
git add apps/iam/internal/service/svcapikey/api_key.go
git commit -m "refactor(iam): replace RSA keypair with random string API Key and AES encryption"
```

---

## Task 4: 验证整体构建

- [ ] **Step 1: 运行 lint**

Run: `make lint APP=iam`
Expected: 无 lint 错误

- [ ] **Step 2: 运行测试**

Run: `make test APP=iam`
Expected: 所有测试通过

- [ ] **Step 3: 构建应用**

Run: `make build APP=iam`
Expected: 构建成功

---

## 实施检查清单

- [x] Task 1: Model 层修改完成
- [x] Task 2: DTO 层修改完成
- [x] Task 3: Service 层修改完成
- [x] Task 4: 整体验证通过

---

## 依赖信息

- **gcrypto AES**: `/Users/morehao/Documents/practice/go/golib/gcrypto/aes.go`
  - `gcrypto.NewAES(key)` - 创建 AES 加密器
  - `EncryptString(plaintext)` - 加密并返回 base64
  - `DecryptString(ciphertext)` - 解密 base64 密文