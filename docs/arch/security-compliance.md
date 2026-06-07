# 安全与合规设计

## 1. 上传安全

| 项 | 规则 | 落点 |
|---|---|---|
| 文件类型 | jpg / jpeg / png / webp | `internal/upload/validator.go` |
| 文件大小 | 默认不超过 5MB | `UPLOAD_MAX_SIZE_MB` |
| 文件名 | UUID 重命名，不信任原文件名 | `internal/upload/storage.go` |
| 路径 | 防路径穿越，只暴露本地 `/static` / `/uploads` | storage / router |
| 目录 | 自拍、场景、生成图、后台导入分目录保存 | `uploads/*` |

## 2. AI 对话合规

- AI 不得声称自己是真实复活的人物。
- 不编造无法确认的确定性史实。
- 方言表达必须附普通话解释。
- 避免地域刻板印象。
- prompt 模板见 `../design/ai-orchestration-detailed-design.md`。

## 3. 生图合规

- 只接受用户主动上传的本人照片。
- 禁止在世公众人物肖像、色情、暴力、侮辱性内容、他人照片恶搞。
- 不做声音克隆。
- 生图任务受每日限额和最大重试次数控制。

## 4. API 与密钥

| 项 | 规则 |
|---|---|
| LLM / IMAGE Key | 只在后端环境变量中读取。 |
| 前端 | 只允许高德公开 Key 使用 `VITE_` 注入。 |
| Admin | `/api/admin/*` 需 `X-Admin-Token` 或 Bearer token。 |
| CORS | `CORS_ALLOW_ORIGINS` 白名单。 |
| AI 限流 | `AI_CHAT_DAILY_LIMIT` / `AI_IMAGE_DAILY_LIMIT`。 |
| 错误响应 | 统一错误结构，不暴露堆栈和密钥。 |

## 5. 红线速查

- 禁止前端出现 LLM / IMAGE Key。
- 禁止提交 `.env`、密钥、`uploads/` 运行时文件。
- 禁止 AI 自称真实复活人物。
- 禁止在世名人肖像生图和不合规内容生成。
- 禁止绕过后台 token 维护内容。

## 6. 代码责任点

| 能力 | 位置 |
|---|---|
| 上传校验 | `backend/internal/upload/validator.go` |
| 文件保存 | `backend/internal/upload/storage.go` |
| Prompt | `backend/internal/ai/prompt_builder.go` |
| AI client | `backend/internal/ai/llm_client.go`、`image_client.go` |
| CORS / 限流 / recover | `backend/internal/middleware/` |
| 后台鉴权 | `backend/internal/api/admin_handler.go` |
