## Phase 1 — 基础与用户（portal-public-site + shell）

- [x] 1.1 初始化 `portal/` Next.js 14 + Semi Design + Tailwind + TypeScript
- [x] 1.2 实现顶栏全局导航与营销布局壳层
- [x] 1.3 登录页 `/login`：主按钮文案「继续」；对接 `POST /auth/login`；支持 expired 查询参数提示
- [x] 1.4 注册页 `/register`：用户名/邮箱/密码；`?aff=` 写入 `invite_code`；对接 `POST /auth/register`
- [x] 1.5 JWT 存储与 `/console/*` 路由守卫；401 → `/login?expired=true`
- [x] 1.6 控制台侧栏 Semi Navigation（菜单项与参考站一致）
- [x] 1.7 （后端）`GET /api/user/self` 或扩展 login 响应含 balance、group、display_name

## Phase 2 — 模型与定价（portal-model-pricing）

- [x] 2.1 （后端）`token_groups` 表与种子数据；用户/令牌关联分组倍率
- [x] 2.2 （后端）`GET /api/models`：供应商、分组、计费类型、端点类型筛选与分页
- [x] 2.3 模型广场 `/pricing`：左侧筛选 + 顶栏搜索/视图切换（卡片/表格 M/L）/价格与倍率开关
- [x] 2.4 （后端）价格计算服务对外暴露分项单价（input/output/cache × multiplier）

## Phase 3 — 令牌与网关（portal-token-management）

- [x] 3.1 （后端）`ApiToken` 扩展字段迁移；`GET/POST/PUT/DELETE /api/tokens`
- [x] 3.2 令牌管理页：表格 CRUD、掩码、复制/二维码、批量操作、启用/禁用
- [x] 3.3 网关鉴权：校验额度/IP/模型白名单（扩展现有 middleware）
- [x] 3.4 「聊天」快捷入口：跳转操练场并预填令牌

## Phase 4 — 控制台数据（dashboard + logs）

- [x] 4.1 （后端）`GET /api/dashboard/stats`、`GET /api/dashboard/charts?range=`
- [x] 4.2 数据看板：四指标卡 + 模型分析 Tab 图表 + API 节点卡（主站/香港/美区，复制/测速）
- [x] 4.3 （后端）`GET /api/logs/usage` 分页筛选
- [x] 4.4 使用日志页：时间/令牌/模型/Request ID/分组筛选；列设置
- [x] 4.5 （后端）`task_logs` 表与 `GET /api/logs/tasks`
- [x] 4.6 任务日志页：空状态与表格列

## Phase 5 — 钱包与运营（wallet + announcements）

- [x] 5.1 钱包页：余额/历史消耗/请求数；兑换码表单对接 `POST /user/redeem`
- [x] 5.2 邀请链接 `register?aff={code}`、待结算收益与划转到余额（后端事务）
- [x] 5.3 在线充值 UI（`RECHARGE_ENABLED` mock 或关闭时展示引导文案）
- [x] 5.4 （后端）公告 CRUD + `GET /api/announcements`；首页/看板 Modal + localStorage 今日关闭

## Phase 6 — 操练场（portal-playground）

- [x] 6.1 （后端）`POST /api/playground/chat` SSE；日志令牌名 `playground-default` 策略
- [x] 6.2 操练场三栏 UI：参数滑条、多轮对话、重新生成/复制/编辑/删除、调试信息
- [x] 6.3 配置导入导出 JSON

## Phase 7 — 多节点、文档、i18n

- [x] 7.1 节点配置环境变量；看板测速 `GET /api/nodes/ping`（可选）
- [x] 7.2 关于页 `/about`；首页营销区块对照截图
- [x] 7.3 顶栏语言切换（简体中文默认，英文预留）
- [x] 7.4 文档入口（外链或占位 `/docs`）

## Phase 8 — 质量与 OpenSpec

- [x] 8.1 更新 `backend/internal/openapi/spec.yaml` 与租户面契约测试
- [x] 8.2 核心路径 E2E：注册 → 登录 → 创建令牌 → 操练场一条消息
- [x] 8.3 `portal/README.md` 本地联调说明
- [x] 8.4 归档前 `npx @fission-ai/openspec@latest status --change gateway-tenant-portal-platform`
