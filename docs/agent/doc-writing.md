# 文档编写约束

> 写任何项目文档前必读。当前代码与运行配置是事实标准，文档只记录已实现能力、明确约束和可执行计划。

## 1. 事实来源

- 接口事实以 `backend/internal/api/router.go`、handler、service 和 `frontend/src/api/` 为准。
- 数据库事实以 `backend/internal/model/`、`backend/migrations/schema.sql` 和 `docs/design/database-detailed-design.md` 对齐。
- 前端形态以 `frontend/src/App.tsx`、`components/`、`pages/`、`store/` 为准。
- 内容事实以数据库为准；仓库 seed 是受控初始化资料，不代表线上内容上限。

## 2. 命名口径

- 项目中文名统一为“任意门”。
- 可在技术上下文中保留仓库名 `random-door`。
- 不再把当前状态描述为旧验证阶段或临时版本。

## 3. 写作要求

- 明确区分“已实现”“运行约束”“待办”。
- 不引用不存在的文件；新增引用前先确认文件已存在。
- 不把历史对标、临时计划、一次性修复记录放入产品文档。
- 接口字段、路径、错误码必须与 `docs/design/api-contract.md` 一致。

## 4. 更新规则

- 改接口：先更新 `docs/design/api-contract.md`。
- 改表结构：同步 `docs/design/database-detailed-design.md` 与 `backend/migrations/schema.sql`。
- 改功能范围：同步 `docs/product/prd.md`、`docs/product/user-flows.md`、`docs/product/acceptance-criteria.md`。
- 发现未完备项：写入根目录 `TODO.md`，不要散落在临时文档中。

## 5. 提交前自检

- [ ] 没有旧项目名作为项目名出现。
- [ ] 没有把当前产品称为旧验证阶段版本。
- [ ] 没有引用不存在的文档。
- [ ] README、CHANGELOG、TODO 与 docs 口径一致。
