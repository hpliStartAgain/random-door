# Random Door MVP 任务清单

## 执行规则

- 每个新会话先完整读取 `input.md`、`plan.md`、`tasks.md`，并用内置计划工具注册本轮小 todo。
- 每轮只执行第一个未完成 task；不得跨 task 提前修改。
- 普通 task 完成前必须基于 diff、日志、类型检查、构建或测试证据反复检查；有代码修改时提交代码。
- 每完成三个普通 task，必须先执行紧随其后的全面检查-debug循环。
- 完成 task 后填写：完成内容、验证结果、剩余风险、下一步、提交。

## 状态说明

- `[ ]` 未开始
- `[-]` 进行中
- `[x]` 已完成

## Tasks

### Task 1 [x] 建立完整基线审计与缺口矩阵

读取全部仓库约束、产品文档、设计文档和现有代码；排除 `node_modules` 后盘点目录、git 状态、配置、测试和契约覆盖情况。将逐条验收缺口、已有实现可信度、后续改动边界写入本文件记录区。

独立验证：形成覆盖验收标准 1~20、API 1~10、部署与安全约束的缺口矩阵；运行只读基线命令并记录结果。

完成内容：

- 完整读取现有 `AGENTS.md`、`CLAUDE.md`、4 份现存 agent 子约束、4 份产品文档、7 份 design 文档、6 份 arch 文档和根配置文件。
- 排除 `frontend/node_modules` 后盘点仓库：后端 50 个 `.go` 文件、12 张表 schema、12 城 seed、11 个成就 seed；前端 8 个页面文件但仅 1 个 API 文件，真实业务联调未完成。
- 记录初始化前已有脏工作树边界：`backend/data/seed/cities.json`、`backend/go.mod`、`docker-compose.yml` 已修改，大量后端文件和 `AGENTS.md` 已存在但未跟踪。后续不得粗暴回滚。
- 建立下方覆盖验收标准 1~20、API 1~10、部署、质量和安全约束的基线缺口矩阵。

验证结果：

- `go test ./...`：失败；`backend/internal/api/router.go` 存在未使用的 `service` import，后端无法编译；所有后端包均无 `_test.go`。
- `go vet ./...`：同样因未使用 import 失败。
- `gofmt -l .`：列出 `backend/cmd/server/main.go`、`backend/internal/service/city_service.go`、`backend/internal/service/game_service.go`、`backend/internal/upload/validator.go`。
- `npm run lint`：失败；前端没有 ESLint 配置文件。
- `npm run build`：失败；存在 Node/Vite 类型缺失、`ImportMeta.env` 类型缺失、未使用变量、overlay/store 类型漂移、`pannellum-react` 声明缺失等错误。构建生成的 4 个派生文件已在本 task 内精确清理。
- `docker compose config`：未执行成功；当前机器找不到 `docker` 可执行文件。
- seed 只读解析：通过；当前 JSON 为 12 城，每城 2 个地标、2 个美食、1 个人物且有方言字段；成就 JSON 为 11 条；schema 检出 12 张表。
- `git diff --check`：通过，仅有现存 LF/CRLF 提示。
- 密钥静态扫描：未发现前端 AI Key 或真实密钥提交；`.env` 与 `frontend/.env` 均处于 ignore 状态。

剩余风险：

- 当前产品不可构建、不可启动、不可部署、不可交互验收。后续必须严格按 task 顺序修缮。
- 本机缺少 Docker，Compose 运行态证据要在 Docker 可用环境补验；这不阻断当前可离线推进的实现任务。
- 文档引用但缺失：`docs/agent/git-workflow.md`、`docs/agent/doc-writing.md`、`docs/arch/observability.md`、多份 product 扩展文档、`.cursorrules` 和 `docs/design/00-detailed-design-index.md`。

下一步：

- 执行 Task 2：对齐 schema、model、seed 与幂等初始化所需数据约束。

提交：

- 本 task 仅更新 goal 控制文件；提交见本轮 git commit。

### Task 2 [ ] 修缮数据库 schema、model 与 12 城 seed

对齐 12 张表、枚举、索引、GORM model、12 城精选内容和成就 seed；确保每城满足地标、美食、人物、方言要求，并保证初始化幂等。

