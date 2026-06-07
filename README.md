# 任意门 · AI 城市漫游（Random Door）v1.0

> **推开门，遇见大美中国。**
> 一个城市文化互动探索 Web 应用：自由探索 + 大富翁式随机漫游 + AI 跨时空人物对话 + AI 赛博打卡海报 + 成就收集。

## 架构概览
```
浏览器 (React + 高德地图 3D + Pannellum 360全景)
        │ HTTP
        ▼
   Go 后端 (Gin + GORM)

   ├─ 业务 / 随机算法 / AI 编排 / 成就判定
        │                 │
        ▼                 ▼
     MySQL          外部 AI API (LLM / 生图)
```
**前端技术栈亮点**：采用 Zustand 进行无缝的跨组件事件流转；宏观视角使用高德 3D API 进行赛博飞越；微观探索利用 WebGL (Pannellum) 渲染高清无损的等距圆柱全景摄影图。
详见 `docs/arch/system-architecture.md`。

## 快速开始
```bash
# 1. 克隆
git clone <repo> && cd random-door

# 2. 配置环境变量
cp .env.example .env
# 编辑 .env：至少修改外部 MySQL 的 DB_*，并填入高德 Key
# LLM_API_KEY / IMAGE_API_KEY 可先留空或 mock，演示真实 AI 时再填

# 3. 一键启动
docker compose up -d        # 或 make up

# 4. 内容数据
# 后端启动会自动建表，但默认不写 seed；数据库是后台内容的事实源。
# 如需给空库补齐演示 seed，先只读盘点，再显式 bootstrap：
make seed-audit
make seed

# 5. 访问
# 前端: http://localhost           （或 Caddy 暴露端口）
# 后端: http://localhost:8080/api
```

## 常用命令（Makefile）
| 命令 | 作用 |
|---|---|
| `make up` / `make down` | 启停所有容器 |
| `make build` | 构建镜像 |
| `make migrate` | (已废弃，启动时自动建表) |
| `make seed-audit` | 只读盘点真实库与 seed 的差异 |
| `make seed` | 显式补齐缺失 seed 行，不覆盖后台已有内容 |
| `make seed-sync` | 谨慎：按 seed 覆盖自然键匹配行 |
| `make seed-landmark-coords` | 仅回填缺失的地标经纬度 |
| `make logs` | 查看后端日志 |
| `make lint` / `make test` | 检查 / 测试 |

## AI 联调

真实 Key 填入 `.env` 后，可先跑不泄露密钥的端点自检：

```bash
python scripts/ai_smoke.py --llm
python scripts/ai_smoke.py --image --selfie path/to/selfie.jpg --confirm-image-cost
```

脚本只输出是否配置、HTTP 状态和结果摘要，不打印 `LLM_API_KEY` / `IMAGE_API_KEY`。

## 目录说明
| 目录 | 内容 |
|---|---|
| `backend/` | Go 后端：cmd/server 入口，internal/ 分层，data/seed 种子，migrations DDL |
| `frontend/` | React 前端：src/pages|components|api|store |
| `docs/design/` | 详细设计（agent 写代码依据） |
| `docs/agent/` | coding agent 约束规则 |
| `docs/arch/` | 架构文档 |
| `docs/product/` | 产品文档 |

## 文档索引
- 入门：`CLAUDE.md`（agent 总约束）
- 设计：`docs/design/00-detailed-design-index.md`
- 架构：`docs/arch/system-architecture.md`、`directory-structure.md`
- 产品：`docs/product/prd.md`、`user-flows.md`
- 验收：`docs/product/acceptance-criteria.md`

## 后台管理

Navbar 右上角齿轮图标 → 输入 `ADMIN_TOKEN`（`.env` 中配置）→ 进入媒体资产管理界面。

功能：为每个城市的**封面图 / 地标图 / 人物头像 / 美食图**批量上传文件或粘贴外链 URL。

## 黄金动线（<10 分钟演示）

见 `docs/product/demo-script.md`：开场品牌页 → 掷骰仪式 → 城市探索 → AI 跨时空对话 → 赛博打卡 → 成就解锁。

每 ~90s 一个 WOW，适合评委演示。

## MVP 范围与约束
- 35 个精选演示城市；匿名用户（无注册登录）；不引入 Redis/MQ；AI 走外部 API；2C2G 单机部署。
- 验收标准见 `docs/product/acceptance-criteria.md`（20 条）。
