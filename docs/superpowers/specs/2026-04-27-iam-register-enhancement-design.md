# IAM 注册接口改进设计

## 1. 背景与目标

### 当前问题

| 问题 | 现状 |
|------|------|
| 强制填租户信息 | RegisterReq 要求 `TenantName`、`TenantCode`，每次注册都创建新租户 |
| 无子域名识别 | 没有根据请求子域名自动识别租户的逻辑 |
| 无层级租户概念 | 所有租户平级，无法匹配多级结构 |
| 审核流程缺失 | `requireApproval` 只改变用户状态为 disabled，无真正的审核机制 |

### 设计目标

1. 用户注册时不再强制创建新租户，而是通过子域名精确匹配已存在的租户
2. 支持多层级租户结构，子域名精确匹配到对应层级的租户
3. 注册用户状态为 pending 时，需要子租户管理员审核后才能激活

---

## 2. 数据模型变更

### 2.1 TenantEntity 新增 Domain 字段

**文件**: `apps/iam/model/tenant.go`

```go
type TenantEntity struct {
    gorm.Model
    // ... 现有字段 ...
    Domain string `gorm:"column:domain;type:varchar(255);;default '';comment: 租户域名(用于注册时子域名匹配)"`
}
```

**用途**: 存储租户的域名，用于注册时通过子域名精确匹配租户。

### 2.2 UserStatus 新增 pending 状态

**文件**: `apps/iam/model/user.go`

```go
const (
    UserStatusEnabled  UserStatus = "enabled"
    UserStatusLocked   UserStatus = "locked"
    UserStatusDisabled UserStatus = "disabled"
    UserStatusPending  UserStatus = "pending"  // 新增：待审核
)
```

**用途**: 表示用户注册后待审核状态，需要管理员审核才能激活。

---

## 3. DTO 变更

### 3.1 RegisterReq 移除租户字段

**文件**: `apps/iam/internal/dto/dtoauth/request.go`

```go
type RegisterReq struct {
    Username string `json:"username" validate:"required" label:"用户名"`
    Password string `json:"password" validate:"required" label:"密码"`
    Mobile   string `json:"mobile" label:"手机号"`
    Email    string `json:"email" validate:"required" label:"邮箱"`
    RealName string `json:"realName" validate:"required" label:"真实姓名"`
}
```

**变更说明**:
- 移除 `TenantName` 字段
- 移除 `TenantCode` 字段
- 用户不再在注册时创建租户

### 3.2 RegisterResp 调整

**文件**: `apps/iam/internal/dto/dtoauth/response.go`

```go
type RegisterResp struct {
    UserID       uint   `json:"userId"`
    PersonID     uint   `json:"personId"`
    Status       string `json:"status"`  // pending/enabled
    PersonExists bool   `json:"personExists"`
    Message      string `json:"message"`
}
```

**变更说明**:
- 移除 `TenantID` 字段
- `Status` 返回 `pending` 或 `enabled`

---

## 4. 注册流程

### 4.1 流程图

```
┌─────────────────────────────────────────────────────────────┐
│                      Register(req)                          │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│  1. getCurrentOrg() → 获取当前组织                          │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│  2. 检查 register.enabled 配置                             │
│     └─ 若禁用 → 返回 AuthRegisterDisabled 错误             │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│  3. resolveDomainFromHost() → 从请求 Host 解析子域名        │
│     例: user.dept1.company.com → "user.dept1.company.com"  │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│  4. getTenantByDomain(orgID, domain) → 精确匹配租户         │
│     └─ 若租户不存在 → 返回"租户不存在"错误                │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│  5. 检查租户状态是否为 enabled                              │
│     └─ 若停用 → 返回"租户已停用"错误                        │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│  6. validateRegisterIdentity() → 验证身份格式(mobile/email)│
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│  7. 检查 requireApproval 配置                              │
│     ├─ 需要审核 → userStatus = "pending"                   │
│     └─ 不需要审核 → userStatus = "enabled"                 │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│  8. 事务创建:                                              │
│     Person → User(关联到已存在的租户和部门) → UserDepartment│
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│  返回 RegisterResp                                          │
└─────────────────────────────────────────────────────────────┘
```