独立验证：schema/model 对照；JSON 解析；断言 12 城、每城内容数量、人物数量和成就规则；运行数据库相关测试。

完成内容：

验证结果：

剩余风险：

下一步：

提交：

### Task 3 [ ] 修缮后端匿名用户、城市与自由访问 API

完成 `POST /api/users/anonymous`、`GET /api/cities`、`GET /api/cities/{city_id}`、`POST /api/visits/free` 的 repository/service/handler 分层和契约一致性。

独立验证：运行相关 Go 测试；对合法输入、缺失用户、缺失城市、非法参数和统一错误包装逐项验证。

完成内容：

验证结果：

剩余风险：

下一步：

提交：

### Checkpoint A [ ] 完成 Task 1~3 后全面检查-debug循环

检查需求偏离、代码 bug、`go test ./...`、`go vet ./...`、前端 lint/build 基线、schema/seed 一致性、安全、文档同步和回滚策略；发现问题必须在本检查点内修复并复测。

检查结果：

修复内容：

验证结果：

剩余风险：

提交：

### Task 4 [ ] 修缮地理随机漫游算法

对齐八方向、距离档位、球面目标点、Haversine 距离、最近城市匹配和兜底策略，覆盖边界条件。

独立验证：为 geo 包运行表驱动测试，覆盖确定性输入、范围、空城市、边界和兜底。

完成内容：

验证结果：

剩余风险：

下一步：

提交：

### Task 5 [ ] 修缮游戏 init/roll API 与访问记录

完成 `POST /api/game/init` 和 `POST /api/game/roll`，确保位置默认值、骰子、方向、距离、目标城市、roll 与 visit 记录一致。

独立验证：运行 service/handler 测试，核对合法流程、默认位置、非法用户、非法坐标、空城市和事务一致性。

完成内容：

验证结果：

剩余风险：

下一步：

提交：

### Task 6 [ ] 修缮成就引擎与成就墙 API

完成通用成就、游戏专属成就、幂等解锁和 `GET /api/users/{user_id}/achievements`，对齐规则字典。

独立验证：运行规则表驱动测试和 service 测试，覆盖重复打卡、自由模式、游戏模式、方向连续性和未知规则。

完成内容：

验证结果：

剩余风险：

下一步：

提交：

### Checkpoint B [ ] 完成 Task 4~6 后全面检查-debug循环

检查需求偏离、算法边界、数据库事务、数据一致性、`go test ./...`、`go vet ./...`、前端 lint/build、安全、文档同步和回滚策略；发现问题必须在本检查点内修复并复测。

检查结果：

修复内容：

验证结果：

剩余风险：

提交：

### Task 7 [ ] 修缮 AI 对话编排与 chat API

完成角色 prompt、免责声明、历史消息、输入限制、超时、失败降级、后端 provider 边界和 `POST /api/chat`。

独立验证：使用 mock provider 运行测试，覆盖成功、超时、敏感输入、历史记录、非法角色和 provider 失败；扫描前端无 AI Key。

完成内容：

验证结果：

剩余风险：

下一步：

提交：

### Task 8 [ ] 修缮上传、生图与赛博照片 API

完成上传校验、安全存储、生图 prompt、合规限制、超时、失败降级和 `POST /api/checkin/generate-image`。

独立验证：使用 mock provider 测试有效图片、格式错误、超限、路径安全、非法 city/landmark、provider 失败和返回 URL。

完成内容：

验证结果：

剩余风险：

下一步：

提交：

### Task 9 [ ] 修缮打卡 API 与成就联动

完成 `POST /api/checkin`，确保打卡记录、模式来源、生成图引用和成就解锁事务行为一致。

独立验证：运行 service/handler 测试，覆盖自由打卡、游戏打卡、重复请求、非法 visit、非法图片引用和成就返回。

完成内容：

验证结果：

剩余风险：

下一步：

提交：

### Checkpoint C [ ] 完成 Task 7~9 后全面检查-debug循环

检查需求偏离、AI 合规、上传安全、事务、数据一致性、`go test ./...`、`go vet ./...`、前端 lint/build、文档同步和回滚策略；发现问题必须在本检查点内修复并复测。

