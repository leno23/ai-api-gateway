## ADDED Requirements

### Requirement: Tenant self-service (P1)

If implemented in the same application or module, a **non-admin** authenticated area MAY expose:

- Creation of API keys via `POST /user/api-keys` with one-time display of the raw key
- Redeem flow via `POST /user/redeem`
- Read-only display of quota and invite code when the backend exposes stable read APIs or safe derivations

#### Scenario: API key shown once

- **WHEN** key creation succeeds
- **THEN** the UI shows the raw key in a modal with copy button
- **AND** closing the modal does not retain the raw key in client state beyond the session policy

#### Scenario: Out of scope for strict admin-only deployments

- **WHEN** the product decision is admin-console only
- **THEN** this capability can be omitted without violating other specs in this change
