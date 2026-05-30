"# Go 后端编码约束

> 写任何 `.go` 文件前必读。父约束：`CLAUDE.md`。

## 1. 分层职责（严格遵守）
```
api (handler)  →  只做：参数解析、校验、调用 service、封装响应。禁止写业务逻辑。
service        →  业务流程编排、事务边界、调用 repository / geo / ai / achievement。
repository     →  只做数据库读写（GORM）。禁止业务判断。
model          →  结构体定义 + GORM tag，对应数据库表。
geo/ai/upload/achievement → 领域工具包，被 service 调用。
middleware     →  cors / logger / recover / rate_limit。
config         →  Viper + 环境变量加载。
```
依赖方向单向：api → service → repository，不得反向或跨层（handler 不得直接调 repository）。

## 2. 命名与风格
- 包名全小写、无下划线；文件名小写下划线（game_handler.go）。
- 导出标识符大驼峰，私有小驼峰。
- 错误变量 err，自定义错误用 errors.New / fmt.Errorf(""...: %w"", err) 包装，保留 wrap 链。
- 遵循 gofmt / goimports，提交前 make lint 通过。

## 3. 错误处理
- service 返回 (result, error)，handler 统一捕获并转为标准错误响应（见 api-contract.md 错误码）。
- 不吞错误；不可恢复错误向上传递；panic 由 recover 中间件兜底。

## 4. GORM 使用规范
- 一律用 GORM 方法，**禁止裸写 SQL 字符串**（迁移 DDL 除外）。
- model 必须带 tag：gorm:""column:...;index""、json:""...""。
- 查询带 context；批量写入用事务（如打卡 + 成就判定）。
- 时间字段统一 created_at / updated_at（GORM 自动维护）。

## 5. 接口纪律
- **新增/修改任何接口前，先改 docs/design/api-contract.md**，再写代码。
- 请求/响应结构体字段与契约的 JSON 字段名一一对应。

## 6. 日志
- 使用 slog（或 zap），结构化输出。关键链路打日志：游戏掷骰、AI 调用、生图、成就解锁、错误。
- 日志类型见 docs/arch/observability.md。

## 7. 配置与密钥
- 所有配置经 internal/config 读取；敏感项（DB 密码、AI Key）只从环境变量来，禁止硬编码。

## 8. 提交前自检
- [ ] 分层正确，handler 无业务逻辑
- [ ] 改了接口 → 已同步 api-contract.md
- [ ] make lint 通过，无裸 SQL，无硬编码密钥"