检查结果：

修复内容：

验证结果：

剩余风险：

提交：

### Task 10 [ ] 修缮后端启动、配置、路由与中间件

检查 config、依赖注入、路由注册、健康检查、CORS、日志、recover、限流、静态上传访问、连接池和 seed 自动导入。

独立验证：运行全量 Go 测试、`go vet ./...`，启动后端并验证 health、路由、错误包装、限流和 seed 日志。

完成内容：

验证结果：

剩余风险：

下一步：

提交：

### Task 11 [ ] 修缮前端 API 封装与 Zustand 状态

读取前端约束；对齐 API 契约，实现匿名用户恢复、城市、访问、游戏、聊天、打卡和成就的类型化封装与状态管理，不在组件中直接 fetch。

独立验证：运行前端 lint/build；静态搜索组件直接 fetch、AI Key 和违规 SDK 调用；核对 TS 类型。

完成内容：

验证结果：

剩余风险：

下一步：

提交：

### Task 12 [ ] 修缮首页、模式选择与路由骨架

完成 HomePage、ModeSelectPage、Navbar、Sidebar、路由和响应式基本布局，确保入口清晰可达。

独立验证：运行前端 lint/build；使用 Browser 验证首页加载、两种模式入口、导航和基础响应式。

完成内容：

验证结果：

剩余风险：

下一步：

提交：

### Checkpoint D [ ] 完成 Task 10~12 后全面检查-debug循环

检查需求偏离、后端启动、路由中间件、前端 lint/build、基础 UI/UX、安全、配置、文档同步和回滚策略；发现问题必须在本检查点内修复并复测。

检查结果：

修复内容：

验证结果：

剩余风险：

提交：

### Task 13 [ ] 修缮自由探索地图与城市详情

完成城市标记、点击访问、详情跳转和城市地标/美食/人物/方言展示；高德 SDK 调用集中到地图组件。

独立验证：运行前端 lint/build；使用 Browser 验证自由探索点击城市进入详情，检查加载态、空态和错误态。

完成内容：

验证结果：

剩余风险：

下一步：

提交：

### Task 14 [ ] 修缮游戏互动地图、骰子与移动动画

完成当前位置/默认位置初始化、掷骰、随机方向距离反馈、目标城市和地图移动动画。

独立验证：运行前端 lint/build；使用 Browser 验证游戏 init、roll、动画结束进入目标城市和错误态。

完成内容：

验证结果：

剩余风险：

下一步：

提交：

### Task 15 [ ] 修缮 AI 人物聊天交互

完成 ChatPage 的人物上下文、消息列表、发送、加载、失败提示和合规说明。

独立验证：运行前端 lint/build；使用 Browser 验证城市人物聊天成功与 provider 降级错误态。

完成内容：

验证结果：

剩余风险：

下一步：

提交：

### Checkpoint E [ ] 完成 Task 13~15 后全面检查-debug循环

检查需求偏离、地图交互、动画、聊天 UI、API 联调、前端 lint/build、后端测试、安全、文档同步和回滚策略；发现问题必须在本检查点内修复并复测。

检查结果：

修复内容：

验证结果：

剩余风险：

提交：

### Task 16 [ ] 修缮赛博打卡与上传交互

完成照片选择、上传校验提示、生图预览、打卡提交、处理中状态和失败恢复。

独立验证：运行前端 lint/build；使用 Browser 验证有效图片、无效格式、超限、provider 降级和成功打卡流程。

完成内容：

验证结果：

剩余风险：

下一步：

提交：

### Task 17 [ ] 修缮成就墙与解锁反馈

完成成就墙、已解锁/未解锁展示以及打卡后的解锁反馈，区分通用和游戏专属成就。

独立验证：运行前端 lint/build；使用 Browser 分别走自由和游戏打卡，验证通用及游戏专属成就反馈。

完成内容：

验证结果：

剩余风险：

下一步：

提交：

### Task 18 [ ] 完成前端细节、可访问性与响应式修缮

检查所有页面的 loading、empty、error、禁用态、键盘可操作性、标签、移动端布局和视觉一致性。

独立验证：运行前端 lint/build；使用 Browser 在桌面和移动视口做主页面巡检并记录证据。

