# IAM 注册接口改进实施计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 改进 IAM 注册接口，实现通过子域名精确匹配已存在租户，支持多级租户结构和注册审核流程。

**Architecture:** 注册接口不再创建新租户，而是通过请求子域名精确匹配已存在的租户。用户注册后状态为 pending 时需要子租户管理员审核才能激活。

**Tech Stack:** Go, GORM, Gin

---

## 1. 文件变更概览

| 任务 | 涉及文件 |
|------|---------|
| Task 1 | `apps/iam/model/tenant.go`, `apps/iam/model/user.go` |
| Task 2 | `apps/iam/dao/tenant.go` |
| Task 3 | `apps/iam/dao/user.go` |
| Task 4 | `apps/iam/internal/dto/dtoauth/request.go`, `apps/iam/internal/dto/dtoauth/response.go` |
| Task 5 | `apps/iam/internal/service/svcauth/auth.go` |
| Task 6 | `apps/iam/internal/controller/ctruser/user.go` |
| Task 7 | `apps/iam/internal/router/user.go` |
| Task 8 | SQL 迁移脚本 |
| Task 9 | Swagger 文档更新 |

---

## 2. 任务详情

### Task 1: 数据模型变更

**Files:**
- Modify: `apps/iam/model/tenant.go:16-35`
- Modify: `apps/iam/model/user.go:9-15`

- [ ] **Step 1: TenantEntity 新增 Domain 字段**

修改 `apps/iam/model/tenant.go`，在 TenantEntity 结构体中添加 Domain 字段：

```go
type TenantEntity struct {
    gorm.Model
    // ... existing fields ...
    Domain string `gorm:"column:domain;type:varchar(255);;default '';comment: 租户域名(用于注册时子域名匹配)"`
    // ... rest of fields ...
}
```

- [ ] **Step 2: UserStatus 新增 pending 状态**

修改 `apps/iam/model/user.go`，在 UserStatus 常量中添加 pending：

```go
const (
    UserStatusEnabled  UserStatus = "enabled"
    UserStatusLocked   UserStatus = "locked"
    UserStatusDisabled UserStatus = "disabled"
    UserStatusPending  UserStatus = "pending"
)
```

- [ ] **Step 3: 提交变更**

```bash
git add apps/iam/model/tenant.go apps/iam/model/user.go
git commit -m "feat(iam): add Domain field to TenantEntity and pending status to UserStatus"
```

---

### Task 2: DAO 层 - TenantCond 新增 Domain 查询条件

**Files:**
- Modify: `apps/iam/dao/tenant.go:1-50`

- [ ] **Step 1: 在 TenantCond 结构体中添加 Domain 字段**

修改 `apps/iam/dao/tenant.go`，在 TenantCond 结构体中添加 Domain 字段：

```go
type TenantCond struct {
    // ... existing fields ...
    Domain string
}

// 在查询条件中添加 Domain 过滤
if c.Domain != "" {
    db.Where(tableName+".domain = ?", c.Domain)
}
```

- [ ] **Step 2: 提交变更**

```bash
git add apps/iam/dao/tenant.go
git commit -m "feat(iam): add Domain field to TenantCond for subdomain matching"
```

---

### Task 3: DAO 层 - 新增 GetPendingUsers 方法

**Files:**
- Modify: `apps/iam/dao/user.go`

- [ ] **Step 1: 在 UserCond 中添加 Status 查询支持**

修改 `apps/iam/dao/user.go`，确保 TenantID 和 Status 可以同时查询：

```go
type UserCond struct {
    // ... existing fields ...
    Status string
}
```

- [ ] **Step 2: 新增 GetPendingUsers 方法**

在 `apps/iam/dao/user.go` 中添加：

```go
func (dao *userDao) GetPendingUsers(ctx context.Context, tenantID uint) ([]*model.UserEntity, error) {
    var users []*model.UserEntity
    err := dao.GetListByCond(ctx, &UserCond{
        TenantID: tenantID,
        Status:   string(model.UserStatusPending),
    }, &users)
    return users, err
}
```

- [ ] **Step 3: 提交变更**

