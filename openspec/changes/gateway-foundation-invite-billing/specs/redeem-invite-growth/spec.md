## ADDED Requirements

### Requirement: Admin batch generation of redeem codes

The system SHALL allow administrators to generate batches of redeem codes using cryptographically secure randomness, formatted with a configurable prefix, guaranteed uniqueness at persistence layer, and associated quota face value and optional expiry.

#### Scenario: Generated code is unpredictable

- **WHEN** the system generates a new redeem code
- **THEN** the generation MUST use a CSPRNG (e.g., crypto/rand) and MUST NOT rely on predictable PRNG seeds

### Requirement: Redeem code redemption is atomic and idempotent

The system SHALL redeem a code in a single database transaction that: validates existence, unused status, and expiry; marks the code used with the redeemer; increases the user quota; writes a quota ledger entry; and prevents concurrent double redemption of the same code.

#### Scenario: Concurrent redemption attempts

- **WHEN** two concurrent requests attempt to redeem the same valid unused code
- **THEN** exactly one succeeds and the other receives an error indicating the code is already used or unavailable

### Requirement: Per-user invite code

Each user SHALL have exactly one stable invite code generated at registration (or first activation) using an unambiguous charset policy, stored uniquely on the user record.

#### Scenario: Invite code uniqueness

- **WHEN** the system assigns an invite code to a user
- **THEN** the invite code MUST be unique across all users

### Requirement: Registration rewards and inviter rewards

When a new user registers with a valid inviter invite code, the system SHALL record the inviter relationship and SHALL grant configured quota bonuses to invitee and inviter in one transaction, writing both quota ledger entries and an invite record with uniqueness on invitee.

#### Scenario: Invalid invite code at registration

- **WHEN** a registrant supplies an invite code that does not match any user
- **THEN** registration MUST fail or proceed without invite linkage according to the product policy, but MUST NOT grant inviter-side rewards

### Requirement: Optional spend commission (async)

If enabled by configuration, the system SHALL enqueue asynchronous processing to credit inviters a configured percentage of invitee spend after successful billable consumption, without blocking the primary request completion path.

#### Scenario: Commission disabled

- **WHEN** commission configuration is disabled
- **THEN** no commission ledger entries are created for invitees’ consumption
