## ADDED Requirements

### Requirement: Model pricing supports token and per-request billing

The system SHALL maintain a `model_prices` (or equivalent) configuration that supports at minimum: per-1K-token prompt and completion prices for token-billed models, and a unit price for per-request billed capabilities (e.g., certain image endpoints).

#### Scenario: Token-priced model

- **WHEN** a model is configured with token billing
- **THEN** the billing engine MUST compute charges using prompt and completion token counts and the configured per-1K rates

#### Scenario: Per-request priced model

- **WHEN** a model is configured with per-request billing
- **THEN** the billing engine MUST charge exactly the configured unit price for each successful billable invocation

### Requirement: Insufficient quota blocks relay

Before forwarding a billable request, the system SHALL estimate maximum quota consumption for the operation and SHALL reject the request if the user does not have sufficient quota.

#### Scenario: Insufficient quota returns 402

- **WHEN** the authenticated user’s available quota is below the estimated requirement
- **THEN** the system MUST NOT forward the request to upstream and MUST return HTTP status 402 with a machine-readable error body

### Requirement: Pre-deduct and reconcile actual usage

For accepted requests, the system SHALL pre-deduct an estimated quota amount in Redis, then after completion SHALL reconcile against actual usage (upstream usage when present, otherwise an approved local estimation method), adjusting the pre-deducted balance to match actual cost.

#### Scenario: Non-stream completion

- **WHEN** a non-stream request completes successfully with known usage counts
- **THEN** the system MUST settle quota to the actual computed cost and MUST emit an auditable quota movement record

#### Scenario: Stream completion

- **WHEN** a stream completes
- **THEN** the system MUST compute completion usage from upstream usage if available, otherwise from the approved local estimator, MUST settle quota against the pre-deduction, and MUST emit an auditable quota movement record

### Requirement: Quota audit log

The system SHALL append an immutable-style ledger entry for each quota change including user reference, signed delta, balance snapshot after change, type classification, reference key, and timestamp.

#### Scenario: Consume after successful relay

- **WHEN** a billable upstream call completes successfully
- **THEN** the system MUST create a ledger record classifying the movement as consumption (or the documented equivalent type code)

### Requirement: Redis to database quota synchronization

The system SHALL synchronize authoritative user quota from Redis to PostgreSQL on a periodic schedule and/or when change thresholds are exceeded, and SHALL rebuild Redis from PostgreSQL on cold start.

#### Scenario: Service restart

- **WHEN** the gateway process restarts and Redis is empty for a user key
- **THEN** the system MUST load quota from PostgreSQL before accepting billable traffic for that user
