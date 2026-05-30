-- Dashboard / usage logs / async task logs

ALTER TABLE request_logs
    ADD COLUMN IF NOT EXISTS request_id VARCHAR(64),
    ADD COLUMN IF NOT EXISTS token_name VARCHAR(128),
    ADD COLUMN IF NOT EXISTS token_group VARCHAR(64),
    ADD COLUMN IF NOT EXISTS time_to_first_ms INT,
    ADD COLUMN IF NOT EXISTS billing_detail JSONB;

CREATE INDEX IF NOT EXISTS idx_request_logs_request_id ON request_logs (request_id);
CREATE INDEX IF NOT EXISTS idx_request_logs_token_name ON request_logs (user_id, token_name);

CREATE TABLE IF NOT EXISTS task_logs (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    task_id         VARCHAR(128) NOT NULL,
    platform        VARCHAR(64) NOT NULL DEFAULT '',
    task_type       VARCHAR(64) NOT NULL DEFAULT '',
    status          VARCHAR(32) NOT NULL DEFAULT 'pending',
    progress        INT NOT NULL DEFAULT 0,
    detail          TEXT,
    submitted_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    finished_at     TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_task_logs_user_submitted ON task_logs (user_id, submitted_at DESC);
CREATE INDEX IF NOT EXISTS idx_task_logs_task_id ON task_logs (user_id, task_id);
