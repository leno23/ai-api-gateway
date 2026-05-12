## ADDED Requirements

### Requirement: Redeem batch generation

The console SHALL submit `POST /admin/redeem/batch` with `count`, `quota`, and optional `expires_days`, and SHALL display returned codes in a copy-friendly way.

#### Scenario: Successful batch shows codes once

- **WHEN** generation succeeds
- **THEN** the UI lists all returned codes and offers copy or CSV export
- **AND** a warning explains codes may not be retrievable later if the backend does not persist plaintext beyond creation context

### Requirement: Redeem code listing

The console SHALL page through `GET /admin/redeem/codes` with `page`, `page_size`, and optional `status` filter.

#### Scenario: Pagination controls

- **WHEN** total exceeds page size
- **THEN** the user can move between pages without losing filter state

### Requirement: Redeem statistics

The console SHALL render `GET /admin/redeem/stats` as a compact summary (e.g. counts per status).

#### Scenario: Zero data

- **WHEN** no codes exist
- **THEN** the stats section still renders without errors
