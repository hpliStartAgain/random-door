# Random Door MVP 可部署交付计划

## 1. 目标

在现有 `random-door` 单体仓库中完成产品文档和设计文档要求的黑客松 MVP v0.1，使其能够通过 Docker Compose 在 2C2G 单机环境部署，并走通：

`首页 -> 模式选择 -> 自由探索或游戏互动 -> 城市详情 -> AI 人物对话 -> 赛博打卡图 -> 打卡 -> 成就解锁`

验收以 `docs/product/acceptance-criteria.md` 的 20 条清单、`docs/design/api-contract.md` 的接口契约和仓库 `AGENTS.md` 为准。

## 2. 当前上下文

- 仓库已经存在 Go 后端、React 前端、MySQL schema、seed、Docker Compose 和设计文档。
- 初始化时工作树并非干净状态：已有大量后端文件、`AGENTS.md` 以及部分配置改动。后续必须将这些视为现存实现，逐项审计，不得粗暴回滚。
- `frontend/node_modules` 已存在，会在文件检索和 diff 中排除。
- 当前仅完成 goal 初始化；尚未修改任何产品代码。

## 3. 范围与约束

- 保持 React + Vite + TypeScript + Tailwind CSS + shadcn/ui + Zustand + 高德地图 JS API。
- 保持 Go + Gin + GORM 单体后端、MySQL 8 / MariaDB 和 Docker Compose 单机部署。
- 不引入 Redis、Kafka、Elasticsearch、向量库、微服务网关或本地 AI 模型。
- AI API Key 只能由后端环境变量读取；前端不能直连 AI API。
- 内容范围固定为 12 个精选城市，每城 1~2 地标、1~2 美食、1 人物。
- 任何 API 变动先同步 `docs/design/api-contract.md`。
- 不提交 `.env`、真实密钥、uploads 生成图或二进制产物。
- 每轮只执行 `tasks.md` 中第一个未完成 task。每完成三个普通 task，执行一个全面检查-debug循环。

## 4. 执行方案

1. 建立基线：阅读全部约束和设计文档，盘点现有实现与缺口，确认工作树来源边界。
2. 先使后端数据层、核心 API、地理算法、成就引擎、上传和 AI 编排达到契约要求，并补足自动测试。
3. 再使前端 API 封装、状态管理、路由、双模式地图、城市详情、聊天、打卡和成就墙完整联通。
4. 完成 12 城 seed、数据库自动导入、环境模板、健康检查和 Docker Compose 一键部署。
5. 启动完整栈，做 API 联调与浏览器主链路交互测试。
6. 最后进行面向 C 端、代码、安全、数据一致性、权限、错误处理、构建、测试、文档和回滚的全量审查。

## 5. 验证方式

- 静态：`git diff --check`、密钥扫描、API 契约对照、seed 数量和字段校验。
- 后端：`go test ./...`、`go vet ./...`，必要时增加针对 geo、achievement、service、handler 的测试。
- 前端：`npm run lint`、`npm run build`，必要时增加可重复的测试脚本。
- 数据库：schema 与 GORM model 对照；容器启动后确认自动 seed 12 城和成就数据。
- 部署：`docker compose config`、`docker compose build`、`docker compose up -d`、健康检查、容器日志检查。
- 联调：按 API 契约验证匿名用户、城市、自由访问、游戏 init/roll、chat、generate-image、checkin、achievement。
- UI/UX：使用 in-app Browser 逐步验证首页、模式选择、地图、动画、详情、聊天、上传、打卡、成就墙和错误态。
- 最终：逐条记录 `docs/product/acceptance-criteria.md` 20 条的证据。

## 6. 风险与应对

| 风险 | 应对 |
|---|---|
| 初始工作树已有未提交改动 | 先审计后改动；避免重置；每轮提交本轮明确变更并记录边界 |
| 外部 LLM、生图或高德 Key 在本机不可用 | 不写入真实 Key；实现可配置 provider、超时和明确错误/演示兜底；真实 provider 以 mock 或可控替身验证编排 |
| Docker 或数据库服务不可用 | 先完成离线构建和配置校验；记录阻断证据；继续推进不依赖外部状态的任务 |
| 2C2G 资源紧张 | 保持单体、限制上传大小、限制日志和连接池、不引入额外基础设施 |
| seed 与 schema 不一致 | 以 schema 和设计文档为准，增加自动导入与幂等校验 |
| 前后端契约漂移 | 所有接口实现和 TS 类型逐项对照 `api-contract.md`；变更先更新契约 |
| AI 内容合规风险 | 强制角色声明、敏感输入处理、历史不确定性提示、生图限制和后端调用边界 |

## 7. 回滚方案

- 所有实现按 task 小步提交；发现回归时优先追加修复提交。
- 不对初始化时已有变更执行 `git reset --hard` 或 `git checkout --`。
- schema 变更必须保持可重复初始化；如需迁移，补充明确回滚 SQL 或重建开发数据库说明。
- Docker 配置修改前后均运行 `docker compose config`，失败时恢复到最近可启动提交。
- 外部 AI 不可用时保留清晰降级路径，避免主链路整体不可演示。

## 8. 默认假设

- 这是黑客松 MVP，不实现完整注册登录、支付、后台、好友、排行榜或多人实时功能。
- 匿名用户机制是 MVP 的权限边界；所有涉及用户资源的 API 仍需校验用户 ID 和输入。
- 本地没有真实外部 AI Key 时，允许使用受控演示降级结果保证主链路可验证，但生产配置仍支持真实后端 provider。
- 高德地图浏览器 Key 可能由部署者配置；无 Key 或 SDK 加载失败时 UI 应给出可理解状态，不能泄漏任何后端密钥。
- goal 控制文件用于无人值守执行记录，不属于产品运行时。
