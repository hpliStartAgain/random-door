# Go 后端编码约束

> 写任何 `.go` 文件前必读。父约束：`AGENTS.md` / `CLAUDE.md`。

## 1. 分层职责

```text
api          handler：参数解析、校验、调用 service、封装响应
service      业务流程编排、事务边界、调用 repository / geo / ai / achievement
repository   只做 GORM 数据访问
model        结构体定义 + GORM tag，对应数据库表
geo/ai/upload/achievement/seed 领域包，被 service 或 cmd 调用
middleware   cors / logger / recover / rate_limit
config       Viper + 环境变量加载
```

依赖方向单向：`api -> service -> repository`。handler 不得直接调 repository。

## 2. 命名与风格

- 包名全小写、无下划线；文件名小写下划线。
- 导出标识符大驼峰，私有标识符小驼峰。
- 错误用 `fmt.Errorf("...: %w", err)` 保留 wrap 链。
- 提交前运行 gofmt；仓库命令为 `make lint` / `make test`。

## 3. 错误处理

- service 返回 `(result, error)`，handler 统一转为标准错误响应。
- 不吞错误；panic 由 recover 中间件兜底。
- 新错误码必须同步 `docs/design/api-contract.md`。

## 4. GORM 使用规范

- 业务代码使用 GORM，迁移 DDL 除外。
- model 必须带 `gorm:"column:...;index"` 与 `json:"..."`。
- 查询带 context；跨表写入用事务，例如掷骰 + 访问、打卡 + 成就。
- 时间字段统一 `created_at` / `updated_at`。

## 5. 接口纪律

- 新增或修改接口前，先改 `docs/design/api-contract.md`。
- 请求/响应结构体字段与契约 JSON 字段名一一对应。

## 6. 日志

- 使用结构化日志，关键链路至少覆盖：游戏掷骰、AI 调用、生图任务、成就解锁、后台导入、错误。
- 当前仓库未单独维护 observability 文档；新增日志规范时同步补文档。

## 7. 配置与密钥

- 所有配置经 `internal/config` 读取。
- DB 密码、AI Key、ADMIN_TOKEN 只从环境变量来，禁止硬编码。

## 8. 提交前自检

- [ ] 分层正确，handler 无业务逻辑。
- [ ] 改了接口则已同步 `api-contract.md`。
- [ ] 改了表则已同步数据库设计与 schema。
- [ ] `make lint` / `make test` 通过或已记录无法执行原因。
