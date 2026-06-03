 # Asgard 性能基线（前后端分离后）

 > T-502 占位文档：手动跑一次的基线样例，作为后续回归参考。
 > 实际基线需要在目标环境（与生产同等规格）跑 `scripts/bench.sh`（见末尾）后填入真实数字。

 ## 测试环境

 - 后端：1× Asgard master + 1× Asgard web + 1× Asgard agent，单机 4C8G，MySQL 8.0 + Redis 7
 - 前端：web-admin build artifact（nginx 1.27 托管），访问端 Chrome 120 + Lighthouse
 - 数据：10 个 agent、500 个 app、200 个 job、100 个 timing

 ## 后端 API 基线

 | 接口 | 期望 P50 | 期望 P95 | 期望 P99 |
 | --- | --- | --- | --- |
 | `POST /api/v1/auth/login` | < 80ms | < 200ms | < 400ms |
 | `GET  /api/v1/apps?page=1&page_size=20` | < 30ms | < 80ms | < 200ms |
 | `GET  /api/v1/monitor/agent?agent_id=1` | < 50ms | < 120ms | < 300ms |
 | `POST /api/v1/apps/:id/start` | < 60ms | < 200ms | < 500ms |
 | `GET  /api/v1/sse/out_log/app`（SSE） | 第一帧 < 1s | 稳态 < 1s | 持续不断流 |

 ## 前端基线

 | 指标 | 期望 |
 | --- | --- |
 | 首屏 FCP | < 1.5s |
 | 首屏 LCP | < 2.5s |
 | TTI（可交互） | < 3.0s |
 | App 列表页切换 | < 200ms（前端） |

 ## 跑分方法

 1. 后端：`wrk -t8 -c200 -d30s http://localhost:12345/api/v1/apps?page=1`
 2. 前端：Chrome DevTools Lighthouse 跑一次 PWA / Performance
 3. SSE：用 `curl -N` 订阅 `/api/v1/sse/out_log/app?app_id=1` 确认每秒有事件

 ## 历史记录

 - 2026-06-03 前后端分离首次落地：后端 11 个 route group / 60+ 接口；前端 17 个页面、9 个公共组件、4 类状态展示；nginx 1× 反代 + 1× 静态托管。

 ## 后续优化点

 - 列表接口加 GORM 预加载（`Preload`）减少 N+1
 - 监控接口按 `created_at` 加分区/索引
 - SSE 端点加 `Last-Event-ID` 重连机制（目前用 EventSource 内置重连）
 - 前端 App/Job/Timing 列表可虚拟滚动（Element Plus 暂未内置）
