# 验收标准

## 1. 主链路验收

| # | 验收条件 | 关联接口 / 页面 |
|---|---|---|
| 1 | 用户可打开首页并看到“任意门”品牌与两个主入口。 | WelcomeOverlay |
| 2 | 首次进入自动创建或恢复匿名用户。 | POST /users/anonymous |
| 3 | 可加载城市列表，显示城市卡片、地图 marker、搜索和筛选。 | GET /cities, Sidebar, MapCanvas |
| 4 | 自由探索点击城市后写访问记录并进入详情。 | POST /visits/free, GET /cities/{id} |
| 5 | 城市详情展示简介、标签、地标、美食、人物、方言。 | CityDetailPanel |
| 6 | 地标或城市风光可进入沉浸图片视图，声景按钮不自动播放。 | StreetViewCanvas, SoundscapeControl |
| 7 | 可与城市人物发起 AI 对话，失败时不阻断浏览。 | POST /chat, RightDrawer |
| 8 | 可发表评论并读取评论列表。 | GET/POST /comments |
| 9 | 任意门模式可初始化起点城市。 | POST /game/init |
| 10 | 任意门模式可随机方向与距离，返回目标城市并写访问。 | POST /game/roll |
| 11 | 掷骰后地图展示移动/降落过程，并进入目标城市详情。 | MapCanvas, RandomCityModal |
| 12 | 可上传自拍创建异步生图任务。 | POST /checkin/generate-image |
| 13 | 可查询/重试生图任务，成功后展示生成图。 | GET/POST /checkin/image-tasks |
| 14 | 可确认打卡并生成海报。 | POST /checkin, CheckinPoster |
| 15 | 访问或打卡可触发新成就解锁，成就墙可展示进度。 | AchievementUnlock, AchievementPage |
| 16 | 我的足迹可展示访问城市、海报和成就进度。 | GET /users/{id}/assets |
| 17 | 用户可编辑匿名资料，注册/登录可保留足迹。 | profile, auth |
| 18 | 可生成猜城市文案和挑战链接，好友可提交答案。 | /guess/* |
| 19 | 后台需通过 ADMIN_TOKEN 验证后才能维护内容。 | /admin/* |
| 20 | Docker Compose 启动后健康检查通过，空库可通过受控 seed 初始化。 | docker-compose.yml, Makefile |

## 2. 工程验收

- `make test` 后端单元测试通过，或明确记录无法执行原因。
- `npm run build` 前端构建通过，或明确记录无法执行原因。
- README、API 契约、数据库设计、前后端详设与当前代码一致。
- `.env`、密钥、运行时上传/生成文件不入库。

## 3. 内容验收

- 当前 seed 可通过 seedtool 校验。
- 每座 seed 城市至少具备城市基础信息、标签、封面、地标、美食、人物、方言。
- 后台 coverage 中 `missing_fields` 应作为内容运营待办处理。

## 4. 主链路定义

```text
进入任意门
  -> 自由探索或随机漫游到城市
  -> 城市详情
  -> 风光浏览 / AI 对话 / 评论
  -> 赛博打卡
  -> 海报与成就
  -> 足迹资产沉淀
```
