 # Asgard Web Admin

 Asgard 管理控制台的前端工程。Vue 3 + Vite + TypeScript（严格模式）+ Element Plus + Pinia + Vue Router 4 + Axios。

 ## 快速开始

 ```bash
 cd web-admin
 npm install
 npm run dev          # 开发服务器：http://localhost:5173
 npm run build        # 生产构建（产物在 dist/）
 npm run preview      # 本地预览 dist
 npm run lint         # ESLint + Prettier
 npm run typecheck    # vue-tsc 严格类型检查
 npm run test         # Vitest 单元测试
 npm run e2e          # Playwright 端到端测试
 ```

 ## 后端联调

 默认通过 Vite 代理把 `/api/*` 转发到 `http://localhost:12345`（即 `Asgard web`）。如需切换：

 ```bash
 VITE_BACKEND_TARGET=http://your-backend:12345 npm run dev
 ```

 鉴权方式：登录成功后 `localStorage` 存 JWT，axios 自动加 `Authorization: Bearer <token>`；同时后端继续用 DES cookie 兜底，迁移期两种方式都可用。

 ## 工程结构

 ```
 src/
   api/           # 类型化 axios 封装（http.ts 拦截器 + index.ts 业务 API）
   assets/        # 全局样式（variables/reset/layout）
   components/    # 通用组件（MonitorChart 图表、TerminalLog xterm）
   layouts/       # MainLayout 框架（侧边栏 + 顶栏 + 内容区）
   router/        # vue-router 配置 + 鉴权守卫
   stores/        # Pinia store（auth）
   types/         # 业务类型（与后端 models 对齐）
   utils/         # 通用工具（status/format）
   views/         # 11 个路由 group 的页面
     auth/        # 登录
     users/       # 用户管理 + 个人设置 + 改密
     agents/      # 实例管理
     groups/      # 分组管理
     apps/        # 应用管理
     jobs/        # 计划任务
     timings/     # 定时任务
     monitor/     # 监控
     archive/     # 归档
     logs/        # 日志
     exception/   # 异常记录
     operation/   # 操作日志
     error/       # 错误页（no_login/auth_fail/admin_only/forbidden/not_found）
   App.vue
   main.ts
 ```

 ## 与后端 API 的对接

 详见 [doc/API.md](../doc/API.md)。

 鉴权 token 默认 TTL 7200s（与后端 `WEB_JWT_TTL` 一致）；过期后由后端返回 401，前端拦截器跳 `/login`。
