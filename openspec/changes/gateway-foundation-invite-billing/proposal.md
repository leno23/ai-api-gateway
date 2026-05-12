## Why

需要实现 **AI API 中转站（OpenAI 兼容统一入口）** 的 MVP 基础能力，并落地文档中已定义的 **邀请码 / 兑换码增长体系** 与 **额度（Quota）计费**，以便多租户按量消费、可审计、可运营；当前仅有设计文档，尚无可执行的规格与任务拆分。

## What Changes

- 定义 **OpenAI 兼容网关面**（对话 / 嵌入 / 图像 / 语音 / 模型列表，含 SSE 流式）及 API Key 鉴权入口。
- 定义 **上游渠道与路由**（多厂商适配、模型映射、负载均衡、熔断、多级限流）。
- 定义 **计费与额度**（内部整数额度单位、`model_prices`、预估与预扣、402 余额不足、实际结算与异步落库、Redis 与 DB 同步策略）。
- 定义 **兑换码与邀请增长**（兑换码批量生成与兑换事务、每用户邀请码、注册/邀请奖励、`quota_logs` / `invite_records`、可选消费返利异步）。
- 定义 **用户与 API Key 生命周期**（注册登录、Key 哈希存储、基础 RBAC）。

## Capabilities

### New Capabilities

- `openai-compatible-gateway`: 对外统一 OpenAI 兼容 HTTP API、流式与非流式行为、认证与请求元数据约定。
- `upstream-channel-routing`: 渠道配置、协议适配、模型解析、调度与容错、限流与熔断。
- `billing-quota`: 定价表、预估/预扣/结算、额度日志与请求日志、402 语义。
- `redeem-invite-growth`: 兑换码与邀请码生命周期、奖励与审计、与额度系统的集成。
- `user-identity-apikeys`: 用户账户、API Key 管理与角色边界（相对管理员能力）。

### Modified Capabilities

- （无）当前 `openspec/specs/` 下尚无已归档基线规格，本次全部为新增能力。

## Impact

- **后端**：Go 网关核心、中间件、适配器、路由引擎、计费与额度服务、管理 API；PostgreSQL / Redis  schema 与迁移。
- **前端**（后续迭代）：Next.js 控制台与后台（本提案聚焦 API 与领域行为，UI 可另开 change）。
- **运维**：Nginx、监控日志、密钥与渠道配置管理。
- **依赖**：与支付渠道（支付宝/微信）对接可在后续 change 展开，本提案仅要求额度入账路径可扩展。

## References

- `docs/technical-proposal.md.md` — 总体架构与模块划分。
- `docs/invitation-and-billing-design.md.md` — 邀请/兑换与计费细节。
- `GET /openapi.yaml` — 运行时 OpenAPI 3 描述（源码 `backend/internal/openapi/spec.yaml`）；README 含本地联调与集成测试说明。
- OpenSpec change：`openspec/changes/gateway-foundation-invite-billing/`（`tasks.md`、`design.md`、`specs/*`）。
- 前端控制台（另 change）：`openspec/changes/gateway-admin-console-frontend/` — 管理端 UI 与 OpenAPI 对齐。
