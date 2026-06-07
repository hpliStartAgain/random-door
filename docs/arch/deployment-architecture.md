# 部署架构

## 1. 部署组件

| 组件 | 说明 | 必需 |
|---|---|---|
| app | Go 后端服务，监听 `SERVER_PORT`，默认 8080。 | 是 |
| caddy | 前端静态托管与 `/api` 反代。 | 是 |
| MySQL / MariaDB | 外部数据库实例，通过 `.env` 的 `DB_*` 连接。 | 是 |
| static assets | 城市封面、地标、美食、人物、勋章、声景等静态资源。 | 是 |
| uploads volume | 用户上传、AI 生成图、后台导入图。 | 是 |

## 2. Docker Compose 拓扑

```text
Host
  :80  -> caddy
          ├── frontend dist
          ├── /api -> app:8080
          ├── /static -> backend/static readonly
          └── /uploads -> uploads volume readonly

  :8080 -> app
          ├── env_file .env
          ├── /static -> backend/static readonly
          └── /uploads -> uploads volume read/write

  MySQL is external
```

## 3. 启动顺序

```text
1. cp .env.example .env
2. 配置 DB_*、VITE_AMAP_KEY、ADMIN_TOKEN、可选 AI Key
3. docker compose up -d
4. app 连接 MySQL -> AutoMigrate -> 按 SEED_MODE 决定是否导入 seed
5. caddy 等待 app healthy 后启动
```

默认 `SEED_MODE=off`。空库内容初始化请运行 `make seed-audit` 和 `make seed`。

## 4. 资源约束

| 组件 | 典型内存 |
|---|---|
| Go app | 50-320MB |
| Caddy | 20-100MB |
| 外部 MySQL | 由外部实例承担 |
| OS 与余量 | 按服务器实际配置预留 |

2C2G 单机部署下必须保持：

- 不部署本地 LLM 或本地生图模型。
- 不新增 Redis / Kafka / Elasticsearch 等常驻中间件。
- 上传文件默认不超过 5MB。
- AI worker 并发默认 1。

## 5. MySQL 建议

```text
innodb_buffer_pool_size=256M
max_connections=50
performance_schema=OFF
```

连接池通过 `.env` 的 `DB_MAX_OPEN_CONNS`、`DB_MAX_IDLE_CONNS` 与数据库实例协调。

## 6. 运维入口

| 命令 | 用途 |
|---|---|
| `make up` / `make down` | 启停容器 |
| `make logs` | 查看 app 日志 |
| `make seed-audit` | 只读检查 seed 与数据库差异 |
| `make seed` | bootstrap 缺失内容 |
| `make seed-sync` | 明确确认后按 seed 覆盖内容 |
| `make seed-landmark-coords` | 回填缺失地标坐标 |
