# Org 模块重构设计文档

## 1. 背景与目标

### 1.1 问题描述

1. **Org Create 接口存在逻辑漏洞**：
   - Domain 唯一性检查缺失
   - OrgCode 唯一性检查缺失
   - Admin 创建时 mobile/email 缺失会失败
   - Configs 校验跳过时无感知
   - Delete 删除关联数据缺少 UserEntity

2. **架构设计问题**：
   - Org.domain 和 Tenant.domain 关系不清晰
   - 注册策略获取流程存在断环

### 1.2 设计决策

| 决策项 | 选择 | 说明 |
|--------|------|------|
| Org 是否需要 domain | **否** | Org 是租户的分组，从 Tenant 反查 Org 更清晰 |
| Code 命名 | `code` + `display_code` | `code` 系统内部生成，`display_code` 手动输入对外展示 |
| 注册策略 | Org 共享配置 | org_config = 所有下属 tenant 共享的配置 |

## 2. 数据库变更

### 2.1 iam_organization 表

**删除字段**：
- `domain VARCHAR(255) COMMENT '组织域名'`

**修改字段**：
- `org_code VARCHAR(32) UNIQUE NOT NULL` → `display_code VARCHAR(32) UNIQUE NOT NULL COMMENT '组织编码(对外展示)'`

**新增字段**：
- `code VARCHAR(32) NOT NULL COMMENT '系统内部唯一标识'`
- 添加 `UNIQUE KEY uk_code (code)`

### 2.2 初始化数据变更

**v1.0.0_2__insert_data.sql**：
- 平台组织的 `org_code` → `code` 和 `display_code`
- `SET @platform_org_code = 'PLATFORM';` → 分别设置 `code` 和 `display_code`

## 3. Go 模型变更

### 3.1 model/organization.go

```go
// 删除
Domain string `gorm:"column:domain;..."`

// 修改
OrgCode string → DisplayCode string

// 新增
Code string `gorm:"column:code;type:varchar(32);not null;default '';comment: 系统内部唯一标识"`
```

### 3.2 dao/organization.go

```go
type OrganizationCond struct {
    // 新增
    Code string
    DisplayCode string
    // ...
}
```

## 4. Service 层修复

### 4.1 Create 接口修复

| 问题 | 修复方案 |
|------|---------|
| Domain 唯一性检查 | **删除**（无 domain） |
| Code 唯一性检查 | 新增，dao 查询校验 |
| Code 系统生成 | 自动生成（如 snowflake 或 UUID） |
| DisplayCode 唯一性检查 | 保留，调整字段名 |
| Admin mobile/email 校验 | 必填校验或调整生成逻辑 |
| Configs 校验失败 | 返回错误而不是 silent continue |
| 平台租户/部门预校验 | 事务外预先检查 |
| AppIDs 事务内校验 | 移到事务开始后 |
| 字段长度校验 | 添加 |

### 4.2 Delete 接口修复

| 问题 | 修复方案 |
|------|---------|
| 缺少 UserEntity 删除 | 补充删除 |
| 事务完整性 | 确认所有操作在事务内 |

## 5. 注册流程修复

### 5.1 修复 domain_strategy.go

**原逻辑（有 bug）**：
```go
// 用同一域名既查 org 又查 tenant —— 架构不一致会断环
orgEntity, err := getCurrentOrg(ctx)  // 用域名查 org
tenant, err := dao.NewTenantDao().GetByCond(ctx, &dao.TenantCond{
    OrgID:  orgEntity.ID,
    Domain: domain,  // 又用同一域名查 tenant
})
```

**新逻辑**：
```go
// 1. 用域名或参数获取 Tenant
tenant, err := dao.NewTenantDao().GetByCond(ctx, &dao.TenantCond{
    Domain: domain,  // 或用 tenant_code/tenant_id
})

// 2. 从 Tenant 反查 Org（获取注册策略）
orgEntity, err := dao.NewOrganizationDao().GetByID(ctx, tenant.OrgID)

// 3. 使用 orgEntity.ID 获取 org_config
```

### 5.2 其他策略文件

检查并修复：
- `sso_strategy.go`
- `open_strategy.go`
- `invite_strategy.go`

### 5.3 auth.go:getCurrentOrg 调整

不再使用 domain 查 org，改为：
- 从请求参数/上下文获取 org_id
- 或从 Tenant 反查 Org

## 6. 实施步骤

### Step 1: 数据库变更
- [ ] 修改 v1.0.0_1__table_init.sql
- [ ] 修改 v1.0.0_2__insert_data.sql

### Step 2: Go 模型变更
- [ ] 修改 model/organization.go
- [ ] 修改 dao/organization.go
- [ ] 修改 DTO request/response

### Step 3: Org Service 修复
- [ ] Create 接口修复
- [ ] Update 接口修复
- [ ] Delete 接口修复

### Step 4: 注册流程修复
- [ ] domain_strategy.go
- [ ] sso_strategy.go
- [ ] open_strategy.go
- [ ] invite_strategy.go
- [ ] auth.go:getCurrentOrg

### Step 5: 全面检查
- [ ] 全局搜索 `OrgCode` → `DisplayCode`
- [ ] 全局搜索 `org_code` → `display_code`
- [ ] 检查 Swagger 注解
- [ ] 运行测试

## 7. Code 生成策略

`code` 字段系统内部生成，推荐方案：

| 方案 | 优点 | 缺点 |
|------|------|------|
| Snowflake | 有序、可读 | 依赖外部组件 |
| UUID | 简单、分布式友好 | 无序、较长 |
| 基于时间戳+随机 | 简单 | 可能冲突 |

**推荐**：使用 snowflake 算法生成，保持有序且可读。
