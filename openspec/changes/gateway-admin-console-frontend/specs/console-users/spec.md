## ADDED Requirements

### Requirement: User status moderation

The console SHALL call `PATCH /admin/users/:id/status` with JSON `{ "status": <int> }` to activate or ban users, with confirmation for destructive changes.

#### Scenario: Ban confirmation

- **WHEN** the operator sets status to disabled (0)
- **THEN** the UI requires confirmation describing impact on sign-in and API access

#### Scenario: Invalid id handling

- **WHEN** the server returns 404 or 500
- **THEN** the UI shows an actionable error without leaking stack traces

### Requirement: Future user directory (optional backlog)

If the backend later exposes `GET /admin/users` with pagination and search, the console SHOULD replace the MVP single-id form with a searchable table. Until then this requirement is **documentation-only** and does not block MVP.

#### Scenario: Placeholder

- **WHEN** no list API exists
- **THEN** the MVP implements only the PATCH-by-id flow described above
