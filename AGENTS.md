# Asgard 项目开发者指南（AGENTS.md）

本文件面向后续参与 Asgard 维护或二次开发的工程师 / Agent。内容基于当前仓库根目录 `/Users/dengjialong/git/github/Asgard` 的实际结构与代码细节总结，请把它当作"先读这份，再去翻代码"的入口。

## 1. 项目是什么

Asgard 是一个**单二进制、多子命令**的分布式作业管理平台，目标是把常驻进程（`app`）、计划任务（`job`/crontab）和定时任务（`timing`/at）统一管理起来。它由三类节点组成：

- **web 节点** — Gin 渲染的管理控制台（`/Users/dengjialong/git/github/Asgard/web`），负责配置、状态查看、启停指令。
- **master 节点** — 中央调度面（`/Users/dengjialong/git/github/Asgard/cmds/master`），接收 agent 注册、转存监控/归档/异常数据、检测 agent 在线状态。
- **agent 节点** — 真正起进程的 worker（`/Users/dengjialong/git/github/Asgard/cmds/agent`），按 master 下发的配置起 `app/job/timing`，并把运行数据回写到 master。

架构图见 [doc/Asgard.png](/Users/dengjialong/git/github/Asgard/doc/Asgard.png)，架构与节点职责说明见 [README.md](/Users/dengjialong/git/github/Asgard/README.md)。

此外还有**单机模式**的两个守护子命令：

- `Asgard guard` — supervisor 风格的常驻进程守护，通过 Unix socket 提供控制面（`/Users/dengjialong/git/github/Asgard/cmds/guard`）。
- `Asgard cron` — crontab 风格的计划任务守护，同样走 Unix socket（`/Users/dengjialong/git/github/Asgard/cmds/cron`）。
- `Asgard debug` — 调试邮件 / RPC 的工具集（`/Users/dengjialong/git/github/Asgard/cmds/debug`）。

所有节点共享同一个二进制 `Asgard`，通过 cobra 派发（`/Users/dengjialong/git/github/Asgard/cmds/root.go`）。

## 2. 仓库结构与职责

| 目录 | 角色 | 关键文件 |
| --- | --- | --- |
| `main.go` | 入口；`recover` 兜底并执行 cobra 根命令 | `main.go` |
| `cmds/` | 子命令集合，每个子包一个 cobra command | `root.go`, `web/`, `master/`, `agent/`, `guard/`, `cron/`, `debug/` |
| `constants/` | 枚举、状态码、配置默认值（运行期由 viper 覆盖） | `constant.go`, `web.go`, `master.go`, `agent.go`, `error_code.go` |
| `models/` | GORM 数据模型 + 通用查询辅助 | `base.go`（`BaseModel/OperatorModel/CmdModel`）, `action.go`（查询函数）, 各业务模型 |
| `services/` | 业务逻辑层，封装 CRUD / 状态机 / 缓存 | 11 个 `*Service` 文件 |
| `providers/` | service 单例容器（`init()` 中实例化）+ agent gRPC 客户端缓存 | `service.go`, `client.go` |
| `managers/` | agent 端"运行期编排"层：进程/协程/crontab 状态机 | `agent.go`, `app.go`, `job.go`, `timing.go` |
| `runtimes/` | agent 端"真正干活"的层：`Command` 抽象 + `App/Job/Timing` 各自生命周期 + `Monitor` + 异常/归档结构体 | `cmd.go`, `app.go`, `job.go`, `timing.go`, `monitor.go`, `archive.go`, `exception.go`, `config.go` |
| `server/` | gRPC server 端实现（`AgentServer/MasterServer/CronServer/GuardServer/TimerServer/baseServer`） | `base.go`（`NewRPCServer`, `GetLog`） + 各 server |
| `clients/` | gRPC client 端封装 | `agent.go`, `master.go`, `guard.go`, `cron.go`, `base.go`（Unix socket dialer） |
| `registry/` | etcd 集成：master 集群 leader 选举 + 客户端 resolver | `registry.go`, `resolver.go` |
| `rpc/` | protoc 生成的 gRPC 代码 + `Build/Format/Parse` 转换函数 | `common.go`（**所有 model↔rpc 转换都集中在这里**） |
| `protos/` | proto 源文件；改完要 `protoc` 重新生成 | `base.proto`, `app.proto`, `job.proto`, `timing.proto`, `agent.proto`, `master.proto` |
| `web/` | web 节点实现：路由 / 控制器 / 中间件 / 模板 | `server.go`, `router.go`, `controllers/`, `middlewares/`, `utils/`, `views/`, `assets/` |
| `scripts/` | 部署脚本 + 初始化 SQL + protoc 脚本 | `Asgard.sql`, `Asgard.service`, `Asgard-{web,master,agent}`, `protoc.sh` |
| `doc/` | 架构图与界面截图 | `Asgard.png`, `page[1-4].png` |

