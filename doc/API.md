 # Asgard JSON API 文档

 > 前后端分离 Phase 1 落地。基础路径：`/api/v1`，所有响应统一为 `{code, message, data}` 结构（列表数据为 `{list, total, page, page_size}`）。
 >
 > 与现有 HTML 路由并存：HTML 路由继续在根路径（`/user/*`、`/app/*` 等），前端通过 Vite 代理到 `http://localhost:12345` 即可。

 ## 1. 通用约定

 ### 1.1 鉴权

 所有接口（除 `/health` 与 `/auth/login`）需要鉴权，二选一：

 ```http
 Authorization: Bearer <jwt>
 ```

 或保持现有 DES cookie（`token=<DES(userID, salt)>`，由 `web.cookie_salt` 加密，8 字节）。

 JWT 签发使用 HS256 + `web.jwt_secret`（`constants.WEB_JWT_SECRET`），TTL 由 `web.jwt_ttl` 控制（默认 7200s）。`sub` 为 userID，`exp` 为过期时间。

 ### 1.2 错误码

 | code | 含义 |
 | --- | --- |
 | 200 | 成功 |
 | 400 | 请求参数异常（参数缺失、格式错误） |
 | 401 | 未登录 / 登录失效 |
 | 403 | 权限不足（禁用 / 非管理员） |
 | 500 | 业务异常 |

 ### 1.3 分页

 列表接口统一返回：

 ```json
 {
   "code": 200,
   "message": "ok",
   "data": {
     "list": [...],
     "total": 100,
     "page": 1,
     "page_size": 20
   }
 }
 ```

 query 参数：`page`（默认 1），`page_size` 由后端常量控制（默认 20）。

 ### 1.4 批量操作

 ```json
 POST /api/v1/apps/batch_start
 { "ids": [1, 2, 3] }
 ```

 ## 2. 鉴权 `/api/v1/auth/*`

 | 方法 | 路径 | 说明 |
 | --- | --- | --- |
| POST | `/auth/login` | 登录，返回 token + user |
| GET | `/auth/info` | 当前用户信息 |
| POST | `/auth/logout` | 登出（清 cookie） |
| POST | `/auth/change_password` | 当前用户改密（会清 cookie） |

 ## 3. 用户 `/api/v1/users`（管理员）

 | 方法 | 路径 | 说明 |
 | --- | --- | --- |
| GET | `/users` | 列表，支持 `nickname/phone/email/status` 过滤 |
| GET | `/users/:id` | 详情 |
| POST | `/users` | 创建（仅管理员） |
| PUT | `/users/:id` | 更新（仅管理员） |
| POST | `/users/:id/forbidden` | 禁用（仅管理员） |
| POST | `/users/:id/reset_password` | 重置密码（仅管理员） |

 ## 4. 实例 `/api/v1/agents`

 | 方法 | 路径 | 说明 |
 | --- | --- | --- |
| GET | `/agents` | 列表 |
| GET | `/agents/:id` | 详情 |
| PUT | `/agents/:id` | 改别名（仅管理员） |
| POST | `/agents/:id/forbidden` | 禁用实例（级联停用该实例上的 app/job/timing） |

 ## 5. 分组 `/api/v1/groups`

 | 方法 | 路径 | 说明 |
 | --- | --- | --- |
| GET | `/groups` | 列表 |
| POST | `/groups` | 创建 |
| PUT | `/groups/:id` | 更新 |
| DELETE | `/groups/:id` | 删除（标记 DELETED） |

 ## 6. 应用 `/api/v1/apps`

 | 方法 | 路径 | 说明 |
 | --- | --- | --- |
