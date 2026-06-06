# CLAUDE.md — AI 城市漫游产品 · Coding Agent 总约束

> 这是 coding agent（Claude Code / Cursor 等）进入本仓库后**首先必须阅读**的文件。
> 对应概要设计：全文。阶段：黑客松 MVP (v0.1)。

## 1. 项目一句话简介
一个将「中国城市文化 + 地图探索 + 大富翁式随机漫游 + AI 人物对话 + AI 赛博打卡 + 成就收集」结合的 Web 互动探索产品。

## 2. 技术栈（不得擅自变更）
- 前端：React + Vite + TypeScript + Tailwind CSS + shadcn/ui + Zustand + 高德地图 JS API
- 后端：Go + Gin + GORM
- 数据库：MySQL 8 / MariaDB
- AI：外部 LLM API + 外部图像生成 API（**仅后端调用**）
- 部署：Docker Compose 单机，目标服务器 2C2G

## 3. 目录结构总览
```
.
├── CLAUDE.md / .cursorrules / README.md / Makefile / .gitignore / .env.example
├── backend/        Go 后端（cmd / internal / data/seed / migrations）
├── frontend/       React 前端（src/pages|components|api|store）
└── docs/
    ├── design/     详细设计（agent 据此写代码）
    ├── agent/      子约束文件
    ├── arch/       架构文档
    └── product/    产品文档
```
完整文件树见 `docs/arch/directory-structure.md`。

## 4. 五条铁律（违反即视为错误）
1. **不引入 Redis / Kafka / Elasticsearch / 向量库 / 微服务网关**。MVP 用单体后端 + MySQL。
2. **AI API Key 只存在于后端环境变量**，前端永远不得出现 Key，前端不得直连任何 AI API。
3. **保持单体轻量**：后端不拆微服务，分层为 api → service → repository。
4. **质量优先于数量**：演示版城市内容做 35 个精选城市，每城 1~2 地标 / 1~2 美食 / 1 人物；seed 校验允许 12~100 城。
5. **以契约为准**：任何接口改动必须先改 `docs/design/api-contract.md`，前后端均以它为唯一真相源。

## 5. 如何启动
```bash
cp .env.example .env      # 填入数据库与 AI 的真实配置
docker compose up -d      # 一键启动 app + mysql (+ 可选 caddy)
# 详见 README.md
```

## 6. 写代码前必读的子约束（按你要改的部分）
| 你要改 | 先读 |
|---|---|
| 任何 `.go` | `docs/agent/go-backend-rules.md` |
| 任何前端文件 | `docs/agent/react-frontend-rules.md` |
| 任何 `.sql` / 表结构 | `docs/agent/sql-rules.md` |
| 任何 AI 相关代码 | `docs/agent/ai-integration-rules.md` |
| 提交代码 | `docs/agent/git-workflow.md` |
| 写文档 | `docs/agent/doc-writing.md` |

## 7. 写代码前必读的设计文档
- 接口：`docs/design/api-contract.md`
- 后端文件职责：`docs/design/backend-detailed-design.md`
- 前端文件职责：`docs/design/frontend-detailed-design.md`
- 数据库：`docs/design/database-detailed-design.md`
- 算法/AI/成就：`docs/design/geo-algorithm-detailed-design.md`、`ai-orchestration-detailed-design.md`、`achievement-engine-detailed-design.md`

## 8. 禁止事项清单
- ❌ 在 handler 里写业务逻辑（必须在 service）。
- ❌ 在前端组件里直接 fetch / 直接调高德 SDK / 存 Key。
- ❌ 部署本地大模型或本地生图模型（2C2G 扛不住）。
- ❌ 提交 .env、密钥、uploads/ 下的生成图、二进制产物。
- ❌ 让 AI 声称自己是「真实复活的人物」或编造确定性史实。
- ❌ 生图使用在世名人肖像 / 生成色情暴力侮辱内容。
- ❌ 擅自新增中间件、改技术栈、扩展到 seed 校验上限以外。

## 9. 完成定义（DoD）
任一功能视为完成需满足 `docs/product/acceptance-criteria.md` 中对应验收条目，且 `docker compose up` 后无需手改数据库即可演示。
