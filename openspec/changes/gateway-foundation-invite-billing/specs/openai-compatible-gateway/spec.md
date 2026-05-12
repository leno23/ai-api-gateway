## ADDED Requirements

### Requirement: OpenAI compatible routes are exposed

The system SHALL expose HTTP endpoints compatible with OpenAI-style clients for: `POST /v1/chat/completions`, `POST /v1/embeddings`, `POST /v1/images/generations`, `POST /v1/audio/transcriptions`, `POST /v1/audio/speech`, and `GET /v1/models`.

#### Scenario: Chat completion non-streaming

- **WHEN** a client sends a valid `POST /v1/chat/completions` request with `stream: false` and a valid user API key
- **THEN** the system returns a JSON body in OpenAI-compatible shape suitable for downstream clients

#### Scenario: Chat completion streaming

- **WHEN** a client sends `POST /v1/chat/completions` with `stream: true`
- **THEN** the system returns an SSE stream using `data: ...` chunks and terminates with a stream end marker consistent with OpenAI-style streaming

### Requirement: API key authentication on gateway routes

The system SHALL authenticate gateway requests using a user-issued API key for all `/v1/*` relay routes unless explicitly documented as public.

#### Scenario: Missing API key

- **WHEN** a client calls a protected `/v1/*` route without credentials
- **THEN** the system rejects the request with an appropriate 4xx error and does not forward to upstream providers

#### Scenario: Invalid or disabled API key

- **WHEN** a client presents an API key that is unknown, revoked, or expired
- **THEN** the system rejects the request with an appropriate 4xx error and does not forward to upstream providers

### Requirement: Models listing reflects routable configuration

The system SHALL implement `GET /v1/models` to return models that the gateway is configured to route for the authenticated tenant according to policy (allowlist on keys vs global catalog as defined in implementation).

#### Scenario: Successful models list

- **WHEN** a client with a valid API key calls `GET /v1/models`
- **THEN** the system returns a JSON payload including at least an array of model descriptors with stable `id` strings usable in subsequent requests
