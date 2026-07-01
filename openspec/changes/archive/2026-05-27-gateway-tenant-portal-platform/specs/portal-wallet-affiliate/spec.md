## ADDED Requirements

### Requirement: Wallet overview

The wallet page at `/console/wallet` SHALL display current balance, historical consumption, and request count.

#### Scenario: View wallet after login

- **WHEN** the user opens `/console/wallet`
- **THEN** balance metrics load from the wallet or user self API

### Requirement: Redeem code recharge

The wallet page SHALL provide a redeem code input that calls the existing redeem API (`POST /user/redeem` or aliased `/api/wallet/redeem`).

#### Scenario: Successful redeem

- **WHEN** the user submits a valid unused code
- **THEN** balance increases and a success toast is shown

### Requirement: Affiliate invite link and transfer

The wallet page SHALL display an invite URL of the form `{portalOrigin}/register?aff={invite_code}`, pending affiliate earnings, and an action to transfer earnings to balance.

#### Scenario: Copy invite link

- **WHEN** the user copies the invite link
- **THEN** the clipboard contains the full URL with the user’s invite code

### Requirement: Online recharge gate

When online recharge is disabled by configuration, the UI SHALL show guidance to contact admin or use redeem codes instead of a payment form.

#### Scenario: Recharge disabled

- **WHEN** `RECHARGE_ENABLED` is false
- **THEN** no payment checkout UI is rendered
