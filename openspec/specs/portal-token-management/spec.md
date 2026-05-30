# portal-token-management Specification

## Purpose
TBD - created by archiving change gateway-tenant-portal-platform. Update Purpose after archive.
## Requirements
### Requirement: API token CRUD UI

The portal SHALL provide `/console/token` to create, list, edit, enable/disable, and delete API keys with fields: name, status, quota (unlimited or fixed), token group, masked `sk-` secret, allowed models, IP whitelist.

#### Scenario: Create token shows secret once

- **WHEN** the user creates a new token
- **THEN** the full `sk-` plaintext is shown in a one-time reveal pattern (copy/QR supported)
- **AND** subsequent list views show masked values only

#### Scenario: Batch delete selected tokens

- **WHEN** the user selects multiple rows and confirms batch delete
- **THEN** the client calls the backend delete API for each selected id

### Requirement: Quick open playground from token row

The token table SHALL offer a 「聊天」 action that opens the playground with the token prefilled.

#### Scenario: Open chat from token

- **WHEN** the user chooses 聊天 on a token row
- **THEN** the app navigates to `/console/playground` with that token context

