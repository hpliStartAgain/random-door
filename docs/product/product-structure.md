# 产品结构与信息架构（product-structure.md）

> 对应概要设计 8 章。说明页面层级、双模式关系与统一抽象。

## 1. 总体产品结构
```text
首页
  ↓
模式选择页
  ├── 自由探索模式
  │     ↓ 城市地图 → 点击城市 → 城市详情页
  └── 游戏互动模式
        ↓ 获取当前位置 → 掷骰子 → 计算目标城市 → 城市详情页

城市详情页
  ├── 城市简介
  ├── 地标景点
  ├── 特色美食
  ├── 代表人物
  ├── 方言样例
  ├── AI 对话
  ├── 赛博打卡
  └── 成就解锁
```

## 2. 页面层级与路由
| 层级 | 页面 | 路由 |
|---|---|---|
| L0 | 首页 HomePage | / |
| L1 | 模式选择 ModeSelectPage | /mode |
| L2a | 自由探索 FreeExplorePage | /explore |
| L2b | 游戏互动 GameModePage | /game |
| L3 | 城市详情 CityPage（两模式共用） | /city/:id |
| L4a | 人物对话 ChatPage | /city/:id/chat/:cid |
| L4b | 赛博打卡 CheckinPage | /city/:id/checkin |
| L4c | 成就墙 AchievementPage | /achievements |

## 3. 双模式关系（核心理念）
自由探索与游戏互动**不是两套独立系统**，而是两种"到达城市"的方式：
```text
自由探索：用户主动选择城市
游戏互动：系统随机生成目标城市
       ↓ 二者最终都进入 ↓
城市详情页 → AI 对话 → 赛博打卡 → 成就系统
```

## 4. 统一抽象：城市访问 Visit
为复用城市探索内核、降低双模式开发成本，系统抽象「城市访问 Visit」：
- 记录用户进入某城的行为；
- 用 visit_mode(free/game) 标识来源；
- 游戏模式额外关联 from_city_id 与 dice_roll_id。

数据建模见 ../arch/data-model-er.md，表设计见 ../design/database-detailed-design.md。

## 5. 跳转关系图
```text
Home →(进入)→ ModeSelect
ModeSelect →(自由探索)→ FreeExplore →(点城市)→ City
ModeSelect →(游戏互动)→ GameMode →(掷骰到城)→ City
City →(选人物)→ Chat
City →(选景点)→ Checkin
City/任意 →(查看成就)→ Achievement
```