# Asgard 任务跟踪（TASKS）

> 本文件是项目内的**单一任务跟踪表**。所有正在做、计划做、被卡住的工程改动都登记在这里。
>
> - **开始一项工作前**：先查本文件，避免重复或冲突。
> - **完成或状态变化时**：直接编辑本文件，把对应 `- [ ]` 改成 `- [x]`、`- [/]` 或 `- [-]`，并在 §4「状态变更日志」追加一行。
> - **AGENTS.md** 是项目规范（架构、约定、构建、坑）；本文件是**动态工作清单**。两份文档互为补充。

## 目录

1. [约定与图例](#1-约定与图例)
2. [当前重点：前后端分离](#2-当前重点前后端分离)
3. [任务详情](#3-任务详情)
   - [Phase 0：规划与决策](#phase-0规划与决策)
   - [Phase 1：后端 API 化](#phase-1后端-api-化)
   - [Phase 2：前端项目初始化](#phase-2前端项目初始化)
   - [Phase 3：前端页面迁移](#phase-3前端页面迁移)
   - [Phase 4：打包与部署](#phase-4打包与部署)
   - [Phase 5：验证与下线](#phase-5验证与下线)
4. [状态变更日志](#4-状态变更日志)

## 1. 约定与图例

**状态标记**（GitHub Flavored Markdown 复选框）：

| 标记 | 含义 | 何时使用 |
| --- | --- | --- |
| `- [ ]` | 待办 | 已识别但未开始 |
| `- [/]` | 进行中 | 有人正在做，建议先完成再开新项 |
| `- [x]` | 已完成 | 当前仓库 main 分支已包含 |
| `- [-]` | 阻塞 | 缺决策 / 缺资源 / 依赖未到位，需在任务说明里写阻塞原因 |

**任务编号**：`T-<phase>-<seq>`，例如 `T-101`、`T-203`。Phase 编号固定，新增任务顺延。

**关联代码**：每个任务在括号里写明主要触达的文件或目录，方便改动时定位。

**更新时机**：

- 开始一项工作 → 把对应行改成 `- [/]`，并在 §4 追加 `YYYY-MM-DD 启动 T-xxx`。
- 完成一项工作 → 改成 `- [x]`，追加 `YYYY-MM-DD 完成 T-xxx`。
- 阻塞时 → 改成 `- [-]`，在任务说明里写阻塞原因，并在 §4 追加。
- 任务被分解 → 旧任务保留作为父任务，新增子任务用 `<parent>.1`、`<parent>.2` 后缀。
- 新识别出的任务 → 加在所属 Phase 末尾，编号顺延。

**评审门槛**：单个任务预计超过 1 天工作量时，必须先拆成子任务再开工。

**示例**（任务从待办到完成的一次生命周期）：

```
- [ ] T-101 在 web/server.go 注册新的 api 路由 group
- [/] T-101 在 web/server.go 注册新的 api 路由 group
- [x] T-101 在 web/server.go 注册新的 api 路由 group
```

对应 §4 追加：

```
2026-06-02 启动  T-101
2026-06-02 完成  T-101  server.go 已加 /api/v1 group；旧 HTML 路由未动
```

## 2. 当前重点：前后端分离

Asgard 的 web 层（`web/`）当前是 **Gin + goview + Tabler + jQuery** 的服务端渲染架构（见 [AGENTS.md §6 Web 层](/Users/dengjialong/git/Asgard/AGENTS.md)）。本次任务目标是把 `web/` 拆成独立的「纯 JSON API 后端 + 独立前端工程」，结束条件：

1. 浏览器请求不再经过 goview 模板，所有 HTML 由前端构建产物提供。
2. 后端 `web/controllers/` 不再依赖 `gin.H` + `ctx.HTML`；只返 JSON。
3. `web/views/`、`web/assets/`、`goview` 依赖可以安全删除。
4. 部署包可以独立发布前端静态资源（CDN / Nginx / 静态站点都行）。
5. 功能等价：现有 11 个路由 group（user / agent / group / app / job / timing / monitor / archive / out_log / err_log / exception / operation）全部有前端对应页面。

执行原则：**先增量、后下线**。Phase 1 期间 HTML 路由与 API 路由并存，确认前端稳定后（Phase 5）再删 HTML。

## 3. 任务详情

### Phase 0：规划与决策

> 在动手前需要先定下来的事情。每条决策落地后，状态改成 `[x]` 并在备注里写明最终选择。

- [x] **T-001** 选定前端框架 → **Vue 3**（Composition API + `<script setup>`，中文社区生态成熟，TS 友好，配合 Element Plus / Vite / Pinia 一站式，无 SSR 需求因此不上 Nuxt）
- [x] **T-002** 选定构建工具 → **Vite 5**（Vue 3 官方推荐、HMR 极快、生产 Rollup 打包成熟、对 pnpm/TypeScript 一等公民支持；Turbopack 仍 preview，Webpack 5 已被 Vite 在新项目里取代）
- [x] **T-003** 选定 UI 组件库 → **Element Plus**（Vue 3 配套最广的桌面端组件库，中文文档完整、Form/Table/Tree/Pagination/Dialog 等直接覆盖 Asgard 后台所需场景）
- [x] **T-004** 选定 HTTP 客户端 → **axios**（拦截器灵活，request/response 钩子成熟，社区文档最多；统一封装在 src/api/http.ts，baseURL 走 vite proxy；拦截器约定：401 跳 /login、code !== 200 走 ElMessage 错误提示、loading 计数管理）
- [x] **T-005** 选定状态管理 → **Pinia**（Vue 3 官方推荐，TS 推断天然友好，store 即组合式 API，devtools 支持完整；Zustand/Redux 是 React 体系；Vuex 4 在 Vue 3 时代已被 Pinia 取代）
- [x] **T-006** TypeScript vs JavaScript → **TypeScript**（严格模式，tsconfig 开启 strict + noUncheckedIndexedAccess + exactOptionalPropertyTypes；后端 API 字段类型契约能由 OpenAPI 后续生成，TS 编译期拦截明显多于运行时排查）
- [x] **T-007** 前端工程位置 → **仓库内子目录 `web-admin/`**（与 Asgard 后端同 repo，统一 review / 版本号 / 部署脚本；前端产物由独立 Dockerfile 或 Asgard 打包脚本单独构建发布；现阶段不引入 monorepo 工具，web-admin/ 保持独立 package.json）
- [x] **T-008** 后端 API 路由前缀 → **`/api/v1`**（版本化路径，约定：所有 JSON 接口挂在 `server.Group("/api/v1")` 下；未来 `/api/v2` 留作破坏性升级；现有 HTML 路由继续挂在根路径，与 API 并存）
- [x] **T-009** 鉴权协议 → **过渡期双轨：保留 DES cookie + 新增 JWT Bearer**（web/middlewares/api_auth.go 同时接受 `Authorization: Bearer <jwt>` 与现有 DES cookie；短期保证旧 HTML 路由继续工作；后端用 golang-jwt/jwt/v5 签发 2h 过期 token；前端 Pinia store 用 localStorage 持久化；T-102 是这一决策的落地锚点）
- [x] **T-010** 实时数据方案 → **SSE（Server-Sent Events）单向 + 可选 WebSocket 双向**（实时监控图表、实时日志走 SSE：HTTP 友好、axios/fetch 原生支持 EventSource、调试简单、断线重连内置；只有需要双向推送时才用 WebSocket：T-117/T-118 走 SSE，T-310/T-312 前端用 EventSource 订阅）

### Phase 1：后端 API 化

> 目标：让前端能拿到与现有 HTML 路由等价的数据。**不删任何 HTML 路由**，先做增量。

- [x] **T-101** 在 web/router.go 注册 api 路由 group（`server.Group("/api/v1")`），与现有 HTML 路由并存；JSON 子路由拆到 web/routers/api_router.go（占位 /health）
- [/] **T-102** 实现 `APIAuth` 中间件（新建 `web/middlewares/api_auth.go`），同时支持 `Authorization: Bearer <token>` 和现有 cookie（保证过渡期兼容）
- [/] **T-103** 引入 CORS 中间件（`common_middlewares.Cors` 或自写），允许前端开发域名访问
- [/] **T-104** 统一 JSON 响应格式：复用 `{code, message, data}` 约定（`web/utils/respose.go`），新增 `APIPage` 工具返 `{code, message, data:{list, total, page, page_size}}`
- [/] **T-105** 鉴权接口：`POST /api/v1/auth/login`、`GET /api/v1/auth/info`、`POST /api/v1/auth/logout`、`POST /api/v1/auth/change_password`
- [/] **T-106** 用户管理 API（`controllers/user.go`）：`GET /users`、`GET /users/:id`、`POST /users`、`PUT /users/:id`、`POST /users/:id/forbidden`、`POST /users/:id/reset_password`
- [/] **T-107** 实例管理 API（`controllers/agent.go`）：`GET /agents`、`PUT /agents/:id`、`POST /agents/:id/forbidden`
- [/] **T-108** 分组管理 API（`controllers/group.go`）：`GET /groups`、`POST /groups`、`PUT /groups/:id`、`DELETE /groups/:id`
- [/] **T-109** 应用管理 API（`controllers/app.go`）：list/show/create/update/copy + start/restart/pause/delete + batch 四个批量接口
- [/] **T-110** 计划任务 API（`controllers/job.go`）：同应用
- [/] **T-111** 定时任务 API（`controllers/timing.go`）：同应用
- [/] **T-112** 监控 API（`controllers/monitor.go`）：`GET /monitor/agent?agent_id=`、`GET /monitor/app?app_id=`、`GET /monitor/job?job_id=`、`GET /monitor/timing?timing_id=`
- [/] **T-113** 归档 API（`controllers/archive.go`）：`GET /archives/app?app_id=`、`GET /archives/job?job_id=`、`GET /archives/timing?timing_id=`
- [/] **T-114** 日志 API（`controllers/log.go`）：`GET /out_logs/app?app_id=&page=`、`GET /out_logs/app/data?app_id=`；err_log 同理
- [/] **T-115** 异常列表 API（`controllers/exception.go`）：`GET /exceptions?page=&type=`
- [/] **T-116** 操作日志 API（`controllers/operation.go`）：`GET /operations?page=&user_id=&type=`
- [/] **T-117** （可选）实时日志流：`GET /out_logs/app/stream?app_id=` 用 SSE；`err_log` 同理
- [/] **T-118** （可选）实时监控数据：`GET /monitor/agent/stream?agent_id=` 用 SSE
- [/] **T-119** 编写 API 文档（`doc/API.md` 或 OpenAPI/Swagger），含鉴权方式、错误码、分页约定

### Phase 2：前端项目初始化

- [/] **T-201** 创建前端项目（`pnpm create vite` 或 `npm create vite@latest`）
- [/] **T-202** 配置 TypeScript（`tsconfig.json` 严格模式）+ ESLint + Prettier + EditorConfig
- [/] **T-203** 配置 Vite 代理：开发期把 `/api/*` 转发到 `http://localhost:12345`
- [/] **T-204** 引入路由库（vue-router 4 / react-router 6 / sveltekit 内置），按 `web/router.go` 11 个 group 建路由表
- [/] **T-205** 引入 HTTP client（axios 推荐），实现拦截器：401 跳登录、统一 `alert`/toast 错误处理
- [/] **T-206** 配置全局状态：用户信息、token、菜单权限
- [/] **T-207** 引入 UI 组件库，按 Phase 0 选定的方案接入
- [/] **T-208** 引入图表库（首页监控）：ApexCharts / ECharts / Chart.js
- [/] **T-209** 引入日志查看组件（按需）：xterm.js / 自实现
- [/] **T-210** 引入时间选择组件（Timing Add/Edit 页需要）
- [/] **T-211** 配单元测试（Vitest / Jest）+ E2E（Playwright / Cypress）

### Phase 3：前端页面迁移

> 与 `web/controllers/` 一一对应。**所有现有 HTML 页面必须在功能上完整复刻**（包括 `assets/js/asgard.js` 里 `ActionBatch / runAction / runBatchAction` 那套交互约定）。

- [/] **T-301** 登录页 + 鉴权守卫（路由 meta + 401 跳转）
- [/] **T-302** 主框架：侧边栏（11 个一级菜单）+ 顶栏（用户信息、退出）+ `<router-view>`
- [/] **T-303** 首页仪表盘（`controllers/index.go` 的 `Index`）：四类计数 + 图表
- [/] **T-304** 用户管理（`controllers/user.go`）：list/info/setting/change_password/add/create/edit/update/forbidden/reset_password
- [/] **T-305** 实例管理（`controllers/agent.go`）：list/edit/update/forbidden
- [/] **T-306** 分组管理（`controllers/group.go`）：list/add/create/edit/update/delete
- [/] **T-307** 应用管理（`controllers/app.go`）：list/show/add/create/edit/update/copy + start/restart/pause/delete + batch 四个
- [/] **T-308** 计划任务管理（`controllers/job.go`）：同应用
- [/] **T-309** 定时任务管理（`controllers/timing.go`）：同应用（多一个时间字段）
- [/] **T-310** 监控视图：4 个二级页（agent/app/job/timing）
- [/] **T-311** 归档列表：3 个二级页
- [/] **T-312** 日志查看：out_log + err_log 各 3 个实体页 + 数据接口
- [/] **T-313** 异常记录
- [/] **T-314** 操作日志
- [/] **T-315** 错误页：no_login / auth_fail / admin_only / forbidden / err

### Phase 4：打包与部署

- [/] **T-401** 前端 Dockerfile（基于 nginx 托管 dist 产物）
- [/] **T-402** 反向代理配置示例（Nginx）：`/api/*` 转发到后端 `Asgard web`，其他走前端静态
- [/] **T-403** 调整 `.goreleaser.yml`：不再把 `web/views/` 和 `web/assets/` 打入发布包
- [/] **T-404** 后端移除 `goview` 依赖（`go.mod` + import）+ HTML 渲染代码
- [/] **T-405** 更新 README.md 部署章节：增加前端构建说明、Nginx 反代示例
- [/] **T-406** 调整 `scripts/Asgard-{web,master,agent}` 启动脚本（如需）：不再要求 `WORKDIR` 下有 `web/views/`

### Phase 5：验证与下线

- [/] **T-501** 端到端冒烟测试：11 个路由 group 全部走通一遍
- [ ] **T-502** 性能基线对比：列表页 TTFB、首页 FCP/LCP
- [/] **T-503** 移除 `web/views/` 与 `web/assets/` 目录
- [/] **T-504** 移除 HTML 控制器（`controllers/*` 里 `Render(...)` 调用清零）
- [/] **T-505** 更新 AGENTS.md §6 Web 层：改写为「纯 API 后端 + 独立前端工程」说明
- [/] **T-506** 移除 `web/middlewares/{login,admin,app,job,timing,...}` 中只服务于 HTML 路由的中间件
- [/] **T-507** 归档旧 `cmd_form.html` / `cmd_info.html` 等模板（如有需要保留为设计参考，可移至 `doc/legacy-templates/`）

## 4. 状态变更日志

> **追加**而非覆盖。每条一行：`YYYY-MM-DD 状态 T-xxx 备注`。
> 状态字段：`启动` / `完成` / `阻塞` / `恢复` / `备注`。

<!-- 在下方追加新条目，最新的放最下面 -->

- 2026-06-02 备注  创建本任务跟踪文档。初始登记前后端分离 Phase 0-5 共 60+ 条任务；待 Phase 0 决策落地后开始 Phase 1。
- 2026-06-02 完成  T-001  前端框架最终选择 Vue 3，理由：中文社区成熟、TS 友好、Element Plus / Pinia / Vite 生态一站式、无 SSR 需求。
- 2026-06-02 完成  T-002  构建工具最终选择 Vite 5，理由：Vue 3 官方推荐、HMR 极快、Rollup 生产构建成熟、pnpm/TS 一等公民。
- 2026-06-02 完成  T-003  UI 组件库最终选择 Element Plus，理由：Vue 3 配套最广、中文文档完整、Table/Form/Tree/Pagination 直接覆盖 Asgard 后台运维场景。
- 2026-06-02 完成  T-004  HTTP 客户端最终选择 axios，理由：拦截器生态最成熟、统一 401/错误码处理、ElMessage 集成方便。
- 2026-06-02 完成  T-005  状态管理最终选择 Pinia，理由：Vue 3 官方推荐、TS 推断友好、store 组合式 API 与 <script setup> 风格一致。
- 2026-06-02 完成  T-006  决定使用 TypeScript 严格模式，理由：编译期拦截 API 字段不匹配、IDE 重构友好、与 Pinia/Vue 3 类型推断完美契合。
- 2026-06-02 完成  T-007  前端工程位置选择仓库内 web-admin/ 子目录，理由：与后端同 repo 统一 review/版本/部署；产物独立打包；暂不引入 monorepo。
- 2026-06-02 完成  T-008  API 路由前缀选择 /api/v1，理由：版本化路径便于未来升级；现有 HTML 路由继续在根路径，与 API 并存。
- 2026-06-02 完成  T-009  鉴权协议决定过渡期双轨：保留 DES cookie 同时新增 JWT Bearer，理由：旧 HTML 路由继续可用、前端 Pinia store 缓存 token、API 中间件同时接受两种形式。
- 2026-06-02 完成  T-010  实时数据方案选择 SSE，理由：HTTP 友好、断线重连内置、EventSource 集成简单；仅在需要双向推送时才用 WebSocket。T-117/T-118/T-310/T-312 均按此约定。
- 2026-06-03 启动  T-101  在 web/router.go 注册 /api/v1 group，与现有 HTML 路由并存。
- 2026-06-03 完成  T-101  server.Group("/api/v1") 已落地，子路由拆到 web/routers/api_router.go（占位 /health），后续 T-102~T-119 逐任务追加。
