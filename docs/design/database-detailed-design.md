# 数据库详细设计（database-detailed-design.md）

> 对应概要设计 14、15 章。SQL DDL（schema.sql）以本文件为准。约束见 sql-rules.md。
> 主键统一 `id BIGINT PK AUTO_INCREMENT`；时间 `created_at DATETIME NOT NULL`，可变表加 `updated_at`。
> 枚举用 VARCHAR + 应用层校验（不用 MySQL ENUM）。

## 0. 表清单（17 张）
users / cities / city_tags / landmarks / foods / characters / city_visits / dice_rolls / checkins / achievements / user_achievements / chat_messages / comments / ai_tasks / ai_usage_logs / guess_challenges / guess_answers

## 0.1 枚举字典
| 字段 | 取值 |
|---|---|
| city_visits.visit_mode | free, game |
| city_visits.source | map_click, search, dice_roll, achievement |
| characters.character_type | history, culture, symbol |
| checkins.checkin_mode | free, game |
| chat_messages.role | user, assistant |
| comments.target_type | landmark, food, character |
| dice_rolls.direction | 北,东北,东,东南,南,西南,西,西北 |
| achievements.rule_type | first_checkin, checkin_count, city_tag, tag_count, game_visit_count, dice_direction, dice_distance |
| ai_tasks.type | checkin_image |
| ai_tasks.status | queued, running, succeeded, failed, retryable |
| ai_usage_logs.usage_type | chat, image |

## 0.2 外键策略
MVP **不建物理外键约束**（避免 seed 导入顺序问题与删除级联复杂度），改为：
1. 所有 *_id 字段**建普通索引**；
2. 一致性由 service 层保证。

---

## 1. users — 用户表
| 字段 | 类型 | 约束 | 说明 |
|---|---|---|---|
| id | BIGINT | PK AI | |
| anonymous_id | VARCHAR(128) | NOT NULL, UNIQUE | 前端 UUID |
| nickname | VARCHAR(64) | NULL | |
| avatar_url | VARCHAR(512) | NULL | |
| age | INT | NULL | 匿名 profile 年龄 |
| home_region | VARCHAR(64) | NULL | 匿名 profile 家乡/地区 |
| current_city_id | BIGINT | NULL | 当前所在城市 |
| created_at | DATETIME | NOT NULL | |
| updated_at | DATETIME | NOT NULL | |

**索引**：UNIQUE(anonymous_id), INDEX(current_city_id)。

---

## 2. cities — 城市表
| 字段 | 类型 | 约束 |
|---|---|---|
| id | BIGINT | PK AI |
| name | VARCHAR(64) | NOT NULL |
| province | VARCHAR(64) | NOT NULL |
| lat | DOUBLE | NOT NULL |
| lng | DOUBLE | NOT NULL |
| intro | TEXT | NULL |
| cover_image_url | VARCHAR(512) | NULL |
| dialect_sample | VARCHAR(255) | NULL |
| dialect_explanation | TEXT | NULL |
| created_at / updated_at | DATETIME | NOT NULL |

**索引**：UNIQUE(name)。城市名是受控 seed 导入的自然键。

---

## 3. city_tags — 城市标签表
| 字段 | 类型 | 约束 |
|---|---|---|
| id | BIGINT | PK AI |
| city_id | BIGINT | NOT NULL, INDEX |
| tag | VARCHAR(64) | NOT NULL |
| created_at | DATETIME | NOT NULL |

**示例 tag**：ancient_capital, dongbei, jiangnan, wuyue, coastal, spicy_food, modern_city, northwest, lingnan。
**索引**：UNIQUE(city_id, tag), INDEX(city_id), INDEX(tag)。

---

## 4. landmarks — 地标表
| 字段 | 类型 | 约束 |
|---|---|---|
| id | BIGINT | PK AI |
| city_id | BIGINT | NOT NULL, INDEX |
| name | VARCHAR(128) | NOT NULL |
| image_url | VARCHAR(512) | NULL（生图参考图） |
| description | TEXT | NULL |
| soundscape_url | VARCHAR(512) | NULL（本地 `/static/soundscapes/...`） |
| created_at | DATETIME | NOT NULL |

