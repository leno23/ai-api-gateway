## ADDED Requirements

### Requirement: Channel list

The console SHALL list channels from `GET /admin/channels` in a table including at least name, provider, base URL, masked API key preview, priority, weight, status, and rate limit.

#### Scenario: Empty state

- **WHEN** the API returns an empty list
- **THEN** the UI shows an empty state with a primary action to create a channel

### Requirement: Channel create and update

The console SHALL support creating and updating channels using `POST /admin/channels` and `PUT /admin/channels/:id` with fields aligned to OpenAPI `ChannelPayload`.

#### Scenario: Model mapping JSON validation

- **WHEN** the operator enters `model_mapping` as JSON text
- **THEN** the client validates JSON shape (object of string to string) before submit
- **AND** server validation errors are shown on the form

#### Scenario: Models array editing

- **WHEN** the operator edits the allowed models list
- **THEN** the payload sends `models` as a string array consistent with the gateway API

### Requirement: Channel delete

The console SHALL delete a channel via `DELETE /admin/channels/:id` only after explicit confirmation.

#### Scenario: Confirm modal prevents accidental delete

- **WHEN** the user clicks delete
- **THEN** a confirmation modal is shown summarizing the channel name or id
- **AND** cancel leaves data unchanged
