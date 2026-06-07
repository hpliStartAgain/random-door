# 部署架构（deployment-architecture.md）

> 对应概要设计 22 章。落地依据为 docker-compose.yml。目标：2C2G 单机 Docker Compose 一键起。

## 1. 部署组件
| 组件 | 说明 | 是否必需 |
|---|---|---|
| app | Go 后端服务（Gin） | 必需 |
| mysql | 外部 MySQL 8 / MariaDB 实例 | 必需 |
| caddy | 静态托管前端 + 反代 /api（可选，亦可由 app 直接托管） | 可选 |
| static assets | 城市地标/美食/人物/勋章静态资源（镜像内置或只读挂载） | 必需 |
| uploads volume | 用户上传与生成图片 | 必需 |

> 注：**不含 Nginx**（理由见 system-architecture.md 第 6 节）。

## 2. Docker Compose 拓扑
```text
2C2G Server
  ├── caddy (可选)   :80/:443  →  反代 app:8080，托管 frontend/dist
  ├── app (Go)       :8080     →  读 .env，连外部 mysql，挂 static/uploads
  ├── mysql          :3306     →  外部实例，按第 5 节低资源配置
  ├── volume static
  └── volume uploads
```

## 3. compose 关键约定（写 docker-compose.yml 时遵守）
- app 配 healthcheck；MySQL 使用外部实例，通过 `.env` 的 DB_* 连接。
- app 通过 `env_file: .env` 注入配置（DB/AI/upload/cors）。
- 外部 mysql 用低资源参数（见第 5 节）。
- 持久卷：uploads；静态资源通过只读挂载或镜像内置提供。
- 统一 `networks` 内部互通；仅 caddy（或 app）对外暴露端口。

## 4. 资源评估（2C2G，对应概要 22.3）
| 组件 | 预估内存 |
|---|---|
| Go Backend | 50MB – 200MB |
| MySQL | 400MB – 900MB |
| Caddy | 20MB – 100MB |
| OS | 400MB – 700MB |
| 预留 | 200MB – 500MB |

结论：2C2G 可承载黑客松 MVP，但必须遵守约束清单（第 6 节）。

## 5. MySQL 低资源配置（对应概要 22.4）
```text
innodb_buffer_pool_size=256M
max_connections=50
performance_schema=OFF
```
通过 compose 的 command 或挂载 my.cnf 注入。

## 6. 2C2G 约束清单（铁律，违反会 OOM）
1. 不部署本地 LLM；2. 不部署本地生图模型；3. 不引 Redis/Kafka/ES 等中间件；4. MySQL 用低资源配置；5. 限制上传图片大小（≤5MB）；6. Go 后端保持单体轻量。

## 7. 启动与初始化顺序
```text
1. cp .env.example .env 并填真实值
2. docker compose up -d        # 起 app (+ caddy)，MySQL 为外部实例
3. app 启动时：连 mysql → AutoMigrate → 按 SEED_MODE 决定是否受控导入 seed（默认 off）
4. 访问：caddy 暴露端口(前端) / app:8080/api(后端)
```
常用命令封装在 Makefile（up/down/build/migrate/seed-audit/seed/logs）。

## 8. 演示稳定性建议
- 演示前预热：先跑一遍掷骰/对话/生图，确认外部 AI 可用。
- 为重点城市准备**预生成样例图**，生图失败时降级展示（见 cyber-checkin-design.md）。
