# WeKnora 迁移 - 知识库模块设计

## 目标

将 WeKnora 的知识库模块迁移到 `apps/ragflow`，代码风格和目录结构参照 `apps/demo`。

## 迁移范围

| 模块 | 功能 |
|------|------|
| model | KnowledgeBase、Document、Chunk 数据模型 |
| dao | 对应的数据访问层 |
| service | 知识库 CRUD、文档上传、解析、分块、向量化 |
| controller | API 控制器 |
| dto | 请求/响应 DTO |
| router | 路由注册 |

## 目录结构

```
apps/ragflow/
├── model/
│   ├── knowledgebase.go    # 知识库实体
│   ├── document.go          # 文档实体
│   └── chunk.go             # 分块实体
├── dao/
│   ├── knowledgebase.go
│   ├── document.go
│   └── chunk.go
├── object/
│   └── objknowledge/
│       └── knowledge.go
├── internal/
│   ├── controller/ctrknowledge/
│   │   └── knowledgebase.go
│   ├── service/svcknowledge/
│   │   ├── knowledgebase.go
│   │   └── document.go
│   ├── dto/dtoknowledge/
│   │   ├── request.go
│   │   └── response.go
│   └── router/
│       └── knowledgebase.go
```

## 数据模型

### KnowledgeBaseEntity

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| name | string | 知识库名称 |
| description | string | 描述 |
| embedding_model | string | embedding 模型 |
| vector_store_type | string | 向量库类型（pgvector） |
| permission_type | string | 权限类型 |
| status | string | 状态 |
| chunk_method | string | 分块策略 |

### DocumentEntity

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| knowledgebase_id | uint | 所属知识库 |
| name | string | 文档名称 |
| type | string | 文档类型 |
| location | string | 存储位置 |
| size | int64 | 文件大小 |
| status | string | 状态 |
| chunk_status | string | 分块状态 |
| vector_status | string | 向量化状态 |

### ChunkEntity

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| document_id | uint | 所属文档 |
| content | text | 分块内容 |
| vector | pgvector | 向量 |
| metadata | json | 元数据 |

## API 设计

| 操作 | 路径 | 方法 |
|------|------|------|
| 创建知识库 | `/v1/ragflow/knowledgebase/create` | POST |
| 删除知识库 | `/v1/ragflow/knowledgebase/delete` | POST |
| 更新知识库 | `/v1/ragflow/knowledgebase/update` | POST |
| 知识库详情 | `/v1/ragflow/knowledgebase/detail` | GET |
| 知识库列表 | `/v1/ragflow/knowledgebase/list` | GET |
| 上传文档 | `/v1/ragflow/knowledgebase/{id}/document/upload` | POST |
| 文档列表 | `/v1/ragflow/knowledgebase/{id}/document/list` | GET |
| 删除文档 | `/v1/ragflow/document/delete` | POST |

## 配置

- 数据库配置：参照 demo 的配置方式
- 其他配置（Redis、存储、AI 模型等）：使用 golib 的 configkv

## 迁移顺序

1. model 层（实体定义）
2. dao 层（数据访问）
3. service 层（业务逻辑）
4. controller/dto 层（API）
5. router 层（路由注册）

## 待后续模块

- 对话模块（chat）
- Agent 模块（agent）
- IM 集成模块（im）
- AI 模型集成模块（ai）