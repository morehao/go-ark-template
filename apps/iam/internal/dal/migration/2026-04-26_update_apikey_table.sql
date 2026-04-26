-- 迁移 API Key 表结构
-- 1. 删除公钥和加密私钥字段
-- 2. 新增加密 API Key 字段

ALTER TABLE `iam_api_key`
    DROP COLUMN IF EXISTS `public_key`,
    DROP COLUMN IF EXISTS `encrypted_private_key`,
    ADD COLUMN IF NOT EXISTS `api_key` TEXT NOT NULL COMMENT 'AES加密的API Key' AFTER `key_prefix`;
