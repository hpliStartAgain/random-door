# 完整目录结构总览（directory-structure.md）

> 汇总全仓库要创建的目录与文件，供 coding agent 一眼看清归属。逐文件职责见 backend/frontend 详设。

## 1. 仓库根
```text
random-door/
├── CLAUDE.md                 # agent 总约束（首读）
├── .cursorrules              # Cursor 规则
├── README.md                 # 项目门面
├── Makefile                  # 常用命令
├── .gitignore
├── .env.example              # 环境变量模板
├── docker-compose.yml        # 容器编排
├── backend/                  # Go 后端
├── frontend/                 # React 前端
└── docs/                     # 文档
```

## 2. backend/
```text
backend/
├── Dockerfile
├── go.mod / go.sum
├── config/config.yaml        # 非敏感默认配置（可选）
├── cmd/server/main.go        # 入口
├── internal/
│   ├── api/                  # router + 6 handler
│   │   ├── router.go
│   │   ├── game_handler.go   city_handler.go   visit_handler.go
│   │   ├── chat_handler.go   checkin_handler.go achievement_handler.go
│   ├── service/              # 6 业务编排
│   │   ├── game_service.go   city_service.go   visit_service.go
│   │   ├── chat_service.go   checkin_service.go achievement_service.go
│   ├── repository/           # 7 数据访问
│   │   ├── user_repo.go  city_repo.go  visit_repo.go  dice_repo.go
│   │   ├── chat_repo.go  checkin_repo.go  achievement_repo.go
│   ├── model/                # 12 数据模型
│   │   ├── user.go city.go city_tag.go landmark.go food.go character.go
│   │   ├── visit.go dice_roll.go checkin.go achievement.go
│   │   ├── user_achievement.go chat_message.go
│   ├── geo/                  # distance.go bearing.go target_point.go city_matcher.go
│   ├── ai/                   # llm_client.go image_client.go prompt_builder.go
│   ├── upload/               # validator.go storage.go
│   ├── seed/                 # seed.go（校验 + 事务化幂等 upsert）
│   ├── achievement/          # evaluator.go rules.go
│   ├── middleware/           # cors.go logger.go recover.go rate_limit.go
│   └── config/config.go
├── data/seed/                # cities.json  achievements.json
└── migrations/schema.sql
```

## 3. frontend/
```text
frontend/
├── Dockerfile
├── package.json / vite.config.ts / tailwind.config.js / tsconfig.json
├── index.html
└── src/
    ├── main.tsx  App.tsx  router.tsx
    ├── pages/        # 8 页
    │   HomePage  ModeSelectPage  FreeExplorePage  GameModePage
    │   CityPage   ChatPage        CheckinPage      AchievementPage
    ├── components/   # 13 组件
    │   MapCanvas  CityMarker  DicePanel  CityCard  LandmarkCard
    │   FoodCard   CharacterCard DialectCard ChatBox ImageUploader
    │   CheckinResult  BadgeWall  ModeEntryCard
    ├── api/          # client.ts city.ts game.ts visit.ts chat.ts checkin.ts achievement.ts
    └── store/        # useUserStore.ts useGameStore.ts useCityStore.ts
```

## 4. docs/
```text
docs/
├── design/   00-index / api-contract / backend-detailed-design /
│             frontend-detailed-design / database-detailed-design /
│             geo-algorithm / ai-orchestration / achievement-engine
├── agent/    go-backend-rules / react-frontend-rules / sql-rules /
│             ai-integration-rules / git-workflow / doc-writing
├── arch/     system-architecture / data-model-er / tech-stack-decision /
│             deployment-architecture / security-compliance /
│             observability / directory-structure
└── product/  prd / user-flows / product-structure / achievement-design /
              city-content-plan / ai-character-design / cyber-checkin-design /
              acceptance-criteria / roadmap
```

## 5. 运行时目录（不入库，.gitignore 忽略生成图）
```text
static/
  landmarks/  foods/  characters/  badges/
uploads/
  selfies/    generated/
```

## 6. 文件统计速览
- 后端：约 50 个 .go + 2 seed + 1 sql
- 前端：8 页 + 13 组件 + 7 api + 3 store + 入口配置
- 文档：design 8 + agent 6 + arch 7 + product 9
- 配置/部署：compose / 2×Dockerfile / .env.example / go.mod / package.json 等