**索引**：UNIQUE(city_id, name), INDEX(city_id)。`(city_id, name)` 是受控 seed 导入的自然键。

---

## 5. foods — 美食表
结构同 landmarks：id / city_id(INDEX) / name / image_url / description / created_at。

**索引**：UNIQUE(city_id, name), INDEX(city_id)。`(city_id, name)` 是受控 seed 导入的自然键。

---

## 6. characters — 人物表
| 字段 | 类型 | 约束 | 说明 |
|---|---|---|---|
| id | BIGINT | PK AI | |
| city_id | BIGINT | NOT NULL, INDEX | |
| name | VARCHAR(128) | NOT NULL | |
| character_type | VARCHAR(32) | NOT NULL | history/culture/symbol |
| avatar_url | VARCHAR(512) | NULL | |
| persona | TEXT | NOT NULL | 人设（**不下发前端**） |
| dialect_style | TEXT | NULL | |
| prompt | TEXT | NOT NULL | 系统 Prompt（**不下发前端**） |
| role_title | VARCHAR(128) | NULL | 人物身份（下发前端） |
| life_span | VARCHAR(64) | NULL | 年代/生卒（下发前端） |
| intro_quote | VARCHAR(255) | NULL | 角色口吻引导语（下发前端，非史实断言） |
| created_at | DATETIME | NOT NULL | |

**索引**：UNIQUE(city_id, name), INDEX(city_id)。`(city_id, name)` 是受控 seed 导入的自然键。

---

## 7. city_visits — 城市访问表（双模式关键表）
| 字段 | 类型 | 约束 | 说明 |
|---|---|---|---|
| id | BIGINT | PK AI | |
| user_id | BIGINT | NOT NULL, INDEX | |
| city_id | BIGINT | NOT NULL, INDEX | |
| visit_mode | VARCHAR(32) | NOT NULL | free/game |
| source | VARCHAR(64) | NULL | map_click/search/dice_roll/achievement |
| from_city_id | BIGINT | NULL | 游戏模式出发城市 |
| dice_roll_id | BIGINT | NULL | 游戏模式关联掷骰 |
| created_at | DATETIME | NOT NULL | |

**索引**：INDEX(user_id, city_id), INDEX(user_id, created_at), INDEX(city_id), INDEX(from_city_id), INDEX(dice_roll_id)。

---

## 8. dice_rolls — 掷骰记录表
| 字段 | 类型 | 约束 |
|---|---|---|
| id | BIGINT | PK AI |
| user_id | BIGINT | NOT NULL, INDEX |
| from_city_id | BIGINT | NULL |
| to_city_id | BIGINT | NOT NULL |
| direction | VARCHAR(32) | NOT NULL |
| distance_km | INT | NOT NULL |
| target_lat | DOUBLE | NULL |
| target_lng | DOUBLE | NULL |
| created_at | DATETIME | NOT NULL |

**索引**：INDEX(user_id, created_at), INDEX(from_city_id), INDEX(to_city_id)。
**关系**：DiceRoll 1-1 CityVisit（通过 city_visits.dice_roll_id 回指）。

---

## 9. checkins — 打卡表
| 字段 | 类型 | 约束 |
|---|---|---|
| id | BIGINT | PK AI |
| user_id | BIGINT | NOT NULL, INDEX |
| city_id | BIGINT | NOT NULL, INDEX |
| landmark_id | BIGINT | NULL |
| visit_id | BIGINT | NULL |
| generated_image_url | VARCHAR(512) | NULL |
| checkin_mode | VARCHAR(32) | NULL | free/game |
| created_at | DATETIME | NOT NULL |

