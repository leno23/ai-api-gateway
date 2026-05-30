-- Token groups + model catalog fields for portal /pricing

CREATE TABLE IF NOT EXISTS token_groups (
    id           BIGSERIAL PRIMARY KEY,
    slug         VARCHAR(64) UNIQUE NOT NULL,
    name         VARCHAR(128) NOT NULL,
    multiplier   NUMERIC(6, 2) NOT NULL DEFAULT 1.00,
    status       SMALLINT NOT NULL DEFAULT 1,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE users
    ADD COLUMN IF NOT EXISTS token_group_id BIGINT REFERENCES token_groups (id);

CREATE INDEX IF NOT EXISTS idx_users_token_group_id ON users (token_group_id);

ALTER TABLE model_prices
    ADD COLUMN IF NOT EXISTS provider VARCHAR(32) NOT NULL DEFAULT 'openai',
    ADD COLUMN IF NOT EXISTS endpoint_type VARCHAR(32) NOT NULL DEFAULT 'openai',
    ADD COLUMN IF NOT EXISTS cache_read_price BIGINT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS cache_write_price BIGINT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS tags TEXT[] NOT NULL DEFAULT '{}',
    ADD COLUMN IF NOT EXISTS display_name VARCHAR(128);

INSERT INTO token_groups (slug, name, multiplier, status)
VALUES
    ('default', 'default', 1.00, 1),
    ('vip', 'vip', 1.00, 1),
    ('claude_code', 'claude_code', 1.50, 1),
    ('aws', 'AWS分组', 3.00, 1),
    ('guanzhuan_max', '官转max', 2.50, 1)
ON CONFLICT (slug) DO NOTHING;

UPDATE users
SET token_group_id = (SELECT id FROM token_groups WHERE slug = 'default' LIMIT 1)
WHERE token_group_id IS NULL;

-- Seed catalog rows when table is empty (dev/demo)
INSERT INTO model_prices (
    model, display_name, provider, endpoint_type,
    prompt_price, completion_price, cache_read_price, cache_write_price,
    unit_price, billing_type, tags
)
SELECT v.model, v.display_name, v.provider, v.endpoint_type,
       v.prompt_price, v.completion_price, v.cache_read_price, v.cache_write_price,
       v.unit_price, v.billing_type, v.tags
FROM (VALUES
    ('gpt-4o', 'GPT-4o', 'openai', 'openai', 9000::bigint, 27000::bigint, 4500::bigint, 0::bigint, 0::bigint, 1::smallint, ARRAY['chat']::text[]),
    ('gpt-4o-mini', 'GPT-4o mini', 'openai', 'openai', 750::bigint, 3000::bigint, 0::bigint, 0::bigint, 0::bigint, 1::smallint, ARRAY['chat']::text[]),
    ('claude-sonnet-4-20250514', 'Claude Sonnet 4', 'anthropic', 'anthropic', 15000::bigint, 75000::bigint, 7500::bigint, 0::bigint, 0::bigint, 1::smallint, ARRAY['chat']::text[]),
    ('claude-4-opus', 'Claude 4 Opus', 'anthropic', 'anthropic', 75000::bigint, 375000::bigint, 0::bigint, 0::bigint, 0::bigint, 1::smallint, ARRAY['chat']::text[]),
    ('deepseek-v3', 'DeepSeek V3', 'deepseek', 'openai', 1350::bigint, 5500::bigint, 0::bigint, 0::bigint, 0::bigint, 1::smallint, ARRAY['chat']::text[]),
    ('dall-e-3', 'DALL·E 3', 'openai', 'openai', 0::bigint, 0::bigint, 0::bigint, 0::bigint, 20000::bigint, 2::smallint, ARRAY['image']::text[])
) AS v(model, display_name, provider, endpoint_type, prompt_price, completion_price, cache_read_price, cache_write_price, unit_price, billing_type, tags)
WHERE NOT EXISTS (SELECT 1 FROM model_prices LIMIT 1);

-- Backfill provider/endpoint for existing rows
UPDATE model_prices SET provider = 'anthropic', endpoint_type = 'anthropic'
WHERE provider = 'openai' AND (model ILIKE 'claude%' OR model ILIKE '%anthropic%');

UPDATE model_prices SET provider = 'deepseek', endpoint_type = 'openai'
WHERE provider = 'openai' AND model ILIKE 'deepseek%';

UPDATE model_prices SET display_name = model WHERE display_name IS NULL OR display_name = '';
