## Why

后端 **`gateway-foundation-invite-billing`** 已提供管理端 HTTP API（渠道、兑换码、用户封禁等）与 **`GET /openapi.yaml`**，但缺少 **可运营的前端控制台**：运营/管理员仍依赖手工调用接口，无法安全、一致地完成日常配置与稽核。本 change 将 **前端能力** 以 OpenSpec 固化，便于评审、拆分任务与和后端契约对齐。

## What Changes

- 定义 **管理控制台（Admin Console）** 的信息架构、页面与交互（对齐现有 `/admin/*` 与 OpenAPI）。
- 定义 **认证与会话**：JWT 存储策略、路由守卫、管理员与普通用户分流（控制台仅管理员可进入核心运营页）。
- 定义 **渠道管理 UI**：列表、创建/编辑、模型列表与 `model_mapping` 编辑、敏感字段脱敏展示。
- 定义 **兑换运营 UI**：批量生成表单、兑换码分页列表、按状态统计看板。
- 定义 **用户治理 UI**：按 ID/邮箱查询（若后端后续扩展则对接）、封禁/解封（`PATCH /admin/users/:id/status`）。
- （可选 P1）定义 **租户自助区**：登录用户查看额度/邀请码、创建 API Key（明文仅一次）、兑换码兑换；依赖 `/user/*` 与扩展只读接口。

## Capabilities

### New Capabilities

- `console-shell-auth`: 布局、登录页、Token 生命周期、403/401 统一处理、环境化 API Base URL。
- `console-channels`: 渠道 CRUD 与模型映射编辑，对接 `GET/POST/PUT/DELETE /admin/channels*`。
- `console-redeem`: 批量生成、列表与统计，对接 `POST /admin/redeem/batch`、`GET /admin/redeem/codes`、`GET /admin/redeem/stats`。
- `console-users`: 用户状态变更，对接 `PATCH /admin/users/:id/status`；列表/搜索以后端扩展为准。
- `console-tenant-portal`（P1）: 租户自助能力与网关 `sk-` 调用说明入口。

### Modified Capabilities

- （无）不修改已归档后端 spec；仅新增前端能力规格。若实现时发现后端缺字段，另开 **gateway API 增量** change。

## Impact

- **前端工程**：本仓库目录为 `frontend/`（见 `frontend/README.md`）；亦可拆至独立仓库；技术选型见 `design.md`。
- **后端**：无强制变更；可选增加「用户分页查询」等管理接口以支撑用户列表页（可另列小 change）。
- **安全**：控制台部署于内网或 SSO 之后；CORS、Cookie `Secure`/`SameSite`、XSS 与 CSRF 策略在实现阶段落地。

## References

- 后端 change：`openspec/changes/gateway-foundation-invite-billing/`
- HTTP 契约：`GET /openapi.yaml`（`backend/internal/openapi/spec.yaml`）
- 架构背景：`docs/technical-proposal.md.md`、`docs/invitation-and-billing-design.md.md`