```bash
git add apps/iam/dao/user.go
git commit -m "feat(iam): add GetPendingUsers method for approval workflow"
```

---

### Task 4: DTO 层变更

**Files:**
- Modify: `apps/iam/internal/dto/dtoauth/request.go:20-28`
- Modify: `apps/iam/internal/dto/dtoauth/response.go:41-48`

- [ ] **Step 1: 修改 RegisterReq，移除 TenantName 和 TenantCode**

修改 `apps/iam/internal/dto/dtoauth/request.go`，RegisterReq 改为：

```go
type RegisterReq struct {
    Username string `json:"username" validate:"required" label:"用户名"`
    Password string `json:"password" validate:"required" label:"密码"`
    Mobile   string `json:"mobile" label:"手机号"`
    Email    string `json:"email" validate:"required" label:"邮箱"`
    RealName string `json:"realName" validate:"required" label:"真实姓名"`
}
```

- [ ] **Step 2: 修改 RegisterResp，移除 TenantID**

修改 `apps/iam/internal/dto/dtoauth/response.go`，RegisterResp 改为：

```go
type RegisterResp struct {
    UserID       uint   `json:"userId"`
    PersonID     uint   `json:"personId"`
    Status       string `json:"status"`
    PersonExists bool   `json:"personExists"`
    Message      string `json:"message"`
}
```

- [ ] **Step 3: 新增 ApproveReq 和 PendingListResp DTO**

在 `apps/iam/internal/dto/dtoauth/request.go` 添加：

```go
type ApproveReq struct {
    UserID   uint   `json:"userId" validate:"required" label:"用户ID"`
    Approved bool   `json:"approved" label:"是否通过"`  // true=通过, false=拒绝
}

type PendingListResp struct {
    Users    []PendingUserInfo `json:"users"`
    Total    int64             `json:"total"`
}

type PendingUserInfo struct {
    UserID    uint   `json:"userId"`
    Username  string `json:"username"`
    RealName  string `json:"realName"`
    Email     string `json:"email"`
    Mobile    string `json:"mobile"`
    CreatedAt string `json:"createdAt"`
}
```

- [ ] **Step 4: 提交变更**

```bash
git add apps/iam/internal/dto/dtoauth/request.go apps/iam/internal/dto/dtoauth/response.go
git commit -m "feat(iam): update RegisterReq/RegisterResp and add approval DTOs"
```

---

### Task 5: Service 层 - 重写 Register 方法

**Files:**
- Modify: `apps/iam/internal/service/svcauth/auth.go:535-679`

- [ ] **Step 1: 新增 resolveDomainFromHost 方法**

在 `apps/iam/internal/service/svcauth/auth.go` 中添加：

```go
func resolveDomainFromHost(ctx *gin.Context) string {
    host := ctx.Request.Host
    if idx := strings.Index(host, ":"); idx != -1 {
        host = host[:idx]
    }
    return strings.TrimSpace(host)
}
```

- [ ] **Step 2: 新增 getTenantByDomain 方法**

在 `apps/iam/internal/service/svcauth/auth.go` 中添加：

```go
func (svc *authSvc) getTenantByDomain(ctx *gin.Context, orgID uint, domain string) (*model.TenantEntity, error) {
    return dao.NewTenantDao().GetByCond(ctx, &dao.TenantCond{
        OrgID:  orgID,
        Domain: domain,
        Status: model.TenantStatusEnabled,
    })
}
```

- [ ] **Step 3: 重写 Register 方法**

修改 `apps/iam/internal/service/svcauth/auth.go` 中的 Register 方法，核心逻辑：

