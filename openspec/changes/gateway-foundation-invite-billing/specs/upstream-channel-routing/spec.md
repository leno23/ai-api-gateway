## ADDED Requirements

### Requirement: Channel configuration drives upstream routing

The system SHALL maintain channel configuration including provider type, base URL, credentials reference, supported models, optional JSON model mapping, priority, weight, status, concurrency and upstream rate limits sufficient to select a target channel for each request.

#### Scenario: Model resolves to an enabled channel

- **WHEN** a request names a model that is mapped to at least one enabled channel passing health and policy filters
- **THEN** the system selects a channel according to configured load balancing rules and forwards the adapted request

#### Scenario: No usable channel

- **WHEN** no channel supports the requested model or all candidates are disabled or circuit-open
- **THEN** the system fails fast with an error response and does not partially charge for upstream usage

### Requirement: Provider adapters normalize protocol differences

The system SHALL use an adapter abstraction to translate unified gateway requests into provider-specific HTTP requests and translate responses (including streaming) back to OpenAI-compatible output.

#### Scenario: Upstream returns success

- **WHEN** the selected upstream returns a successful response
- **THEN** the adapter produces a unified response or stream chunks consumable by the gateway client contract

### Requirement: Circuit breaker protects failing upstreams

The system SHALL implement a circuit breaker that temporarily stops routing to an upstream channel when error thresholds are exceeded and SHALL attempt recovery via a half-open probing strategy.

#### Scenario: Channel enters open circuit

- **WHEN** observed error rate in the configured window exceeds the configured threshold for a channel
- **THEN** the system stops selecting that channel until cooldown or probe success criteria are met

### Requirement: Tiered rate limiting

The system SHALL enforce rate limits at global, authenticated user (or key), and per-channel scopes using Redis-backed algorithms as described in the technical proposal.

#### Scenario: User exceeds per-minute limit

- **WHEN** a user exceeds the configured per-user rate limit
- **THEN** the system rejects the request without forwarding to upstream and returns an appropriate 429-class response