| GET | `/apps` | 列表（普通用户自动按 creator 过滤） |
| GET | `/apps/:id` | 详情 |
| POST | `/apps` | 创建 |
| PUT | `/apps/:id` | 更新 |
| POST | `/apps/:id/copy` | 复制 |
| POST | `/apps/:id/start` | 启动 |
| POST | `/apps/:id/restart` | 重启 |
| POST | `/apps/:id/pause` | 暂停 |
| DELETE | `/apps/:id` | 删除（先暂停） |
| POST | `/apps/batch_start` | 批量启动 |
| POST | `/apps/batch_restart` | 批量重启 |
| POST | `/apps/batch_pause` | 批量暂停 |
| POST | `/apps/batch_delete` | 批量删除 |

 启动/重启/暂停/删除 都会通过 gRPC 通知 agent；agent 不在线时返回 400。

 ## 7. 计划任务 `/api/v1/jobs`

 接口与 `/apps` 一一对应，外加 `spec`（cron 表达式）、`timeout`（秒）字段。

 ## 8. 定时任务 `/api/v1/timings`

 接口与 `/apps` 一一对应，外加 `time`（执行时刻）、`timeout`（秒）字段。

 ## 9. 监控 `/api/v1/monitor`

 | 方法 | 路径 | 参数 | 说明 |
 | --- | --- | --- | --- |
| GET | `/monitor/agent` | `agent_id, size` | CPU/Memory 时间序列 |
| GET | `/monitor/app` | `app_id, size` | 同上 |
| GET | `/monitor/job` | `job_id, size` | 同上 |
| GET | `/monitor/timing` | `timing_id, size` | 同上 |

 ## 10. 归档 `/api/v1/archives`

 | 方法 | 路径 | 参数 | 说明 |
 | --- | --- | --- | --- |
| GET | `/archives/app` | `app_id, page` | 应用归档列表 |
| GET | `/archives/job` | `job_id, page` | 计划任务归档 |
| GET | `/archives/timing` | `timing_id, page` | 定时任务归档 |

 ## 11. 日志 `/api/v1/{out,err}_logs/*/data`

 | 方法 | 路径 | 参数 | 说明 |
 | --- | --- | --- | --- |
| GET | `/out_logs/app/data` | `app_id, lines` | 应用 stdout |
| GET | `/err_logs/app/data` | `app_id, lines` | 应用 stderr |
| GET | `/out_logs/job/data` | `job_id, lines` | 计划任务 stdout |
| GET | `/err_logs/job/data` | `job_id, lines` | 计划任务 stderr |
| GET | `/out_logs/timing/data` | `timing_id, lines` | 定时任务 stdout |
| GET | `/err_logs/timing/data` | `timing_id, lines` | 定时任务 stderr |

 ## 12. 异常 + 操作日志

 | 方法 | 路径 | 参数 | 说明 |
 | --- | --- | --- | --- |
| GET | `/exceptions` | `page, type` | 异常记录 |
| GET | `/operations` | `page, user_id, type` | 操作日志 |

 ## 13. SSE 实时数据 `/api/v1/sse/*`

 实时日志/监控走 SSE，浏览器用 `EventSource` 订阅。统一 `interval` query 参数控制推送频率（秒）。

 | 方法 | 路径 | 事件 | 说明 |
 | --- | --- | --- | --- |
| GET | `/sse/out_log/app` | `log / ping` | app_id, interval |
| GET | `/sse/err_log/app` | `log / ping` | app_id, interval |
| GET | `/sse/out_log/job` | `log / ping` | job_id, interval |
| GET | `/sse/err_log/job` | `log / ping` | job_id, interval |
| GET | `/sse/out_log/timing` | `log / ping` | timing_id, interval |
| GET | `/sse/err_log/timing` | `log / ping` | timing_id, interval |
| GET | `/sse/monitor/agent` | `point / ping` | agent_id, interval |
| GET | `/sse/monitor/app` | `point / ping` | app_id, interval |
| GET | `/sse/monitor/job` | `point / ping` | job_id, interval |
| GET | `/sse/monitor/timing` | `point / ping` | timing_id, interval |

 ## 14. 开发环境

 ### 14.1 Vite 代理

 ```js
 // web-admin/vite.config.ts
 server: {
   proxy: {
     '/api': {
       target: 'http://localhost:12345',
       changeOrigin: true,
     },
   },
 },
 ```

 ### 14.2 CORS

 开发期由 `web/middlewares/cors.go` 全局放开；生产建议由 Nginx 收敛。

 ## 15. 后续

 - Phase 2 起：前端 `web-admin/` 工程使用本文件作为唯一接口契约。
 - Phase 5：HTML 路由全量下线，本文件成为后端唯一对外接口。