```go
func (svc *authSvc) Register(ctx *gin.Context, req *dtoauth.RegisterReq) (*dtoauth.RegisterResp, error) {
    // 1. 获取当前组织
    orgEntity, err := svc.getCurrentOrg(ctx)
    if err != nil {
        return nil, err
    }

    // 2. 检查注册是否开放
    registerEnabled, err := svc.getOrgConfigBool(ctx, orgEntity.ID, model.OrgConfigKeyRegisterEnabled)
    if err != nil || !registerEnabled {
        return nil, code.GetError(code.AuthRegisterDisabled)
    }

    // 3. 从请求 Host 解析子域名
    domain := resolveDomainFromHost(ctx)

    // 4. 通过子域名精确匹配租户
    tenantEntity, err := svc.getTenantByDomain(ctx, orgEntity.ID, domain)
    if err != nil || tenantEntity == nil {
        return nil, code.GetError(code.TenantNotFoundError)  // 需新增错误码
    }

    // 5. 检查租户状态
    if tenantEntity.Status != model.TenantStatusEnabled {
        return nil, code.GetError(code.TenantDisabledError)  // 需新增错误码
    }

    // 6. 验证身份格式
    identityType, _ := svc.getOrgConfigString(ctx, orgEntity.ID, model.OrgConfigKeyRegisterIdentityType)
    if identityType == "" {
        identityType = string(model.RegisterIdentityTypeEmail)
    }
    if err := svc.validateRegisterIdentity(ctx, req, model.RegisterIdentityType(identityType)); err != nil {
        return nil, err
    }

    // 7. 确定用户状态
    requireApproval, _ := svc.getOrgConfigBool(ctx, orgEntity.ID, model.OrgConfigKeyRegisterRequireApproval)
    userStatus := model.UserStatusEnabled
    message := "注册成功"
    if requireApproval {
        userStatus = model.UserStatusPending
        message = "注册成功，等待管理员审核"
    }

    // 8. 事务中创建: Person → User → UserDepartment
    // (不再创建租户，直接关联到 subdomain 匹配的租户)
    passwordHash, err := gcrypto.GeneratePasswordHash(req.Password)
    if err != nil {
        glog.Errorf(ctx, "[svcauth.Register] GeneratePasswordHash fail, err:%v", err)
        return nil, code.GetError(code.AuthRegisterError)
    }

    email := strings.TrimSpace(req.Email)
    personEntity, _ := dao.NewPersonDao().GetByCond(ctx, &dao.PersonCond{Email: email})
    personExists := personEntity != nil && personEntity.ID > 0

    var userID, personID uint
    txErr := dbclient.IamDB(ctx).Transaction(func(tx *gorm.DB) error {
        var personID uint
        if personExists {
            personID = personEntity.ID
        } else {
            newPerson := &model.PersonEntity{
                Mobile:       strings.TrimSpace(req.Mobile),
                Email:        email,
                RealName:     req.RealName,
                PasswordHash: passwordHash,
                CreatedBy:    0,
                UpdatedBy:    0,
            }
            if err := dao.NewPersonDao().WithTx(tx).Insert(ctx, newPerson); err != nil {
                return err
            }
            personID = newPerson.ID
        }

        // 创建用户，关联到租户
        userEntity := &model.UserEntity{
            TenantID:  tenantEntity.ID,
            PersonID:  personID,
            Username:  req.Username,
            UserType:  model.UserTypeNormal,  // 普通用户，非租户管理员
            Status:    userStatus,
            CreatedBy: 0,
            UpdatedBy: 0,
        }
        if err := dao.NewUserDao().WithTx(tx).Insert(ctx, userEntity); err != nil {
            return err
        }
        userID = userEntity.ID

        return nil
    })
    if txErr != nil {
        glog.Errorf(ctx, "[svcauth.Register] Transaction fail, err:%v", txErr)
        return nil, code.GetError(code.AuthRegisterError)
    }

    return &dtoauth.RegisterResp{
        UserID:       userID,
        PersonID:     personID,
        Status:       string(userStatus),
        PersonExists: personExists,
        Message:      message,
    }, nil
}
```

- [ ] **Step 4: 提交变更**

```bash
git add apps/iam/internal/service/svcauth/auth.go
git commit -m "feat(iam): rewrite Register method for subdomain-based tenant matching"
```

---

### Task 6: Controller 层 - 新增审核方法

**Files:**
- Modify: `apps/iam/internal/controller/ctruser/user.go`

- [ ] **Step 1: 新增 PendingList 方法**