## 3. 运行时拓扑与数据流

### 3.1 启动路径

每个子命令都遵循统一模式（见 `cmds/<cmd>/<cmd>.go`）：

1. `cmd.Flag("conf").Value.String()` 拿到配置目录（默认 `conf`）。
2. `runtimes.ParseConfig(confPath)` 用 viper 读 `app.yaml`，把 `system.moniter / system.timer` 写入 `constants`。
3. 注册所需的 avenger 组件（`logger.Register / db.Register / cache.Register / mail.Register`）。
4. 初始化对应 `Manager`（`agent → AgentManager`，`web → 各种 Controller`，`master → registry+providers`）。
5. 启动协程 + `runtimes.Wait(<stopFunc>)` 阻塞监听 `SIGKILL/SIGHUP/SIGINT/SIGQUIT/SIGTERM`（`runtimes/cmd.go`）。

### 3.2 配置 → 进程 数据流

```
Web (Gin) ──> master (DB) ──> agent gRPC ──> appManager / jobManager / timingManager ──> os/exec
   │                                                                          │
   └─> controllers ──> providers.*Service ──> models.* (GORM)                 └─> runtimes.Command
                                                                                       │
                       ←── Archive / Monitor / Exception chan ←── Master Client ←───┘
```

- 启动/停止/重启：web 端经 controller 调 `providers.GetAgent(...)` 拿到 gRPC 客户端（`providers/client.go` 里带缓存），直接 push 到 agent；agent 端 `managers.<X>Manager.Add/Update/Remove` 会操作 `runtimes` 层的 `App/Job/Timing`。
- 状态回写：master 端 `cmds/master/master.go` 的 `MoniterMaster` 用 `time.Ticker` 周期扫库 + TCP 探活 agent 端口（`checkPort`），用 `markAppStatus/markJobStatus/markTimigStatus`（注意 `Timig` 拼写）回写状态。
- 运行数据上报：agent 端 `Command` 完成后回调 `ArchiveReport / MonitorReport / ExceptionReport`（在 `managers` 里注入），`Master.Report()` 单协程消费 9 个 channel 后调用 master 的 gRPC（`MasterServer.*ArchiveReport` 等）。

### 3.3 状态机

`constants/constant.go` 里集中定义了所有状态的 int64 值和对应的中文展示名。任何修改状态语义的工作请同时更新：

- `constant.go` 里的常量 + `XXX_STATUS` 切片
- `models/*` 默认值（部分表 `status` 默认为 `0`/停）
- `services/*` 中的状态转换函数（如 `AppService.ChangeAPPStatus`）
- `web/utils/format.go` 里的 `GetObjectName`（按 `TYPE_*` 拿名字）

常见状态：

- `AGENT_*`：`ONLINE=1 / OFFLINE=0 / FORBIDDEN=-1`
- `APP/JOB_STATUS_*`：`RUNNING=1 / PAUSE=2 / STOP=0 / UNKNOWN=-2 / DELETED=-1`
- `TIMING_STATUS_*`：比上面多一个 `FINISHED=3`
- `TYPE_*`：1=agent, 2=app, 3=job, 4=timing, 5=group, 6=user（**给 operation 表和 exception 表用**）

## 4. gRPC 协议