完成内容：

验证结果：

剩余风险：

下一步：

提交：

### Checkpoint F [ ] 完成 Task 16~18 后全面检查-debug循环

检查需求偏离、端到端 UI/UX、可访问性、响应式、前后端联调、前端 lint/build、后端测试、安全、文档同步和回滚策略；发现问题必须在本检查点内修复并复测。

检查结果：

修复内容：

验证结果：

剩余风险：

提交：

### Task 19 [ ] 修缮 Docker Compose、环境模板与自动 seed

对齐 Dockerfile、Compose、`.env.example`、健康检查、依赖启动顺序、MySQL 低资源参数、uploads 持久化和自动导入。

独立验证：运行 `docker compose config`、构建镜像、启动完整栈、检查健康状态和日志；确认无需手改数据库。

完成内容：

验证结果：

剩余风险：

下一步：

提交：

### Task 20 [ ] 执行完整 API 联调

在 Compose 栈上按契约依次走通匿名用户、城市列表详情、自由访问、游戏 init/roll、chat、生图、打卡和成就墙，核验错误路径。

独立验证：保存命令和响应摘要；确认统一包装、状态码、数据记录和降级路径符合契约。

完成内容：

验证结果：

剩余风险：

下一步：

提交：

### Task 21 [ ] 执行 Browser 主链路交互测试

使用 in-app Browser 从首页完整走通自由探索和游戏互动两条链路，覆盖城市详情、聊天、上传、生图、打卡、成就墙与刷新恢复。

独立验证：记录 20 条验收条件逐条 UI 证据和发现的问题；在本 task 内修复并复测交互问题。

完成内容：

验证结果：

剩余风险：

下一步：

提交：

### Checkpoint G [ ] 完成 Task 19~21 后全面检查-debug循环

检查需求偏离、Compose 部署、容器日志、API 联调、完整 UI/UX、类型检查、构建、测试、安全、数据一致性、文档同步和回滚策略；发现问题必须在本检查点内修复并复测。

检查结果：

修复内容：

验证结果：

剩余风险：

提交：

### Task 22 [ ] 同步 README、部署与演示文档

根据实际实现同步 README、环境变量、启动说明、AI 降级说明、演示路径、健康检查和常见故障排查。

独立验证：从干净配置视角逐条执行文档命令；检查文档无真实密钥、无过期接口和无手改数据库步骤。

完成内容：

验证结果：

剩余风险：

下一步：

提交：

### Task 23 [ ] 执行安全与合规专项审计

检查密钥、上传、路径、CORS、限流、输入校验、匿名用户资源、日志泄漏、AI prompt、免责声明、生图限制、依赖和生成文件忽略规则。

独立验证：运行密钥扫描、静态搜索、恶意输入 API 测试和上传边界测试；修复并复测全部高风险问题。

完成内容：

验证结果：

剩余风险：

下一步：

提交：

### Task 24 [ ] 执行交付前构建与回归测试

运行完整静态检查、测试、镜像构建、Compose 启动、API 冒烟、Browser 双模式回归和验收清单预审。

独立验证：记录每条命令结果和验收标准 1~20 当前证据；不得以局部检查代替完整回归。

完成内容：

验证结果：

剩余风险：

下一步：

提交：

### Checkpoint H [ ] 完成 Task 22~24 后全面检查-debug循环

检查需求偏离、文档、安全、Compose 部署、容器日志、API、完整 UI/UX、类型检查、构建、测试、数据一致性和回滚策略；发现问题必须在本检查点内修复并复测。

检查结果：

修复内容：

验证结果：

剩余风险：

提交：

### Final Review [ ] 最终全量审查、修缮与交付证明

从 C 端体验、产品需求、20 条验收、API 契约、代码质量、安全合规、匿名权限、错误处理、AI 降级、数据一致性、seed、2C2G 部署、日志、测试、构建、文档和回滚逐项审查。发现问题必须在本阶段修复并重复全量验证，直到没有已知高风险问题。

独立验证：完成逐项交付证明；运行最终全量命令；确认可直接部署；将 goal 标记完成。

完成内容：

验证结果：

剩余风险：

提交：

## Task 1 基线缺口矩阵记录区

