# 租户门户（Portal）

蓝移类产品形态的 **用户门户 + 控制台**，技术栈：Next.js + Semi Design + Tailwind。

与仓库内 `frontend/`（Ant Design 管理端 `/admin`）分离运行，默认端口 **3001**。

## 本地联调

### 1. 数据库

在 `001_init.sql` 之后按顺序执行迁移：

```bash
export DATABASE_URL='postgres://gateway:gateway@localhost:5432/ai_gateway?sslmode=disable'
psql "$DATABASE_URL" -f backend/migrations/003_token_groups_catalog.sql
psql "$DATABASE_URL" -f backend/migrations/004_api_tokens_extend.sql
psql "$DATABASE_URL" -f backend/migrations/005_dashboard_logs.sql
psql "$DATABASE_URL" -f backend/migrations/006_wallet_announcements.sql
```

### 2. 后端网关

```bash
cd backend
cp .env.example .env   # 按需修改 DATABASE_DSN、JWT_SECRET、UPSTREAM_*
go run ./cmd/server    # 默认 http://localhost:8080
```

常用环境变量见 `backend/.env.example`（`PORTAL_API_NODES_JSON`、`RECHARGE_ENABLED`、`PORTAL_ORIGIN` 等）。

### 3. 门户前端

```bash
cd portal
cp .env.example .env.local
# NEXT_PUBLIC_GATEWAY_API_URL=http://localhost:8080
# NEXT_PUBLIC_PORTAL_ORIGIN=http://localhost:3001
# NEXT_PUBLIC_DOCS_URL=https://docs.example.com   # 可选外链文档
npm install
npm run dev
```

浏览器打开 [http://localhost:3001](http://localhost:3001)。

### 4. 管理端（可选）

```bash
cd frontend
npm install
npm run dev   # 默认 http://localhost:3000
```

## API 契约

租户面 REST 与 OpenAI 网关路径见 `backend/internal/openapi/spec.yaml`，运行时可通过 `GET http://localhost:8080/openapi.yaml` 获取。

## 集成测试（后端 E2E）

需本机 Postgres（及 Redis，默认 `localhost:6379`）：

```bash
export INTEGRATION_DATABASE_DSN='host=localhost user=gateway password=gateway dbname=ai_gateway port=5432 sslmode=disable'
export INTEGRATION_REDIS_ADDR=localhost:6379   # 可选，默认与 REDIS_ADDR 一致

cd backend
go test ./internal/openapi/... -count=1
go test ./internal/integration/... -count=1 -v
```

`TestIntegration_PortalCoreFlow` 覆盖：**注册 → 登录 → 创建令牌 → 操练场请求 → 看板/节点测速**（操练场在无上游时可能返回 502，属预期）。

## 功能概览

| 模块 | 路由 / API |
|------|------------|
| 营销站 | `/`、`/about`、`/docs`、`/pricing` |
| 认证 | `/login`、`/register?aff=` |
| 控制台 | `/console/*`（看板、令牌、操练场、日志、钱包、设置） |
| 网关调用 | `Authorization: Bearer sk-...` → `/v1/chat/completions` |

OpenSpec：现行规格 `openspec/specs/portal-*/`；归档提案 `openspec/changes/archive/2026-05-27-gateway-tenant-portal-platform/`

## UI 还原验证（ai-ui-verification）

```bash
# 参考截图 baseline（来自 lanyiapi-site-audit）
# portal/.ai-verification/baselines/{home,pricing,playground}.png

bun run ai:verify:static    # lint / tsc
bun run ai:verify:ui -- --baseline ./.ai-verification/baselines/home.png --url http://localhost:3001/

# 本地截图对比（需 dev :3001）
bun -e "import { chromium } from 'playwright'; ..."
```

设计对齐要点：品牌「蓝移 API」、首页渐变 Hero、模型广场顶栏渐变、控制台侧栏 240px + 收起、操练场「模型配置 | AI 对话」双栏。
