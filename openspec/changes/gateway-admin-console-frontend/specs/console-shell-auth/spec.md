## ADDED Requirements

### Requirement: Admin authentication UI

The console SHALL provide a sign-in view that exchanges credentials for a JWT using `POST /auth/login` and SHALL NOT persist passwords in localStorage.

#### Scenario: Successful login redirects to admin home

- **WHEN** a user submits valid email and password
- **THEN** the client stores the session per `design.md` (prefer httpOnly cookie pattern or agreed token storage)
- **AND** the user is navigated to the protected admin area

#### Scenario: Invalid credentials show server message

- **WHEN** the server returns 401
- **THEN** the UI displays a non-leaky error message and does not clear unrelated fields unnecessarily

### Requirement: Session guard for admin routes

All routes under `/admin` (or equivalent prefix) SHALL require an authenticated session and SHALL handle 403 as «insufficient role».

#### Scenario: Unauthenticated access redirects to login

- **WHEN** no valid session exists and the user opens a protected route
- **THEN** the application redirects to the login page

#### Scenario: Normal user receives forbidden treatment

- **WHEN** the backend returns 403 on an admin endpoint
- **THEN** the UI shows a dedicated forbidden page with guidance

### Requirement: Configurable API base URL

The console SHALL read the gateway base URL from build-time or runtime public configuration (e.g. `NEXT_PUBLIC_GATEWAY_API_URL`) so the same build can target staging and production.

#### Scenario: Missing base URL fails fast in development

- **WHEN** the base URL is not configured in development
- **THEN** the app surfaces a clear configuration error
