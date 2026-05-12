-- Optional consume rebate queue (7.4); processor credits quota asynchronously.

CREATE TABLE IF NOT EXISTS rebate_tasks (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users (id),
    amount          BIGINT NOT NULL CHECK (amount > 0),
    status          SMALLINT NOT NULL DEFAULT 0,
    ref             VARCHAR(128),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    processed_at    TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_rebate_tasks_pending ON rebate_tasks (status, created_at);
