# 产品结构与信息架构

## 1. 总体结构

```text
任意门 App
├── 首屏欢迎层
│   ├── 自由探索
│   ├── 任意门随机漫游
│   └── 注册 / 登录
├── 主工作台
│   ├── Navbar：品牌、我的足迹、成就墙、后台入口
│   ├── Sidebar：城市列表、搜索筛选、任意门入口、城市详情面板
│   ├── MapCanvas：高德地图、城市/地标 marker、飞行动画
│   ├── StreetViewCanvas：城市/地标风光浏览、声景、猜一猜、打卡入口
│   └── RightDrawer：人物 AI 对话
├── 覆盖页
│   ├── AchievementPage：成就墙
│   ├── AssetPage：个人资产
│   ├── AdminPage：后台 CMS
│   └── GuessChallengePage：好友猜城市挑战
└── 全局反馈
    ├── Toast
    └── AchievementUnlock
```

## 2. 前端视图状态

当前前端不是多路由页面栈，而是由 `App.tsx` 组合的单页工作台。主要视图状态来自 `useViewStore`：

| 状态 | 用途 |
|---|---|
| `HOME` | 首屏 / 默认视图。 |
| `CITY_DETAIL` | Sidebar 上方滑入城市详情。 |
| `GAME_DICE` | 地图区域显示任意门随机城市弹窗。 |
| `ACHIEVEMENT` | 全屏成就墙覆盖层。 |
| `ASSETS` | 全屏资产页覆盖层。 |
| `canvasMode=map/street` | 主画布在地图与风光浏览间切换。 |
| `profileOpen` | 右侧个人足迹面板。 |

`/?guess=<code>` 是例外入口，会直接渲染 `GuessChallengePage`。

## 3. 到达城市的两种方式

```text
自由探索：搜索 / 筛选 / 地图点击 / 城市卡片点击
任意门：定位或默认起点 -> 随机方向 -> 随机距离 -> 匹配目标城市
                         ↓
                 统一进入城市详情
                         ↓
       风光浏览 / AI 对话 / 评论 / 赛博打卡 / 成就
```

## 4. 统一抽象：CityVisit

`city_visits` 记录用户进入城市的事实，用 `visit_mode` 区分来源：

- `free`：自由探索、搜索、地图点击。
- `game`：任意门随机漫游，关联 `dice_roll_id` 与 `from_city_id`。

成就、资产和用户足迹都基于访问记录与打卡记录聚合。

## 5. 主要数据对象

| 对象 | 用途 |
|---|---|
| City | 城市基础信息、坐标、封面、方言。 |
| CityTag | 城市筛选、成就规则。 |
| Landmark / Food / Character | 城市详情内容、评论目标、打卡参考。 |
| DiceRoll | 任意门随机方向、距离、目标点和目标城市。 |
| Checkin | 打卡与海报资产。 |
| Achievement / UserAchievement | 成就定义与用户解锁记录。 |
| AITask / AIUsageLog | 异步生图任务与每日用量限制。 |
| GuessChallenge / GuessAnswer | 猜城市挑战链接与答案。 |

详细字段见 `../design/database-detailed-design.md`。
