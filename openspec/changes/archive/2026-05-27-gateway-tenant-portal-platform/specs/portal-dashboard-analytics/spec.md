## ADDED Requirements

### Requirement: Dashboard summary cards

The data dashboard SHALL show four summary cards: account balance/consumption, usage statistics, resource consumption, and performance metrics (including sparklines), backed by `GET /api/dashboard/stats`.

#### Scenario: Load dashboard on entry

- **WHEN** the user opens `/console`
- **THEN** the four cards load with loading states and render numeric metrics from the API

### Requirement: Model analytics charts

The dashboard SHALL provide tabs for consumption distribution, call trend, call count distribution, and ranking charts aggregated by hour via `GET /api/dashboard/charts`.

#### Scenario: Switch chart tab

- **WHEN** the user selects 「调用趋势」
- **THEN** a time-series chart renders for the selected range query parameter

### Requirement: Multi-region API node card

The dashboard SHALL list configured API nodes (primary, Hong Kong, US) with copy URL, latency test, and external link actions.

#### Scenario: Copy primary base URL

- **WHEN** the user clicks copy on the primary node row
- **THEN** the configured base URL is copied to the clipboard

### Requirement: Announcements FAQ and availability

The dashboard side column SHALL list system announcements, collapsible FAQ, and service availability indicators (e.g. Claude reachability).

#### Scenario: FAQ expand

- **WHEN** the user expands an FAQ item
- **THEN** the answer content is shown without leaving the page
