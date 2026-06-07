# 总体系统架构

## 1. 架构总览

任意门采用前后端分离 + Go 单体后端架构。前端构建为静态资源，由 Caddy 托管并反向代理 `/api`；后端是单一 Gin 进程，使用 GORM 连接外部 MySQL，AI 能力通过外部 API 调用，上传和生成文件落本地卷。

```text
Browser
  React / Vite / TypeScript / Tailwind / Zustand / AMap
        |
        | HTTP
        v
Caddy
  static frontend dist
  /api -> app:8080
  /static and /uploads file serving
        |
        v
Go app
  Gin router
  api -> service -> repository
  geo / ai / upload / achievement / seed
        |                      |
        v                      v
External MySQL           External AI APIs
  users/cities/...         LLM / image generation
        |
        v
Local volumes
  static assets / uploads
```

## 2. 前端职责

- 首屏欢迎、自由探索、任意门随机漫游入口。
- 高德地图加载、城市/地标 marker、飞行动画。
- 城市列表搜索、区域筛选、标签筛选。
- 城市详情、风光浏览、声景控制、评论、AI 对话、打卡流程。
- 成就墙、资产页、个人资料面板、猜城市挑战页。
- 后台 CMS 界面。

## 3. 后端职责

- 用户匿名身份、注册/登录、资料维护。
- 城市 catalog 查询与后台维护。
- 自由访问、随机漫游、地理算法、访问记录。
- AI 对话 prompt 编排、LLM 调用、消息落库。
- 异步生图任务、上传校验、生成图保存、用量限制。
- 打卡、资产聚合、成就评估。
- 评论、猜城市挑战与答案。
- seed audit/bootstrap/sync 工具与启动时受控导入。

## 4. 典型数据流

```text
自由探索：点击城市 -> POST /visits/free -> city_visits -> GET /cities/{id}
任意门：POST /game/init -> POST /game/roll -> dice_rolls + city_visits -> 地图动画
对话：POST /chat -> prompt_builder -> LLM -> chat_messages -> reply
生图：上传自拍 -> ai_tasks queued -> worker 调 image API -> uploads/generated
打卡：POST /checkin -> checkins -> achievement evaluator -> user_achievements
猜城市：POST /guess/caption -> POST /guess/challenges -> /?guess=code -> answers
后台：ADMIN_TOKEN -> /api/admin/* -> catalog / media / achievements
```

## 5. 架构原则

- 单体优先：不拆微服务，不引入 Redis / MQ / ES / 网关。
- 数据库为内容事实源：seed 默认关闭，只用于受控初始化或维护。
- AI 后端隔离：LLM / IMAGE Key 不下发前端。
- 前端体验集中：地图、动效、抽屉、覆盖层由 React/Zustand 编排。
- 资源受控：上传限大小，AI 限每日次数，2C2G 单机可运行。

## 6. 为什么使用 Caddy

Caddy 当前负责前端静态托管、`/api` 反代、静态/上传文件暴露和健康检查。项目不引入 Nginx，是为了保持部署组件少、配置简单、资源占用低。

## 7. 文档衔接

- 部署：`deployment-architecture.md`
- 数据模型：`data-model-er.md`
- 安全合规：`security-compliance.md`
- 文件树：`directory-structure.md`
