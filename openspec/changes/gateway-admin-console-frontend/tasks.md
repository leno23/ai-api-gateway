## 1. 工程与契约

- [x] 1.1 初始化前端工程（目录/仓库）与 `NEXT_PUBLIC_GATEWAY_API_URL`（或等价）环境约定
- [x] 1.2 自 `GET /openapi.yaml` 或仓库内嵌 spec 生成 TypeScript 类型 / API 客户端（openapi-typescript 或手写最小 client）
- [x] 1.3 全局请求拦截：注入 JWT、`401` 跳转登录、`403` 统一无权限页

## 2. 认证与壳层

- [x] 2.1 登录页：对接 `POST /auth/login`，错误提示与加载态
- [x] 2.2 Admin 布局：侧栏/顶栏、退出登录（清除 Token）
- [x] 2.3 （可选）注册页或「仅管理员由 DB 开通」说明文案

## 3. 渠道管理

- [x] 3.1 渠道列表：`GET /admin/channels`，表格列与空态
- [x] 3.2 新建/编辑表单：`POST` / `PUT /admin/channels/:id`，`models` 与 `model_mapping` 编辑与校验
- [x] 3.3 删除渠道：确认模态 + `DELETE /admin/channels/:id`

## 4. 兑换运营

- [x] 4.1 批量生成：`POST /admin/redeem/batch`，结果码列表展示与导出
- [x] 4.2 兑换码列表：`GET /admin/redeem/codes` 分页与 `status` 筛选
- [x] 4.3 统计：`GET /admin/redeem/stats` 可视化（卡片/表）

## 5. 用户治理

- [x] 5.1 封禁/解封：`PATCH /admin/users/:id/status`，表单校验与确认
- [ ] 5.2 （二期）用户列表：依赖后端 `GET /admin/users` 等接口另开 change

## 6. P1 租户自助（可选）

- [ ] 6.1 额度/邀请码只读展示（需后端只读接口或复用注册响应字段策略）
- [ ] 6.2 API Key 创建：`POST /user/api-keys`，明文一次性 Modal
- [ ] 6.3 兑换：`POST /user/redeem`

## 7. 质量与交付

- [ ] 7.1 核心路径 E2E 或 RTL 测试（登录、列表、创建渠道）
- [x] 7.2 README：本地联调步骤（网关地址、管理员账号如何设置 `role`）
- [ ] 7.3 归档前 `openspec status --change gateway-admin-console-frontend`
