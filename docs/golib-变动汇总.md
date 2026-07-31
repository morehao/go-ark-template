# golib 变动汇总

本文件记录 `morehao/golib` 重构后的破坏性变更，供其他依赖 golib 的项目（goark 及其后续应用）升级参考。

## 一、包路径迁移

| 旧路径 | 新路径 | 说明 |
|--------|--------|------|
| `golib/biz/genericdao` | `golib/dbaccess/gormdao` | DAO 层，含 `BaseCond` / `Dao` |
| `golib/biz/gconstant` | `golib/gconstant` | 顶层，含错误码 Map 与常量 |
| — | `golib/glog/driver/slog`、`golib/glog/driver/zap` | glog 注册驱动子包 |
| — | `golib/dbaccess/dbgorm/driver/mysql`、`/postgres`、`/sqlite` | dbgorm 注册驱动子包 |

> 注：`golib/biz/gserver`、`golib/biz/gmiddleware`、`golib/biz/gcontext`、`golib/biz/gobject`、`golib/biz/testkit` 路径未变。

## 二、API 变更

### 1. genericdao → gormdao（API 改名）

| 旧 | 新 |
|----|----|
| `genericdao.GenericDao[T, L]` | `gormdao.Dao[T, L]` |
| `genericdao.NewGenericDao[T,L](tableName, daoName, getDB)` | `gormdao.NewDao[T,L](tableName, daoName, getDB DBGetter, opts...)` |
| `genericdao.GenericDao:`（内嵌字段名） | `gormdao.Dao:` |
| `genericdao.BaseCond` | `gormdao.BaseCond` |

- `DBGetter` 类型：`func(context.Context) *gorm.DB`，与 `dbclient.DemoDB` 等签名一致。
- 新增可选参 `gormdao.WithoutSoftDelete()`。
- **重要**：`genericdao.DBErrorMsgMap` 已移到顶层 **`golib/gconstant`**，代码中需改用 `gconstant.DBErrorMsgMap`。

### 2. ginserver 路由分组

| 旧 | 新 |
|----|----|
| `ginserver.Version{Name: ...}` | `ginserver.VersionGroup{Version: ...}` |

`NewRouterGroups(engine, appName, versions ...VersionGroup)` 参数类型已变更，字段 `Name` 改为 `Version`。

### 3. 常量位置

> `golib/gconstant` 中**没有** `ApiVersionV1` 等版本常量，它们位于 `golib/biz/gserver/ginserver`（如 `ginserver.ApiVersionV1 = "v1"`）。

## 三、初始化注册机制（重点）

### 1. glog 注册机制

glog 改为通过 driver 子包的 `init()` 注册 LoggerType 实现：

```go
// 顶层 glog
func RegisterLoggerType(t LoggerType, factory LoggerFactory)

// 内部依据 cfg.LoggerType 从 registeredFactories 匹配
// 空值默认为 LoggerTypeSlog；未注册对应 driver 会报错：
// "unknown LoggerType X, import glog/driver/slog or glog/driver/zap to register"
```

升级时**必须** blank import 对应驱动，否则 `glog.InitLogger` / `glog.GetDefaultLogger` 报错：

```go
import _ "github.com/morehao/golib/glog/driver/slog"  // 默认 slog
// 或
import _ "github.com/morehao/golib/glog/driver/zap"   // zap（配置需设 logger_type: zap）
```

- `InitLogger(cfg *LogConfig, ...)` 接口不变，`LogConfig` 新增 `LoggerType` 字段（yaml 可配 `logger_type`）。
- **陷阱**：`cmd` 入口 blank import 只覆盖应用进程；**测试初始化的路径（如 testsetup / testkit）也必须 blank import**，否则测试跑 `glog.InitLogger` 时 panic。

### 2. dbgorm 注册机制

dbgorm 改为通过 driver 子包的 `init()` 注册 `DialectorFactory`，`New` 依据 URL 前缀匹配已注册的 dialector：

```go
func Register(name string, factory DialectorFactory)  // DialectorFactory: Name / MatchURL / Dialector / ParseURL

// 未注册匹配驱动报错：
// "no registered dialector matches url, make sure to import the driver (e.g. _ \".../dbaccess/dbgorm/driver/mysql\")"
```

升级时**必须** blank import 对应数据库驱动（例如 MySQL）：

```go
import _ "github.com/morehao/golib/dbaccess/dbgorm/driver/mysql"
```

## 四、本仓库升级记录（goark）

goark 本次迁移涉及的文件分类：

- **路径迁移**：`biz/genericdao` → `dbaccess/gormdao`、`biz/gconstant` → `golib/gconstant`。
- **类型改名**：`dbgorm.GormConfig` → `dbgorm.Config`；`ginserver.Version{Name}` → `ginserver.VersionGroup{Version}`。
- **符号迁移**：`genericdao.DBErrorMsgMap` → `gconstant.DBErrorMsgMap`；`gconstant.ApiVersionV1` → `ginserver.ApiVersionV1`。
- **新增 blank import**：`glog/driver/slog`（cmd 与 testsetup 两处）、`dbgorm/driver/mysql`（dbclient）。

### 附带修复（非迁移引起）

- `pkg/code/ragforge.go` 与 `apikey.go` 错误码在 `101101-101108` 段重叠导致的 code 包 init panic，已把 ragforge 段整体平移至 `101300-101352`。
- `httpbingo` 测试依赖全局 `config.Conf` 未初始化，已在测试内 `config.LoadConfig` 修复。
- 清理了一批 pre-existing 的 golangci-lint 问题（unused 死代码、errcheck、staticcheck、grpc.Dial 弃用等）。
