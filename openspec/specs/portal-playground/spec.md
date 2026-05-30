# portal-playground Specification

## Purpose
TBD - created by archiving change gateway-tenant-portal-platform. Update Purpose after archive.
## Requirements
### Requirement: Playground layout and parameters

The playground at `/console/playground` SHALL provide a left configuration panel (custom body toggle, group, model, image URL for multimodal, Temperature/TopP/Frequency/Presence penalties with enable switches) and a right conversation panel (user/assistant bubbles).

#### Scenario: Send a chat message

- **WHEN** the user sends a message with valid group and model selected
- **THEN** the client calls `POST /api/playground/chat` with SSE streaming
- **AND** assistant content appends incrementally in the thread

### Requirement: Playground message actions

The conversation UI SHALL support regenerate, copy, edit, delete, and show-debug for messages.

#### Scenario: Regenerate last assistant reply

- **WHEN** the user triggers regenerate on the latest assistant message
- **THEN** a new completion replaces or appends per product rules while logging usage

### Requirement: Playground config import export

The playground SHALL export and import session parameters as JSON.

#### Scenario: Export config

- **WHEN** the user exports configuration
- **THEN** a JSON file or clipboard payload contains group, model, and slider values

