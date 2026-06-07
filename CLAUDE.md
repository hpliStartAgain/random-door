# CLAUDE.md - 任意门 Coding Agent 总约束

> 与 `AGENTS.md` 保持一致，供 Claude Code 等工具读取。

## 1. 项目简介

任意门是一个将中国城市文化、地图探索、随机漫游、AI 人物对话、AI 赛博打卡、猜城市挑战、足迹资产和成就收集结合的 Web 互动探索产品。

## 2. 技术栈

- 前端：React + Vite + TypeScript + Tailwind CSS + Zustand + 高德地图 JS API
- 后端：Go + Gin + GORM
- 数据库：MySQL 8 / MariaDB
- AI：外部 LLM API + 外部图像生成 API，仅后端调用
- 部署：Docker Compose 单机，Caddy 托管前端并反代 API

## 3. 核心约束

1. 不引入 Redis / Kafka / Elasticsearch / 向量库 / 微服务网关。
2. AI API Key 只存在于后端环境变量；前端不得出现 LLM / IMAGE Key。
3. 后端保持 `api -> service -> repository` 分层。
4. 数据库是内容事实源；seed 仅用于受控导入。
5. 接口变更必须先改 `docs/design/api-contract.md`。

## 4. 变更前阅读

| 变更范围 | 先读 |
|---|---|
| Go 后端 | `docs/agent/go-backend-rules.md` |
| React 前端 | `docs/agent/react-frontend-rules.md` |
| SQL / 表结构 | `docs/agent/sql-rules.md` |
| AI | `docs/agent/ai-integration-rules.md` |
| 文档 | `docs/agent/doc-writing.md` |
| 提交 | `docs/agent/git-workflow.md` |

设计文档索引见 `docs/design/00-detailed-design-index.md`。

## 5. 启动

```bash
cp .env.example .env
docker compose up -d
```

内容初始化：`make seed-audit` 后执行 `make seed`。
