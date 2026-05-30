# 数据库详细设计（database-detailed-design.md）

> 对应概要设计 14、15 章。SQL DDL（schema.sql）以本文件为准。约束见 sql-rules.md。
> 主键统一 `id BIGINT PK AUTO_INCREMENT`；时间 `created_at DATETIME NOT NULL`，可变表加 `updated_at`。
> 枚举用 VARCHAR + 应用层校验（不用 MySQL ENUM）。

## 0. 表清单（12 张）
users / cities / city_tags / landmarks / foods / characters / city_visits / dice_rolls / checkins / achievements / user_achievements / chat_messages

## 0.1 枚举字典
| 字段 | 取值 |
|---|---|
| city_visits.visit_mode | free, game |
| city_visits.source | map_click, search, dice_roll, achievement |
| characters.character_type | history, culture, symbol |
| checkins.checkin_mode | free, game |
| chat_messages.role | user, assistant |
| dice_rolls.direction | 北,东北,东,东南,南,西南,西,西北 |
| achievements.rule_type | first_checkin, checkin_count, city_tag, tag_count, game_visit_count, dice_direction, dice_distance |

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
| current_city_id | BIGINT | NULL | 当前所在城市 |
| created_at | DATETIME | NOT NULL | |
| updated_at | DATETIME | NOT NULL | |

**索引**：UNIQUE(anonymous_id)。

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

**索引**：INDEX(name)。

---

## 3. city_tags — 城市标签表
| 字段 | 类型 | 约束 |
|---|---|---|
| id | BIGINT | PK AI |
| city_id | BIGINT | NOT NULL, INDEX |
| tag | VARCHAR(64) | NOT NULL |
| created_at | DATETIME | NOT NULL |

**示例 tag**：ancient_capital, dongbei, jiangnan, wuyue, coastal, spicy_food, modern_city, northwest, lingnan。
**索引**：INDEX(city_id), INDEX(tag)。

---

## 4. landmarks — 地标表
| 字段 | 类型 | 约束 |
|---|---|---|
| id | BIGINT | PK AI |
| city_id | BIGINT | NOT NULL, INDEX |
| name | VARCHAR(128) | NOT NULL |
| image_url | VARCHAR(512) | NULL（生图参考图） |
| description | TEXT | NULL |
| created_at | DATETIME | NOT NULL |

---

## 5. foods — 美食表
结构同 landmarks：id / city_id(INDEX) / name / image_url / description / created_at。

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
| created_at | DATETIME | NOT NULL | |

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

**索引**：INDEX(user_id, city_id), INDEX(user_id, created_at)。

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

**索引**：INDEX(user_id, created_at)。
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

**索引**：INDEX(user_id, created_at), INDEX(user_id, city_id)。
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

**索引**：UNIQUE(user_id, achievement_id)（防重复解锁）。

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

**索引**：INDEX(user_id, character_id, created_at)（拉取对话历史）。

---

## 13. GORM model 映射提醒
- 每个 model struct 对应一表，含 `gorm:"column:...;index"` 与 `json:"..."` tag。
- json 字段名需与 api-contract.md 一致。
- created_at/updated_at 用 GORM 自动维护（autoCreateTime/autoUpdateTime）。