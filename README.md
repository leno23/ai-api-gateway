# AI API Gateway

Go 实现的 OpenAI 兼容 API 网关（多租户 API Key、额度、兑换/邀请、渠道路由）。本仓库为 **monorepo**：**`backend/`** 为网关服务，**`frontend/`** 为管理控制台（Next.js + Ant Design）。

设计说明见 `docs/` 与 OpenSpec：**后端** `openspec/changes/gateway-foundation-invite-billing/`；**管理控制台前端提案** `openspec/changes/gateway-admin-console-frontend/`。

## 目录结构

| 路径 | 说明 |
|------|------|
| `backend/` | Go 网关：`go.mod`、`cmd/server`、`internal/`、`migrations/` |
| `frontend/` | 管理控制台（`npm run dev`，详见 `frontend/README.md`） |
| `docker-compose.yml` | 本地 Postgres + Redis（仓库根，与后端 `.env` 配合） |
| `openspec/` | 变更提案与规格 |

## 运行（后端）

1. 在仓库根执行：`docker compose up -d`（Postgres + Redis，若已自备可跳过）
2. 应用迁移：`backend/migrations/001_init.sql`
3. 若启用消费返利：追加 `backend/migrations/002_rebate_tasks.sql`，并设置 `REBATE_ENABLED=true`、`REBATE_BPS`（万分比，如 100=1%）
4. 复制 `backend/.env.example` 为 `backend/.env` 并按需填写
5. `cd backend && go run ./cmd/server`

## 运行（管理控制台）

1. `cd frontend && cp .env.example .env.local`，设置 `NEXT_PUBLIC_GATEWAY_API_URL`（与网关监听地址一致，无尾斜杠）
2. `npm install && npm run dev`，默认 [http://localhost:3000](http://localhost:3000)
3. 使用具备管理员 `role` 的账号登录（见 `frontend/README.md`）

## API 文档

- 运行时：`GET /openapi.yaml`（与仓库内嵌 spec 同源：`backend/internal/openapi/spec.yaml`）
- 可将该文件导入 Postman / Stoplight / Swagger UI

## 测试

在 **`backend/`** 目录下执行：

- 单元测试（默认）：`go test ./...`
- 集成测试（需可用 Postgres + Redis）：  
  `INTEGRATION_DATABASE_DSN=... INTEGRATION_REDIS_ADDR=127.0.0.1:6379 go test ./internal/integration/... -count=1`  
  未设置 `INTEGRATION_DATABASE_DSN` 时集成用例会自动 `Skip`。

## OpenSpec

归档前在仓库根目录执行：

`npx @fission-ai/openspec@latest status --change gateway-foundation-invite-billing`
