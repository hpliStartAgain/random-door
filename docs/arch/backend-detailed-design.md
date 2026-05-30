# 后端详细设计（backend-detailed-design.md）

> 对应概要设计 12 章。这是 agent 创建所有 `.go` 文件的直接依据。约束见 go-backend-rules.md。
> 分层：api(handler) → service → repository；geo/ai/upload/achievement 为领域包。
> 每个文件给出：职责 / 关键导出 / 主逻辑 / 输入输出 / 错误处理 / 依赖。

## 0. 完整后端文件树
```
backend/
  cmd/server/main.go
  internal/
    api/        router.go game_handler.go city_handler.go visit_handler.go
                chat_handler.go checkin_handler.go achievement_handler.go
    service/    game_service.go city_service.go visit_service.go
                chat_service.go checkin_service.go achievement_service.go
    repository/ user_repo.go city_repo.go visit_repo.go dice_repo.go
                chat_repo.go checkin_repo.go achievement_repo.go
    model/      user.go city.go city_tag.go landmark.go food.go character.go
                visit.go dice_roll.go checkin.go achievement.go
                user_achievement.go chat_message.go
    geo/        distance.go bearing.go target_point.go city_matcher.go
    ai/         llm_client.go image_client.go prompt_builder.go
    upload/     validator.go storage.go
    seed/       seed.go
    achievement/ evaluator.go rules.go
    middleware/ cors.go logger.go recover.go rate_limit.go
    config/     config.go
  data/seed/    cities.json achievements.json
  migrations/   schema.sql
```

---

## 1. cmd/server/main.go
- 职责：程序入口，组装与启动。
- 主逻辑：
  1. `config.Load()` 读取环境变量/配置；
  2. 初始化 GORM（连接 MySQL，配置连接池 max_open=50）；
  3. 执行 AutoMigrate；通过 `internal/seed` 校验 `SEED_DIR` 下的 JSON，并在事务内幂等 upsert 内容；
  4. 构造各 repository → service → handler 依赖（手动注入）；
  5. `api.NewRouter(...)` 装配 Gin 引擎与中间件；
  6. 监听 `SERVER_PORT` 启动。
- 错误：启动期任何失败直接 `log.Fatal`。

---

## 2. internal/api（HTTP 层，禁止业务逻辑）

### 2.1 router.go
- 职责：注册全部路由 + 全局中间件。
- 导出：`func NewRouter(handlers ...) *gin.Engine`。
- 主逻辑：use(cors, logger, recover, rate_limit)；按 api-contract.md 注册 10 条路由到对应 handler；挂载 /static 与 /uploads 静态目录。

### 2.2 各 handler 通用模式
每个 handler 持有对应 service，方法签名 `func(c *gin.Context)`。统一流程：解析/校验参数 → 调 service → 成功返回 JSON / 失败用统一错误响应（错误码见 api-contract.md 0.4）。

| 文件 | 处理接口 | 关键方法 |
|---|---|---|
| city_handler.go | GET /cities、GET /cities/{id} | List、Detail |
| visit_handler.go | POST /visits/free | CreateFree |
| game_handler.go | POST /game/init、POST /game/roll | Init、Roll |
| chat_handler.go | POST /chat | Chat |
| checkin_handler.go | POST /checkin/generate-image、POST /checkin | GenerateImage、Create |
| achievement_handler.go | GET /users/{id}/achievements | Wall |

注：用户初始化 POST /users/anonymous 可放在 visit_handler 或单独 user_handler；本设计放 visit_handler.go 的 CreateAnonymousUser（避免新增文件）。
- checkin_handler.GenerateImage：读取 multipart 表单，调 upload.validator 校验，再交 service。

---

## 3. internal/service（业务编排，事务边界）

### 3.1 city_service.go
- City 列表/详情聚合：List() 返回城市+tags；Detail(id) 聚合 cities+tags+landmarks+foods+characters（**剔除 persona/prompt**）。依赖 city_repo。

### 3.2 visit_service.go
- CreateAnonymousUser(anonymousID)：查/建 user。
- CreateFreeVisit(userID, cityID, source)：写 city_visits(mode=free)。依赖 user_repo、visit_repo。

### 3.3 game_service.go（核心）
- Init(userID, lat, lng)：用 geo.city_matcher 找最近城市作起点（lat/lng 缺省用默认北京）。
- Roll(userID, fromCityID, lat, lng)：
  1. geo.RandomDirection() → 方向+角度；
  2. geo.RandomDistance() → 距离档；
  3. geo.TargetPoint(lat,lng,bearing,dist) → 目标经纬度；
  4. geo.MatchNearestCity(target, exclude=from, preferUnvisited) → 目标城市（含兜底）；
  5. 事务：写 dice_rolls，写 city_visits(mode=game, source=dice_roll, from_city_id, dice_roll_id)；
  6. 返回 visit_id/dice_roll_id/direction/distance/target_point/target_city。
