-- 租户表重命名为产品线表
RENAME TABLE iam_tenant TO iam_organization;

-- 租户配置表重命名为产品线配置表
RENAME TABLE iam_tenant_config TO iam_organization_config;

-- 公司表字段重命名：tenant_id → organization_id
ALTER TABLE iam_company 
    CHANGE COLUMN tenant_id organization_id BIGINT NOT NULL COMMENT '所属产品线ID';

-- 产品线配置表字段重命名：tenant_id → organization_id
ALTER TABLE iam_organization_config 
    CHANGE COLUMN tenant_id organization_id BIGINT NOT NULL COMMENT '产品线ID';

-- 更新索引名称
ALTER TABLE iam_organization_config 
    DROP INDEX uk_tenant_config,
    ADD UNIQUE INDEX uk_organization_config (organization_id, config_key);

ALTER TABLE iam_organization_config 
    DROP INDEX idx_tenant_id,
    ADD INDEX idx_organization_id (organization_id);

ALTER TABLE iam_company 
    DROP INDEX uk_tenant_code,
    ADD UNIQUE INDEX uk_organization_code (organization_id, company_code);

ALTER TABLE iam_company 
    DROP INDEX idx_tenant_id,
    ADD INDEX idx_organization_id (organization_id);

ALTER TABLE iam_company 
    DROP INDEX idx_tenant_status,
    ADD INDEX idx_organization_status (organization_id, status);