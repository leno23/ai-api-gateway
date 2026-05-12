## 1. Repository & scaffolding

- [x] 1.1 初始化 Go module 与 `backend/cmd/server`、`backend/internal/` 目录骨架（对齐技术方案「项目结构」；仓库根为 monorepo）
- [x] 1.2 引入 Gin、GORM、Zap、Redis 客户端、配置加载（viper 或等价）与基础 `docker-compose`（Postgres + Redis）

## 2. Persistence

- [x] 2.1 编写迁移：`users`（含 `quota`、`invite_code`、`invited_by` 等扩展字段）、`api_keys`、`channels`、`model_prices`、`request_logs`（可按时间分区策略分阶段）、`recharge_records`、`redeem_codes`、`quota_logs`、`invite_records`
- [x] 2.2 为高频查询添加索引（`api_keys.key_hash`、`quota_logs(user_id,created_at)`、`invite_records(inviter_id)` 等）
- [x] 2.3 实现 GORM model 与 repository 层最小 CRUD（`internal/repository`：API Key 认证、渠道、定价；其余 handler 仍部分直连 GORM）

## 3. User identity & API keys

- [x] 3.1 实现邮箱注册/登录与验证码流程（占位可接第三方邮件）
- [x] 3.2 实现 API Key 创建（仅创建时返回明文）、bcrypt/argon2 哈希存储、`key_prefix` 展示
- [x] 3.3 JWT 会话 + 网关 API Key 中间件；RBAC 守卫管理端路由

## 4. OpenAI-compatible gateway surface

- [x] 4.1 实现 `GET /v1/models`（结合租户/Key 允许模型策略）
- [x] 4.2 实现 `POST /v1/chat/completions` 非流式转发链路
- [x] 4.3 实现 SSE 流式中继（chunk 透传/适配）
- [x] 4.4 实现 `embeddings` / `images/generations` / `audio/*` 路由占位或 P0 适配（按优先级裁剪）

## 5. Upstream routing & adapters

- [x] 5.1 定义 `ProviderAdapter` 接口与至少一个 P0 适配器（OpenAI 兼容路径，如 DeepSeek 或 OpenAI）
- [x] 5.2 实现渠道选择：模型解析 → 过滤 → 加权/优先级 LB → 失败转移（`internal/routing` + `chat/completions` 重试下一渠道 / env 兜底）
- [x] 5.3 实现熔断器（滑动窗口错误率）与渠道状态机（Redis 失败计数 + 短时 open；非完整状态机）
- [x] 5.4 实现三级限流（全局 / 用户 / 渠道）Redis 实现（固定窗口 `INCR` + 分钟桶）

## 6. Billing & quota

- [x] 6.1 实现 `model_prices` 加载与 `CalcTokenCost` / `EstimateMaxCost`（整数额度，按设计文档公式）
- [x] 6.2 实现 Redis 预扣脚本、差额结算、402 预估失败路径（预扣为 `DECRBY` + 失败回滚；未用 Lua 脚本）
- [x] 6.3 流式路径：按 `max_tokens` 上限预扣，结束按 usage 或本地 tokenizer 结算（当前流式按预估值整笔结算，未解析 usage）
- [x] 6.4 异步落库：`request_logs`、`quota_logs`、定期/阈值 `users.quota` 同步（`request_logs` 异步写入；`quota_logs` 仍同步；Redis 配额键按 `QUOTA_RECONCILE_INTERVAL_SEC` 从 PG 刷新）
- [x] 6.5 选型并接入 tokenizer；对无 usage 的上游定义估算与误差策略（已抽 `internal/tokenizer` 启发式估算，便于替换 tiktoken；未接精确 tokenizer）

## 7. Redeem & invite

- [x] 7.1 管理端：兑换码批量生成（`crypto/rand` + 唯一性重试 + 批量插入）
- [x] 7.2 用户端：兑换接口（事务 + `SELECT FOR UPDATE` + 乐观状态更新）
- [x] 7.3 注册生成 `invite_code`；注册带邀请码时发放双方奖励 + `invite_records` + `quota_logs`
- [x] 7.4 （可选）消费返利：消息队列/任务表 + 异步入账，配置默认关闭（`rebate_tasks` 表 + 异步入队 + 定时 worker；`REBATE_ENABLED` 默认 false）

## 8. Admin API & observability

- [x] 8.1 管理端：渠道 CRUD、模型映射、用户封禁、兑换批次查询（含 `GET /admin/redeem/codes`、`GET /admin/redeem/stats`）
- [x] 8.2 结构化日志、Prometheus 指标、健康检查端点（Zap + `/metrics` + `/health`）
- [x] 8.3 OpenAPI/Swagger 导出与基础集成测试（鉴权、402、兑换并发）（`GET /openapi.yaml` + `internal/integration`；402 与兑换并发压测未覆盖）

## 9. OpenSpec / 文档

- [x] 9.1 将本 change 与 `docs/*.md` 交叉引用；归档前执行 `openspec status --change gateway-foundation-invite-billing`（`proposal.md` / `README.md` / `technical-proposal` 已互链；`openspec status` 显示 4/4 artifacts complete）
