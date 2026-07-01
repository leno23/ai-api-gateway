# portal-user-settings Specification

## Purpose
TBD - created by archiving change gateway-tenant-portal-platform. Update Purpose after archive.
## Requirements
### Requirement: Personal settings tabs

The settings page at `/console/setting` SHALL provide account management (third-party binding placeholders: email, WeChat, GitHub, Discord, OIDC, Telegram, LinuxDO) and other settings (notification channels, balance alert threshold, notification email, price/privacy/sidebar sub-tabs).

#### Scenario: Save notification threshold

- **WHEN** the user updates balance alert threshold and saves
- **THEN** the client persists settings via the backend notification config API

### Requirement: Language preference

The portal SHALL default to Simplified Chinese and allow switching language from the top bar, persisting preference per user or local storage until backend i18n is available.

#### Scenario: Switch to English

- **WHEN** the user selects English from the language control
- **THEN** UI strings switch to English catalog where translated

