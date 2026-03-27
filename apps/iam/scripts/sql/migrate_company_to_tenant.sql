-- 将公司表重命名为租户表
RENAME TABLE iam_company TO iam_tenant;

-- 租户表字段重命名：company_code -> tenant_code, company_name -> tenant_name
ALTER TABLE iam_tenant 
    CHANGE COLUMN company_code tenant_code VARCHAR(32) NOT NULL COMMENT '租户编码',
    CHANGE COLUMN company_name tenant_name VARCHAR(128) NOT NULL COMMENT '租户名称';

-- 更新租户表索引名称
ALTER TABLE iam_tenant 
    DROP INDEX uk_organization_code,
    ADD UNIQUE INDEX uk_organization_tenant_code (organization_id, tenant_code);

ALTER TABLE iam_tenant 
    DROP INDEX idx_organization_status,
    ADD INDEX idx_organization_tenant_status (organization_id, status);

-- 用户表字段重命名：company_id -> tenant_id
ALTER TABLE iam_user 
    CHANGE COLUMN company_id tenant_id BIGINT NOT NULL DEFAULT 0 COMMENT '所属租户ID, 0表示平台管理员账号';

-- 更新用户表索引名称
ALTER TABLE iam_user 
    DROP INDEX uk_company_username,
    ADD UNIQUE INDEX uk_tenant_username (tenant_id, username);

ALTER TABLE iam_user 
    DROP INDEX uk_company_employee_no,
    ADD UNIQUE INDEX uk_tenant_employee_no (tenant_id, employee_no);

ALTER TABLE iam_user 
    DROP INDEX idx_company_dept,
    ADD INDEX idx_tenant_dept (tenant_id, dept_id);

ALTER TABLE iam_user 
    DROP INDEX idx_company_status,
    ADD INDEX idx_tenant_status (tenant_id, status);

ALTER TABLE iam_user 
    DROP INDEX idx_company_user_type,
    ADD INDEX idx_tenant_user_type (tenant_id, user_type);

-- 部门表字段重命名：company_id -> tenant_id
ALTER TABLE iam_department 
    CHANGE COLUMN company_id tenant_id BIGINT NOT NULL COMMENT '所属租户ID';

-- 更新部门表索引名称
ALTER TABLE iam_department 
    DROP INDEX uk_company_code,
    ADD UNIQUE INDEX uk_tenant_code (tenant_id, dept_code);

ALTER TABLE iam_department 
    DROP INDEX idx_company_parent,
    ADD INDEX idx_tenant_parent (tenant_id, parent_id);

ALTER TABLE iam_department 
    DROP INDEX idx_company_path,
    ADD INDEX idx_tenant_path (tenant_id, dept_path(100));

ALTER TABLE iam_department 
    DROP INDEX idx_company_status,
    ADD INDEX idx_tenant_status (tenant_id, status);

-- 角色表字段重命名：company_id -> tenant_id
ALTER TABLE iam_role 
    CHANGE COLUMN company_id tenant_id BIGINT NOT NULL COMMENT '所属租户ID';

-- 更新角色表索引名称
ALTER TABLE iam_role 
    DROP INDEX uk_company_code,
    ADD UNIQUE INDEX uk_tenant_code (tenant_id, role_code);

ALTER TABLE iam_role 
    DROP INDEX idx_company_status,
    ADD INDEX idx_tenant_status (tenant_id, status);

-- 菜单表字段重命名：company_id -> tenant_id
ALTER TABLE iam_menu 
    CHANGE COLUMN company_id tenant_id BIGINT NOT NULL COMMENT '所属租户ID';

-- 更新菜单表索引名称
ALTER TABLE iam_menu 
    DROP INDEX uk_company_code,
    ADD UNIQUE INDEX uk_tenant_code (tenant_id, menu_code);

ALTER TABLE iam_menu 
    DROP INDEX idx_company_parent,
    ADD INDEX idx_tenant_parent (tenant_id, parent_id);

ALTER TABLE iam_menu 
    DROP INDEX idx_company_status,
    ADD INDEX idx_tenant_status (tenant_id, status);

-- 用户角色关联表字段重命名：company_id -> tenant_id
ALTER TABLE iam_user_role 
    CHANGE COLUMN company_id tenant_id BIGINT NOT NULL COMMENT '租户ID(冗余)';

-- 更新用户角色关联表索引名称
ALTER TABLE iam_user_role 
    DROP INDEX uk_user_role,
    ADD UNIQUE INDEX uk_tenant_user_role (tenant_id, user_id, role_id);

ALTER TABLE iam_user_role 
    DROP INDEX idx_company_user,
    ADD INDEX idx_tenant_user (tenant_id, user_id);

ALTER TABLE iam_user_role 
    DROP INDEX idx_company_role,
    ADD INDEX idx_tenant_role (tenant_id, role_id);

-- 用户部门关联表字段重命名：company_id -> tenant_id
ALTER TABLE iam_user_department 
    CHANGE COLUMN company_id tenant_id BIGINT NOT NULL COMMENT '租户ID(冗余)';

-- 更新用户部门关联表索引名称
ALTER TABLE iam_user_department 
    DROP INDEX idx_company_dept,
    ADD INDEX idx_tenant_dept (tenant_id, dept_id);

-- 角色菜单关联表字段重命名：company_id -> tenant_id
ALTER TABLE iam_role_menu 
    CHANGE COLUMN company_id tenant_id BIGINT NOT NULL COMMENT '租户ID(冗余)';

-- 更新角色菜单关联表索引名称
ALTER TABLE iam_role_menu 
    DROP INDEX uk_role_menu,
    ADD UNIQUE INDEX uk_tenant_role_menu (tenant_id, role_id, menu_id);

ALTER TABLE iam_role_menu 
    DROP INDEX idx_company_role,
    ADD INDEX idx_tenant_role (tenant_id, role_id);

ALTER TABLE iam_role_menu 
    DROP INDEX idx_company_menu,
    ADD INDEX idx_tenant_menu (tenant_id, menu_id);

-- 操作日志表字段重命名：company_id -> tenant_id
ALTER TABLE iam_operation_log 
    CHANGE COLUMN company_id tenant_id BIGINT NOT NULL COMMENT '租户ID';

-- 更新操作日志表索引名称
ALTER TABLE iam_operation_log 
    DROP INDEX idx_company_user,
    ADD INDEX idx_tenant_user (tenant_id, user_id);

ALTER TABLE iam_operation_log 
    DROP INDEX idx_company_created,
    ADD INDEX idx_tenant_created (tenant_id, created_at);

ALTER TABLE iam_operation_log 
    DROP INDEX idx_company_module,
    ADD INDEX idx_tenant_module (tenant_id, module);

-- 登录日志表字段重命名：company_id -> tenant_id
ALTER TABLE iam_login_log 
    CHANGE COLUMN company_id tenant_id BIGINT NOT NULL COMMENT '租户ID';

-- 更新登录日志表索引名称
ALTER TABLE iam_login_log 
    DROP INDEX idx_company_user,
    ADD INDEX idx_tenant_user (tenant_id, user_id);

ALTER TABLE iam_login_log 
    DROP INDEX idx_company_created,
    ADD INDEX idx_tenant_created (tenant_id, created_at);