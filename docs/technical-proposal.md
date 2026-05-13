# AI API Gateway 中转站 — 技术方案

## 1. 项目概述

构建一个类似 [lanyiapi.com](https://lanyiapi.com/) 的 **AI API 中转站**，为开发者提供统一的 OpenAI 兼容接口，聚合多家 AI 模型提供商（OpenAI、Claude、Gemini、DeepSeek、通义千问等），实现一个 API Key 调用全网主流大模型。

### 核心价值

| 价值点 | 说明 |
|--------|------|
| 统一接口 | 所有模型统一为 OpenAI 兼容格式，降低接入成本 |
| 国内直连 | 优化网络路由，无需科学上网即可调用海外模型 |
| 成本优化 | 渠道聚合 + 按量计费，价格低于官方直调 |
| 高可用 | 多渠道负载均衡 + 自动故障转移，保障服务稳定性 |
| 多租户 | 支持用户/组织隔离、独立计费、配额管控 |

---

## 2. 功能模块

### 2.1 核心功能

```
┌─────────────────────────────────────────────────────────────┐
│                      AI API Gateway                         │
├──────────┬──────────┬──────────┬──────────┬────────────────┤
│  统一API  │  渠道管理  │  用户体系  │  计费系统  │  管理后台     │
│  接口层   │  & 路由   │  & 认证   │  & 配额   │  (Dashboard)  │
└──────────┴──────────┴──────────┴──────────┴────────────────┘
```

#### API 接口层（OpenAI 兼容）
- `POST /v1/chat/completions` — 对话补全（支持流式 SSE）
- `POST /v1/embeddings` — 文本向量化
- `POST /v1/images/generations` — 图像生成
- `POST /v1/audio/transcriptions` — 语音转文本
- `POST /v1/audio/speech` — 文本转语音
- `GET  /v1/models` — 可用模型列表

#### 渠道管理
- 多上游提供商配置（API Key、Base URL、协议类型）
- 模型映射：将统一模型名映射到各渠道的实际模型名
- 渠道优先级 & 权重配置
- 自动健康检查 & 熔断机制
- 故障自动转移到备用渠道

#### 用户体系
- 用户注册/登录（邮箱 + 验证码）
- API Key 生命周期管理（创建/禁用/删除/过期）
- 用户组/组织管理
- RBAC 权限控制（管理员、普通用户、渠道管理员）

#### 计费系统
- 按 Token 用量计费（输入/输出分别计价）
- 按次计费（图像/音频类）
- 余额管理 & 充值（对接支付宝/微信支付）
- 兑换码系统
- 用量统计 & 消费明细

#### 管理后台
- 数据看板（请求量、Token 消耗、活跃用户、收入）
- 渠道管理界面
- 用户管理界面
- 模型配置界面
- 系统设置

### 2.2 高级功能（二期）

- 模型灰度发布 & A/B 测试
- 请求内容审核（敏感词过滤）
- 对话日志记录 & 审计
- Webhook 通知
- 多币种支持
- 邀请返利系统

---

## 3. 技术架构

### 3.1 整体架构图

```
                            ┌─────────────┐
                            │   Nginx     │
                            │  (反向代理)  │
                            └──────┬──────┘
                                   │
                    ┌──────────────┼──────────────┐
                    │              │              │
              ┌─────┴─────┐ ┌─────┴─────┐ ┌─────┴─────┐
              │  Frontend  │ │  API 网关  │ │  Admin API │
              │  Next.js   │ │  /v1/*    │ │  /api/*    │
              │  (Web UI)  │ │  (Go)     │ │  (Go)      │
              └────────────┘ └─────┬─────┘ └─────┬─────┘
                                   │              │
                            ┌──────┴──────────────┘
                            │
                    ┌───────┴────────┐
                    │   Core Engine  │
                    │   (Go 服务)    │
                    ├────────────────┤
                    │ • 路由调度      │
                    │ • 协议适配      │
                    │ • 负载均衡      │
                    │ • 熔断限流      │
                    │ • 计费统计      │
                    └───────┬────────┘
                            │
              ┌─────────────┼─────────────┐
              │             │             │
        ┌─────┴─────┐ ┌────┴────┐ ┌──────┴──────┐
        │ PostgreSQL │ │  Redis  │ │ 上游 AI 厂商  │
        │  (主存储)  │ │ (缓存)  │ │ OpenAI/Claude│
        │            │ │ (限流)  │ │ Gemini/...   │
        └────────────┘ └─────────┘ └─────────────┘
```

### 3.2 技术栈选型

| 层级 | 技术选型 | 选型理由 |
|------|---------|---------|
| **后端框架** | Go + Gin | 高并发性能优异，适合 API 网关场景；生态成熟 |
| **前端框架** | Next.js 15 + React 19 + TypeScript | SSR 支持好，开发体验佳，生态丰富 |
| **UI 组件库** | Shadcn/ui + Tailwind CSS | 现代化设计，高度可定制 |
| **数据库** | PostgreSQL 16 | 可靠的关系型数据库，支持 JSONB，扩展性强 |
| **缓存** | Redis 7 | 高性能缓存 + 分布式限流 + 消息队列 |
| **ORM** | GORM | Go 生态最流行的 ORM，功能完善 |
| **认证** | JWT + API Key | 无状态认证，适合 API 网关场景 |
| **API 文档** | Swagger/OpenAPI 3.0 | 标准化接口文档 |
| **容器化** | Docker + Docker Compose | 标准化部署，环境一致性 |
| **反向代理** | Nginx | 高性能，SSL 终止，负载均衡 |
| **监控** | Prometheus + Grafana | 指标采集 + 可视化告警 |
| **日志** | Zap (结构化日志) | Go 生态高性能日志库 |

### 3.3 为什么选 Go 而不是 Node.js/Python？

| 维度 | Go | Node.js | Python |
|------|-----|---------|--------|
| 并发性能 | ⭐⭐⭐ goroutine 轻量 | ⭐⭐ 单线程事件循环 | ⭐ GIL 限制 |
| 内存占用 | 低 | 中 | 高 |
| SSE 流式支持 | 原生优秀 | 良好 | 一般 |
| 部署简单性 | 单二进制 | 需 runtime | 需 runtime |
| 类型安全 | 编译期检查 | TypeScript 可选 | 弱 |

Go 在 API 网关场景下的高并发、低延迟、低资源占用特性使其成为最佳选择。one-api、new-api 等成功项目均采用 Go 技术栈。

---

## 4. 核心模块设计

### 4.1 适配器模式 — 多厂商协议适配

```go
// 适配器接口定义
type ProviderAdapter interface {
    // 将统一请求转换为厂商特定格式
    ConvertRequest(req *ChatCompletionRequest) (*http.Request, error)
    // 将厂商响应转换为统一格式
    ConvertResponse(resp *http.Response) (*ChatCompletionResponse, error)
    // 处理流式响应
    ConvertStreamResponse(resp *http.Response) (<-chan *ChatCompletionChunk, error)
    // 计算 Token 用量
    CountTokens(req *ChatCompletionRequest, resp *ChatCompletionResponse) (promptTokens, completionTokens int)
}
```

已支持/计划支持的适配器：

| 厂商 | 优先级 | 协议格式 | 支持模型 |
|------|--------|---------|---------|
| OpenAI | P0 | OpenAI 原生 | GPT-4o, GPT-4.1, o3, o4-mini |
| Anthropic | P0 | Anthropic Messages API | Claude 4 Opus/Sonnet |
| Google | P0 | Gemini API | Gemini 2.5 Pro/Flash |
| DeepSeek | P0 | OpenAI 兼容 | DeepSeek-V3, DeepSeek-R1 |
| 阿里云 | P1 | OpenAI 兼容 | Qwen-Max, Qwen-Plus |
| 字节跳动 | P1 | OpenAI 兼容 | Doubao-Pro |
| 百度 | P1 | 百度自定义 | ERNIE-4.0 |
| 智谱 | P1 | OpenAI 兼容 | GLM-4 |
| Mistral | P2 | OpenAI 兼容 | Mistral Large |
| xAI | P2 | OpenAI 兼容 | Grok |

### 4.2 路由调度引擎

```
请求进入
    │
    ▼
┌──────────────┐
│ 解析模型名称   │
└──────┬───────┘
       │
       ▼
┌──────────────┐     ┌──────────────┐
│ 查找可用渠道   │────▶│ 过滤条件：     │
│              │     │ • 模型支持     │
└──────┬───────┘     │ • 渠道状态     │
       │             │ • 配额未超限   │
       ▼             │ • 未被熔断     │
┌──────────────┐     └──────────────┘
│ 负载均衡选择   │
│ • 加权轮询     │
│ • 优先级优先   │
└──────┬───────┘
       │
       ▼
┌──────────────┐     失败
│ 发送请求      │────────┐
└──────┬───────┘        │
       │ 成功            ▼
       ▼          ┌──────────────┐
┌──────────────┐  │ 故障转移到     │
│ 返回响应      │  │ 下一个渠道     │
└──────────────┘  └──────────────┘
```

**负载均衡策略**：
- **加权轮询（Weighted Round Robin）**：按渠道权重分配流量
- **优先级优先**：高优先级渠道优先使用，满载/不可用时降级
- **最少连接**：选择当前活跃连接最少的渠道
- **成本优先**：选择单位 Token 成本最低的渠道

**熔断机制**：
- 滑动窗口统计错误率（默认 10s 窗口，50% 错误率触发）
- 熔断后自动进入半开状态探测恢复
- 支持手动熔断/恢复

### 4.3 限流设计

```
三级限流：

1. 全局限流（保护系统）
   └─ Redis + 令牌桶算法
   └─ 默认: 10000 req/min

2. 用户级限流（公平使用）
   └─ Redis + 滑动窗口
   └─ 默认: 60 req/min (可按用户等级调整)

3. 渠道级限流（保护上游）
   └─ 本地 + 令牌桶
   └─ 根据上游 API 的 Rate Limit 配置
```

### 4.4 流式响应（SSE）处理

```
客户端                    网关                      上游 AI 厂商
  │                       │                           │
  │  POST /v1/chat/...    │                           │
  │  stream: true         │                           │
  │──────────────────────▶│                           │
  │                       │  转换为厂商格式请求          │
  │                       │──────────────────────────▶│
  │                       │                           │
  │                       │  SSE: data: {...chunk1}   │
  │                       │◀──────────────────────────│
  │  SSE: data: {...}     │                           │
  │◀──────────────────────│  (实时协议转换)             │
  │                       │                           │
  │                       │  SSE: data: {...chunk2}   │
  │                       │◀──────────────────────────│
  │  SSE: data: {...}     │                           │
  │◀──────────────────────│                           │
  │                       │                           │
  │                       │  SSE: data: [DONE]        │
  │                       │◀──────────────────────────│
  │  SSE: data: [DONE]    │                           │
  │◀──────────────────────│  (异步记录用量 & 扣费)      │
  │                       │                           │
```

---

## 5. 数据库设计

### 5.1 核心表结构

```sql
-- 用户表
CREATE TABLE users (
    id            BIGSERIAL PRIMARY KEY,
    username      VARCHAR(64) UNIQUE NOT NULL,
    email         VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role          SMALLINT DEFAULT 1,        -- 1:普通用户 10:管理员 100:超级管理员
    status        SMALLINT DEFAULT 1,        -- 1:正常 2:封禁
    balance       DECIMAL(16,6) DEFAULT 0,   -- 账户余额（元）
    quota         BIGINT DEFAULT 0,          -- 剩余配额
    used_quota    BIGINT DEFAULT 0,          -- 已用配额
    created_at    TIMESTAMPTZ DEFAULT NOW(),
    updated_at    TIMESTAMPTZ DEFAULT NOW()
);

-- API Key 表
CREATE TABLE api_keys (
    id            BIGSERIAL PRIMARY KEY,
    user_id       BIGINT REFERENCES users(id),
    name          VARCHAR(128) NOT NULL,
    key_hash      VARCHAR(255) NOT NULL,     -- bcrypt 加密存储
    key_prefix    VARCHAR(16) NOT NULL,      -- 用于展示 "sk-xxxx..."
    status        SMALLINT DEFAULT 1,
    models        TEXT[],                    -- 允许的模型列表，NULL 表示全部
    rate_limit    INT DEFAULT 60,            -- 每分钟请求上限
    expires_at    TIMESTAMPTZ,
    created_at    TIMESTAMPTZ DEFAULT NOW()
);

-- 渠道表（上游 AI 厂商配置）
CREATE TABLE channels (
    id            BIGSERIAL PRIMARY KEY,
    name          VARCHAR(128) NOT NULL,
    provider      VARCHAR(32) NOT NULL,      -- openai/anthropic/google/...
    base_url      VARCHAR(512) NOT NULL,
    api_key       VARCHAR(512) NOT NULL,     -- 加密存储
    models        TEXT[] NOT NULL,            -- 支持的模型列表
    model_mapping JSONB,                     -- 模型名映射 {"gpt-4o": "gpt-4o-2024-08-06"}
    priority      INT DEFAULT 0,
    weight        INT DEFAULT 1,
    status        SMALLINT DEFAULT 1,        -- 1:启用 2:禁用 3:熔断
    max_concurrent INT DEFAULT 100,
    rate_limit    INT DEFAULT 1000,
    created_at    TIMESTAMPTZ DEFAULT NOW(),
    updated_at    TIMESTAMPTZ DEFAULT NOW()
);

-- 请求日志表（用于统计和审计）
CREATE TABLE request_logs (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT,
    api_key_id      BIGINT,
    channel_id      BIGINT,
    model           VARCHAR(64),
    request_method  VARCHAR(16),
    request_path    VARCHAR(256),
    prompt_tokens   INT DEFAULT 0,
    completion_tokens INT DEFAULT 0,
    total_tokens    INT DEFAULT 0,
    cost            DECIMAL(16,8) DEFAULT 0,
    latency_ms      INT,
    status_code     INT,
    error_message   TEXT,
    ip_address      INET,
    created_at      TIMESTAMPTZ DEFAULT NOW()
) PARTITION BY RANGE (created_at);  -- 按时间分区

-- 模型定价表
CREATE TABLE model_prices (
    id              BIGSERIAL PRIMARY KEY,
    model           VARCHAR(64) UNIQUE NOT NULL,
    prompt_price    DECIMAL(16,8) NOT NULL,     -- 每 1K tokens 价格
    completion_price DECIMAL(16,8) NOT NULL,
    unit_price      DECIMAL(16,8),              -- 按次计费价格
    billing_type    SMALLINT DEFAULT 1,          -- 1:按token 2:按次
    currency        VARCHAR(8) DEFAULT 'CNY',
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

-- 充值记录表
CREATE TABLE recharge_records (
    id            BIGSERIAL PRIMARY KEY,
    user_id       BIGINT REFERENCES users(id),
    amount        DECIMAL(16,2) NOT NULL,
    payment_method VARCHAR(32),              -- alipay/wechat/card
    trade_no      VARCHAR(128) UNIQUE,       -- 外部交易号
    status        SMALLINT DEFAULT 0,        -- 0:待支付 1:成功 2:失败
    created_at    TIMESTAMPTZ DEFAULT NOW(),
    completed_at  TIMESTAMPTZ
);

-- 兑换码表
CREATE TABLE redeem_codes (
    id            BIGSERIAL PRIMARY KEY,
    code          VARCHAR(64) UNIQUE NOT NULL,
    quota         BIGINT NOT NULL,
    used_by       BIGINT REFERENCES users(id),
    status        SMALLINT DEFAULT 1,        -- 1:未使用 2:已使用 3:已过期
    expires_at    TIMESTAMPTZ,
    created_at    TIMESTAMPTZ DEFAULT NOW(),
    used_at       TIMESTAMPTZ
);
```

### 5.2 索引策略

```sql
CREATE INDEX idx_api_keys_user_id ON api_keys(user_id);
CREATE INDEX idx_api_keys_key_hash ON api_keys(key_hash);
CREATE INDEX idx_request_logs_user_id ON request_logs(user_id, created_at DESC);
CREATE INDEX idx_request_logs_channel_id ON request_logs(channel_id, created_at DESC);
CREATE INDEX idx_channels_provider ON channels(provider, status);
```

---

## 6. 项目结构

```
ai-api-gateway/
├── cmd/
│   └── server/
│       └── main.go                 # 程序入口
├── internal/
│   ├── config/
│   │   └── config.go               # 配置加载
│   ├── middleware/
│   │   ├── auth.go                  # API Key 认证中间件
│   │   ├── ratelimit.go             # 限流中间件
│   │   ├── cors.go                  # CORS 中间件
│   │   └── logger.go                # 请求日志中间件
│   ├── handler/
│   │   ├── relay/
│   │   │   ├── chat.go              # /v1/chat/completions
│   │   │   ├── embeddings.go        # /v1/embeddings
│   │   │   ├── images.go            # /v1/images/generations
│   │   │   └── models.go            # /v1/models
│   │   └── admin/
│   │       ├── user.go              # 用户管理 API
│   │       ├── channel.go           # 渠道管理 API
│   │       ├── apikey.go            # API Key 管理 API
│   │       ├── dashboard.go         # 数据看板 API
│   │       └── settings.go          # 系统设置 API
│   ├── adapter/
│   │   ├── adapter.go               # 适配器接口定义
│   │   ├── openai/
│   │   │   └── openai.go            # OpenAI 适配器
│   │   ├── anthropic/
│   │   │   └── anthropic.go         # Anthropic 适配器
│   │   ├── google/
│   │   │   └── gemini.go            # Gemini 适配器
│   │   ├── deepseek/
│   │   │   └── deepseek.go          # DeepSeek 适配器
│   │   └── ...
│   ├── router/
│   │   ├── router.go                # 路由调度引擎
│   │   ├── balancer.go              # 负载均衡器
│   │   └── circuit_breaker.go       # 熔断器
│   ├── model/
│   │   ├── user.go                  # 用户模型
│   │   ├── channel.go               # 渠道模型
│   │   ├── api_key.go               # API Key 模型
│   │   └── request_log.go           # 请求日志模型
│   ├── service/
│   │   ├── user_service.go          # 用户业务逻辑
│   │   ├── channel_service.go       # 渠道业务逻辑
│   │   ├── billing_service.go       # 计费业务逻辑
│   │   └── relay_service.go         # 中继业务逻辑
│   └── pkg/
│       ├── jwt/                     # JWT 工具
│       ├── crypto/                  # 加密工具
│       └── sse/                     # SSE 流式处理
├── web/                             # 前端项目 (Next.js)
│   ├── src/
│   │   ├── app/
│   │   │   ├── (auth)/
│   │   │   │   ├── login/
│   │   │   │   └── register/
│   │   │   ├── (dashboard)/
│   │   │   │   ├── dashboard/       # 数据看板
│   │   │   │   ├── channels/        # 渠道管理
│   │   │   │   ├── api-keys/        # API Key 管理
│   │   │   │   ├── usage/           # 用量统计
│   │   │   │   ├── billing/         # 充值 & 账单
│   │   │   │   └── settings/        # 个人设置
│   │   │   ├── (admin)/
│   │   │   │   ├── users/           # 用户管理（管理员）
│   │   │   │   ├── models/          # 模型配置
│   │   │   │   └── system/          # 系统设置
│   │   │   └── layout.tsx
│   │   ├── components/
│   │   │   ├── ui/                  # Shadcn/ui 组件
│   │   │   └── ...
│   │   └── lib/
│   │       ├── api.ts               # API 客户端
│   │       └── utils.ts
│   ├── package.json
│   └── tailwind.config.ts
├── deploy/
│   ├── docker-compose.yml           # 开发环境编排
│   ├── docker-compose.prod.yml      # 生产环境编排
│   ├── Dockerfile                   # Go 后端镜像
│   ├── Dockerfile.web               # 前端镜像
│   └── nginx/
│       └── nginx.conf               # Nginx 配置
├── scripts/
│   ├── migrate.sh                   # 数据库迁移脚本
│   └── seed.sh                      # 种子数据脚本
├── go.mod
├── go.sum
├── .env.example
├── Makefile
└── README.md
```

---

## 7. API 设计

### 7.1 中继 API（对外，OpenAI 兼容）

所有中继 API 使用 `Authorization: Bearer sk-xxxxx` 认证。

#### Chat Completions

```http
POST /v1/chat/completions
Content-Type: application/json
Authorization: Bearer sk-xxxxx

{
  "model": "gpt-4o",
  "messages": [
    {"role": "system", "content": "You are a helpful assistant."},
    {"role": "user", "content": "Hello!"}
  ],
  "temperature": 0.7,
  "stream": true
}
```

**响应（非流式）**：
```json
{
  "id": "chatcmpl-xxx",
  "object": "chat.completion",
  "created": 1234567890,
  "model": "gpt-4o",
  "choices": [{
    "index": 0,
    "message": {"role": "assistant", "content": "Hello! How can I help you?"},
    "finish_reason": "stop"
  }],
  "usage": {
    "prompt_tokens": 20,
    "completion_tokens": 10,
    "total_tokens": 30
  }
}
```

**响应（流式 SSE）**：
```
data: {"id":"chatcmpl-xxx","choices":[{"delta":{"content":"Hello"}}]}

data: {"id":"chatcmpl-xxx","choices":[{"delta":{"content":"!"}}]}

data: [DONE]
```

### 7.2 管理 API（后台）

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/auth/login` | 登录 |
| POST | `/api/auth/register` | 注册 |
| GET | `/api/user/profile` | 获取个人信息 |
| GET | `/api/user/api-keys` | 获取 API Key 列表 |
| POST | `/api/user/api-keys` | 创建 API Key |
| DELETE | `/api/user/api-keys/:id` | 删除 API Key |
| GET | `/api/user/usage` | 获取用量统计 |
| GET | `/api/user/billing` | 获取账单记录 |
| POST | `/api/user/recharge` | 创建充值订单 |
| POST | `/api/user/redeem` | 兑换码兑换 |
| **管理员接口** | | |
| GET | `/api/admin/channels` | 渠道列表 |
| POST | `/api/admin/channels` | 创建渠道 |
| PUT | `/api/admin/channels/:id` | 更新渠道 |
| DELETE | `/api/admin/channels/:id` | 删除渠道 |
| GET | `/api/admin/users` | 用户列表 |
| PUT | `/api/admin/users/:id` | 更新用户 |
| GET | `/api/admin/dashboard` | 看板数据 |
| GET | `/api/admin/models` | 模型定价列表 |
| PUT | `/api/admin/models/:id` | 更新模型定价 |

---

## 8. 部署方案

### 8.1 开发环境

```yaml
# deploy/docker-compose.yml
services:
  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_DB: ai_gateway
      POSTGRES_USER: gateway
      POSTGRES_PASSWORD: gateway_dev
    ports:
      - "5432:5432"
    volumes:
      - pgdata:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

  backend:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgres://gateway:gateway_dev@postgres:5432/ai_gateway?sslmode=disable
      - REDIS_URL=redis://redis:6379
      - JWT_SECRET=dev-secret-change-in-prod
    depends_on:
      - postgres
      - redis

  frontend:
    build:
      context: ./web
      dockerfile: ../deploy/Dockerfile.web
    ports:
      - "3000:3000"
    environment:
      - NEXT_PUBLIC_API_URL=http://localhost:8080

volumes:
  pgdata:
```

### 8.2 生产环境

- **计算**：2C4G 云服务器起步（阿里云/腾讯云 ECS）
- **数据库**：云数据库 RDS PostgreSQL（或自建主从）
- **缓存**：云 Redis（或自建哨兵模式）
- **CDN**：前端静态资源走 CDN
- **SSL**：Let's Encrypt 免费证书 + Nginx 终止
- **监控**：Prometheus + Grafana + AlertManager

### 8.3 扩展方案

```
阶段一（单机）:
  1 台服务器 → Go 后端 + PostgreSQL + Redis + Nginx

阶段二（读写分离）:
  应用服务器 × 2 + PG 主从 + Redis 哨兵

阶段三（水平扩展）:
  应用服务器 × N + PG 读写分离 + Redis Cluster
  + 请求日志写入 ClickHouse（OLAP 分析）
```

---

## 9. 安全设计

| 安全措施 | 说明 |
|---------|------|
| API Key 加密存储 | bcrypt 哈希，不可逆 |
| 传输加密 | 全站 HTTPS |
| 上游 Key 加密 | AES-256-GCM 加密渠道 API Key |
| SQL 注入防护 | ORM 参数化查询 |
| XSS 防护 | Content-Security-Policy + 输入转义 |
| CSRF 防护 | SameSite Cookie + Token 验证 |
| IP 限流 | 未认证请求的 IP 级限流 |
| SSRF 防护 | 渠道 Base URL 白名单校验 |
| 敏感信息脱敏 | 日志中不记录完整 API Key 和请求内容 |

---

## 10. 开发路线图

### Phase 1 — MVP（核心中继 + 基础管理）

- [x] 项目初始化 & 技术选型
- [ ] Go 后端骨架（Gin + GORM + PostgreSQL）
- [ ] OpenAI 适配器 + Chat Completions 中继（含流式）
- [ ] API Key 认证中间件
- [ ] 基础渠道管理（增删改查）
- [ ] Token 用量记录 & 余额扣减
- [ ] 前端：登录/注册 + API Key 管理页面
- [ ] Docker Compose 开发环境
- [ ] 基础单元测试

### Phase 2 — 多模型 + 管理增强

- [ ] Anthropic / Gemini / DeepSeek 适配器
- [ ] 负载均衡 & 故障转移
- [ ] 熔断器实现
- [ ] Redis 分布式限流
- [ ] 前端：管理后台（渠道管理、用户管理、看板）
- [ ] 模型定价管理
- [ ] Embeddings / Images API

### Phase 3 — 商业化

- [ ] 支付宝/微信支付对接
- [ ] 兑换码系统
- [ ] 邀请返利
- [ ] 更多国产模型适配器
- [ ] 内容审核
- [ ] 对话日志 & 审计
- [ ] 监控告警系统

### Phase 4 — 规模化

- [ ] ClickHouse 日志分析
- [ ] 水平扩展支持
- [ ] 多区域部署
- [ ] 模型灰度发布
- [ ] SLA 保障体系

---

## 11. 环境变量配置

```bash
# .env.example

# 服务配置
PORT=8080
GIN_MODE=debug                        # debug/release

# 数据库
DATABASE_URL=postgres://gateway:password@localhost:5432/ai_gateway?sslmode=disable

# Redis
REDIS_URL=redis://localhost:6379

# JWT
JWT_SECRET=your-jwt-secret-here
JWT_EXPIRY=24h

# 加密
CRYPTO_SECRET=your-crypto-secret-32chars

# 管理员初始账号
ADMIN_USERNAME=admin
ADMIN_PASSWORD=admin123456

# 支付（Phase 3）
# ALIPAY_APP_ID=
# ALIPAY_PRIVATE_KEY=
# WECHAT_MCH_ID=
# WECHAT_API_KEY=
```

---

## 12. 技术风险 & 应对

| 风险 | 影响 | 应对策略 |
|------|------|---------|
| 上游 API 不稳定 | 服务可用性下降 | 多渠道冗余 + 自动故障转移 + 熔断保护 |
| 流量突增 | 服务过载 | 三级限流 + 弹性扩容 + 请求排队 |
| API Key 泄露 | 安全风险 | 加密存储 + 定期轮转 + 异常检测 |
| 计费不准确 | 财务损失 | 双重校验（本地计算 + 上游用量对比）|
| 上游价格调整 | 利润率波动 | 定价解耦 + 定期比价 + 成本优先路由 |

---

## 总结

本方案采用 **Go + Next.js + PostgreSQL + Redis** 技术栈，通过**适配器模式**实现多厂商协议统一，通过**路由调度引擎**实现智能负载均衡和故障转移。系统从 MVP 开始迭代，逐步实现完整的 AI API 中转站功能。

核心设计原则：
1. **OpenAI 兼容优先** — 最大程度降低用户迁移成本
2. **高可用** — 多渠道冗余 + 熔断 + 限流
3. **可扩展** — 适配器模式，新增厂商只需实现接口
4. **安全第一** — 加密存储 + 传输加密 + 访问控制
