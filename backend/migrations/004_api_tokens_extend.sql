-- Extend api_keys for portal token management

ALTER TABLE api_keys
    ADD COLUMN IF NOT EXISTS token_group_id BIGINT REFERENCES token_groups (id),
    ADD COLUMN IF NOT EXISTS quota_limit BIGINT,
    ADD COLUMN IF NOT EXISTS used_quota BIGINT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS ip_whitelist TEXT[] NOT NULL DEFAULT '{}';

UPDATE api_keys
SET token_group_id = (SELECT id FROM token_groups WHERE slug = 'default' LIMIT 1)
WHERE token_group_id IS NULL;
