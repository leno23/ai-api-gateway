## ADDED Requirements

### Requirement: User registration and authentication

The system SHALL support user registration and login flows using email-based verification as described in the technical proposal, storing password material using an approved slow hashing algorithm.

#### Scenario: Successful registration

- **WHEN** a new user completes registration with valid email verification
- **THEN** the system creates a user record in active status suitable for subsequent API key issuance

### Requirement: API key lifecycle management

The system SHALL allow authenticated users to create, list, disable, and delete API keys. Raw key material MUST only be shown once at creation; persisted form MUST be a one-way hash with a display prefix.

#### Scenario: Revoked key cannot authenticate

- **WHEN** a user disables an API key
- **THEN** subsequent gateway requests using that key MUST be rejected

### Requirement: Role-based access for administration

The system SHALL distinguish at minimum normal users and administrative roles capable of managing channels, redeem batches, and user moderation as defined in the technical proposal RBAC section.

#### Scenario: Non-admin cannot access admin routes

- **WHEN** a normal user calls an administrative HTTP endpoint
- **THEN** the system denies access with a 403-class response
