-- Enable pgvector extension
CREATE EXTENSION IF NOT EXISTS vector;

-- ============================================================
-- rg_tenant
-- ============================================================
CREATE TABLE IF NOT EXISTS rg_tenant (
    id              BIGSERIAL PRIMARY KEY,
    name            VARCHAR(255) NOT NULL,
    status          VARCHAR(50) NOT NULL DEFAULT 'active',
    storage_config  JSONB NOT NULL DEFAULT '{}',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

-- ============================================================
-- rg_user
-- ============================================================
CREATE TABLE IF NOT EXISTS rg_user (
    id              BIGSERIAL PRIMARY KEY,
    tenant_id       BIGINT NOT NULL REFERENCES rg_tenant(id),
    username        VARCHAR(255) NOT NULL,
    email           VARCHAR(255),
    password_hash   VARCHAR(255),
    role            VARCHAR(50) NOT NULL DEFAULT 'viewer',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_user_tenant ON rg_user(tenant_id);

-- ============================================================
-- rg_knowledge_base
-- ============================================================
CREATE TABLE IF NOT EXISTS rg_knowledge_base (
    id              BIGSERIAL PRIMARY KEY,
    tenant_id       BIGINT NOT NULL,
    name            VARCHAR(255) NOT NULL,
    description     TEXT,
    kb_type         VARCHAR(50) NOT NULL DEFAULT 'normal',
    parser_engine   VARCHAR(100),
    embedding_config JSONB NOT NULL DEFAULT '{}',
    index_strategy  JSONB NOT NULL DEFAULT '{}',
    creator_id      BIGINT NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_kb_tenant ON rg_knowledge_base(tenant_id);

-- ============================================================
-- rg_knowledge
-- ============================================================
CREATE TABLE IF NOT EXISTS rg_knowledge (
    id              BIGSERIAL PRIMARY KEY,
    kb_id           BIGINT NOT NULL,
    tenant_id       BIGINT NOT NULL,
    type            VARCHAR(50) NOT NULL,
    title           VARCHAR(500),
    content         TEXT,
    file_url        TEXT,
    source_url      TEXT,
    parse_status    VARCHAR(50) NOT NULL DEFAULT 'pending',
    file_size       BIGINT NOT NULL DEFAULT 0,
    creator_id      BIGINT NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_knowledge_kb_tenant ON rg_knowledge(kb_id, tenant_id);
CREATE INDEX IF NOT EXISTS idx_knowledge_parse_status ON rg_knowledge(parse_status);

-- ============================================================
-- rg_chunk (with pgvector support)
-- ============================================================
CREATE TABLE IF NOT EXISTS rg_chunk (
    id              BIGSERIAL PRIMARY KEY,
    knowledge_id    BIGINT NOT NULL,
    kb_id           BIGINT NOT NULL,
    tenant_id       BIGINT NOT NULL,
    content         TEXT NOT NULL,
    seq_id          INT NOT NULL DEFAULT 0,
    tokens          INT NOT NULL DEFAULT 0,
    vector          vector(1536),
    meta_info       JSONB NOT NULL DEFAULT '{}',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_chunk_kb_tenant ON rg_chunk(kb_id, tenant_id);
CREATE INDEX IF NOT EXISTS idx_chunk_knowledge ON rg_chunk(knowledge_id);

-- HNSW index for pgvector similarity search
-- Uses vector_cosine_ops for cosine distance (<-> operator)
-- m=16 (max connections per node), ef_construction=64 (build speed/quality trade-off)
CREATE INDEX IF NOT EXISTS idx_chunk_vector ON rg_chunk
USING hnsw (vector vector_cosine_ops)
WITH (m = 16, ef_construction = 64);

-- ============================================================
-- rg_faq
-- ============================================================
CREATE TABLE IF NOT EXISTS rg_faq (
    id                BIGSERIAL PRIMARY KEY,
    kb_id             BIGINT NOT NULL,
    tenant_id         BIGINT NOT NULL,
    question          TEXT NOT NULL,
    answer            TEXT NOT NULL,
    similar_questions JSONB NOT NULL DEFAULT '[]',
    tags              JSONB NOT NULL DEFAULT '[]',
    status            VARCHAR(50) NOT NULL DEFAULT 'active',
    creator_id        BIGINT NOT NULL DEFAULT 0,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at        TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_faq_kb_tenant ON rg_faq(kb_id, tenant_id);

-- ============================================================
-- rg_session
-- ============================================================
CREATE TABLE IF NOT EXISTS rg_session (
    id              BIGSERIAL PRIMARY KEY,
    tenant_id       BIGINT NOT NULL,
    user_id         BIGINT NOT NULL,
    kb_id           BIGINT NOT NULL DEFAULT 0,
    title           VARCHAR(500),
    description     TEXT,
    is_pinned       BOOLEAN NOT NULL DEFAULT false,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_session_tenant_user ON rg_session(tenant_id, user_id);

-- ============================================================
-- rg_message
-- ============================================================
CREATE TABLE IF NOT EXISTS rg_message (
    id              BIGSERIAL PRIMARY KEY,
    session_id      BIGINT NOT NULL,
    tenant_id       BIGINT NOT NULL,
    role            VARCHAR(50) NOT NULL,
    content         TEXT NOT NULL,
    metadata        JSONB NOT NULL DEFAULT '{}',
    token_count     INT NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_message_session ON rg_message(session_id);
CREATE INDEX IF NOT EXISTS idx_message_tenant ON rg_message(tenant_id);

-- ============================================================
-- rg_model
-- ============================================================
CREATE TABLE IF NOT EXISTS rg_model (
    id              BIGSERIAL PRIMARY KEY,
    tenant_id       BIGINT NOT NULL,
    name            VARCHAR(255) NOT NULL,
    model_type      VARCHAR(50) NOT NULL,
    provider        VARCHAR(100) NOT NULL,
    model_name      VARCHAR(255) NOT NULL,
    config          JSONB NOT NULL DEFAULT '{}',
    status          VARCHAR(50) NOT NULL DEFAULT 'active',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_model_tenant ON rg_model(tenant_id);

-- ============================================================
-- rg_vector_store
-- ============================================================
CREATE TABLE IF NOT EXISTS rg_vector_store (
    id              BIGSERIAL PRIMARY KEY,
    tenant_id       BIGINT NOT NULL,
    name            VARCHAR(255) NOT NULL,
    engine_type     VARCHAR(100) NOT NULL,
    config          JSONB NOT NULL DEFAULT '{}',
    status          VARCHAR(50) NOT NULL DEFAULT 'active',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_vector_store_tenant ON rg_vector_store(tenant_id);

-- ============================================================
-- rg_tag
-- ============================================================
CREATE TABLE IF NOT EXISTS rg_tag (
    id              BIGSERIAL PRIMARY KEY,
    kb_id           BIGINT NOT NULL,
    tenant_id       BIGINT NOT NULL,
    name            VARCHAR(255) NOT NULL,
    color           VARCHAR(50),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_tag_kb_tenant ON rg_tag(kb_id, tenant_id);

-- ============================================================
-- rg_migration (migration tracking table)
-- ============================================================
CREATE TABLE IF NOT EXISTS rg_migration (
    id              BIGSERIAL PRIMARY KEY,
    filename        VARCHAR(255) NOT NULL UNIQUE,
    applied_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