- 协议定义在 [protos/](/Users/dengjialong/git/github/Asgard/protos/)，生成代码在 [rpc/](/Users/dengjialong/git/github/Asgard/rpc/)。`scripts/protoc.sh` 给出重新生成命令：

  ```sh
  protoc -I protos protos/*.proto --go_out=plugins=grpc:./rpc/
  ```

  注意：此命令需要老版本 `protoc-gen-go`（`plugins=grpc` 语法只支持 `v1.x`）。在 Go 1.17+ 工具链上请用 `protoc-gen-go-grpc` 插件并调整参数，否则需固定旧版工具。

- 四个 gRPC service：`Agent / Master / Guard / Cron`（proto 中还有一个 `Timer`，对应 `server/timer.go` 但目前**未被注册到任何 agent 端**——新增时记得在 `cmds/agent/agent.go` 的 `StartRpcServer` 中加 `rpc.RegisterTimerServer`）。

- 所有"业务对象 ↔ protobuf 消息"的转换都在 [rpc/common.go](/Users/dengjialong/git/github/Asgard/rpc/common.go)，不要在 controllers / managers / services 里再次手动映射：

  - `FormatApp/Job/Timing`：从 `models` 走向 proto（master 给 agent 看）
  - `BuildApp/Job/Timing`：从 `runtimes` 走向 proto（agent 内 RPC 用）
  - `BuildArchive / BuildAgentMonitor / BuildAppArchive / ...`：runtimes 内部数据上送 master
  - `ParseMonitor / ParseArchive / ParseException`：master 把 proto 落地成 `models` 记录
  - `BuildAppConfig / BuildJobConfig / BuildTimingConfig`：proto → `map[string]interface{}` 用于 `managers.Register`（**新增字段时记得同步**）

## 5. agent 端运行期：managers + runtimes

`runtimes/cmd.go` 中的 `Command` 是核心抽象，所有 `App/Job/Timing` 都内嵌它。`Configure(config map[string]interface{})` 负责校验基础字段（name/dir/program/args/stdout/stderr/is_monitor），具体类型再补字段（`App` 补 `auto_restart`，`Job` 补 `spec/timeout`，`Timing` 补 `time/timeout`）。

生命周期（具体到 `runtimes/cmd.go` 的方法）：

1. `build()` — 拆 `args`、建 `exec.Cmd`、准备 stdout/stderr 文件（自动建目录）。
2. `start()` — `cmd.Start()`，记 `UUID/Begin/Pid/Running=true`，如开启监控则 `Monitor.Add`。
3. `wait(callback)` — `cmd.Wait()` 后填 `Status/Signal/Successed`，回调 `ArchiveReport`，再触发 `callback`（App 用作 `restart`，Job/Timing 用作 `record`）。
4. `Kill()` — `Process.Kill()` + `ArchiveReport`。
5. `finish()` — 关锁、清监控、置 `Running=false`。

三类的差异（请改之前先读对应文件）：

- **App**（`runtimes/app.go`）：可选 `auto_restart`；用 `mcache` 计数 5 分钟内重启次数，超 5 次记 `Dead` 不再启。
- **Job**（`runtimes/job.go`）：外层包了 `robfig/cron/v3`；每次触发 `Run()`，可选 `timeout` 强杀。
- **Timing**（`runtimes/timing.go`）：`Time time.Time` + `Executed bool`；`managers.TimingManager.Run` 周期检查 `time.Now().Unix()` 到点即 `UnRegister`（一次性）。

agent 的 `AgentManager.StartAll()` 顺序：上报协程 → 自身 `AgentMonitorReport` → `AgentRegister` → `AppsRegister/JobsRegister/TimingsRegister` → `StartAll(true)` 三个子 manager。

## 6. Web 层（Gin + goview）

### 6.1 装配

`web/server.go` 中 `Init()` 负责：路由 / 模板 / 静态资源（`web/assets` 静态文件由 `goreleaser` 打包进去）。`Run()` 调 `setupController()`（注入 viper → `constants.WEB_OUT_DIR`）后 `setupRouter()`。

### 6.2 路由与中间件分层