**索引**：INDEX(user_id, created_at), INDEX(user_id, city_id), INDEX(city_id), INDEX(landmark_id), INDEX(visit_id)。
**用途**：成就判定的主要数据来源。

---

## 10. achievements — 成就表
| 字段 | 类型 | 约束 | 说明 |
|---|---|---|---|
| id | BIGINT | PK AI | |
| code | VARCHAR(64) | NOT NULL, UNIQUE | 业务唯一标识 |
| name | VARCHAR(128) | NOT NULL | |
| description | TEXT | NULL | |
| rule_type | VARCHAR(64) | NOT NULL | 见枚举字典 |
| rule_value | VARCHAR(255) | NOT NULL | 规则参数（见下） |
| badge_url | VARCHAR(512) | NULL | |
| created_at | DATETIME | NOT NULL | |

**rule_value 约定**（与 achievement-engine 对齐）：
| rule_type | rule_value 示例 | 含义 |
|---|---|---|
| first_checkin | "" | 任意首次打卡 |
| checkin_count | "5" | 累计打卡 5 次 |
| city_tag | "ancient_capital" | 打卡任意该标签城市 |
| tag_count | "jiangnan:3" | 打卡 3 座该标签城市 |
| game_visit_count | "5" | 游戏模式打卡 5 城 |
| dice_direction | "北:2" | 连续 2 次掷出该方向 |
| dice_distance | "1200" | 单次掷出 1200km |

---

## 11. user_achievements — 用户成就表
| 字段 | 类型 | 约束 |
|---|---|---|
| id | BIGINT | PK AI |
| user_id | BIGINT | NOT NULL, INDEX |
| achievement_id | BIGINT | NOT NULL, INDEX |
| unlocked_at | DATETIME | NOT NULL |

**索引**：UNIQUE(user_id, achievement_id)（防重复解锁）, INDEX(achievement_id)。

---

## 12. chat_messages — 对话消息表
| 字段 | 类型 | 约束 |
|---|---|---|
| id | BIGINT | PK AI |
| user_id | BIGINT | NOT NULL, INDEX |
| city_id | BIGINT | NOT NULL |
| character_id | BIGINT | NOT NULL, INDEX |
| role | VARCHAR(16) | NOT NULL | user/assistant |
| content | TEXT | NOT NULL |
| created_at | DATETIME | NOT NULL |

**索引**：INDEX(user_id, character_id, created_at)（拉取对话历史）, INDEX(city_id), INDEX(character_id)。

---

## 13. comments — 游客评论 / 弹幕表

| 字段 | 类型 | 约束 | 说明 |
|---|---|---|---|
| id | BIGINT | PK AI | |
| target_type | VARCHAR(32) | NOT NULL, INDEX | landmark/food/character |
| target_id | BIGINT | NOT NULL, INDEX | 对应资源 id |
| user_id | BIGINT | NULL, INDEX | 匿名用户，可为空 |
| nickname | VARCHAR(64) | NOT NULL | 留空时应用层写“游客” |
| content | VARCHAR(500) | NOT NULL | 展示在评论区与弹幕层 |
| created_at | DATETIME | NOT NULL | |

**索引**：INDEX(target_type, target_id, created_at), INDEX(user_id, created_at)。
**一致性**：target_id 指向何表由 target_type 决定，service 层校验目标存在；MVP 不建多态物理外键。

---

## 14. ai_tasks — AI 异步任务表

| 字段 | 类型 | 约束 | 说明 |
|---|---|---|---|
| id | BIGINT | PK AI | |
| user_id | BIGINT | NOT NULL, INDEX | 匿名用户 |
| type | VARCHAR(32) | NOT NULL | checkin_image |
| status | VARCHAR(32) | NOT NULL, INDEX | queued/running/succeeded/failed/retryable |
| input_json | JSON | NOT NULL | 城市、地标、自拍路径、prompt 输入 |
| result_url | VARCHAR(512) | NULL | 成功后本地 `/uploads/...` |
| error | TEXT | NULL | 最近一次失败原因 |
| attempts | INT | NOT NULL DEFAULT 0 | 已尝试次数 |
| created_at / updated_at | DATETIME | NOT NULL | |

