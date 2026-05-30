# API 契约（api-contract.md）

> 前后端唯一真相源。对应概要设计 16 章。任何接口改动**必须先改本文件**。
> 通用约定见下方「0. 全局约定」。

## 0. 全局约定

### 0.1 基础路径
所有接口前缀 `/api`。后端默认端口 `8080`，前端经 vite proxy 或 Caddy 反代。

### 0.2 请求头
| 头 | 说明 |
|---|---|
| Content-Type | `application/json`（除生图接口为 `multipart/form-data`） |
| X-User-Id | 可选。前端从 store 注入匿名 user_id，便于日志追踪；以请求体 user_id 为准 |

### 0.3 统一响应包装
成功响应直接返回业务对象（见各接口）。错误统一结构：
```json
{
  "error": {
    "code": "INVALID_PARAM",
    "message": "city_id is required"
  }
}
```

### 0.4 标准错误码
| HTTP | code | 含义 |
|---|---|---|
| 400 | INVALID_PARAM | 参数缺失或格式错误 |
| 404 | NOT_FOUND | 资源不存在（城市/用户/人物） |
| 413 | FILE_TOO_LARGE | 上传文件超过 5MB |
| 415 | UNSUPPORTED_MEDIA | 文件类型不支持 |
| 502 | AI_UPSTREAM_ERROR | 外部 LLM / 生图 API 失败 |
| 504 | AI_TIMEOUT | 外部 AI 调用超时 |
| 500 | INTERNAL_ERROR | 服务器内部错误 |

### 0.5 字段类型约定
- id 类：整型（int64）
- lat / lng：浮点（double）
- 时间：ISO8601 字符串（如 `2026-05-30T09:39:23+08:00`）

---

## 1. POST /api/users/anonymous — 创建/恢复匿名用户

**请求**
```json
{ "anonymous_id": "browser_generated_uuid" }
```
| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| anonymous_id | string | 是 | 前端生成的 UUID，存 localStorage |

**响应 200**
```json
{
  "user_id": 1,
  "anonymous_id": "browser_generated_uuid",
  "current_city_id": null
}
```
**逻辑**：按 anonymous_id 查 users，存在则返回，不存在则创建。

---

## 2. GET /api/cities — 城市列表

**请求**：无参数（MVP 全量 12 城）。

**响应 200**
```json
{
  "cities": [
    {
      "id": 1,
      "name": "北京",
      "province": "北京",
      "lat": 39.9042,
      "lng": 116.4074,
      "cover_image_url": "/static/landmarks/beijing_forbidden_city.jpg",
      "tags": ["ancient_capital", "north_china"]
    }
  ]
}
```
**用途**：自由探索地图打点 + 游戏模式候选城市集合。

---

## 3. GET /api/cities/{city_id} — 城市详情

**路径参数**：city_id（int64）。

**响应 200**
```json
{
  "id": 3,
  "name": "西安",
  "province": "陕西",
  "lat": 34.3416,
  "lng": 108.9398,
  "intro": "西安是中国重要古都之一……",
  "cover_image_url": "/static/landmarks/xian.jpg",
  "dialect_sample": "嘹咋咧",
  "dialect_explanation": "陕西关中方言，表示很好、很棒。",
  "tags": ["ancient_capital", "northwest"],
  "landmarks": [
    { "id": 12, "name": "兵马俑", "image_url": "/static/landmarks/bmy.jpg", "description": "..." }
  ],
  "foods": [
    { "id": 5, "name": "肉夹馍", "image_url": "/static/foods/rjm.jpg", "description": "..." }
  ],
  "characters": [
    { "id": 8, "name": "唐代长安书生", "character_type": "history", "avatar_url": "/static/characters/c8.jpg", "dialect_style": "关中话" }
  ]
}
```
**错误**：city 不存在 → 404 NOT_FOUND。
**注意**：响应中 characters 不返回 persona / prompt（敏感，仅后端组装 Prompt 用）。

---

## 4. POST /api/visits/free — 自由探索访问城市

**请求**
```json
{ "user_id": 1, "city_id": 3, "source": "map_click" }
```
| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| user_id | int64 | 是 | |
| city_id | int64 | 是 | |
| source | string | 否 | map_click / search，默认 map_click |

**响应 200**
```json
{ "visit_id": 1001, "city_id": 3, "visit_mode": "free" }
```
**逻辑**：写 city_visits，visit_mode=free，dice_roll_id=null。

---

## 5. POST /api/game/init — 游戏模式初始化

**请求**
```json
{ "user_id": 1, "lat": 39.9042, "lng": 116.4074 }
```
lat/lng 可选；缺省时用默认位置（北京）。

**响应 200**
```json
{
  "nearest_city": {
    "id": 1, "name": "北京", "province": "北京",
    "lat": 39.9042, "lng": 116.4074
  }
}
```
**逻辑**：用 geo.city_matcher 找离 (lat,lng) 最近城市作为起点。不写 visit。

