# portal-backend-extensions Specification

## Purpose
TBD - created by archiving change gateway-tenant-portal-platform. Update Purpose after archive.
## Requirements
### Requirement: Tenant REST API surface

The gateway SHALL expose tenant-authenticated JSON APIs documented in OpenAPI under paths aligned with `/api/user/*`, `/api/tokens`, `/api/models`, `/api/dashboard/*`, `/api/logs/*`, `/api/playground/chat`, and wallet/announcement endpoints defined in `design.md`.

#### Scenario: OpenAPI includes new paths

- **WHEN** a client fetches `GET /openapi.yaml`
- **THEN** new tenant routes appear with JWT security scheme

### Requirement: Token group multiplier in billing

The billing service SHALL apply a per-token or per-user group multiplier when calculating final quota debit for a request.

#### Scenario: VIP group 1.5x

- **WHEN** a request uses a token in group `claude_code` with multiplier 1.5
- **THEN** the settled cost equals base cost × 1.5 rounded per integer quota rules

### Requirement: Gateway key policy enforcement

Before forwarding upstream, the gateway SHALL enforce per-key quota limit, enabled flag, allowed model list, and IP whitelist when configured.

#### Scenario: Disabled key rejected

- **WHEN** a disabled key calls `POST /v1/chat/completions`
- **THEN** the gateway returns 401 with a non-leaky error

### Requirement: Usage log fields for console

Usage log API responses SHALL include fields required by the console: timestamp, token name, group, type, model, time-to-first-token, input tokens, and billing detail breakdown.

#### Scenario: Paginated usage logs

- **WHEN** the client calls `GET /api/logs/usage` with page and filters
- **THEN** the response includes total count and rows with billing detail JSON

### Requirement: Playground chat logging

Playground requests SHALL write usage logs identifiable by a playground token name (e.g. `playground-default`) when using the shared playground key strategy.

#### Scenario: Playground completion billed

- **WHEN** a playground SSE session completes successfully
- **THEN** a usage log row is persisted and quota is debited