### 4.2 域名匹配规则

- 子域名精确匹配租户的 `domain` 字段
- 匹配范围: `user.dept1.company.com` → Tenant.domain = "user.dept1.company.com"
- 多级租户示例:
  - `company.com` → 顶级租户 (TenantLevel = 1)
  - `dept1.company.com` → 二级租户 (TenantLevel = 2, ParentID = 顶级租户ID)
  - `user.dept1.company.com` → 三级租户 (TenantLevel = 3, ParentID = 二级租户ID)

---

## 5. 新增方法

### 5.1 getTenantByDomain

**文件**: `apps/iam/internal/service/svcauth/auth.go`

```go
func (svc *authSvc) getTenantByDomain(ctx *gin.Context, orgID uint, domain string) (*model.TenantEntity, error) {
    return dao.NewTenantDao().GetByCond(ctx, &dao.TenantCond{
        OrgID:  orgID,
        Domain: domain,
        Status: model.TenantStatusEnabled,
    })
}
```

### 5.2 resolveDomainFromHost

**文件**: `apps/iam/internal/service/svcauth/auth.go`

```go
func resolveDomainFromHost(ctx *gin.Context) string {
    host := ctx.Request.Host
    if idx := strings.Index(host, ":"); idx != -1 {
        host = host[:idx]
    }
    return strings.TrimSpace(host)
}
```

---

## 6. 审核流程（新增）

### 6.1 接口列表

| 操作 | 接口 | 方法 | 说明 |
|------|------|------|------|
| 查看待审核用户 | GET /v1/iam/user/pendingList | - | 子租户管理员查看本租户待审核用户 |
| 审核用户 | POST /v1/iam/user/approve | approve/reject | 通过或拒绝用户注册 |

### 6.2 审核状态变化

```
pending → (approve) → enabled
pending → (reject) → disabled
```

---

## 7. 需修改的文件清单

| 文件 | 改动内容 |
|------|---------|
| `apps/iam/model/tenant.go` | TenantEntity 新增 Domain 字段 |
| `apps/iam/model/user.go` | UserStatus 新增 UserStatusPending 状态 |
| `apps/iam/internal/dto/dtoauth/request.go` | RegisterReq 移除 TenantName, TenantCode |
| `apps/iam/internal/dto/dtoauth/response.go` | RegisterResp 移除 TenantID |
| `apps/iam/internal/service/svcauth/auth.go` | 重写 Register 方法，移除创建租户逻辑 |
| `apps/iam/dao/tenant.go` | TenantCond 新增 Domain 查询条件 |
| `apps/iam/dao/user.go` | 新增 GetPendingUsers 方法 |
| `apps/iam/internal/controller/ctruser/user.go` | 新增 Approve 方法 |
| `apps/iam/internal/router/user.go` | 注册审核相关路由 |
| `scripts/sql/v1.x.x_*__table_init.sql` | 租户表新增 domain 字段 |
| `docs/iam_docs.go` | 更新 swagger 文档 |

---

## 8. 数据库变更

### 8.1 租户表新增 domain 字段

```sql
ALTER TABLE iam_tenant ADD COLUMN domain VARCHAR(255) DEFAULT '' COMMENT '租户域名(用于注册时子域名匹配)';
```

### 8.2 索引

```sql
ALTER TABLE iam_tenant ADD UNIQUE INDEX idx_org_domain (org_id, domain);
```

---

## 9. 兼容性考虑

1. **平滑过渡**: 旧接口仍可正常使用，只是行为改变
2. **数据库兼容**: 新增 nullable 字段，不影响现有数据
3. **配置兼容**: register.enabled、register.requireApproval 配置保留