`web/router.go` 是一张平铺的对照表，业务上把每类资源（user/agent/group/app/job/timing/monitor/archive/out_log/err_log/exception/operation）分成一个 group：

- `*Init`（如 `AppInit/AppAgentInit/BatchAppAgentInit`）— 解析 query/post 里的 `id`，把 `*models.App` 之类的对象塞到 `gin.Context` 上，并把 `agent` 也一并塞好（带 `AGENT_ONLINE` 校验）。
- `Admin` — 角色校验。
- `Login` — 解析 `token` cookie（DES 加密 userID，盐来自 `web.cookie_salt`，**必须 8 字节字符串**）。
- `CmdConfigVerify` — 提交时的非空校验。

> 修改任何 `*Init` 的语义请同时核对 `web/utils/request.go` 里的 `GetApp/GetAppAgent/...` 提取函数。

### 6.3 控制器写法

以 `controllers.AppController` 为代表（`/Users/dengjialong/git/github/Asgard/web/controllers/app.go`）：

- 列表页用 `where := map[string]interface{}{}` + `providers.XService.GetXPageList`，**普通用户自动按 `creator` 过滤**。
- 操作类接口（start/restart/pause/delete/copy）的标准顺序是：`utils.GetApp(ctx)` → `utils.GetAgent(ctx)`（中间件塞好的）→ `providers.GetAgent(agent)` → 通过 gRPC 推 agent → `providers.XService.ChangeXStatus` 改库 → `utils.OpetationLog(...)` 写操作日志 → `utils.APIOK(ctx)`。

> 注意：批量接口的中间件把对象装到 `app_agent / job_agent / timing_agent` map，controller 用 `utils.GetAppAgent(ctx)` 拿。

### 6.4 响应与模板

`web/utils/respose.go` 提供了 `Render / APIOK / APIError / Warning / JumpWarning / JumpError` 一组小工具；HTML 模板用 goview 的"layouts + partials"机制（`web/views/layouts/master.html` + `web/views/templates/*.html`）。`web/utils/html.go` 的 `PagerHtml` 是后端拼好的分页 HTML。

> 注意：`web/utils/opetation.go` 写的是 `OpetationLog`（**注意拼写**：是 *opetation*，不是 *operation*）；`utils.OpetationLog` 是控制器里实际调用的名字，不要去"修正"它。

### 6.5 资源视图

`web/views/` 的每个子目录对应控制器的实体：agent / app / archive / exception / group / job / log / monitor / operation / timing / user + 顶层 `index.html` 和 `warning.html`。`web/assets/js/asgard.js` 里有 `Asgard.getData / postData` 这种统一 ajax 包装，约定成功码是 `200`。

## 7. 配置与数据

### 7.1 app.yaml 关键项

viper 读 `app.yaml`，按节点类型只读对应前缀：

- 公共：`system.moniter / system.timer`、`component.{db,redis,log,mail}.*`
- `web.*`：`port / mode / domain / cookie_salt`（详见 `cmds/web/web.go` 和 `constants/web.go`）
- `master.*`：`port / cluster / cluster_registry / cluster_name / cluster_id / cluster_ip / moniter / notify / receiver`（详见 `cmds/master/master.go`）
- `agent.*`：`moniter / master.{ip,port} | master.{cluster,cluster_registry,cluster_name}` / `rpc.{ip,port}`（详见 `cmds/agent/agent.go`）

> viper 读不到配置就直接 `panic`，所以部署前请保证 `conf/app.yaml` 存在。`.gitignore` 已经忽略 `conf/` 目录。

### 7.2 数据库与缓存

- MySQL：表结构见 [scripts/Asgard.sql](/Users/dengjialong/git/github/Asgard/scripts/Asgard.sql)。`monitors_YYYYMM` 按月分表，**部署时记得提前建好**（README 也强调过）。`models.Monitor.TableName()` 会动态取 `time.Now().Format("200601")`。
- Redis：仅用于 `UserService / GroupService` 缓存（`services/cache.go` → `constants.CACHE_KEY_*`）。
- GORM 1.x：所有 `Where(...).Find(list)` 风格。

### 7.3 鉴权

