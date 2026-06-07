# AGENTS.md - 任意门 Coding Agent 总约束

> Coding agent 进入本仓库后先读本文件，再按变更范围阅读对应子文档。
> 当前代码是唯一事实标准；文档用于约束后续开发，不得反向覆盖代码事实。

## 1. 项目简介

任意门是一个将中国城市文化、地图探索、随机漫游、AI 人物对话、AI 赛博打卡、猜城市挑战、足迹资产和成就收集结合的 Web 互动探索产品。

## 2. 技术栈

- 前端：React + Vite + TypeScript + Tailwind CSS + Zustand + 高德地图 JS API
- 后端：Go + Gin + GORM
- 数据库：MySQL 8 / MariaDB
- AI：外部 LLM API + 外部图像生成 API，仅后端调用
- 部署：Docker Compose 单机，Caddy 托管前端并反代 API

## 3. 目录结构

```text
.
├── AGENTS.md / CLAUDE.md / README.md / Makefile / .env.example
├── backend/        Go 后端：cmd / internal / data/seed / migrations / static
├── frontend/       React 前端：src/pages|components|api|store|lib
├── docs/           产品、架构、设计与 agent 约束
├── scripts/        seed 与 AI 自检辅助脚本
├── TODO.md         当前未完备功能与整理事项
└── CHANGELOG.md    版本记录
```

完整文件树见 `docs/arch/directory-structure.md`。

## 4. 五条铁律

1. 不引入 Redis / Kafka / Elasticsearch / 向量库 / 微服务网关；保持 Go 单体 + MySQL。
2. AI API Key 只存在于后端环境变量；前端不得出现 LLM / IMAGE Key，不得直连 AI API。
3. 后端分层保持 `api -> service -> repository`，handler 不写业务逻辑。
4. 内容库以数据库为事实源；seed 仅用于受控 bootstrap / sync，默认不在启动时写库。
5. 接口以 `docs/design/api-contract.md` 为唯一契约；接口改动必须先改契约。

## 5. 启动

```bash
cp .env.example .env
docker compose up -d
```

空库初始化内容时先执行 `make seed-audit`，再执行 `make seed`。

## 6. 写代码前必读

| 变更范围 | 先读 |
|---|---|
| 任意 `.go` | `docs/agent/go-backend-rules.md` |
| 任意前端文件 | `docs/agent/react-frontend-rules.md` |
| 任意 `.sql` / 表结构 | `docs/agent/sql-rules.md` |
| 任意 AI 相关代码 | `docs/agent/ai-integration-rules.md` |
| 任意文档 | `docs/agent/doc-writing.md` |
| 提交代码 | `docs/agent/git-workflow.md` |

## 7. 设计文档

- 接口：`docs/design/api-contract.md`
- 后端：`docs/design/backend-detailed-design.md`
- 前端：`docs/design/frontend-detailed-design.md`
- 数据库：`docs/design/database-detailed-design.md`
- 算法 / AI / 成就：`docs/design/geo-algorithm-detailed-design.md`、`docs/design/ai-orchestration-detailed-design.md`、`docs/design/achievement-engine-detailed-design.md`

## 8. 禁止事项

- 在 handler 里写业务逻辑。
- 在前端组件里直接 fetch / axios，或在 `MapCanvas` 外直接操作高德 SDK。
- 部署本地大模型或本地生图模型。
- 提交 `.env`、密钥、`uploads/` 下的运行时文件、二进制产物。
- 让 AI 声称自己是“真实复活的人物”或编造确定性史实。
- 使用在世名人肖像生图，或生成色情、暴力、侮辱性内容。
- 擅自更换技术栈或扩展为微服务架构。

## 9. 完成定义

功能完成需满足 `docs/product/acceptance-criteria.md` 对应验收项；Docker Compose 启动后应能在不手工改表的前提下完成主链路。