### A. 验收标准 1~20 基线

状态含义：`失败` 表示当前证据直接反驳可验收；`未证实` 表示存在骨架但没有足够运行证据。

| # | 基线状态 | 证据与缺口 | 后续 task |
|---|---|---|---|
| 1 | 失败 | `HomePage` 文件存在，但 `main.tsx` 直接挂载 `App`，未使用 `RouterProvider`；前端 build 失败。 | 12 |
| 2 | 失败 | `ModeSelectPage` 有双入口源码，但路由未挂载，运行不可达。 | 12 |
| 3 | 失败 | `/explore` 路由声明存在但未挂载；页面还是本地列表演示。 | 12, 13 |
| 4 | 失败 | `MapCanvas` 有高德 marker 骨架，但自由探索页不使用真实地图；marker 点击仅 `flyTo`，没有访问 API。 | 13 |
| 5 | 失败 | 前端没有 `POST /visits/free` 封装或调用；城市详情读取 Mock。 | 11, 13 |
| 6 | 失败 | `/game` 路由声明存在但未挂载；页面未初始化真实游戏状态。 | 12, 14 |
| 7 | 失败 | 前端无浏览器定位和 `POST /game/init` 调用；仅硬编码北京坐标。 | 11, 14 |
| 8 | 未证实 | 有掷骰按钮 UI，但调用本地随机 Mock，不是后端 roll。 | 14 |
| 9 | 未证实 | 后端 geo 骨架含八方向和六档距离，但没有测试，且后端整体无法编译。 | 4 |
| 10 | 失败 | 后端 roll 骨架存在；前端仍随机抽取本地 Mock 城市。 | 5, 14 |
| 11 | 失败 | `useMapStore.flyTo` 有视角变化骨架，但游戏主链路未联通地图移动过程。 | 14 |
| 12 | 失败 | 后端详情聚合骨架存在；当前可见城市页只展示 Mock 简介，不展示契约要求的完整内容。 | 3, 13 |
| 13 | 失败 | `ChatPage` 为静态 UI，没有 `POST /chat` 调用。 | 7, 11, 15 |
| 14 | 失败 | `CheckinPage` 没有真实文件选择或生图调用，仅 `setTimeout`。 | 8, 11, 16 |
| 15 | 失败 | `CheckinPage` 仅 alert 模拟成功，没有 `POST /checkin`。 | 9, 11, 16 |
| 16 | 未证实 | 成就引擎源码和 `first_checkin` seed 存在，但无测试、无可运行链路。 | 6, 9, 17 |
| 17 | 未证实 | 游戏专属成就 seed 和规则源码存在，但无测试、无可运行链路。 | 6, 9, 17 |
| 18 | 未证实 | 后端 free/game visit 写入骨架存在，但前端没有真实调用，service 校验不足。 | 3, 5 |
| 19 | 失败 | `backend/Dockerfile` 缺失；Compose 偏离设计，移除了内置 MySQL；本机也无 Docker 命令。 | 19 |
| 20 | 失败 | 启动 seed 骨架存在，但后端不可编译、容器不可构建；导入失败仅 warn，不能证明无需手工干预。 | 2, 10, 19 |

### B. API 1~10 基线

| API | 基线状态 | 主要缺口 | 后续 task |
|---|---|---|---|
| `POST /api/users/anonymous` | 部分骨架 | handler/service/repo 存在；后端不编译；匿名 ID 缺少格式/长度边界测试。 | 3 |
| `GET /api/cities` | 部分骨架 | 聚合源码存在，但关联查询错误被忽略，未运行验证。 | 3 |
| `GET /api/cities/{id}` | 部分骨架 | 聚合源码存在，但关联查询错误被忽略，错误分类依赖字符串比较。 | 3 |
| `POST /api/visits/free` | 部分骨架 | 未校验 user、city、source 枚举；前端无调用。 | 3, 11 |
| `POST /api/game/init` | 部分骨架 | 未校验 user 和坐标范围；前端无调用。 | 5, 11 |
| `POST /api/game/roll` | 部分骨架 | 未校验 user、from city 和坐标；事务骨架存在但无测试；前端仍 Mock。 | 5, 11 |
| `POST /api/chat` | 部分骨架 | prompt 合规文本存在；未校验 user、消息长度、人物归属城市；历史消息未传 LLM；落库错误被吞；前端静态。 | 7, 11, 15 |
| `POST /api/checkin/generate-image` | 失败 | 上传后传给 AI client 的是 URL 路径而非文件系统路径；seed 参考图 URL 也没有映射为本地路径；静态资源目录缺失；仅扩展名校验。 | 8, 16 |
| `POST /api/checkin` | 部分骨架 | 未校验 user/city/landmark/visit/图片引用；未设置 `checkin_mode`；打卡与成就解锁不在同一事务。 | 9 |
| `GET /api/users/{id}/achievements` | 部分骨架 | 墙数据源码存在；未校验 user；统计查询错误被忽略；前端没有 API 封装。 | 6, 11, 17 |

