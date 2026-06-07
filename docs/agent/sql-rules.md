# SQL / 数据库约束

> 写任何 `.sql` 或改表结构前必读。父约束：`AGENTS.md` / `CLAUDE.md`。

## 1. 命名

- 表名：复数、小写、下划线，如 `city_visits`。
- 字段名：小写下划线，如 `from_city_id`。

## 2. 主键与时间

- 主键统一 `id BIGINT PRIMARY KEY AUTO_INCREMENT`。
- 时间字段统一 `created_at DATETIME NOT NULL`；可变表加 `updated_at DATETIME NOT NULL`。

## 3. 枚举

- 枚举值用 `VARCHAR` 存储，在应用层校验，不使用 MySQL ENUM。
- 枚举全集见 `docs/design/database-detailed-design.md`。

## 4. 索引与外键

- 所有 `*_id` 关系字段建普通索引。
- 唯一自然键按设计文档维护，例如 `users.anonymous_id`、`cities.name`、`achievements.code`。
- 高频查询建复合索引，例如 `checkins(user_id, created_at)`、`city_visits(user_id, city_id)`。
- 当前实现不建物理外键，一致性由 service 层校验。

## 5. 变更流程

- 改表必须同步 `docs/design/database-detailed-design.md` 与 `backend/migrations/schema.sql`。
- seed 与 schema 分离：schema 只建表，种子数据在 `backend/data/seed/*.json`。
- 数据库是内容事实源；seed 导入默认关闭，维护前先 audit。

## 6. MySQL 低资源配置

```text
innodb_buffer_pool_size=256M
max_connections=50
performance_schema=OFF
```

## 7. 提交前自检

- [ ] 命名规范、主键/时间字段齐全。
- [ ] 关系字段有索引。
- [ ] schema 与数据库设计一致。
- [ ] seed 变更不会覆盖后台已维护内容，除非显式 sync。