在 `apps/iam/internal/controller/ctruser/user.go` 中添加：

```go
func (ctr *userCtr) PendingList(ctx *gin.Context) {
    var req dtoauth.PageListReq
    if err := ctx.ShouldBindQuery(&req); err != nil {
        gincontext.Fail(ctx, err)
        return
    }

    // 从 context 获取当前用户和租户信息
    // 调用 svc 获取 pending 用户列表
    // 返回结果
}
```

- [ ] **Step 2: 新增 Approve 方法**

在 `apps/iam/internal/controller/ctruser/user.go` 中添加：

```go
func (ctr *userCtr) Approve(ctx *gin.Context) {
    var req dtoauth.ApproveReq
    if err := ctx.ShouldBindJSON(&req); err != nil {
        gincontext.Fail(ctx, err)
        return
    }

    res, err := ctr.userSvc.Approve(ctx, &req)
    if err != nil {
        gincontext.Fail(ctx, err)
        return
    }
    gincontext.Success(ctx, res)
}
```

- [ ] **Step 3: 提交变更**

```bash
git add apps/iam/internal/controller/ctruser/user.go
git commit -m "feat(iam): add PendingList and Approve controller methods"
```

---

### Task 7: Router 层 - 注册审核路由

**Files:**
- Modify: `apps/iam/internal/router/user.go`

- [ ] **Step 1: 添加 pendingList 和 approve 路由**

修改 `apps/iam/internal/router/user.go`：

```go
func userRouter(groups *ginserver.RouterGroups) {
    v1Group := groups.MustGetGroup(constant.ApiVersionV1)
    iamGroup := v1Group.Group("/iam")

    // 现有路由...
    // userRouter(iamGroup)

    // 新增审核相关路由
    iamGroup.GET("/user/pendingList", userCtr.PendingList)
    iamGroup.POST("/user/approve", userCtr.Approve)
}
```

- [ ] **Step 2: 提交变更**

```bash
git add apps/iam/internal/router/user.go
git commit -m "feat(iam): add pendingList and approve routes"
```

---

### Task 8: 数据库变更

**Files:**
- Create: `apps/iam/scripts/sql/v1.1.0_1__add_tenant_domain.sql`

- [ ] **Step 1: 创建 SQL 迁移脚本**

创建 `apps/iam/scripts/sql/v1.1.0_1__add_tenant_domain.sql`：

```sql
-- 租户表新增 domain 字段
ALTER TABLE iam_tenant ADD COLUMN domain VARCHAR(255) DEFAULT '' COMMENT '租户域名(用于注册时子域名匹配)';

-- 为 org_id 和 domain 组合添加唯一索引
ALTER TABLE iam_tenant ADD UNIQUE INDEX idx_org_domain (org_id, domain);
```

- [ ] **Step 2: 提交变更**

```bash
git add apps/iam/scripts/sql/v1.1.0_1__add_tenant_domain.sql
git commit -m "feat(iam): add tenant domain column migration script"
```

---

### Task 9: 错误码新增

**Files:**
- Modify: `pkg/code/code.go`

- [ ] **Step 1: 添加 TenantNotFoundError 和 TenantDisabledError**

在错误码文件中添加：

```go
const (
    TenantNotFoundError ErrorCode = 100301
    TenantDisabledError ErrorCode = 100302
)
```

- [ ] **Step 2: 提交变更**

```bash
git add pkg/code/code.go
git commit -m "feat(iam): add tenant not found and disabled error codes"
```

---

## 3. 自检清单

- [ ] 所有 RegisterReq.TenantName 和 RegisterReq.TenantCode 引用已移除
- [ ] resolveDomainFromHost 能正确从 Host 中提取子域名
- [ ] getTenantByDomain 能精确匹配租户
- [ ] 注册用户正确关联到匹配的租户（而非创建新租户）
- [ ] requireApproval=true 时用户状态为 pending
- [ ] 审核接口能正确修改用户状态
- [ ] SQL 脚本中的字段与 model 一致

---

**Plan complete and saved to `docs/superpowers/plans/2026-04-27-iam-register-enhancement-implementation.md`.**