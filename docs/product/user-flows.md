# 用户流程文档（user-flows.md）

> 对应概要设计 9 章。每条流程含步骤与对应接口（接口契约见 api-contract.md）。

## 1. 首次进入流程
```text
用户打开网页
  ↓ 前端初始化应用
  ↓ 检查 localStorage 是否已有匿名 user_id
  ↓ 没有 → 调后端创建匿名用户 (POST /api/users/anonymous)
  ↓ 进入模式选择页
  ↓ 用户选择 自由探索 或 游戏互动
```
对应接口：POST /api/users/anonymous。前端 store：useUserStore.initUser()。

## 2. 自由探索模式流程
```text
选择自由探索
  ↓ 前端加载城市地图 (GET /api/cities)
  ↓ 地图展示城市 Marker
  ↓ 用户点击目标城市
  ↓ 后端写 city_visits, visit_mode=free (POST /api/visits/free)
  ↓ 进入城市详情页 (GET /api/cities/{id})
  ↓ 查看内容 → AI 对话 / 赛博打卡
  ↓ 完成打卡 (POST /api/checkin)
  ↓ 系统判断是否解锁成就
```
对应接口：GET /api/cities、POST /api/visits/free、GET /api/cities/{id}、POST /api/checkin。

## 3. 游戏互动模式流程
```text
选择游戏互动
  ↓ 前端请求浏览器定位授权
  ↓ 获取当前位置(失败用默认北京)
  ↓ 游戏初始化 (POST /api/game/init) → 确定起点城市
  ↓ 用户点击掷骰子
  ↓ 后端生成随机方向 + 随机距离
  ↓ 后端算目标坐标点 → 匹配最近目标城市
  ↓ 写 dice_rolls + 写 city_visits(visit_mode=game) (POST /api/game/roll)
  ↓ 前端展示移动动画
  ↓ 进入目标城市详情页
  ↓ 探索 / 对话 / 打卡 / 成就解锁
```
对应接口：POST /api/game/init、POST /api/game/roll。算法见 geo-algorithm-detailed-design.md。

## 4. 城市探索流程
```text
进入城市详情页
  ↓ 城市简介 → 地标景点 → 特色美食 → 代表人物 → 方言样例
  ↓ 选择与人物对话
  ↓ 选择生成赛博游客照
  ↓ 确认打卡
  ↓ 触发成就判断
```

## 5. AI 人物对话流程
```text
选择城市人物
  ↓ 输入问题
  ↓ 后端读 城市/人物/可选方言风格和历史消息
  ↓ 组装 Prompt → 调外部 LLM API
  ↓ 保存 user 与 assistant 两条消息
  ↓ 返回 AI 回复 → 前端展示
```
对应接口：POST /api/chat。失败兜底：友好提示，不阻断浏览。

## 6. 赛博打卡流程
```text
选择景点
  ↓ 上传自拍/个人照
  ↓ 后端校验文件类型与大小 → 保存上传图
  ↓ 读景点参考图 → 调外部生图 API → 保存生成图
  ↓ 返回生成图 URL (POST /api/checkin/generate-image)
  ↓ 用户确认打卡 → 写 checkins (POST /api/checkin)
  ↓ 触发成就判断 → 返回新解锁成就
```
对应接口：POST /api/checkin/generate-image、POST /api/checkin。注意：生图接口不写库，确认打卡才落 checkins。

## 7. 流程与接口/表对照
| 流程 | 主要接口 | 主要写表 |
|---|---|---|
| 首次进入 | /users/anonymous | users |
| 自由探索 | /cities、/visits/free | city_visits |
| 游戏互动 | /game/init、/game/roll | dice_rolls、city_visits |
| AI 对话 | /chat | chat_messages |
| 赛博打卡 | /checkin/generate-image、/checkin | checkins、user_achievements |
