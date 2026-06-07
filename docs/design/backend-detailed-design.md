# 后端详细设计

## 1. 分层结构

```text
cmd/server      组装配置、数据库、迁移、seed mode、依赖注入、启动 Gin
cmd/seedtool    seed audit/bootstrap/sync/backfill 工具
internal/api    HTTP handler 与路由
internal/service 业务编排与事务边界
internal/repository GORM 数据访问
internal/model  GORM 模型
internal/geo    任意门随机漫游算法
internal/ai     LLM / 图像生成 client 与 prompt
internal/upload 文件校验与保存
internal/seed   seed 校验和受控导入
internal/achievement 成就规则与评估
internal/middleware CORS / logger / recover / rate limit
internal/config 环境变量配置
```

依赖方向保持 `api -> service -> repository`，领域包由 service 调用。

## 2. `cmd/server/main.go`

启动流程：

1. 读取 `.env` 和环境变量。
2. 连接外部 MySQL，配置连接池。
3. AutoMigrate 17 张表。
4. 按 `SEED_MODE=off|bootstrap|sync` 决定是否导入 seed，默认 `off`。
5. 手动依赖注入 repository、service、handler。
6. 启动异步生图 worker。
7. 注册 Gin router 并监听 `SERVER_PORT`。

启动期错误直接 fail fast。

## 3. API 层

| 文件 | 职责 |
|---|---|
| `router.go` | 注册中间件、静态目录、健康检查和全部 `/api` 路由。 |
| `response.go` | 统一错误响应与 service error 映射。 |
| `visit_handler.go` | 匿名用户、自由探索访问。 |
| `user_handler.go` | 注册、登录、profile 读写。 |
| `city_handler.go` | 城市列表与详情。 |
| `game_handler.go` | 任意门初始化与掷骰。 |
| `chat_handler.go` | AI 人物对话。 |
| `comment_handler.go` | 评论列表与创建。 |
| `guess_handler.go` / `guess_challenge_handler.go` | 猜城市文案、挑战和答案。 |
| `checkin_handler.go` | 生图任务、任务查询/重试、确认打卡。 |
| `achievement_handler.go` | 成就墙。 |
| `asset_handler.go` | 用户资产聚合。 |
| `admin_handler.go` | ADMIN_TOKEN 鉴权与后台 CMS。 |

Handler 只做参数解析、基础校验、调用 service、返回 JSON。

## 4. Service 层

| 文件 | 核心职责 |
|---|---|
| `city_service.go` | 聚合城市、标签、地标、美食、人物和内容计数。 |
| `visit_service.go` | 匿名用户创建/恢复，自由访问写入与成就评估。 |
| `user_service.go` | profile、注册、登录、匿名足迹保留。 |
| `game_service.go` | 初始化最近城市，随机方向/距离，目标点和最近城市匹配，事务写掷骰与访问。 |
| `chat_service.go` | 消息校验、每日限额、历史消息、prompt、LLM、消息落库。 |
| `comment_service.go` | 评论目标校验、内容校验、列表读取。 |
| `guess_service.go` | 猜城市文案生成与兜底。 |
| `guess_challenge_service.go` | data URL 保存、code 生成、挑战读取、答案判定。 |
| `checkin_service.go` | 自拍/场景任务入队、worker 生图、任务状态、确认打卡。 |
| `achievement_service.go` | Evaluate、Wall、进度计算。 |
| `asset_service.go` | visited cities、posters、achievement progress 聚合。 |
| `admin_service.go` | coverage、标签、内容、成就、图片导入与本地路径校验。 |
| `errors.go` | service error 类型。 |

## 5. Repository 层

Repository 只负责 GORM 查询/写入，不做业务判断。当前包括：

```text
admin_repo, achievement_repo, ai_task_repo, ai_usage_repo,
checkin_repo, checkin_store, chat_repo, city_repo, comment_repo,
dice_repo, game_store, guess_challenge_repo, user_repo, visit_repo
```

跨表原子写入通过 store 类 repository 封装事务，例如 `game_store`、`checkin_store`。

## 6. Model 层

每个 model 对应一张表，字段与 `database-detailed-design.md` 对齐：

```text
User, City, CityTag, Landmark, Food, Character,
CityVisit, DiceRoll, Checkin,
Achievement, UserAchievement,
ChatMessage, Comment,
AITask, AIUsageLog,
GuessChallenge, GuessAnswer
```

`Character.Persona` 与 `Character.Prompt` 不下发前端。

## 7. 领域包

| 包 | 职责 |
|---|---|
| `geo` | Haversine、方向、距离、目标点、最近城市匹配。 |
| `ai` | LLM client、image client、prompt builder、配置 guard。 |
| `upload` | 文件类型/大小校验、UUID 保存、防路径穿越。 |
| `seed` | 解析 seed、校验 12-100 城、bootstrap/sync/backfill。 |
| `achievement` | rule_type 解析、UserStats 匹配、Evaluate。 |
| `middleware` | CORS、结构化日志、panic recover、内存令牌桶限流。 |

## 8. 配置

`internal/config/config.go` 读取：

- `DB_*`
- `SERVER_PORT`、`STATIC_DIR`、`UPLOAD_DIR`
- `SEED_DIR`、`SEED_MODE`
- `LLM_*`、`IMAGE_*`
- `UPLOAD_*`
- `CORS_ALLOW_ORIGINS`
- `AI_TIMEOUT_SECONDS`、AI 每日限额、worker 配置
- `ADMIN_TOKEN`
- `LOG_LEVEL`

敏感项禁止硬编码。
