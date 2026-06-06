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
| 429 | AI_QUOTA_EXCEEDED | 匿名用户当日 AI 调用次数超限 |
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
**错误**：anonymous_id 不是标准 UUID → 400 INVALID_PARAM。

---

## 2. GET /api/cities — 城市列表

**请求**：无参数（演示版全量 35 城，seed 校验允许 12~100 城）。

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
**错误**：city_id 不是正整数 → 400 INVALID_PARAM。
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
**错误**：user / city 不存在 → 404 NOT_FOUND；user_id / city_id 不是正整数或 source 不属于 map_click / search → 400 INVALID_PARAM。

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
**逻辑**：读 city/character/可选 dialect_style 与历史消息 → prompt_builder 组装简短角色扮演 system prompt → llm_client 调用 → 落库 user 与 assistant 两条 chat_messages → 返回 reply。

---

## 7.1 评论与弹幕

游客可对地标、人物、美食发表评论；前端读取同一列表用于评论区与弹幕层。

### GET /api/comments

**Query**
| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| target_type | string | 是 | landmark / character / food |
| target_id | int64 | 是 | 对应资源 id |
| limit | int | 否 | 1~100，默认 50 |

**响应 200**
```json
{
  "comments": [
    {
      "id": 1,
      "target_type": "landmark",
      "target_id": 12,
      "user_id": 1,
      "nickname": "北京游客",
      "content": "这个角度很适合打卡。",
      "created_at": "2026-05-30T09:39:23+08:00"
    }
  ]
}
```

### POST /api/comments

**请求**
```json
{
  "target_type": "landmark",
  "target_id": 12,
  "user_id": 1,
  "nickname": "北京游客",
  "content": "这个角度很适合打卡。"
}
```
| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| target_type | string | 是 | landmark / character / food |
| target_id | int64 | 是 | 对应资源 id |
| user_id | int64 | 否 | 匿名用户 id，未初始化时可为空 |
| nickname | string | 否 | 留空默认“游客” |
| content | string | 是 | 1~200 字 |

**响应 201**
```json
{
  "id": 2,
  "target_type": "landmark",
  "target_id": 12,
  "user_id": 1,
  "nickname": "北京游客",
  "content": "这个角度很适合打卡。",
  "created_at": "2026-05-30T09:40:23+08:00"
}
```

**错误**：target_type 非法、target_id 非正整数、content 为空或超长 → 400；目标资源不存在 → 404。

---

## 7.2 POST /api/guess/caption — 全景截图猜一猜文案

**请求**
```json
{
  "user_id": 1,
  "city_id": 3,
  "target_name": "兵马俑",
  "scene_hint": "全景截图里有城墙、砖色建筑和游客动线"
}
```
| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| user_id | int64 | 否 | 匿名用户 id，便于后续统计 |
| city_id | int64 | 是 | 当前全景所属城市 |
| target_name | string | 否 | 当前全景/地标名，缺省用城市名 |
| scene_hint | string | 否 | 前端截图或用户视角描述，最长 200 字 |

**响应 200**
```json
{
  "weibo": "猜猜我在哪座城？砖色城墙、古都风物和一阵西北风都在画面里。#任意门漫游# #西安#",
  "moments": "把视角停在这座古都的一角，像是突然闯进一段历史现场。猜猜这是哪里？",
  "hashtags": ["任意门漫游", "西安"]
}
```

**逻辑**：后端读取城市/地标/美食/人物摘要 → 组装简短 prompt → 调 LLM 生成微博与朋友圈文案；LLM 未配置或失败时返回模板兜底，不阻断全景体验。

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
| scene_file | file | 否 | 全景 VR 当前视角截图，jpg/jpeg/png/webp，≤5MB；传入时优先作为生图参考图 |

**响应 202**
```json
{
  "task_id": 9001,
  "status": "queued"
}
```
**错误**：类型不符 → 415；超过 5MB → 413；当日次数超限 → 429。
**逻辑**：upload.validator 校验 → storage 以 UUID 落 uploads/selfies；如有 scene_file 则落 uploads/scenes → 校验城市/地标 → 写 `ai_tasks(status=queued,type=checkin_image)` → 返回任务。**本接口不阻塞等待外部生图，也不写 checkins**（确认后由接口 9 写）。生图参考图优先使用 scene_file，其次使用 landmark.image_url。

---

## 8.1 GET /api/checkin/image-tasks/{task_id} — 查询生图任务

**路径参数**：task_id（int64）。

**响应 200**
```json
{
  "task_id": 9001,
  "status": "succeeded",
  "result_url": "/uploads/generated/img_123.png",
  "error": null,
  "attempts": 1,
  "created_at": "2026-05-30T09:39:23+08:00",
  "updated_at": "2026-05-30T09:39:33+08:00"
}
```
**status**：queued / running / succeeded / failed / retryable。
**错误**：任务不存在或不属于该 user_id（请求头 `X-User-Id` 或 query `user_id`）→ 404 NOT_FOUND。

