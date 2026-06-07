# 详细设计文档索引

## 1. 阅读顺序

1. `api-contract.md`：前后端接口唯一契约。
2. `database-detailed-design.md`：表结构、枚举、索引、seed 导入规则。
3. `../arch/data-model-er.md`：实体关系。
4. `backend-detailed-design.md`：后端文件职责。
5. `frontend-detailed-design.md`：前端文件职责。
6. `geo-algorithm-detailed-design.md`：任意门随机漫游算法。
7. `ai-orchestration-detailed-design.md`：AI 对话与生图编排。
8. `achievement-engine-detailed-design.md`：成就规则与评估。

## 2. 文档清单

| 文档 | 内容 |
|---|---|
| `api-contract.md` | 路由、请求/响应、错误码。 |
| `database-detailed-design.md` | 17 张表设计、seed 策略。 |
| `backend-detailed-design.md` | Go 后端分层与文件职责。 |
| `frontend-detailed-design.md` | React 单页工作台与组件/store/API 职责。 |
| `geo-algorithm-detailed-design.md` | 随机方向、距离、目标点、最近城市匹配。 |
| `ai-orchestration-detailed-design.md` | LLM / 生图 client、prompt、任务与兜底。 |
| `achievement-engine-detailed-design.md` | rule_type 字典、UserStats、Evaluate。 |

## 3. 交叉约束

- 改接口：先改 `api-contract.md`，再改前后端。
- 改表：同步 `database-detailed-design.md` 与 `backend/migrations/schema.sql`。
- 改前端字段类型：同步 `frontend/src/api/types.ts` 与接口契约。
- 改产品范围：同步 `../product/prd.md`、`../product/user-flows.md`、`../product/acceptance-criteria.md`。
