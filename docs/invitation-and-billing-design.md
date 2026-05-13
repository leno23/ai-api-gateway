# 邀请码 & 额度计费 — 详细设计

## 目录

1. [邀请码生成机制](#1-邀请码生成机制)
2. [邀请码兑换额度流程](#2-邀请码兑换额度流程)
3. [用量计费与额度扣减机制](#3-用量计费与额度扣减机制)

---

## 1. 邀请码生成机制

### 1.1 邀请码的两种类型

系统中存在两种不同用途的码，容易混淆，这里先明确区分：

| 类型 | 用途 | 谁生成 | 示例 |
|------|------|--------|------|
| **兑换码 (Redeem Code)** | 管理员批量生成，用户输入后直接获得额度 | 管理员 | `GIFT-A3X9-K7M2-P5Q1` |
| **邀请码 (Invite Code)** | 用户自带的专属码，邀请新用户注册，双方获得奖励 | 系统为每个用户自动生成 | `INV-8F3K2D` |

### 1.2 兑换码生成

#### 生成算法

```
┌─────────────────────────────────────────────────┐
│              兑换码生成流程                        │
├─────────────────────────────────────────────────┤
│                                                 │
│  1. 管理员发起 → 指定数量 + 面额 + 过期时间         │
│                    │                             │
│  2. 生成随机字符串   ▼                             │
│     ┌──────────────────────────────┐             │
│     │ 前缀 + crypto/rand 随机字节   │             │
│     │ Base32 编码 → 分段格式化      │             │
│     └──────────────────────────────┘             │
│                    │                             │
│  3. 唯一性校验      ▼                             │
│     数据库 UNIQUE 约束 + 冲突重试                  │
│                    │                             │
│  4. 批量写入数据库   ▼                             │
│     状态 = unused, 额度 = N                       │
│                                                 │
└─────────────────────────────────────────────────┘
```

#### Go 代码示例

```go
package service

import (
    "crypto/rand"
    "encoding/base32"
    "fmt"
    "strings"
)

// 兑换码格式: GIFT-XXXX-XXXX-XXXX (16字符, 分4段)
func GenerateRedeemCode(prefix string) (string, error) {
    // 12 字节随机数 → Base32 编码 → 取前12个字符
    b := make([]byte, 12)
    if _, err := rand.Read(b); err != nil {
        return "", fmt.Errorf("failed to generate random bytes: %w", err)
    }

    // Base32 编码，去掉填充字符，取前12位
    encoded := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(b)
    code := strings.ToUpper(encoded[:12])

    // 格式化为 PREFIX-XXXX-XXXX-XXXX
    formatted := fmt.Sprintf("%s-%s-%s-%s",
        prefix,           // "GIFT"
        code[0:4],
        code[4:8],
        code[8:12],
    )
    return formatted, nil
}

// 批量生成兑换码
type CreateRedeemCodesRequest struct {
    Count     int     `json:"count"`      // 生成数量
    Quota     int64   `json:"quota"`      // 每张面额（内部额度单位）
    Prefix    string  `json:"prefix"`     // 前缀，默认 "GIFT"
    ExpiresIn int     `json:"expires_in"` // 过期天数，0 = 永不过期
}

func (s *RedeemService) BatchCreate(req CreateRedeemCodesRequest) ([]RedeemCode, error) {
    if req.Prefix == "" {
        req.Prefix = "GIFT"
    }

    codes := make([]RedeemCode, 0, req.Count)
    for i := 0; i < req.Count; i++ {
        var code string
        var err error

        // 冲突重试，最多 3 次
        for retry := 0; retry < 3; retry++ {
            code, err = GenerateRedeemCode(req.Prefix)
            if err != nil {
                return nil, err
            }
            // 检查数据库唯一性（也可以依赖 DB UNIQUE 约束 + 捕获冲突错误）
            exists, _ := s.repo.ExistsByCode(code)
            if !exists {
                break
            }
        }

        redeemCode := RedeemCode{
            Code:      code,
            Quota:     req.Quota,
            Status:    StatusUnused,
            ExpiresAt: calcExpiry(req.ExpiresIn),
        }
        codes = append(codes, redeemCode)
    }

    // 批量插入数据库
    if err := s.repo.BatchInsert(codes); err != nil {
        return nil, err
    }
    return codes, nil
}
```

#### 为什么用 `crypto/rand` 而不是 `math/rand`？

- `math/rand` 是伪随机，种子可预测 → 兑换码可被猜测/暴力破解
- `crypto/rand` 是密码学安全随机 → 不可预测
- 兑换码等价于"钱"，必须用密码学安全随机数

### 1.3 邀请码生成

邀请码是用户的专属推广标识，在用户注册时自动生成。

```go
// 邀请码格式: INV-XXXXXX (6位字母数字)
// 每个用户只有一个，终身不变
func GenerateInviteCode() (string, error) {
    b := make([]byte, 6)
    if _, err := rand.Read(b); err != nil {
        return "", err
    }
    // 使用可读字符集 (去掉 0/O/I/L 等易混淆字符)
    const charset = "23456789ABCDEFGHJKMNPQRSTUVWXYZ"
    code := make([]byte, 6)
    for i := range code {
        code[i] = charset[int(b[i])%len(charset)]
    }
    return "INV-" + string(code), nil
}
```

#### 邀请码存储

邀请码直接存在 `users` 表中：

```sql
ALTER TABLE users ADD COLUMN invite_code VARCHAR(16) UNIQUE;
ALTER TABLE users ADD COLUMN invited_by  BIGINT REFERENCES users(id);
```

---

## 2. 邀请码兑换额度流程

### 2.1 兑换码兑换流程

```
用户输入兑换码
      │
      ▼
┌──────────────┐     ┌──────────────────────┐
│ 1. 查询兑换码 │────▶│ 验证:                │
│              │     │ • 码是否存在？         │
└──────┬───────┘     │ • 状态是否为"未使用"？  │
       │             │ • 是否已过期？          │
       │ 验证通过     └──────────────────────┘
       ▼                     │ 验证失败
┌──────────────┐             ▼
│ 2. 开启事务   │      返回错误信息
└──────┬───────┘
       │
       ▼
┌──────────────────────────────────┐
│ 3. 在同一事务中执行:               │
│                                  │
│  a. 更新兑换码状态 → "已使用"       │
│     SET status=2, used_by=用户ID  │
│     SET used_at=NOW()            │
│     WHERE status=1 (乐观锁)       │
│                                  │
│  b. 增加用户额度                   │
│     UPDATE users                 │
│     SET quota = quota + 兑换额度   │
│     WHERE id = 用户ID             │
│                                  │
│  c. 记录额度变动日志               │
│     INSERT INTO quota_logs       │
│                                  │
└──────────────────────────────────┘
       │
       ▼
  提交事务, 返回成功
```

#### Go 代码示例

```go
func (s *RedeemService) Redeem(userID int64, code string) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        // 1. 查找并锁定兑换码 (SELECT ... FOR UPDATE)
        var redeemCode RedeemCode
        if err := tx.Where("code = ?", code).
            Set("gorm:query_option", "FOR UPDATE").
            First(&redeemCode).Error; err != nil {
            if errors.Is(err, gorm.ErrRecordNotFound) {
                return ErrCodeNotFound // "兑换码不存在"
            }
            return err
        }

        // 2. 状态校验
        if redeemCode.Status != StatusUnused {
            return ErrCodeAlreadyUsed // "兑换码已被使用"
        }
        if redeemCode.ExpiresAt != nil && redeemCode.ExpiresAt.Before(time.Now()) {
            return ErrCodeExpired // "兑换码已过期"
        }

        // 3. 标记为已使用
        now := time.Now()
        if err := tx.Model(&redeemCode).Updates(map[string]interface{}{
            "status":  StatusUsed,
            "used_by": userID,
            "used_at": now,
        }).Error; err != nil {
            return err
        }

        // 4. 增加用户额度
        if err := tx.Model(&User{}).
            Where("id = ?", userID).
            Update("quota", gorm.Expr("quota + ?", redeemCode.Quota)).
            Error; err != nil {
            return err
        }

        // 5. 记录额度变动日志
        log := QuotaLog{
            UserID:    userID,
            Delta:     redeemCode.Quota,
            Type:      QuotaTypeRedeem,
            Reference: fmt.Sprintf("redeem_code:%d", redeemCode.ID),
            Remark:    fmt.Sprintf("兑换码 %s 充值", maskCode(code)),
            CreatedAt: now,
        }
        return tx.Create(&log).Error
    })
}
```

#### 关键设计点

| 要点 | 说明 |
|------|------|
| **事务隔离** | 兑换码状态更新和用户额度增加必须在同一事务内 |
| **行级锁** | `SELECT FOR UPDATE` 防止并发兑换同一码 |
| **幂等性** | 状态字段做前置校验，已使用的码不可重复兑换 |
| **审计日志** | 每次额度变动都记录 `quota_logs`，便于对账 |

### 2.2 邀请码奖励流程

```
新用户注册时填写邀请码
         │
         ▼
┌──────────────────┐
│ 1. 查找邀请人      │
│ SELECT * FROM     │
│ users WHERE       │
│ invite_code = ?   │
└────────┬─────────┘
         │ 找到邀请人
         ▼
┌──────────────────────────────────────┐
│ 2. 注册事务中一并处理:                  │
│                                      │
│  a. 创建新用户                         │
│     invited_by = 邀请人ID              │
│                                      │
│  b. 新用户获得注册奖励                  │
│     quota += REGISTER_BONUS (如 500)  │
│                                      │
│  c. 邀请人获得邀请奖励                  │
│     quota += INVITE_BONUS (如 1000)   │
│                                      │
│  d. 记录邀请关系日志                    │
│     INSERT INTO invite_records        │
│                                      │
│  e. 双方各记录一条 quota_log            │
│                                      │
└──────────────────────────────────────┘
```

#### 邀请奖励规则配置

```go
// 可在系统设置中动态调整
type InviteRewardConfig struct {
    Enabled          bool  `json:"enabled"`            // 是否开启邀请功能
    RegisterBonus    int64 `json:"register_bonus"`     // 新用户注册奖励额度
    InviterBonus     int64 `json:"inviter_bonus"`      // 邀请人奖励额度
    MaxInviteRewards int   `json:"max_invite_rewards"` // 每人最多邀请奖励次数 (0=不限)
    BonusPercentage  int   `json:"bonus_percentage"`   // 被邀请人消费返利比例 (0-100)
}
```

#### 高级: 消费返利模式

除了一次性奖励，还可支持**持续返利**——被邀请人每次消费时，邀请人获得一定比例的额度奖励：

```
被邀请人发起请求 → 消费 100 额度
                    │
                    ▼
             邀请人获得返利
         100 × 5% = 5 额度（异步）
```

这部分通过异步队列处理，不影响请求响应延迟。

### 2.3 额度变动日志表

```sql
-- 额度变动日志（完整审计）
CREATE TABLE quota_logs (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT NOT NULL REFERENCES users(id),
    delta       BIGINT NOT NULL,               -- 变动量（正=增加，负=消耗）
    balance     BIGINT NOT NULL,               -- 变动后余额（快照）
    type        SMALLINT NOT NULL,             -- 变动类型
    reference   VARCHAR(128),                  -- 关联标识
    remark      VARCHAR(256),                  -- 备注
    created_at  TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_quota_logs_user ON quota_logs(user_id, created_at DESC);

-- type 枚举:
-- 1 = recharge    (充值)
-- 2 = redeem      (兑换码)
-- 3 = consume     (API 消费)
-- 4 = invite      (邀请奖励)
-- 5 = register    (注册奖励)
-- 6 = admin_grant (管理员手动调整)
-- 7 = refund      (退款/补偿)
-- 8 = commission  (邀请返利)
```

### 2.4 邀请记录表

```sql
CREATE TABLE invite_records (
    id             BIGSERIAL PRIMARY KEY,
    inviter_id     BIGINT NOT NULL REFERENCES users(id),
    invitee_id     BIGINT NOT NULL REFERENCES users(id),
    inviter_bonus  BIGINT NOT NULL,             -- 邀请人获得的额度
    invitee_bonus  BIGINT NOT NULL,             -- 被邀请人获得的额度
    status         SMALLINT DEFAULT 1,          -- 1:有效 2:撤销（作弊检测）
    created_at     TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(invitee_id)                          -- 每人只能被邀请一次
);

CREATE INDEX idx_invite_records_inviter ON invite_records(inviter_id);
```

---

## 3. 用量计费与额度扣减机制

### 3.1 额度单位体系

系统内部使用统一的**额度（Quota）**单位，而非直接用人民币，好处是：

1. 避免浮点精度问题
2. 灵活设定"汇率"（1元 = 多少额度），便于促销
3. 与上游成本解耦，独立定价

```
换算关系:
  1 元人民币 = 500,000 额度 (quota)
  即 1 额度 ≈ 0.000002 元

为什么这么大的数字？
  → GPT-4o 输入 $2.5/1M tokens ≈ ¥18/1M tokens
  → 1 token 消耗 ≈ 0.000018 元 ≈ 9 额度
  → 用整数运算避免浮点，精度到单 token 级别
```

### 3.2 模型定价表

```sql
-- 定价示例（额度单位 / 1K tokens）
INSERT INTO model_prices (model, prompt_price, completion_price, billing_type) VALUES
-- OpenAI
('gpt-4o',           9000,    27000,   1),   -- 输入 9/1K, 输出 27/1K
('gpt-4o-mini',      750,     3000,    1),
('gpt-4.1',          10000,   40000,   1),
('o3',               50000,   200000,  1),
('o4-mini',          5500,    22000,   1),
-- Anthropic
('claude-sonnet-4-20250514', 15000, 75000, 1),
('claude-4-opus',    75000,   375000,  1),
-- Google
('gemini-2.5-pro',   6250,    50000,   1),
('gemini-2.5-flash', 375,     6000,    1),
-- DeepSeek
('deepseek-v3',      1350,    5500,    1),   -- ¥0.27/1M in, ¥1.1/1M out
('deepseek-r1',      2750,    11000,   1),
-- 图像生成（按次计费）
('dall-e-3',         0,       0,       2),   -- billing_type=2
('dall-e-3-hd',      0,       0,       2);

-- 按次计费的价格存在 unit_price 字段
UPDATE model_prices SET unit_price = 20000 WHERE model = 'dall-e-3';      -- 每张 0.04元
UPDATE model_prices SET unit_price = 40000 WHERE model = 'dall-e-3-hd';   -- 每张 0.08元
```

### 3.3 计费流程（核心）

#### 非流式请求

```
          ┌────────────────┐
          │ 请求到达网关     │
          └───────┬────────┘
                  │
                  ▼
          ┌────────────────┐
          │ ① 预估消耗       │  根据输入 tokens 预估最大消耗
          │ 余额 ≥ 预估？    │  (prompt_tokens × 价格 + max_tokens × 输出价格)
          └───────┬────────┘
             Yes  │  No → 返回 402 "余额不足"
                  ▼
          ┌────────────────┐
          │ ② 预扣额度       │  Redis: DECRBY user:{id}:quota {预估额度}
          │ (乐观预扣)       │
          └───────┬────────┘
                  │
                  ▼
          ┌────────────────┐
          │ ③ 转发到上游      │
          │ 等待完整响应      │
          └───────┬────────┘
                  │
                  ▼
          ┌────────────────────────────────────────┐
          │ ④ 计算实际消耗                            │
          │                                        │
          │  if billing_type == TOKEN:              │
          │    actual_cost =                        │
          │      prompt_tokens × prompt_price/1000  │
          │    + completion_tokens × compl_price/1000│
          │                                        │
          │  if billing_type == PER_REQUEST:        │
          │    actual_cost = unit_price             │
          │                                        │
          └───────────────┬────────────────────────┘
                          │
                          ▼
          ┌────────────────────────────────────────┐
          │ ⑤ 结算差额                               │
          │                                        │
          │  diff = 预扣额度 - 实际消耗                │
          │  if diff > 0:                          │
          │    Redis: INCRBY user:{id}:quota {diff} │ ← 退还多扣的
          │  if diff < 0:                          │
          │    Redis: DECRBY user:{id}:quota {-diff}│ ← 补扣不足的
          │                                        │
          └───────────────┬────────────────────────┘
                          │
                          ▼
          ┌────────────────────────────────────────┐
          │ ⑥ 异步持久化                              │
          │                                        │
          │  • 更新 DB users.quota (定期同步)         │
          │  • 写入 request_logs                    │
          │  • 写入 quota_logs                      │
          │  • 计算邀请返利（如有）                     │
          │                                        │
          └────────────────────────────────────────┘
```

#### 流式请求（SSE）

流式请求的特殊性在于：**响应期间无法知道最终 token 数**。

```
          ┌────────────────┐
          │ 请求到达网关     │
          └───────┬────────┘
                  │
                  ▼
          ┌────────────────┐
          │ ① 预估消耗       │  prompt_tokens × 价格
          │ + max_tokens    │  + max_tokens × 输出价格（全额预扣）
          │   全额预扣       │
          └───────┬────────┘
                  │
                  ▼
          ┌─────────────────────────────────┐
          │ ② 逐 chunk 转发给客户端            │
          │                                 │
          │  同时累计:                        │
          │  • completion_tokens_count++     │
          │  (通过 tiktoken 本地计算           │
          │   或累加上游返回的 usage 字段)      │
          │                                 │
          └────────────┬────────────────────┘
                       │ 收到 [DONE]
                       ▼
          ┌────────────────────────────────────────┐
          │ ③ 流结束，计算实际消耗                      │
          │                                        │
          │  优先使用上游返回的 usage (如果有)           │
          │  否则使用本地 tiktoken 计算                │
          │                                        │
          │  actual_cost = prompt × price           │
          │              + completion × price       │
          │                                        │
          └───────────────┬────────────────────────┘
                          │
                          ▼
                   ④ 结算差额（同非流式）
                   ⑤ 异步持久化
```

#### Go 代码示例 — 计费核心

```go
// 额度单位常量
const QuotaPerYuan int64 = 500_000 // 1 元 = 500,000 额度

// 计算单次请求的消耗额度
func CalcTokenCost(model string, promptTokens, completionTokens int) (int64, error) {
    price, err := GetModelPrice(model)
    if err != nil {
        return 0, fmt.Errorf("model %s not found in price table", model)
    }

    switch price.BillingType {
    case BillingTypeToken:
        // prompt_price 和 completion_price 单位是 "额度/1K tokens"
        promptCost := int64(promptTokens) * price.PromptPrice / 1000
        completionCost := int64(completionTokens) * price.CompletionPrice / 1000
        return promptCost + completionCost, nil

    case BillingTypePerRequest:
        return price.UnitPrice, nil

    default:
        return 0, fmt.Errorf("unknown billing type: %d", price.BillingType)
    }
}

// 预估最大消耗（用于预扣）
func EstimateMaxCost(model string, promptTokens, maxTokens int) (int64, error) {
    price, err := GetModelPrice(model)
    if err != nil {
        return 0, err
    }

    if price.BillingType == BillingTypePerRequest {
        return price.UnitPrice, nil
    }

    // 如果用户未指定 max_tokens，使用模型默认上限
    if maxTokens == 0 {
        maxTokens = GetModelDefaultMaxTokens(model) // 如 4096
    }
    // 但预扣不要太激进，设一个合理上限
    if maxTokens > 8192 {
        maxTokens = 8192
    }

    promptCost := int64(promptTokens) * price.PromptPrice / 1000
    completionCost := int64(maxTokens) * price.CompletionPrice / 1000
    return promptCost + completionCost, nil
}
```

### 3.4 Redis 额度缓存 & DB 同步

为了高性能，额度的实时扣减在 Redis 中完成，定期同步到 PostgreSQL：

```
┌─────────────┐                    ┌─────────────┐
│   Redis     │                    │ PostgreSQL  │
│             │    定期同步          │             │
│ user:1:quota│ ──────────────────▶ │ users.quota │
│   = 48350   │    (每30秒 或       │  = 48350    │
│             │     变动超100次)     │             │
└─────────────┘                    └─────────────┘

写入路径: 请求计费 → Redis DECRBY → 快速返回
同步路径: 后台 goroutine → 批量 UPDATE → DB
恢复路径: 服务重启 → 从 DB 加载 → 写入 Redis
```

```go
// Redis 额度操作（原子性）
func (s *QuotaService) PreDeduct(ctx context.Context, userID, amount int64) error {
    key := fmt.Sprintf("user:%d:quota", userID)

    // 原子操作：扣减并检查余额
    script := redis.NewScript(`
        local balance = redis.call('GET', KEYS[1])
        if balance == false then
            return -1  -- key 不存在，需要从 DB 加载
        end
        balance = tonumber(balance)
        local amount = tonumber(ARGV[1])
        if balance < amount then
            return -2  -- 余额不足
        end
        redis.call('DECRBY', KEYS[1], amount)
        return balance - amount
    `)

    result, err := script.Run(ctx, s.redis, []string{key}, amount).Int64()
    if err != nil {
        return err
    }

    switch result {
    case -1:
        // Redis 中无缓存，从 DB 加载
        if err := s.loadFromDB(ctx, userID); err != nil {
            return err
        }
        return s.PreDeduct(ctx, userID, amount) // 重试
    case -2:
        return ErrInsufficientQuota
    default:
        return nil
    }
}

// 退还多扣的额度
func (s *QuotaService) Refund(ctx context.Context, userID, amount int64) error {
    key := fmt.Sprintf("user:%d:quota", userID)
    return s.redis.IncrBy(ctx, key, amount).Err()
}
```

### 3.5 Token 计算方式

Token 数量的获取有三种来源，按优先级排列：

| 优先级 | 来源 | 准确性 | 说明 |
|--------|------|--------|------|
| 1 (最高) | 上游 `usage` 字段 | 100% 准确 | 上游返回的官方计数 |
| 2 | tiktoken 本地计算 | 99%+ 准确 | 使用与模型匹配的 tokenizer |
| 3 | 字符数估算 | ~80% 准确 | 英文 ÷ 4，中文 ÷ 2（仅兜底） |

```go
// Token 计算器
type TokenCounter struct {
    // tiktoken 编码器缓存
    encoders map[string]*tiktoken.Encoding
}

func (tc *TokenCounter) Count(model string, text string) int {
    // 根据模型选择对应的编码器
    enc := tc.getEncoder(model)
    if enc != nil {
        tokens := enc.Encode(text, nil, nil)
        return len(tokens)
    }

    // 兜底估算
    return estimateTokens(text)
}

func estimateTokens(text string) int {
    // 中文字符数 × 0.6 + 英文单词数 × 1.3
    // 简化为：UTF-8 字节数 / 3
    return len([]byte(text)) / 3
}
```

### 3.6 完整请求计费时序图

```
Client          Gateway           Redis            DB             Upstream
  │                │                │               │                │
  │  POST /v1/...  │                │               │                │
  │───────────────▶│                │               │                │
  │                │                │               │                │
  │                │  计算 prompt    │               │                │
  │                │  tokens 数量    │               │                │
  │                │                │               │                │
  │                │  预估最大消耗    │               │                │
  │                │  = 2,500 额度   │               │                │
  │                │                │               │                │
  │                │  DECRBY 2500   │               │                │
  │                │───────────────▶│               │                │
  │                │  OK (余额45850) │               │                │
  │                │◀───────────────│               │                │
  │                │                │               │                │
  │                │  转发请求        │               │                │
  │                │────────────────────────────────────────────────▶│
  │                │                │               │                │
  │                │  返回响应 (usage: prompt=120, completion=85)     │
  │                │◀────────────────────────────────────────────────│
  │                │                │               │                │
  │                │  实际消耗 =      │               │                │
  │                │  120×9/1000     │               │                │
  │                │  + 85×27/1000   │               │                │
  │                │  = 1 + 2 = 3   │               │                │
  │                │  (≈ 3 额度)     │               │                │
  │                │                │               │                │
  │                │  退还 2500-3    │               │                │
  │                │  = 2497        │               │                │
  │                │  INCRBY 2497   │               │                │
  │                │───────────────▶│               │                │
  │                │                │               │                │
  │  响应返回       │                │               │                │
  │◀───────────────│                │               │                │
  │                │                │               │                │
  │                │  异步写入日志     │               │                │
  │                │  ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─▶│                │
  │                │                │               │                │
```

### 3.7 计费边界情况处理

| 场景 | 处理策略 |
|------|---------|
| **请求失败（上游报错）** | 全额退还预扣额度，不计费 |
| **请求超时** | 退还预扣额度；如果流式已部分返回，按已消耗 token 计费 |
| **上游未返回 usage** | 使用本地 tiktoken 计算 |
| **余额刚好够预估但不够实际** | 允许小额透支（< 1000 额度），下次请求前补扣 |
| **Redis 宕机** | 降级为直接操作 DB（性能下降但不中断服务） |
| **并发请求导致余额负数** | Lua 脚本保证原子性；允许极小额负数，下次充值时扣回 |
| **上游返回的 usage 明显异常** | 对比本地计算结果，偏差 > 50% 则记录告警，使用本地计算值 |

### 3.8 用户侧查看消费

```json
// GET /api/user/usage?start=2026-05-01&end=2026-05-10
{
  "total_requests": 1523,
  "total_tokens": {
    "prompt": 2450000,
    "completion": 890000,
    "total": 3340000
  },
  "total_cost_quota": 156800,
  "total_cost_yuan": "0.31",
  "by_model": [
    {
      "model": "gpt-4o-mini",
      "requests": 1200,
      "prompt_tokens": 1800000,
      "completion_tokens": 650000,
      "cost_quota": 42750
    },
    {
      "model": "claude-sonnet-4-20250514",
      "requests": 323,
      "prompt_tokens": 650000,
      "completion_tokens": 240000,
      "cost_quota": 114050
    }
  ],
  "daily": [
    {"date": "2026-05-01", "requests": 180, "cost_quota": 18500},
    {"date": "2026-05-02", "requests": 210, "cost_quota": 21300}
  ]
}
```

---

## 4. 总结 — 三个核心流程的关系

```
                    ┌─────────────────────────────┐
                    │         额度池 (Quota)        │
                    │     用户的"虚拟钱包"            │
                    └──────────┬──────────────────┘
                               │
           ┌───────────────────┼───────────────────┐
           │                   │                   │
     ┌─────┴─────┐      ┌─────┴─────┐      ┌─────┴─────┐
     │   入口 1    │      │   入口 2    │      │   入口 3    │
     │  兑换码充值  │      │  邀请奖励    │      │  支付充值   │
     │  管理员生成  │      │  注册+邀请   │      │  支付宝等   │
     │  → 增加额度  │      │  → 增加额度  │      │  → 增加额度  │
     └────────────┘      └────────────┘      └────────────┘

                               │
                         ┌─────┴─────┐
                         │   出口     │
                         │  API 消费   │
                         │  按量扣减   │
                         │  预扣→结算  │
                         └────────────┘

                               │
                         ┌─────┴─────┐
                         │  审计日志   │
                         │ quota_logs │
                         │ 每笔进出   │
                         │ 完整记录   │
                         └────────────┘
```

**设计原则**：
1. **所有额度变动走同一套 quota_logs** → 任何入口/出口都有完整审计链
2. **Redis 做热路径，DB 做持久化** → 高并发下不成为瓶颈
3. **预扣 + 结算两步走** → 不会出现消费超过余额的情况
4. **整数运算** → 避免浮点精度问题