---

## 8.2 POST /api/checkin/image-tasks/{task_id}/retry — 重试生图任务

**请求**
```json
{ "user_id": 1 }
```

**响应 202**
```json
{ "task_id": 9001, "status": "queued" }
```
**逻辑**：仅 `failed/retryable` 任务可重试；重试会递增后续 attempts，但不重新上传自拍。

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

## 10.1 GET /api/users/{user_id}/assets — 匿名用户资产页

**路径参数**：user_id。

**响应 200**
```json
{
  "visited_cities": [
    { "id": 3, "name": "西安", "province": "陕西", "visited_at": "2026-05-30T09:39:23+08:00" }
  ],
  "posters": [
    {
      "checkin_id": 3001,
      "city_id": 3,
      "city_name": "西安",
      "landmark_name": "兵马俑",
      "generated_image_url": "/uploads/generated/img_123.png",
      "created_at": "2026-05-30T09:45:00+08:00"
    }
  ],
  "achievement_progress": [
    { "code": "ancient_capital_tour", "current": 1, "target": 3 }
  ]
}
```

---

## 11. Admin 内容 CMS（ADMIN_TOKEN）

Admin 接口均需 `X-Admin-Token` 或 `Authorization: Bearer <token>`。

### 11.1 GET /api/admin/catalog/coverage

**响应 200**
```json
{
  "total_cities": 35,
  "complete_cities": 35,
  "items": [
    {
      "city_id": 1,
      "city_name": "北京",
      "has_cover_image": true,
      "tag_count": 2,
      "landmark_count": 2,
      "food_count": 2,
      "character_count": 1,
      "missing_fields": []
    }
  ]
}
```

### 11.2 PATCH 内容字段

- `PATCH /api/admin/cities/{id}`
- `PATCH /api/admin/landmarks/{id}`
- `PATCH /api/admin/foods/{id}`
- `PATCH /api/admin/characters/{id}`

请求仅传需要修改的字段。图片字段必须为本地 `/static/...` 或 `/uploads/...` 路径；远程图片需要先通过上传/导入接口落本地。

**示例**
```json
{ "intro": "新的城市简介", "tags": ["ancient_capital", "north_china"] }
```

### 11.3 创建/删除城市下属内容

- `POST /api/admin/cities/{city_id}/landmarks`
- `POST /api/admin/cities/{city_id}/foods`
- `POST /api/admin/cities/{city_id}/characters`
- `DELETE /api/admin/landmarks/{id}`
- `DELETE /api/admin/foods/{id}`
- `DELETE /api/admin/characters/{id}`

地标/美食创建请求：
```json
{ "name": "故宫", "description": "明清宫城。", "image_url": "/uploads/admin_imports/a.png" }
```

人物创建请求：
```json
{
  "name": "朱棣",
  "character_type": "history",
  "avatar_url": "/uploads/admin_imports/a.png",
  "persona": "明成祖朱棣的城市导览角色。",
  "dialect_style": "官话表达为主",
  "prompt": "你在和用户进行角色扮演的游戏，你扮演的人物是朱棣。不声称真实复活，不编史。"
}
```
`name` 必填；图片字段必须为本地 `/static/...` 或 `/uploads/...` 路径。人物 `character_type` 缺省为 `history`，`persona/prompt` 缺省时由后端生成合规默认值。

创建响应 201 返回新对象；删除响应 200：
```json
{ "status": "deleted" }
```

### 11.4 图片上传/导入

保留现有上传接口：
- `POST /api/admin/cities/{city_id}/cover-image`
- `POST /api/admin/landmarks/{landmark_id}/image`
- `POST /api/admin/foods/{food_id}/image`
- `POST /api/admin/characters/{character_id}/avatar`

URL 绑定接口不再保存远程 URL。`PATCH .../image` 或 `PATCH .../cover-image` 接收 `{ "url": "https://..." }` 后由后端下载到本地 `/uploads/admin_imports/...`，再保存本地路径。

---

## 12. 接口与数据表对照速查
| 接口 | 主要读写表 |
|---|---|
| /users/anonymous | users |
| /cities | cities, city_tags |
| /cities/{id} | cities, city_tags, landmarks, foods, characters |
| /visits/free | city_visits |
| /game/init | cities |
| /game/roll | dice_rolls, city_visits, cities |
| /chat | chat_messages, characters, cities |
| /comments | comments, landmarks, foods, characters |
| /guess/caption | cities, landmarks, foods, characters |
| /checkin/generate-image | ai_tasks, ai_usage_logs, uploads/selfies |
| /checkin/image-tasks/{id} | ai_tasks |
| /checkin | checkins, user_achievements, achievements |
| /users/{id}/achievements | achievements, user_achievements |
| /users/{id}/assets | city_visits, checkins, cities, landmarks, achievements |
| /admin/* | cities, city_tags, landmarks, foods, characters |
