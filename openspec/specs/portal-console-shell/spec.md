# portal-console-shell Specification

## Purpose
TBD - created by archiving change gateway-tenant-portal-platform. Update Purpose after archive.
## Requirements
### Requirement: Console layout with Semi sidebar

Authenticated console routes SHALL use a fixed top bar, left Semi Navigation sidebar (~240px, collapsible), and main content on background `#f6f7f9`.

#### Scenario: User opens data dashboard

- **WHEN** the user navigates to `/console` with a valid session
- **THEN** the sidebar shows groups 聊天 / 控制台 / 个人中心 with items matching the reference IA (操练场, 数据看板, 令牌管理, 使用日志, 任务日志, 钱包管理, 个人设置)

### Requirement: Session guard for console routes

All routes under `/console` SHALL require authentication.

#### Scenario: Session expired

- **WHEN** the backend returns 401 on a console API call
- **THEN** the UI shows a toast equivalent to 「未登录或登录已过期」
- **AND** redirects to `/login?expired=true`

### Requirement: Configurable gateway API base URL

The portal SHALL read the gateway base URL from `NEXT_PUBLIC_GATEWAY_API_URL` (no trailing slash).

#### Scenario: Missing configuration in development

- **WHEN** the base URL is unset in development
- **THEN** the app surfaces a clear configuration error

