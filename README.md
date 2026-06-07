# 任意门

> 推开门，遇见大美中国。

任意门是一个城市文化互动探索 Web 产品：用户可以在中国城市地图上自由探索，也可以通过随机方向与随机距离开启一次“大富翁式”城市漫游。到达城市后，产品提供城市内容、地标风光、人物 AI 对话、赛博打卡海报、猜城市挑战、足迹资产和成就收集。

## 当前形态

- 前端：React + Vite + TypeScript + Tailwind CSS + Zustand + 高德地图 JS API。
- 后端：Go + Gin + GORM 单体服务，分层为 `api -> service -> repository`。
- 数据库：外部 MySQL 8 / MariaDB。服务启动自动建表，内容库以数据库为事实源。
- AI：外部 LLM 与外部图像生成 API，全部由后端调用；前端只持有高德公开 Key。
- 部署：Docker Compose 单机部署，`app` 提供 API，`caddy` 托管前端并反代 `/api`。

## 快速开始

```bash
git clone <repo> && cd random-door
cp .env.example .env
# 编辑 .env：配置 DB_*、VITE_AMAP_KEY、ADMIN_TOKEN；真实 AI 能力需配置 LLM_* / IMAGE_*。
docker compose up -d
```

访问地址：

- 前端：`http://localhost`
- 后端健康检查：`http://localhost:8080/health`
- API 前缀：`http://localhost:8080/api`

## 内容初始化

应用启动会执行 GORM AutoMigrate，但默认 `SEED_MODE=off`，不会覆盖数据库内容。空库初始化或维护回填请显式执行：

```bash
make seed-audit
make seed
```

`make seed` 只补齐缺失行；`make seed-sync` 会按 seed 覆盖自然键匹配行，执行前必须先确认差异。

当前仓库 seed 提供 35 座精选城市、每城 1-2 个地标、1-2 个美食、1 个人物，以及 11 个成就。线上内容扩展应优先走后台 CMS。

## 常用命令

| 命令 | 作用 |
|---|---|
| `make up` / `make down` | 启停容器 |
| `make build` | 构建镜像 |
| `make seed-audit` | 只读盘点数据库与 seed 差异 |
| `make seed` | 只插入缺失 seed 内容 |
| `make seed-sync` | 谨慎覆盖自然键匹配内容 |
| `make seed-landmark-coords` | 回填缺失地标经纬度 |
| `make logs` | 查看后端日志 |
| `make lint` | 后端 gofmt/vet + 前端 lint |
| `make test` | 后端单元测试 |
| `make fe-dev` / `make fe-build` | 前端开发 / 构建 |

## AI 联调

```bash
python scripts/ai_smoke.py --llm
python scripts/ai_smoke.py --image --selfie path/to/selfie.jpg --confirm-image-cost
```

脚本只输出配置状态、HTTP 状态和结果摘要，不打印 `LLM_API_KEY` / `IMAGE_API_KEY`。

## 后台管理

Navbar 右上角齿轮图标进入后台管理。后台由 `ADMIN_TOKEN` 保护，支持：

- 城市、标签、地标、美食、人物内容维护。
- 城市封面、地标图、美食图、人物头像、成就勋章上传。
- 远程图片 URL 导入到本地 `/uploads/admin_imports/...` 后绑定。
- 成就规则与勋章维护。

## 文档索引

- Agent 总约束：`AGENTS.md`、`CLAUDE.md`
- API 契约：`docs/design/api-contract.md`
- 数据库设计：`docs/design/database-detailed-design.md`
- 后端详设：`docs/design/backend-detailed-design.md`
- 前端详设：`docs/design/frontend-detailed-design.md`
- 产品总纲：`docs/product/prd.md`
- 用户流程：`docs/product/user-flows.md`
- 验收标准：`docs/product/acceptance-criteria.md`
- 待办事项：`TODO.md`
- 版本记录：`CHANGELOG.md`

## 关键约束

- 不引入 Redis / Kafka / Elasticsearch / 向量库 / 微服务网关。
- AI API Key 只存在于后端环境变量；前端不得直连 LLM 或图像生成 API。
- 接口变更必须先改 `docs/design/api-contract.md`。
- 表结构变更必须同步 `docs/design/database-detailed-design.md` 与 `backend/migrations/schema.sql`。
- 运行时上传与生成文件不入库，不提交 `.env` 或任何密钥。
