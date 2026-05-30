## ADDED Requirements

### Requirement: Model catalog page

The portal SHALL provide `/pricing` listing models with provider filter, token group multiplier filter, billing type, tags, endpoint type (openai/anthropic), fuzzy search, copy-list action, price visibility toggle, multiplier visibility toggle, and grid/table view with size M/L.

#### Scenario: Filter by provider

- **WHEN** the user selects provider Anthropic in the sidebar filter
- **THEN** the model list refreshes to matching models from `GET /api/models`

#### Scenario: Toggle price display

- **WHEN** the user turns off price display
- **THEN** per-token prices are hidden in cards/table while model names remain visible

### Requirement: Model card pricing fields

Each model entry SHALL display input, completion, and cache read/write prices per 1M tokens when prices are visible, plus billing mode label (usage-based).

#### Scenario: Card view with prices on

- **WHEN** prices are visible in card view
- **THEN** input and output unit prices are shown per reference screenshot layout