**索引**：INDEX(user_id, created_at), INDEX(status, updated_at), INDEX(type, status)。

---

## 15. ai_usage_logs — AI 用量流水表

| 字段 | 类型 | 约束 | 说明 |
|---|---|---|---|
| id | BIGINT | PK AI | |
| user_id | BIGINT | NOT NULL, INDEX | |
| usage_type | VARCHAR(32) | NOT NULL | chat/image |
| usage_date | DATE | NOT NULL | 按匿名用户本地日期限流 |
| count | INT | NOT NULL DEFAULT 0 | 当日计数 |
| created_at / updated_at | DATETIME | NOT NULL | |

**索引**：UNIQUE(user_id, usage_type, usage_date), INDEX(usage_date)。

---

## 16. guess_challenges — 猜位置挑战表

| 字段 | 类型 | 约束 | 说明 |
|---|---|---|---|
| id | BIGINT | PK AI | |
| code | VARCHAR(16) | NOT NULL, UNIQUE | 对外分享 code，不暴露 user_id |
| user_id | BIGINT | NULL, INDEX | 匿名创建者，可为空 |
| city_id | BIGINT | NOT NULL, INDEX | 正确城市 |
| target_name | VARCHAR(128) | NULL | 地标/场景名 |
| image_url | VARCHAR(512) | NULL | 本地截图或地标图片 |
| caption | VARCHAR(300) | NULL | 分享文案 |
| created_at | DATETIME | NOT NULL | |
| expires_at | DATETIME | NOT NULL, INDEX | 过期后不可访问 |

**索引**：UNIQUE(code), INDEX(user_id, created_at), INDEX(city_id), INDEX(expires_at)。

---

## 17. guess_answers — 猜位置答案表

| 字段 | 类型 | 约束 | 说明 |
|---|---|---|---|
| id | BIGINT | PK AI | |
| challenge_code | VARCHAR(16) | NOT NULL, INDEX | 对应 guess_challenges.code |
| answer_text | VARCHAR(64) | NOT NULL | 用户提交答案 |
| is_correct | TINYINT(1) | NOT NULL | 是否命中城市或 target_name |
| created_at | DATETIME | NOT NULL | |

**索引**：INDEX(challenge_code, created_at)。

---

## 18. GORM model 映射提醒
- 每个 model struct 对应一表，含 `gorm:"column:...;index"` 与 `json:"..."` tag。
- json 字段名需与 api-contract.md 一致。
- created_at/updated_at 用 GORM 自动维护（autoCreateTime/autoUpdateTime）。

---

## 19. Seed 受控导入

- `backend/internal/seed` 在写库前完整解析并校验 `cities.json` 与 `achievements.json`。
- 城市数量允许 12~100 个；演示版 seed 固定为 35 个精选城市。每城 1~2 地标、1~2 美食、1 人物，且必须含方言、静态资源 URL 和合规人物 Prompt。
- 地标 `soundscape_url` 可为空；非空时必须是 `/static/soundscapes/...` 本地路径，并确保音频文件存在或演示前补齐。
- 服务启动默认 `SEED_MODE=off`，不写入 seed；数据库是后台 catalog 的事实源。
- `bootstrap` 模式在单一事务内只插入缺失行，不更新已有后台内容；适用于空库或人工确认后的补齐。
- `sync` 模式在单一事务内按自然键覆盖匹配行；仅用于明确需要 seed 覆盖数据库的维护场景，执行前必须先 audit。
- 自然键：`cities.name`、`city_tags(city_id,tag)`、`landmarks(city_id,name)`、`foods(city_id,name)`、`characters(city_id,name)`、`achievements.code`。
- 显式工具：`go run ./cmd/seedtool -mode audit|bootstrap|sync`；`sync` 需额外确认参数。