登录走 `controllers/user.go` 的 `DoLogin` → `SetTokenCookie(token)` → cookie 形如 `DES(userID, WEB_COOKIE_SALT)`，中间件 `Login` 解密并 `UserService.GetUserByID` 注入到 context。

> 没有显式登出失效机制（除修改密码外），token 寿命依赖 cookie 的 maxAge=7200s。

## 8. 高可用（master / agent 集群）

- 仅当 `master.cluster: true` 时启用 etcd：
  - `cmds/master/master.go` 中 `RegisterRpcServer` 通过 `registry.Register` 把 `MASTER_CLUSTER_NAME/ID/IP:PORT` 写进 etcd，lease 由 `MASTER_CLUSTER_TTL=10s` 续约。
  - `Campaign("/Asgard/leader", MASTER_CLUSTER_ID)` 进行 leader 选举；`registry.IsLeader()` 决定 `MoniterMaster` 是否做探活（避免多 master 重复探测）。
- agent 端 `cluster: true` 时用 `registry.NewResolver` + `grpc.WithDefaultServiceConfig({loadBalancingPolicy:"round_robin"})` 实现客户端负载均衡。
- 单点（非集群）模式 `clients.Master` 走直连 `ip:port`。

## 9. 命名与代码约定

- 包名：单数、无下划线，与目录同名；`Asgard/...` 全路径 import。
- 错误处理：业务层多用 `logger.Error` 吞错并 `return nil/0`；RPC 边界用 `*Response{Code: 500/404/200}`。
- 中文 / 英文混排：常量键名与 SQL 注释大量使用中文，**保持现状**，不要再回译。
- 状态值：int64，硬编码集中在 `constants/constant.go`；不要在 services / controllers 里出现裸数字。
- 文件/方法命名偏向动词或领域词（`ChangeAPPStatus / ReStart / BatchPause`），不强行 Go-style 短名。
- `go.mod` 是单一 `module Asgard`，没有内部 module；新增子包直接放顶层目录。

## 10. 构建、运行与发布

### 10.1 工具链

- Go 1.16（`go.mod` 声明）。本地用 1.26 时**会因 `gopsutil v2.19.11` 不再支持新 Darwin syscalls 而编译失败**（`process_darwin.go` 找不到 `KinfoProc`）。要本地构建请：① 用 Go 1.16~1.20 工具链；或 ② 临时把 `gopsutil` 升到 v3+。
- 第三方包关键依赖：`github.com/dalonghahaha/avenger`（自研组件库，提供 `db/cache/logger/mail`），`gin-gonic/gin`，`jinzhu/gorm`，`coreos/etcd`（`clientv3` + `concurrency`），`robfig/cron/v3`，`shirou/gopsutil`，`patrickmn/go-cache`。

### 10.2 常用命令

```sh
# 本地构建
go build -o Asgard .

# 单节点（带 conf 目录，里面放 app.yaml）
./Asgard web     -c conf
./Asgard master  -c conf
./Asgard agent   -c conf
./Asgard guard   -c conf [-s runtime/asgard_guard]
./Asgard cron    -c conf [-s runtime/asgard_cron]
./Asgard agent status [-p 27149]    # 看本机 agent 状态

# 重新生成 gRPC 代码
bash scripts/protoc.sh             # 注意：需要旧版 protoc-gen-go

# 发布多平台二进制（参考 .goreleaser.yml）
goreleaser release --snapshot --clean
```

### 10.3 部署素材

- `scripts/Asgard-{web,master,agent}` 是 `/etc/init.d/` 风格启停脚本，默认 `WORKDIR=/data/Asgard`、日志在 `runtime/`。`Asgard.service` 是 systemd unit 模板。
- 端口默认：web=12345 / master=9527 / agent=27149。
- `web/views/`、`web/assets/`、`scripts/*` 都被 `.goreleaser.yml` 打进发布包。

## 11. 常见改动手册

下面这些是最容易踩坑的几类改动，动手前先读相关章节：

### 11.1 新增一种监控对象

