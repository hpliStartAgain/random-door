# 详细设计文档索引（00-detailed-design-index.md）

> 这是所有详细设计文档的导航入口。coding agent 在写代码前，按本文件指引的顺序阅读对应文档。
> 父约束：CLAUDE.md。所有文档头部均标注对应概要设计章节号。

## 1. 阅读顺序（建议）
1. `api-contract.md` —— 接口契约（前后端唯一真相源，最先读）
2. `database-detailed-design.md` —— 12 张表逐字段
3. `../arch/data-model-er.md` —— ER 图（理解实体关系）
4. `backend-detailed-design.md` —— 后端逐文件职责（写 .go 依据）
5. `geo-algorithm-detailed-design.md` —— 随机漫游算法
6. `ai-orchestration-detailed-design.md` —— AI 对话/生图编排
7. `achievement-engine-detailed-design.md` —— 成就引擎
8. `frontend-detailed-design.md` —— 前端逐文件职责（写 .ts 依据）

## 2. 文档清单与对应概要章节
| 文档 | 内容 | 概要章节 |
|---|---|---|
| api-contract.md | 10 个接口请求/响应/错误码 | 16 |
| database-detailed-design.md | 12 张表设计 | 14、15 |
| backend-detailed-design.md | 后端文件树 + 每文件职责/函数签名 | 12 |
| frontend-detailed-design.md | 前端文件树 + 每文件职责 | 13 |
| geo-algorithm-detailed-design.md | 8 方向/6 距离/目标点/最近城市/兜底 | 17 |
| ai-orchestration-detailed-design.md | 对话 Prompt/生图 Prompt/超时兜底 | 19、20 |
| achievement-engine-detailed-design.md | 规则类型/判定流程 | 18 |

## 3. 关键交叉约束
- 改任何接口 → 先改 api-contract.md，再改前后端代码。
- 改任何表 → 同步改 database-detailed-design.md 与 backend/migrations/schema.sql。
- 后端与前端的字段名以 api-contract.md 的 JSON 字段为准，三方对齐。

## 4. 文件创建依据
agent 需要创建哪些目录与文件，见 `../arch/directory-structure.md`（完整文件树）。本目录下的文档只描述"每个文件该写什么"，不含真实代码。