---

## 6. POST /api/game/roll — 掷骰漫游

**请求**
```json
{ "user_id": 1, "from_city_id": 1, "lat": 39.9042, "lng": 116.4074 }
```

**响应 200**
```json
{
  "visit_id": 1002,
  "dice_roll_id": 2001,
  "direction": "西南",
  "distance_km": 800,
  "target_point": { "lat": 34.8, "lng": 109.1 },
  "target_city": {
    "id": 3, "name": "西安", "province": "陕西",
    "lat": 34.3416, "lng": 108.9398
  }
}
```
**逻辑**（详见 geo-algorithm-detailed-design.md）：随机方向(8) + 随机距离(6档) → 算目标点 → 匹配最近城市(过滤当前、优先未访问、含兜底) → 写 dice_rolls → 写 city_visits(visit_mode=game, source=dice_roll, from_city_id, dice_roll_id)。

---

## 7. POST /api/chat — AI 人物对话

**请求**
```json
{
  "user_id": 1,
  "city_id": 3,
  "character_id": 8,
  "message": "你觉得西安最值得看的地方是哪里？"
}
```

**响应 200**
```json
{ "reply": "若问长安风物，那自然绕不开城墙、大雁塔与兵马俑……" }
```
**错误**：character 不存在 → 404；LLM 失败 → 502 AI_UPSTREAM_ERROR；超时 → 504 AI_TIMEOUT（前端给友好提示，不阻断浏览）。
**逻辑**：读 city/character/landmarks/foods/dialect → prompt_builder 组装 → llm_client 调用 → 落库 user 与 assistant 两条 chat_messages → 返回 reply。

---

## 8. POST /api/checkin/generate-image — 生成赛博打卡图

**Content-Type**：`multipart/form-data`

**表单字段**
| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| user_id | int64 | 是 | |
| city_id | int64 | 是 | |
| landmark_id | int64 | 是 | 景点（取参考图） |
| selfie_file | file | 是 | jpg/jpeg/png/webp，≤5MB |

**响应 200**
```json
{
  "status": "success",
  "generated_image_url": "/uploads/generated/img_123.png"
}
```
**错误**：类型不符 → 415；超过 5MB → 413；生图失败 → 502；超时 → 504。
**逻辑**：upload.validator 校验 → storage 以 UUID 落 uploads/selfies → 取景点参考图 → image_client 调外部生图 → 落 uploads/generated → 返回 url。**本接口不写 checkins**（确认后由接口 9 写）。

---

## 9. POST /api/checkin — 完成打卡

**请求**
```json
{
  "user_id": 1,
  "city_id": 3,
  "landmark_id": 12,
  "visit_id": 1002,
  "generated_image_url": "/uploads/generated/img_123.png"
}
```
| 字段 | 必填 | 说明 |
|---|---|---|
| visit_id | 否 | 关联本次访问（游戏模式有） |
| generated_image_url | 否 | 没生图也允许纯打卡 |

**响应 200**
```json
{
  "checkin_id": 3001,
  "unlocked_achievements": [
    { "code": "ancient_capital_first", "name": "古都初见", "description": "首次打卡一座古都城市" }
  ]
}
```
**逻辑**：事务内写 checkins → 调 achievement.evaluator 判定 → 写 user_achievements → 返回新解锁列表（无则空数组）。

---

## 10. GET /api/users/{user_id}/achievements — 成就墙

**路径参数**：user_id。

**响应 200**
```json
{
  "unlocked": [
    { "code": "first_checkin", "name": "初次打卡", "description": "...", "badge_url": "/static/badges/first.png", "unlocked_at": "2026-05-30T09:00:00+08:00" }
  ],
  "locked": [
    { "code": "ancient_capital_tour", "name": "古都巡礼", "description": "打卡3座古都城市", "badge_url": "/static/badges/ac3.png" }
  ],
  "progress": [
    { "code": "ancient_capital_tour", "current": 1, "target": 3 }
  ]
}
```
**逻辑**：读 achievements 全集 + user_achievements → 拆分已解锁/未解锁 → 对可量化成就计算 progress。

---

## 11. 接口与数据表对照速查
| 接口 | 主要读写表 |
|---|---|
| /users/anonymous | users |
| /cities | cities, city_tags |
| /cities/{id} | cities, city_tags, landmarks, foods, characters |
| /visits/free | city_visits |
| /game/init | cities |
| /game/roll | dice_rolls, city_visits, cities |
| /chat | chat_messages, characters, cities, landmarks, foods |
| /checkin/generate-image | (文件存储，不写库) |
| /checkin | checkins, user_achievements, achievements |
| /users/{id}/achievements | achievements, user_achievements |