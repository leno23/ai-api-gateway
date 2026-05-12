-- Core schema for AI API Gateway MVP (see docs/technical-proposal.md.md & invitation-and-billing-design.md.md)

CREATE TABLE IF NOT EXISTS users (
    id              BIGSERIAL PRIMARY KEY,
    username        VARCHAR(64) UNIQUE NOT NULL,
    email           VARCHAR(255) UNIQUE NOT NULL,
    password_hash   VARCHAR(255) NOT NULL,
    role            SMALLINT NOT NULL DEFAULT 1,
    status          SMALLINT NOT NULL DEFAULT 1,
    balance         DECIMAL(16,6) NOT NULL DEFAULT 0,
    quota           BIGINT NOT NULL DEFAULT 0,
    used_quota      BIGINT NOT NULL DEFAULT 0,
    invite_code     VARCHAR(32) UNIQUE,
    invited_by      BIGINT REFERENCES users (id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS api_keys (
    id            BIGSERIAL PRIMARY KEY,
    user_id       BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    name          VARCHAR(128) NOT NULL,
    key_hash      VARCHAR(255) NOT NULL,
    key_prefix    VARCHAR(16) NOT NULL,
    status        SMALLINT NOT NULL DEFAULT 1,
    models        TEXT[],
    rate_limit    INT NOT NULL DEFAULT 60,
    expires_at    TIMESTAMPTZ,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_api_keys_user_id ON api_keys (user_id);
CREATE INDEX IF NOT EXISTS idx_api_keys_key_hash ON api_keys (key_hash);

CREATE TABLE IF NOT EXISTS channels (
    id             BIGSERIAL PRIMARY KEY,
    name           VARCHAR(128) NOT NULL,
    provider       VARCHAR(32) NOT NULL,
    base_url       VARCHAR(512) NOT NULL,
    api_key        VARCHAR(512) NOT NULL,
    models         TEXT[] NOT NULL DEFAULT '{}',
    model_mapping  JSONB,
    priority       INT NOT NULL DEFAULT 0,
    weight         INT NOT NULL DEFAULT 1,
    status         SMALLINT NOT NULL DEFAULT 1,
    max_concurrent INT NOT NULL DEFAULT 100,
    rate_limit     INT NOT NULL DEFAULT 1000,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_channels_provider ON channels (provider, status);

CREATE TABLE IF NOT EXISTS model_prices (
    id                 BIGSERIAL PRIMARY KEY,
    model              VARCHAR(128) UNIQUE NOT NULL,
    prompt_price       BIGINT NOT NULL DEFAULT 0,
    completion_price   BIGINT NOT NULL DEFAULT 0,
    unit_price         BIGINT NOT NULL DEFAULT 0,
    billing_type       SMALLINT NOT NULL DEFAULT 1,
    currency           VARCHAR(8) NOT NULL DEFAULT 'CNY',
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS request_logs (
    id                 BIGSERIAL PRIMARY KEY,
    user_id            BIGINT REFERENCES users (id),
    api_key_id         BIGINT REFERENCES api_keys (id),
    channel_id         BIGINT REFERENCES channels (id),
    model              VARCHAR(128),
    request_method     VARCHAR(16),
    request_path       VARCHAR(256),
    prompt_tokens      INT NOT NULL DEFAULT 0,
    completion_tokens  INT NOT NULL DEFAULT 0,
    total_tokens       INT NOT NULL DEFAULT 0,
    cost_quota         BIGINT NOT NULL DEFAULT 0,
    latency_ms         INT,
    status_code        INT,
    error_message      TEXT,
    ip_address         INET,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_request_logs_user_id ON request_logs (user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_request_logs_channel_id ON request_logs (channel_id, created_at DESC);

CREATE TABLE IF NOT EXISTS recharge_records (
    id             BIGSERIAL PRIMARY KEY,
    user_id        BIGINT NOT NULL REFERENCES users (id),
    amount         DECIMAL(16,2) NOT NULL,
    payment_method VARCHAR(32),
    trade_no       VARCHAR(128) UNIQUE,
    status         SMALLINT NOT NULL DEFAULT 0,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at   TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS redeem_codes (
    id         BIGSERIAL PRIMARY KEY,
    code       VARCHAR(64) UNIQUE NOT NULL,
    quota      BIGINT NOT NULL,
    used_by    BIGINT REFERENCES users (id),
    status     SMALLINT NOT NULL DEFAULT 1,
    expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    used_at    TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_redeem_codes_status ON redeem_codes (status);

CREATE TABLE IF NOT EXISTS quota_logs (
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT NOT NULL REFERENCES users (id),
    delta      BIGINT NOT NULL,
    balance    BIGINT NOT NULL,
    type       SMALLINT NOT NULL,
    reference  VARCHAR(128),
    remark     VARCHAR(512),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_quota_logs_user ON quota_logs (user_id, created_at DESC);

CREATE TABLE IF NOT EXISTS invite_records (
    id             BIGSERIAL PRIMARY KEY,
    inviter_id     BIGINT NOT NULL REFERENCES users (id),
    invitee_id     BIGINT NOT NULL REFERENCES users (id),
    inviter_bonus  BIGINT NOT NULL,
    invitee_bonus  BIGINT NOT NULL,
    status         SMALLINT NOT NULL DEFAULT 1,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (invitee_id)
);

CREATE INDEX IF NOT EXISTS idx_invite_records_inviter ON invite_records (inviter_id);