### C. 数据、部署、质量与安全基线

| 领域 | 已确认事实 | 缺口与风险 | 后续 task |
|---|---|---|---|
| 数据内容 | JSON 解析得到 12 城、每城 2 地标/2 美食/1 人物、11 成就；schema 有 12 表。 | schema 对部分 `*_id` 缺单列索引；model/schema/seed 仍需严格对照；静态图片文件缺失。 | 2 |
| 自动 seed | `main.go` 有 AutoMigrate 和空表 seed 骨架。 | seed 使用相对路径；错误仅 warn；部分写入不检查错误、不做事务；Makefile 的 `-seed` 子命令不存在。 | 2, 10, 19 |
| 后端质量 | 已有 api -> service -> repo 大体分层和 AI 包边界。 | 编译失败；4 个 gofmt 文件；无测试；多个 handler 向前端泄漏内部 `err.Error()`；多处 DB 错误被吞。 | 3~10, 23 |
| 前端架构 | React/Vite/Zustand/高德依赖存在；高德构造集中在 `MapCanvas`。 | `router.tsx` 未挂载；只有 `api/client.ts`，缺 6 个业务 API 文件；城市、骰子、聊天、打卡、成就大量 Mock；期望组件多数缺失；build/lint 失败。 | 11~18 |
| 上传安全 | UUID 文件名和扩展名白名单骨架存在。 | 只信任文件名扩展名和 multipart header 大小；没有内容嗅探/流式上限；URL 与物理路径混淆；生成图引用未约束；限流桶 map 无清理。 | 8, 10, 23 |
| AI 合规 | AI Key 从后端 env 读取；聊天 prompt 含“不声称复活、不编史、方言解释、150 字”。 | 对话输入限制、历史上下文、稳定降级、生图合规提示和 provider mock 测试缺失。 | 7, 8, 23 |
| 密钥隔离 | 静态扫描未发现前端 LLM/IMAGE Key；`.env` 被 ignore。 | `.env.example` 仅占位符可保留；最终仍需全量密钥扫描。 | 23 |
| 部署 | Caddyfile 和前端 Dockerfile 存在。 | `backend/Dockerfile` 缺失；Compose 不含 MySQL、healthcheck 或 mysql 依赖；README/Makefile 与 Compose 矛盾；本机无 Docker。 | 19, 22 |
| 文档 | 核心产品/设计/架构文档存在。 | 多个被引用文档缺失，README 仍写“后端待开发”并包含手工 migrate/seed 步骤。 | 22 |

### D. 工作树边界

- 初始化前已有修改：`backend/data/seed/cities.json`、`backend/go.mod`、`docker-compose.yml`。
- 初始化前已有未跟踪实现：`AGENTS.md`、`backend/cmd/server/main.go`、`backend/data/seed/achievements.json`、`backend/go.sum`、`backend/internal/**`。
- Task 1 新增并维护：`goal-1/input.md`、`goal-1/plan.md`、`goal-1/tasks.md`。
- Task 1 基线构建临时生成并已清理：`frontend/tsconfig.node.tsbuildinfo`、`frontend/tsconfig.tsbuildinfo`、`frontend/vite.config.d.ts`、`frontend/vite.config.js`。

## 最终验收证据记录区

> Final Review 执行时逐条填写 `docs/product/acceptance-criteria.md` 1~20 的最终证据。
