# portal-public-site Specification

## Purpose
TBD - created by archiving change gateway-tenant-portal-platform. Update Purpose after archive.
## Requirements
### Requirement: Marketing home page

The portal SHALL serve a public home page at `/` with hero copy equivalent to the reference product (high-concurrency gateway, cost efficiency, multi-model), feature cards, partner/logo strip, and CTAs to register and login.

#### Scenario: Guest views home

- **WHEN** an unauthenticated user opens `/`
- **THEN** the page renders without requiring login
- **AND** primary CTAs link to `/register` and `/login`

### Requirement: Global top navigation

The portal SHALL show a persistent top bar on public and console layouts with links: Home, Console, Model pricing (`/pricing`), Docs, About; plus notification affordance, theme/display control, language switch, and auth actions (login/register or user menu).

#### Scenario: Logged-in user sees console link

- **WHEN** a valid JWT session exists
- **THEN** the Console link navigates to `/console` without forcing re-login

### Requirement: Login page primary action label

The sign-in form primary submit control SHALL display the label **继续** (not 「登录」).

#### Scenario: User submits credentials

- **WHEN** the user clicks the primary button on `/login`
- **THEN** the visible label is 「继续」

### Requirement: Registration with affiliate query param

The registration flow SHALL accept `?aff=` on `/register` and pass the invite code to the backend register API.

#### Scenario: Invite link registration

- **WHEN** the user opens `/register?aff=ABC123` and completes registration
- **THEN** the invite code is included in the register request per backend contract

### Requirement: System announcement modal on home

The home page SHALL support a dismissible announcement modal (e.g. enterprise contact QR) with optional «do not show again today» stored in `localStorage`.

#### Scenario: User dismisses for today

- **WHEN** the user chooses today’s dismiss option
- **THEN** the modal does not reappear until the next calendar day (per implementation key)