- 依赖 geo、dice_repo、visit_repo、city_repo。

### 3.4 chat_service.go
- Chat(userID, cityID, characterID, message)：
  1. 读 city/character/landmarks/foods/dialect；character 不存在→NOT_FOUND；
  2. ai.prompt_builder 组装；
  3. ai.llm_client 调用（超时/重试，失败→AI_UPSTREAM_ERROR/AI_TIMEOUT）；
  4. 落库 user 与 assistant 两条 chat_messages；
  5. 返回 reply。

### 3.5 checkin_service.go
- GenerateImage(userID,cityID,landmarkID,file)：upload.storage 落 selfie → 取 landmark.image_url 参考图 → ai.image_client 生图 → 落 generated → 返回 url。**不写库**。
- Create(req)：事务内写 checkins → 调 achievement_service.Evaluate → 返回 checkin_id + 新解锁成就。

### 3.6 achievement_service.go
- Evaluate(userID)：调 achievement.evaluator 计算新解锁，写 user_achievements（去重），返回新解锁列表。
- Wall(userID)：拆分 unlocked/locked + 计算 progress。

---

## 4. internal/repository（仅 GORM 读写，无业务判断）
| 文件 | 主要方法 |
|---|---|
| user_repo.go | FindByAnonymousID、Create、UpdateCurrentCity |
| city_repo.go | ListAll、FindByID、ListTags/Landmarks/Foods/Characters(byCityID) |
| visit_repo.go | Create、CountByUser、CountByUserMode、ListByUser |
| dice_repo.go | Create、ListRecentByUser(用于"连续方向"成就) |
| chat_repo.go | Create、ListByUserCharacter |
| checkin_repo.go | Create、ListByUser、CountByUser、CountByUserCityTag |
| achievement_repo.go | ListAll、ListUserAchievements、CreateUserAchievement(去重) |

所有方法接收 context，返回 (data, error)。

---

## 5. internal/model（结构体 + GORM tag，对应 15 章 12 表）
每文件一个 struct，字段与 database-detailed-design.md 一致，带 `gorm:"column:...;index"` 与 `json:"..."`。created_at/updated_at 用 autoCreateTime/autoUpdateTime。
文件：user.go city.go city_tag.go landmark.go food.go character.go visit.go(CityVisit) dice_roll.go checkin.go achievement.go user_achievement.go chat_message.go。
注意：Character 的 Persona/Prompt 字段 json tag 设为 `json:"-"`（不下发）。

---

## 6. internal/geo（详见 geo-algorithm-detailed-design.md）
- distance.go：Haversine(lat1,lng1,lat2,lng2) float64 km。
- bearing.go：8 方向枚举 + RandomDirection() (name,deg)。
- target_point.go：TargetPoint(lat,lng,bearingDeg,distKm) (lat2,lng2)；RandomDistance() int。
- city_matcher.go：MatchNearestCity(cities, target, opts) 含过滤当前/优先未访问/兜底 4 步。

## 7. internal/ai（详见 ai-orchestration-detailed-design.md）
- llm_client.go：Chat(ctx, systemPrompt, userMsg) (string,error)，超时+重试。
- image_client.go：Generate(ctx, selfiePath, refImagePath, prompt) (string,error)。
- prompt_builder.go：BuildChatPrompt(ctx数据)、BuildImagePrompt(city,landmark)。

## 8. internal/upload
- validator.go：校验类型(jpg/jpeg/png/webp)与大小(≤5MB)，否则 415/413。
- storage.go：UUID 命名，分目录落 uploads/selfies、uploads/generated；防路径穿越。

## 8.1 internal/seed
- seed.go：读取 `SEED_DIR` 下的 cities.json / achievements.json；校验 12 城内容、成就可达性与 AI Prompt 合规提醒；事务内按自然键幂等 upsert。

## 9. internal/achievement（详见 achievement-engine-detailed-design.md）
- rules.go：rule_type 解析器（checkin_count/city_tag/tag_count/game_visit_count/dice_direction/dice_distance/first_checkin）。
- evaluator.go：Evaluate(userID, repos) 返回满足且未解锁的成就。

## 10. internal/middleware
- cors.go：按 CORS_ALLOW_ORIGINS 白名单。
- logger.go：结构化请求日志（路径/方法/状态/耗时）。
- recover.go：panic 兜底 → 500 INTERNAL_ERROR。
- rate_limit.go：上传/AI 接口限流（简单内存令牌桶）。

## 11. internal/config/config.go
- Load() 读 env（Viper）：DB_*、SERVER_PORT、SEED_DIR、LLM_*、IMAGE_*、UPLOAD_*、CORS_*、AI_TIMEOUT_SECONDS。敏感项禁硬编码。