1. `models/<name>.go` 写 model（嵌入 `BaseModel`）。
2. `services/<name>.go` 写 Service（参考 `MonitorService`）。
3. 在 `services/cache.go` 旁的 `providers/service.go` 增加全局变量和 `init()` 注入。
4. `constants/constant.go` 加 `TYPE_*` 和状态常量。
5. `protos/*.proto` 加消息 + service，`scripts/protoc.sh` 重生 `rpc/*.pb.go`。
6. `rpc/common.go` 写 `Format/Build/Parse` 转换。
7. `server/<name>.go` 实现 gRPC server，必要时在 `cmds/*/StartRpcServer` 注册。
8. `managers/<name>.go`（如果是 agent 端在跑）和 `runtimes/<name>.go` 写运行期。
9. `web/controllers/<name>.go` + `web/middlewares/<name>.go` + `web/views/<name>/` 模板。
10. `scripts/Asgard.sql` 加表。

### 11.2 改一种状态

1. `constants/constant.go` 改值 / 改展示名。
2. 搜 `WhereAndOrder / ChangeXStatus / markXStatus`（按命名模糊搜）确认没漏。
3. `web/views/.../list.html` 里 `StatusList` 是从 `constants.X_STATUS` 渲染的，**多数情况无需改模板**。

### 11.3 加一个新 web 页面

1. `web/views/<entity>/<page>.html`（继承 `layouts/master.html` 的 `{{define "content"}}`）。
2. `web/controllers/<entity>.go` 写方法，模板里出现的 `Subtitle / List / Pagination` 等 key 一一对应。
3. `web/router.go` 在对应 group 加 `server.GET/POST(...)`。

### 11.4 改 gRPC 消息字段

- `protos/*.proto` 改字段后 **必须** 重新生成 `rpc/*.pb.go`，且**新增/重命名字段**时同步 `rpc/common.go` 中的 `Format/Build/Parse/BuildXConfig`。
- proto 中尽量用 `bool` / `int64` 简化跨语言语义；不要引入 `oneof`（项目里没有先例）。

### 11.5 加新中间件

- 在 `web/middlewares/` 放文件，`web/router.go` 里 `group.Use(...)` 或具体路由上加。
- 若要往 `gin.Context` 塞对象，记得在 `web/utils/request.go` 里加对应的 `GetX(ctx)` 提取函数。

## 12. 已知坑 / 风险提示

- `gopsutil` 旧版本不兼容新 macOS SDK（如上文）。
- `protoc` 生成脚本依赖的 `protoc-gen-go` 旧版插件，新工具链不会直接跑通。
- `providers/client.go` 里 `AgentClients` 是按 `agent.ID` 缓存的 gRPC 连接，**永远不会失效**（agent 改 IP/端口后必须重启 master 节点或清缓存）。
- `models/monitor.go` 的分表名按当前月生成，跨月写入当月表；历史表由 SQL 脚本预创建。
- `services/agent.go` 在缓存里 `GetAgentByIPAndPort` 但 `providers/client.go` 是按 ID 缓存，两者并不冲突但容易混淆。
- master 的 `MoniterMaster` 是 `time.NewTicker` 同步循环，规模大了是瓶颈，但当前实现没有并发控制。
- web `cookie_salt` 必须是 8 字节（DES 限制），改完启动时不会校验，启动后才会在登录时炸。
- `services/*` 大量用 `WhereAndOrder / PageListbyWhereString` 手搓 SQL，**没做参数化**，对用户可控字段（`name/nickname/...`）做模糊查询时按 `%v` 直接拼。新加接口请注意注入风险。

---

读到这里应该已经能动手改了。如果还需要更深的细节，建议按下面顺序再去翻代码：

1. `cmds/<cmd>/<cmd>.go` — 看启动顺序和 viper key。
2. `managers/<x>.go` + `runtimes/<x>.go` — 看 agent 端运行期语义。
3. `services/<x>.go` + `models/<x>.go` — 看 master/web 端持久化语义。
4. `server/<x>.go` + `rpc/common.go` — 看 gRPC 边界和 model↔proto 转换。
5. `web/router.go` + `web/controllers/<x>.go` — 看 web 端业务流。

Good luck.
