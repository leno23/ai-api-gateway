## Why

[蓝移 API（lanyiapi.com）](https://lanyiapi.com/) 类产品的完整形态是 **公开营销站 + 租户控制台 + OpenAI/Anthropic 兼容网关 + 按量计费与运营体系**，而非仅管理员后台。本仓库已完成 **`gateway-foundation-invite-billing`**（Go 网关、额度、`sk-` 鉴权、兑换/邀请）与 **`gateway-admin-console-frontend`**（Ant Design 管理端 `/admin/*`），但 **终端开发者仍无** 模型广场、数据看板、令牌管理、操练场、钱包/邀请等自助能力。

依据 `zhencai/lanyiapi-site-audit/IMPLEMENTATION_PROMPT.md` 及同目录 `screenshots/` 巡站结果，需要将 **业务等价的多模型 API 聚合平台（用户门户 + 控制台）** 固化为可评审、可拆任务的 OpenSpec change，并与现有后端能力对齐、显式列出 API 缺口。

## What Changes

### 公开站点（未登录）

- 路由：`/`（首页营销）、`/pricing`（模型广场）、`/about`、`/register?aff=`、`/login`。
- 顶栏全局导航：首页 | 控制台 | 模型广场 | 文档 | 关于；通知、主题/语言、登录/注册。
- 首页：卖点三列、模型 Logo 墙、系统公告 Modal（支持「今日不再提示」`localStorage`）。
- 登录主按钮文案为 **「继续」**（非「登录」）；支持用户名或邮箱（与后端契约协商后实现）。

### 租户控制台（登录后 `/console`）

- 布局：顶栏 + 左侧 Semi Navigation（~240px，可收起）+ 主内容区（背景 `#f6f7f9`）。
- **操练场** `/console/playground`：分组/模型/采样参数、多轮对话、流式 SSE、调试信息、配置导入导出。
- **数据看板** `/console`：四指标卡（余额/消耗/请求/RPM·TPM）、模型消耗图表、API 多区域节点（主站/香港/美区，复制与测速）、公告/FAQ/服务可用性。
- **令牌管理** `/console/token`：`sk-` CRUD、分组倍率、额度、模型白名单、IP 限制、掩码/一次性展示、批量操作。
- **使用日志** `/console/log`、**任务日志** `/console/task`：分页筛选、列设置、流式首字耗时与计费明细。
- **钱包管理** `/console/wallet`：余额/历史消耗、在线充值（可关）、兑换码、邀请链接与收益划转。
- **个人设置** `/console/setting`：账户绑定、通知阈值、语言与显示偏好。

### 后端/API 增量（相对现有 OpenAPI）

- 用户面 REST 风格对齐参考站：`/api/user/*`（或映射到现有 `/auth/*`、`/user/*` 并扩展）。
- 模型目录：`GET /api/models`（供应商/分组/端点类型筛选、分页）。
- 令牌：`/api/tokens` CRUD（扩展现有 `POST /user/api-keys` 为完整生命周期）。
- 看板：`GET /api/dashboard/stats`、`GET /api/dashboard/charts?range=`。
- 日志：`GET /api/logs/usage`、`GET /api/logs/tasks`。
- 操练场：`POST /api/playground/chat`（SSE，专用或默认令牌）。
- 运营：公告 CRUD、通知配置、钱包充值/划转（支付可 mock）。
- 网关：补充 Anthropic 兼容 `POST /v1/messages`（P1）；多区域节点配置与测速接口。

### 技术栈（本 change 前端约定）

- **Next.js 14 App Router** + **Semi Design**（`@douyinfe/semi-ui`）+ TypeScript + Tailwind 辅助 + ECharts/Recharts + Lucide。
- 与现有 `frontend/`（Ant Design 管理端）**共存**：建议子应用 `portal/` 或 `frontend-portal/` 目录，避免与 `/admin` 设计系统混用；详见 `design.md`。

## Capabilities

### New Capabilities

- `portal-public-site`: 营销首页、关于、顶栏 IA、公告弹窗、注册/登录（含 `aff` 邀请参数）。
- `portal-console-shell`: 控制台壳层、侧栏菜单、路由守卫、401 跳转 `/login?expired=true`。
- `portal-model-pricing`: 模型广场筛选/搜索/卡片与表格视图、价格与倍率展示开关。
- `portal-token-management`: 令牌全生命周期 UI 与 `sk-` 安全展示策略。
- `portal-dashboard-analytics`: 数据看板指标卡、图表、API 节点与可用性探测展示。
- `portal-usage-task-logs`: 使用日志与异步任务日志列表与筛选。
- `portal-wallet-affiliate`: 钱包、兑换码、邀请链接与收益划转。
- `portal-playground`: 操练场三栏 UI 与流式对话调试。
- `portal-user-settings`: 个人设置、通知阈值、第三方绑定占位（按阶段 mock）。
- `portal-backend-extensions`: 租户面 HTTP API 与领域模型增量（令牌分组、定价展示、日志字段等）。

### Modified Capabilities

- `user-identity-apikeys`（`gateway-foundation-invite-billing`）：扩展令牌字段（分组、额度上限、模型白名单、IP 白名单、启用状态）及管理 API。
- `billing-quota`：对外暴露按模型分项价格（input/output/cache）与分组倍率，供广场与日志展示。
- `redeem-invite-growth`：钱包页邀请链接格式 `register?aff={code}`、待结算收益与划转到余额（若尚未实现则本 change 定义）。

### Out of Scope（另开 change 或 P2）

- 完整支付渠道（支付宝/微信）生产对接；本期允许环境变量 mock 充值回调。
- Cloudflare Turnstile 生产接入；开发环境 mock `?turnstile=`。
- 独立 VitePress 文档站（可先外链或静态 MD）。
- 管理后台渠道/兑换能力（已由 `gateway-admin-console-frontend` 覆盖）。

## Impact

- **前端**：新增 Semi 租户门户工程（与 Ant Design 管理端分离）；UI 对照 `lanyiapi-site-audit/screenshots/` 验收。
- **后端**：Go 网关新增/扩展 handler、迁移表（`token_groups`、`announcements`、`task_logs` 等）；OpenAPI 与 `GET /openapi.yaml` 同步更新。
- **运维**：多区域 base URL 配置、CORS 增加门户 Origin；可选独立域名（主站/香港/美区）。
- **依赖**：`gateway-foundation-invite-billing` 必须先可用；`gateway-admin-console-frontend` 无冲突。

## Phased Delivery（对齐 IMPLEMENTATION_PROMPT）

| Phase | 范围 | 本 change 对应 |
|-------|------|----------------|
| 1 | 用户注册登录、路由守卫、公开认证页 | `portal-public-site` + `portal-console-shell`（壳层） |
| 2 | 模型与定价、模型广场 | `portal-model-pricing` + `portal-backend-extensions` |
| 3 | 令牌 CRUD、网关鉴权增强 | `portal-token-management` + 后端扩展 |
| 4 | 看板、使用/任务日志 | `portal-dashboard-analytics` + `portal-usage-task-logs` |
| 5 | 钱包、公告、通知 | `portal-wallet-affiliate` + 后端扩展 |
| 6 | 操练场 SSE | `portal-playground` |
| 7 | 多节点、文档、i18n | 看板节点卡 + 顶栏语言（P1） |
| 8 | 管理后台增强 | 不纳入；沿用现有 admin change |

## References

- **需求来源**：`zhencai/lanyiapi-site-audit/IMPLEMENTATION_PROMPT.md`（Phase 1–8、域模型、API 路径风格、UI 规范）。
- **截图索引**：同目录 `screenshots/`（`01_console_dashboard.png` … `14_playground_url.png`）。
- **已知抓包**：`POST /api/user/login?turnstile=`；分组示例 `default`、`vip`、`claude_code`（1.5x）、`AWS分组`（3x）。
- **本仓库后端**：`openspec/changes/gateway-foundation-invite-billing/`、`GET /openapi.yaml`。
- **本仓库管理端**：`openspec/changes/gateway-admin-console-frontend/`（Ant Design `/admin`，与本 change 租户门户分离）。

## Non-goals

- 1:1 复刻蓝移品牌与域名；实现 **业务等价** 与可替换白标。
- 替换现有 Go 技术栈为 Node/Prisma（IMPLEMENTATION_PROMPT 中的备选栈不采用；后端延续 Go + PostgreSQL + Redis）。
- 在单页内合并 Semi 门户与 Ant Design 管理端为一套 UI 库。
