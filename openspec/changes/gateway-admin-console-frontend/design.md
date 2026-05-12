## Context

- 网关已暴露 **JWT 管理端**（`Authorization: Bearer <access_token>`，`role` 为管理员）与 **OpenAPI 3** 描述。
- 本 change 描述 **B 端控制台** 的产品与技术设计，供实现与评审；不绑定单一 monorepo 路径（可与 `zhencai` 中 Next + Ant Design + Tailwind 栈对齐）。

## Goals / Non-Goals

**Goals**

- **MVP**：管理员登录 → 渠道管理、兑换码生成/列表/统计、用户封禁（按用户 ID）。
- **契约驱动**：所有表单字段与后端 DTO/OpenAPI 对齐；错误码与文案映射（401/403/409/502）一致。
- **可观测入口**：提供跳转或内嵌说明（健康检查 URL、Prometheus `/metrics` 为只读运维面，控制台可仅放链接）。
- **无障碍与审计友好**：关键操作二次确认（批量生成、删除渠道、封禁）；操作结果 Toast/消息可追踪。

**Non-Goals（本期）**

- 不接 **支付网关** 自动充值 UI（`recharge_records` 可二期）。
- 不做 **实时对话调试台**（Playground）除非单独 change。
- **返利任务队列看板**（`rebate_tasks`）后端若无列表 API，本期仅预留文案/二期任务。
- 不替代 **Prometheus/Grafana**；控制台不做重型监控绘图。

## 技术选型（建议）

| 区域 | 建议 | 说明 |
|------|------|------|
| 框架 | Next.js（App Router）或等价 SPA | SSR 可选；管理端多为 CSR + BFF 亦可 |
| UI | Ant Design + Tailwind | 与团队 `antd-best-practices`、设计 token 一致 |
| 请求 | `fetch` + 轻量封装 或 TanStack Query | 统一 `baseURL`、`Authorization` 注入 |
| 表单 | Ant Design Form | 渠道 `models[]`、`model_mapping` 用动态表单项或 JSON 编辑器（校验 JSON） |
| 配置 | `NEXT_PUBLIC_GATEWAY_API_URL` | 禁止把管理员密码写入仓库；生产用环境变量 |

## 信息架构（IA）

```
/login
/admin
  /channels          渠道列表 + 新建/编辑抽屉或子页
  /redeem            子路由：生成 | 列表 | 统计（或 Tab）
  /users             用户治理（MVP：按 ID 封禁表单；二期表格）
/settings            可选：展示当前 OpenAPI 链接、环境名
```

普通 **租户自助**（P1）建议独立前缀 `/portal` 或与主站账号体系合并后再定。

## 认证与权限

1. `POST /auth/login` → 存 `access_token`（**httpOnly Cookie** 优先；若纯 Bearer 存 **memory + sessionStorage** 需评估 XSS）。
2. 路由守卫：请求 `GET /admin/channels` 或轻量 **whoami**（若后端暂无则首次管理请求 403 即视为非管理员）。
3. **403**：统一「无权限」页，提示联系管理员开通 `role`。

## 关键交互

### 渠道

- 列表展示：`name`、`provider`、`base_url`、`api_key_preview`、优先级、权重、状态、`rate_limit`。
- 编辑：`models` 多选或标签输入（字符串数组）；`model_mapping` JSON 文本区 + 校验。
- 删除：二次确认；调用 `DELETE /admin/channels/:id`。

### 兑换码

- **批量生成**：`count`、`quota`、`expires_days`；成功后展示 **一次性列表**（可复制导出 CSV/文本）；提示码仅此时全量可见（与后端一致：生成响应含 `codes`）。
- **列表**：分页、`status` 筛选；展示 `code`（可掩码中间段）、`quota`、`status`、`expires_at`、`used_by`（若有）。
- **统计**：`by_status` 卡片或简单表格。

### 用户

- MVP：`用户 ID` + `status`（1/0）表单调用 `PATCH`。
- 二期：用户分页、搜索需后端 `GET /admin/users`（另 spec）。

## Risks

| 风险 | 缓解 |
|------|------|
| 管理员 Token 泄露 | 短 TTL、HTTPS、HttpOnly、操作审计日志（后端后续） |
| `model_mapping` JSON 输入错误 | 客户端 schema 校验 + 后端 400 提示展示 |
| CORS | 网关配置允许控制台 Origin；或 BFF 同源代理 |

## References

- `openspec/changes/gateway-admin-console-frontend/proposal.md`
- `backend/internal/openapi/spec.yaml`
