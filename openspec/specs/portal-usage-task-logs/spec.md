# portal-usage-task-logs Specification

## Purpose
TBD - created by archiving change gateway-tenant-portal-platform. Update Purpose after archive.
## Requirements
### Requirement: Usage log list page

The portal SHALL provide `/console/log` with filters for time range, token name, model, request ID, and group; table columns for time, token, group, type, model, time-to-first-token, input tokens, and billing details; plus compact mode, column settings, and pagination defaulting to 10 rows per page.

#### Scenario: Filter by model name

- **WHEN** the user enters a model filter and applies search
- **THEN** `GET /api/logs/usage` is called with the model query parameter

### Requirement: Task log list page

The portal SHALL provide `/console/task` listing async tasks with columns: submit time, end time, duration, platform, type, task ID, status, progress, and detail action.

#### Scenario: Empty task log state

- **WHEN** the user has no task records
- **THEN** the page shows an empty state matching reference screenshot behavior

