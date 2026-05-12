## Context

- 目标系统为 **AI API Gateway**：统一 OpenAI 兼容面，聚合多上游（OpenAI、Anthropic、Gemini、DeepSeek、通义等），多租户配额与计费。
- 已有两份设计文档：`technical-proposal.md.md`（架构与模块）、`invitation-and-billing-design.md.md`（兑换/邀请与额度扣减）。
- 技术栈约定：**Go + Gin** 网关核心，**PostgreSQL 16** 持久化，**Redis 7** 缓存/限流/额度热路径，**GORM**，**JWT + API Key** 认证（详细见技术方案文档）。

## Goals / Non-Goals

**Goals:**

- 网关对外提供稳定 **OpenAI 兼容** 契约，支持 **SSE 流式** 与中继计费闭环。
- **渠道层**可配置、可切换、可熔断，模型名可映射到上游真实模型。
- **额度（Quota）** 使用内部整数单位；支持 **按 Token** 与 **按次** 定价；请求路径上 **预估 → 预扣 → 按实际结算**，不足时 **HTTP 402**。
- **兑换码**（管理员批量）与 **邀请码**（每用户唯一）及奖励写入 **事务 + 审计日志**；兑换使用 **行级锁 / 乐观更新** 防并发重复使用。
- **Redis 配额缓存 + 定期/阈值同步 PostgreSQL**，重启可从 DB 回填 Redis。

**Non-Goals:**

- 本期不强制实现 **支付宝/微信自动充值** 全链路（仅保留 `recharge_records` 等扩展点与手工/后续对接）。
- 不实现完整 **管理后台 UI**（仅 API 与领域行为；Dashboard 另列 change）。
- **内容审核、对话持久化审计、Webhook、多币种** 列为二期（技术方案已标注）。

## Decisions

| 决策 | 选择 | 理由 |
|------|------|------|
| 对外协议 | OpenAI 兼容 `/v1/*` | 降低客户端接入成本，与行业工具链一致（技术方案 2.1）。 |
| 网关语言 | Go | 高并发、低延迟、SSE 友好；与 one-api/new-api 类项目一致（技术方案 3.3）。 |
| 上游集成 | 适配器接口 `ProviderAdapter` | 隔离各厂商请求/响应/流式/计数字段差异（技术方案 4.1）。 |
| 路由 | 模型解析 → 渠道过滤 → 加权/优先级 LB → 失败转移 | 满足高可用与成本策略（技术方案 4.2）。 |
| 限流 | 全局 / 用户 / 渠道 三级 | 保护系统与上游（技术方案 4.3）。 |
| 额度热路径 | Redis 原子脚本预扣 + 差额结算 | 高性能；与设计文档 3.3、3.4 一致。 |
| 兑换码随机源 | `crypto/rand` + Base32/可读字符集 | 防猜测；与设计文档 1.2、1.3 一致。 |
| 兑换一致性 | 单事务：`SELECT FOR UPDATE` + 更新码状态 + `quota` 增加 + `quota_logs` | 防止重复兑换与额度丢失（设计文档 2.1）。 |
| 流式计费 | 预估按 `max_tokens` 上限预扣，结束按 usage 或本地 tokenizer 校准 | 上游未结束前无法得知 completion tokens（设计文档 3.3）。 |
| 邀请返利 | 异步队列处理 | 不阻塞请求主路径（设计文档 2.2 高级能力）。 |

## Risks / Trade-offs

| 风险 | 缓解 |
|------|------|
| Redis 与 DB 额度漂移 | 定期同步 + 大变动触发同步；启动从 DB 灌 Redis；关键操作可回源校验。 |
| 流式场景预扣过多/过少 | 结束按实际 usage 差额回补或补扣；封顶 `max_tokens` 避免无限预扣（设计文档伪代码上限）。 |
| 上游不返回 usage | 本地 tiktoken/等价库估算，需在 spec 中声明误差与上限策略。 |
| 兑换码暴力尝试 | 限流 + 审计 + 可选锁定策略（管理配置）。 |

## Migration Plan

1. 先落 **DDL**（`users` 扩展字段、`redeem_codes`、`quota_logs`、`invite_records`、`model_prices` 等与方案对齐）。
2. 部署 **只读定价数据** 与 **空渠道** → 灰度开启路由。
3. 开启 **计费与预扣**（可先 shadow 模式记录不实扣，再切换实扣）。
4. 开放 **兑换/邀请** 接口前完成事务与并发压测。

## Open Questions

- 各模型 **默认 max_tokens** 与 **预扣封顶** 的最终数值表是否由运营配置表驱动？
- **Tokenizer** 选型（按模型分库）与 **非英文** 内容的估算误差接受范围？
- **邀请返利** 一期是否启用，还是仅配置项默认关闭？

## References

- `docs/technical-proposal.md.md`
- `docs/invitation-and-billing-design.md.md`
- `openspec/changes/gateway-foundation-invite-billing/proposal.md`
