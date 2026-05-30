"# SQL / 数据库约束

> 写任何 `.sql` 或改表结构前必读。父约束：`CLAUDE.md`。

## 1. 命名
- 表名：复数、小写、下划线（city_visits、chat_messages）。
- 字段名：小写下划线（from_city_id）。

## 2. 主键与时间
- 主键统一 id BIGINT PRIMARY KEY AUTO_INCREMENT。
- 时间字段统一 created_at DATETIME NOT NULL，需更新的表加 updated_at DATETIME NOT NULL。

## 3. 枚举
- 枚举值用 VARCHAR 存储，**在应用层校验**（不用 MySQL ENUM）。
  - visit_mode：free / game
  - source：map_click / search / dice_roll / achievement
  - role（chat）：user / assistant
  - character_type / checkin_mode / rule_type 见 database-detailed-design.md。

## 4. 索引与外键
- 所有 *_id 外键字段建索引。
- users.anonymous_id 唯一索引；achievements.code 唯一索引。
- 高频查询建复合索引：如 checkins(user_id, created_at)、city_visits(user_id, city_id)。
- 外键关系在文档登记；MVP 可用应用层保证一致性（是否加物理外键见 database-detailed-design.md）。

## 5. 变更流程
- 改表必须**同步更新** docs/design/database-detailed-design.md 与 backend/migrations/schema.sql。
- **seed 与 schema 分离**：schema.sql 只建表，种子数据在 data/seed/*.json。

## 6. MySQL 低资源配置（2C2G）
```
innodb_buffer_pool_size=256M
max_connections=50
performance_schema=OFF
```

## 7. 提交前自检
- [ ] 命名规范、主键/时间字段齐全
- [ ] 外键字段有索引
- [ ] schema.sql 与 database-detailed-design.md 一致"
