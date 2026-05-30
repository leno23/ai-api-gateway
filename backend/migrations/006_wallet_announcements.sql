-- Wallet affiliate pending + system announcements

ALTER TABLE users
    ADD COLUMN IF NOT EXISTS affiliate_pending BIGINT NOT NULL DEFAULT 0;

CREATE TABLE IF NOT EXISTS announcements (
    id          BIGSERIAL PRIMARY KEY,
    title       VARCHAR(256) NOT NULL,
    content     TEXT NOT NULL,
    level       VARCHAR(32) NOT NULL DEFAULT 'info',
    placement   VARCHAR(64) NOT NULL DEFAULT 'home',
    status      SMALLINT NOT NULL DEFAULT 1,
    starts_at   TIMESTAMPTZ,
    ends_at     TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO announcements (title, content, level, placement, status)
SELECT
    '欢迎使用 AI API Gateway',
    '企业级多模型接入网关已上线。注册后即可创建 API 令牌，在模型广场查看价格，于控制台查看用量。',
    'info',
    'home',
    1
WHERE NOT EXISTS (SELECT 1 FROM announcements LIMIT 1);
