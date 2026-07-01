# AI API Gateway — 管理控制台

Next.js（App Router）+ Ant Design + Tailwind。契约对齐 `../backend/internal/openapi/spec.yaml` 与运行时 `GET /openapi.yaml`。

## 本地开发

1. 复制环境变量：

   ```bash
   cp .env.example .env.local
   ```

   将 `NEXT_PUBLIC_GATEWAY_API_URL` 设为网关根地址（无尾斜杠），例如 `http://127.0.0.1:8080`。

2. 启动网关（见仓库根 `README.md`：`cd ../backend && go run ./cmd/server`）及 Postgres/Redis。

3. 安装依赖并启动前端：

   ```bash
   npm install
   npm run dev
   ```

4. 浏览器打开 [http://localhost:3000](http://localhost:3000)，使用**管理员**账号登录（`POST /auth/login`）。普通用户调用管理接口会返回 403 并跳转无权限页。

## 脚本

| 命令 | 说明 |
|------|------|
| `npm run dev` | 开发 |
| `npm run build` / `npm start` | 生产构建与启动 |
| `npm run lint` | ESLint |
| `npm run test:e2e` | Playwright 冒烟测试（需网关与 `npm run dev` 已启动） |

### E2E 冒烟测试

先启动 Postgres/Redis、网关（默认 `http://127.0.0.1:8081`）与本前端 dev，再执行：

```bash
npx playwright install chromium   # 首次
GATEWAY_BASE_URL=http://127.0.0.1:8081 FRONTEND_BASE_URL=http://localhost:3000 npm run test:e2e
```

可选环境变量：`ADMIN_EMAIL`、`ADMIN_PASSWORD`（默认 `admin@example.com` / `Admin@12345`）。

## 规格

OpenSpec：`../openspec/changes/gateway-admin-console-frontend/`。

## 说明

- 会话：`access_token` 存 **sessionStorage**（MVP）；生产可改为 BFF 写入 **httpOnly Cookie**（见 `design.md`）。
- 未配置 `NEXT_PUBLIC_GATEWAY_API_URL` 时，应用会展示明确错误提示（开发契约）。
